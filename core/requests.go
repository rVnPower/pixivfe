package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/i18n"
	"codeberg.org/vnpower/pixivfe/v2/server/audit"
	"codeberg.org/vnpower/pixivfe/v2/server/request_context"
	"codeberg.org/vnpower/pixivfe/v2/server/token_manager"
	"codeberg.org/vnpower/pixivfe/v2/server/utils"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/tidwall/gjson"
)

// SimpleHTTPResponse represents a simplified HTTP response
type SimpleHTTPResponse struct {
	StatusCode int
	Body       string
}

var retryClient *retryablehttp.Client

func init() {
	retryClient = retryablehttp.NewClient()
	retryClient.RetryMax = config.GlobalConfig.APIMaxRetries
	retryClient.RetryWaitMin = config.GlobalConfig.APIBaseTimeout
	retryClient.RetryWaitMax = config.GlobalConfig.APIMaxBackoffTime
	retryClient.HTTPClient = utils.HttpClient
	retryClient.Logger = nil // Disables the default logger in go-retryablehttp
}

// Helper function to handle common request logic
func makeRequest(ctx context.Context, reqFunc func(context.Context, string) (*retryablehttp.Request, error), token *token_manager.Token, url string) (*SimpleHTTPResponse, error) {
	req, err := reqFunc(ctx, token.Value)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	resp, err := retryClient.Do(req)
	end := time.Now()

	if err != nil {
		return nil, i18n.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	audit.LogAPIRoundTrip(audit.APIRequestSpan{
		StartTime: start,
		EndTime:   end,
		RequestId: request_context.GetFromContext(ctx).RequestId,
		Response:  resp,
		Error:     err,
		Method:    req.Method,
		Url:       url,
		Token:     token.Value,
		Body:      string(body),
	})

	return &SimpleHTTPResponse{
		StatusCode: resp.StatusCode,
		Body:       string(body),
	}, nil
}

// retryRequest performs a request with automatic retries and token management
func retryRequest(
	ctx context.Context,
	reqFunc func(context.Context, string) (*retryablehttp.Request, error),
	userToken string,
	isPost bool,
	url string, // used for logging in audit.LogAPIRoundTrip only
) (*SimpleHTTPResponse, error) {
	var lastErr error
	tokenManager := config.GlobalConfig.TokenManager

	if isPost {
		// For POST requests, perform the request once without retrying
		token := &token_manager.Token{Value: userToken}
		return makeRequest(ctx, reqFunc, token, url)
	}

	// For GET requests, use the retry logic
	for i := 0; i < config.GlobalConfig.APIMaxRetries; i++ {
		var token *token_manager.Token
		if userToken != "" {
			token = &token_manager.Token{Value: userToken}
		} else {
			token = tokenManager.GetToken()
		}

		if token == nil {
			tokenManager.ResetAllTokens()
			return nil, i18n.Errorf(
				`All tokens (%d) are timed out, resetting all tokens to their initial good state.
Consider providing additional tokens in PIXIVFE_TOKEN or reviewing API request level backoff configuration.
Please refer the following documentation for additional information:
- https://pixivfe-docs.pages.dev/hosting/obtaining-pixivfe-token/
- https://pixivfe-docs.pages.dev/hosting/environment-variables/#exponential-backoff-configuration`,
				len(config.GlobalConfig.Token))
		}

		resp, err := makeRequest(ctx, reqFunc, token, url)
		if err != nil {
			lastErr = err
			tokenManager.MarkTokenStatus(token, token_manager.TimedOut)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			tokenManager.MarkTokenStatus(token, token_manager.Good)
			return resp, nil
		}

		lastErr = i18n.Errorf("HTTP status code: %d", resp.StatusCode)
		tokenManager.MarkTokenStatus(token, token_manager.TimedOut)

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// Continue to next iteration
		}
	}

	return nil, i18n.Errorf("Max retries reached for GET request. Last error: %v", lastErr)
}

// API_GET performs a GET request to the Pixiv API with automatic retries
func API_GET(ctx context.Context, url string, userToken string) (*SimpleHTTPResponse, error) {
	return retryRequest(ctx, func(ctx context.Context, token string) (*retryablehttp.Request, error) {
		req, err := retryablehttp.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req = req.WithContext(ctx)
		req.Header.Add("User-Agent", config.GetRandomUserAgent())
		req.Header.Add("Accept-Language", config.GlobalConfig.AcceptLanguage)
		req.AddCookie(&http.Cookie{
			Name:  "PHPSESSID",
			Value: token,
		})
		return req, nil
	}, userToken, false, url)
}

// API_GET_UnwrapJson performs a GET request and unwraps the JSON response
func API_GET_UnwrapJson(ctx context.Context, url string, userToken string) (string, error) {
	resp, err := API_GET(ctx, url, userToken)
	if err != nil {
		return "", err
	}

	if !gjson.Valid(resp.Body) {
		return "", i18n.Errorf("Invalid JSON: %v", resp.Body)
	}

	err2 := gjson.Get(resp.Body, "error")

	if !err2.Exists() {
		return "", i18n.Error("Incompatible request body")
	}

	if err2.Bool() {
		return "", errors.New(gjson.Get(resp.Body, "message").String())
	}

	return gjson.Get(resp.Body, "body").String(), nil
}

// createMultipartFormData is a helper function to create multipart form data
func createMultipartFormData(fields map[string]string) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, value := range fields {
		err := writer.WriteField(key, value)
		if err != nil {
			return nil, "", err
		}
	}
	err := writer.Close()
	if err != nil {
		return nil, "", err
	}
	return body, writer.FormDataContentType(), nil
}

// API_POST performs a POST request to the Pixiv API
func API_POST(
	ctx context.Context,
	url string,
	payload interface{},
	userToken, csrf string,
	contentType string,
) (*SimpleHTTPResponse, error) {
	if userToken == "" {
		return nil, i18n.Error("userToken is required for POST requests")
	}

	resp, err := retryRequest(ctx, func(ctx context.Context, token string) (*retryablehttp.Request, error) {
		var req *retryablehttp.Request
		var err error

		switch v := payload.(type) {
		case string:
			req, err = retryablehttp.NewRequest("POST", url, bytes.NewBuffer([]byte(v)))
		case map[string]string:
			body, formContentType, err := createMultipartFormData(v)
			if err != nil {
				return nil, err
			}
			req, err = retryablehttp.NewRequest("POST", url, body)
			if err == nil {
				contentType = formContentType
			}
		default:
			return nil, i18n.Error("Unsupported payload type")
		}

		if err != nil {
			return nil, err
		}

		req = req.WithContext(ctx)
		req.Header.Add("User-Agent", config.GetRandomUserAgent())
		req.Header.Add("Accept", "application/json")
		req.Header.Add("x-csrf-token", csrf)
		req.AddCookie(&http.Cookie{
			Name:  "PHPSESSID",
			Value: token,
		})
		req.Header.Add("Content-Type", contentType)

		return req, nil
	}, userToken, true, url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, resp.Body)
	}

	return resp, nil
}

// ProxyRequest forwards an HTTP request to the target server and copies the response back
func ProxyRequest(w http.ResponseWriter, req *http.Request) error {
	resp, err := utils.HttpClient.Do(req)
	if err != nil {
		return i18n.Errorf("failed to proxy request: %w", err)
	}
	defer resp.Body.Close()

	// Copy response headers
	header := w.Header()
	for k, v := range resp.Header {
		header[k] = v
	}

	// Set the status code
	w.WriteHeader(resp.StatusCode)

	// Copy the body from the response to the original writer
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return i18n.Errorf("failed to copy response body: %w", err)
	}

	return nil
}

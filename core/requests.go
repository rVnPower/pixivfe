package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	config "codeberg.org/vnpower/pixivfe/v2/config"
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
}

// retryRequest performs a request with automatic retries and token management
func retryRequest(ctx context.Context, reqFunc func(context.Context, string) (*retryablehttp.Request, error), userToken string, isPost bool) (SimpleHTTPResponse, error) {
	var lastErr error
	tokenManager := config.GlobalConfig.TokenManager

	for i := 0; i < config.GlobalConfig.APIMaxRetries; i++ {
		var token *token_manager.Token
		if userToken != "" {
			token = &token_manager.Token{Value: userToken}
		} else if !isPost {
			token = tokenManager.GetToken()
		}

		if token == nil && !isPost {
			tokenManager.ResetAllTokens()
			return SimpleHTTPResponse{}, fmt.Errorf("All tokens (%d) are timed out, resetting all tokens to their initial good state.\n"+
				"Consider providing additional tokens in PIXIVFE_TOKEN or reviewing API request level backoff configuration.\n"+
				"Please refer the following documentation for additional information:\n"+
				"- https://pixivfe-docs.pages.dev/hosting/obtaining-pixivfe-token/\n"+
				"- https://pixivfe-docs.pages.dev/hosting/environment-variables/#exponential-backoff-configuration",
				len(config.GlobalConfig.Token))
		}

		tokenValue := ""
		if token != nil {
			tokenValue = token.Value
		}

		req, err := reqFunc(ctx, tokenValue)
		if err != nil {
			return SimpleHTTPResponse{}, err
		}

		start := time.Now()
		resp, err := retryClient.Do(req)
		end := time.Now()

		audit.LogAPIRoundTrip(audit.APIRequestSpan{
			RequestId: request_context.GetFromContext(ctx).RequestId,
			Response:  resp,
			Error:     err,
			Method:    req.Method,
			Token:     tokenValue,
			Body:      "",
			StartTime: start,
			EndTime:   end,
		})

		if err == nil && resp.StatusCode == http.StatusOK {
			if userToken == "" && !isPost {
				tokenManager.MarkTokenStatus(token, token_manager.Good)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return SimpleHTTPResponse{}, err
			}
			return SimpleHTTPResponse{
				StatusCode: resp.StatusCode,
				Body:       string(body),
			}, nil
		}

		lastErr = err
		if err == nil {
			lastErr = fmt.Errorf("HTTP status code: %d", resp.StatusCode)
		}

		if userToken == "" && !isPost {
			tokenManager.MarkTokenStatus(token, token_manager.TimedOut)
		}

		select {
		case <-ctx.Done():
			return SimpleHTTPResponse{}, ctx.Err()
		default:
			// Continue to next iteration
		}
	}

	return SimpleHTTPResponse{}, fmt.Errorf("Max retries reached. Last error: %v", lastErr)
}

// API_GET performs a GET request to the Pixiv API with automatic retries
func API_GET(ctx context.Context, url string, userToken string) (SimpleHTTPResponse, error) {
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
	}, userToken, false)
}

// API_GET_UnwrapJson performs a GET request and unwraps the JSON response
func API_GET_UnwrapJson(ctx context.Context, url string, userToken string) (string, error) {
	resp, err := API_GET(ctx, url, userToken)
	if err != nil {
		return "", err
	}

	if !gjson.Valid(resp.Body) {
		return "", fmt.Errorf("Invalid JSON: %v", resp.Body)
	}

	err2 := gjson.Get(resp.Body, "error")

	if !err2.Exists() {
		return "", errors.New("Incompatible request body")
	}

	if err2.Bool() {
		return "", errors.New(gjson.Get(resp.Body, "message").String())
	}

	return gjson.Get(resp.Body, "body").String(), nil
}

// API_POST performs a POST request to the Pixiv API with automatic retries
func API_POST(ctx context.Context, url, payload, userToken, csrf string, isJSON bool) error {
	if userToken == "" {
		return errors.New("userToken is required for POST requests")
	}

	_, err := retryRequest(ctx, func(ctx context.Context, token string) (*retryablehttp.Request, error) {
		req, err := retryablehttp.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
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
		if isJSON {
			req.Header.Add("Content-Type", "application/json; charset=utf-8")
		} else {
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
		}
		return req, nil
	}, userToken, true)

	return err
}

// ProxyRequest forwards an HTTP request to the target server and copies the response back
func ProxyRequest(w http.ResponseWriter, req *http.Request) error {
	resp, err := utils.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	header := w.Header()
	for k, v := range resp.Header {
		header[k] = v
	}

	_, err = io.Copy(w, resp.Body)
	return err
}

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

// ProxyJob represents a job for the proxy worker pool
type ProxyJob struct {
	Request  *http.Request
	Response http.ResponseWriter
	Done     chan error
}

var proxyJobQueue chan ProxyJob
var numWorkers = 10 // TODO: allow number of workers to be configured

func init() {
	retryClient = retryablehttp.NewClient()
	retryClient.RetryMax = config.GlobalConfig.APIMaxRetries
	retryClient.RetryWaitMin = config.GlobalConfig.APIBaseTimeout
	retryClient.RetryWaitMax = config.GlobalConfig.APIMaxBackoffTime
	retryClient.HTTPClient = utils.HttpClient

	// Initialize the proxy worker pool
	proxyJobQueue = make(chan ProxyJob, 100) // TODO: allow buffer size to be configured
	for i := 0; i < numWorkers; i++ {
		go proxyWorker(proxyJobQueue)
	}
}

// retryRequest performs a request with automatic retries and token management
func retryRequest(ctx context.Context, reqFunc func(context.Context, string) (*retryablehttp.Request, error), userToken string, isPost bool) (*SimpleHTTPResponse, error) {
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
			return nil, fmt.Errorf("All tokens (%d) are timed out, resetting all tokens to their initial good state.\n"+
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
			return nil, err
		}

		start := time.Now()
		resp, err := retryClient.Do(req)
		end := time.Now()
		// Unwrap the body here so that we could log stuff correctly
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		audit.LogAPIRoundTrip(audit.APIRequestSpan{
			RequestId: request_context.GetFromContext(ctx).RequestId,
			Response:  resp,
			Error:     err,
			Method:    req.Method,
			Token:     tokenValue,
			Body:      string(body),
			StartTime: start,
			EndTime:   end,
		})

		if resp.StatusCode == http.StatusOK {
			if userToken == "" && !isPost {
				tokenManager.MarkTokenStatus(token, token_manager.Good)
			}
			return &SimpleHTTPResponse{
				StatusCode: resp.StatusCode,
				Body:       string(body),
			}, nil
		}

		lastErr = fmt.Errorf("HTTP status code: %d", resp.StatusCode)

		if userToken == "" && !isPost {
			tokenManager.MarkTokenStatus(token, token_manager.TimedOut)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// Continue to next iteration
		}
	}

	return nil, fmt.Errorf("Max retries reached. Last error: %v", lastErr)
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
func ProxyRequest(w http.ResponseWriter, req *http.Request) {
	done := make(chan error, 1)
	job := ProxyJob{
		Request:  req,
		Response: w,
		Done:     done,
	}

	proxyJobQueue <- job

	err := <-done
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// proxyWorker processes jobs from the proxyJobQueue
func proxyWorker(jobs <-chan ProxyJob) {
	for job := range jobs {
		err := processProxyJob(job)
		job.Done <- err
	}
}

// processProxyJob handles the actual proxying of the request
func processProxyJob(job ProxyJob) error {
	resp, err := utils.HttpClient.Do(job.Request)
	if err != nil {
		return fmt.Errorf("failed to process request: %w", err)
	}
	defer resp.Body.Close()

	// Copy response headers
	header := job.Response.Header()
	for k, v := range resp.Header {
		header[k] = v
	}

	// Set the status code
	job.Response.WriteHeader(resp.StatusCode)

	// Copy the body from the response to the original writer
	_, err = io.Copy(job.Response, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy response body: %w", err)
	}

	return nil
}

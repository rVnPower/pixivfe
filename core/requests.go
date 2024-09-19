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
	"github.com/tidwall/gjson"
)

// SimpleHTTPResponse represents a simplified HTTP response
type SimpleHTTPResponse struct {
	StatusCode int
	Body       string
}

// retryRequest performs a request with automatic retries and token management
func retryRequest(ctx context.Context, reqFunc func(context.Context, string) (SimpleHTTPResponse, *http.Response, error)) (SimpleHTTPResponse, error) {
	var lastErr error
	tokenManager := config.GlobalConfig.TokenManager

	for i := 0; i < config.GlobalConfig.APIMaxRetries; i++ {
		// Get a token from the token manager
		token := tokenManager.GetToken()
		if token == nil {
			return SimpleHTTPResponse{}, errors.New("All tokens are timed out")
		}

		// Perform the request using the provided function
		res, resp, err := reqFunc(ctx, token.Value)
		if err == nil && res.StatusCode == http.StatusOK {
			tokenManager.MarkTokenStatus(token, token_manager.Good)
			return res, nil
		}

		// Handle errors and prepare for retry
		lastErr = err
		if err == nil {
			lastErr = fmt.Errorf("HTTP status code: %d", res.StatusCode)
		}

		tokenManager.MarkTokenStatus(token, token_manager.TimedOut)

		// Calculate backoff duration for exponential backoff
		backoffMilliseconds := config.GlobalConfig.APIBaseTimeout.Milliseconds() * (1 << uint(i))
		backoffDuration := time.Duration(backoffMilliseconds) * time.Millisecond
		if backoffDuration > config.GlobalConfig.APIMaxBackoffTime {
			backoffDuration = config.GlobalConfig.APIMaxBackoffTime
		}

		// Wait for backoff duration or context cancellation
		select {
		case <-ctx.Done():
			return SimpleHTTPResponse{}, ctx.Err()
		case <-time.After(backoffDuration):
		}

		// Log the API request for auditing purposes
		audit.LogAPIRoundTrip(audit.APIRequestSpan{
			RequestId: request_context.GetFromContext(ctx).RequestId,
			Response:  resp,
			Error:     err,
			Method:    "GET",
			Token:     token.Value,
			Body:      res.Body,
			StartTime: time.Now(),
			EndTime:   time.Now().Add(backoffDuration),
		})
	}

	return SimpleHTTPResponse{}, fmt.Errorf("Max retries reached. Last error: %v", lastErr)
}

// API_GET performs a GET request to the Pixiv API with automatic retries
func API_GET(ctx context.Context, url string, _ string) (SimpleHTTPResponse, error) {
	return retryRequest(ctx, func(ctx context.Context, token string) (SimpleHTTPResponse, *http.Response, error) {
		return _API_GET(ctx, url, token)
	})
}

// _API_GET is the internal function to perform a GET request to the Pixiv API
func _API_GET(ctx context.Context, url string, token string) (SimpleHTTPResponse, *http.Response, error) {
	var res SimpleHTTPResponse

	// Create a new request with the provided context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return res, nil, err
	}

	// Set headers for the request
	req.Header.Add("User-Agent", config.GlobalConfig.UserAgent)
	req.Header.Add("Accept-Language", config.GlobalConfig.AcceptLanguage)

	// Add the token as a cookie
	req.AddCookie(&http.Cookie{
		Name:  "PHPSESSID",
		Value: token,
	})

	// Make the request
	resp, err := utils.HttpClient.Do(req)
	if err != nil {
		return res, nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, resp, err
	}

	// Construct the SimpleHTTPResponse
	res = SimpleHTTPResponse{
		StatusCode: resp.StatusCode,
		Body:       string(body),
	}
	return res, resp, nil
}

// API_GET_UnwrapJson performs a GET request and unwraps the JSON response
func API_GET_UnwrapJson(ctx context.Context, url string, _ string) (string, error) {
	resp, err := API_GET(ctx, url, "")
	if err != nil {
		return "", err
	}

	// Validate JSON response
	if !gjson.Valid(resp.Body) {
		return "", fmt.Errorf("Invalid JSON: %v", resp.Body)
	}

	// Check for errors in the JSON response
	err2 := gjson.Get(resp.Body, "error")

	if !err2.Exists() {
		return "", errors.New("Incompatible request body")
	}

	if err2.Bool() {
		return "", errors.New(gjson.Get(resp.Body, "message").String())
	}

	// Return the "body" field from the JSON response
	return gjson.Get(resp.Body, "body").String(), nil
}

// API_POST performs a POST request to the Pixiv API with automatic retries
func API_POST(ctx context.Context, url, payload, _, csrf string, isJSON bool) error {
	var lastErr error
	for i := 0; i < config.GlobalConfig.APIMaxRetries; i++ {
		// Get a token from the token manager
		token := config.GlobalConfig.TokenManager.GetToken()
		if token == nil {
			return errors.New("All tokens are timed out")
		}

		// Perform the POST request
		start_time := time.Now()
		resp, err := _API_POST(ctx, url, payload, token.Value, csrf, isJSON)
		end_time := time.Now()

		// Log the API request for auditing purposes
		audit.LogAPIRoundTrip(audit.APIRequestSpan{
			RequestId: request_context.GetFromContext(ctx).RequestId,
			Response:  resp,
			Error:     err,
			Method:    "POST",
			Url:       url,
			Token:     token.Value,
			Body:      "",
			StartTime: start_time,
			EndTime:   end_time,
		})

		// Check for successful response
		if err == nil && resp.StatusCode == http.StatusOK {
			config.GlobalConfig.TokenManager.MarkTokenStatus(token, token_manager.Good)
			return nil
		}

		// Handle errors and prepare for retry
		lastErr = err
		if err == nil {
			lastErr = fmt.Errorf("HTTP status code: %d", resp.StatusCode)
		}

		config.GlobalConfig.TokenManager.MarkTokenStatus(token, token_manager.TimedOut)

		// Calculate backoff duration for exponential backoff
		backoffMilliseconds := config.GlobalConfig.APIBaseTimeout.Milliseconds() * (1 << uint(i))
		backoffDuration := time.Duration(backoffMilliseconds) * time.Millisecond
		if backoffDuration > config.GlobalConfig.APIMaxBackoffTime {
			backoffDuration = config.GlobalConfig.APIMaxBackoffTime
		}

		// Wait for backoff duration or context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoffDuration):
		}
	}

	return fmt.Errorf("Max retries reached. Last error: %v", lastErr)
}

// _API_POST is the internal function to perform a POST request to the Pixiv API
func _API_POST(ctx context.Context, url, payload, token, csrf string, isJSON bool) (*http.Response, error) {
	requestBody := []byte(payload)

	// Create a new POST request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	// Set common headers
	req.Header.Add("User-Agent", "Mozilla/5.0")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-csrf-token", csrf)
	req.AddCookie(&http.Cookie{
		Name:  "PHPSESSID",
		Value: token,
	})

	// Set content type header based on isJSON flag
	if isJSON {
		req.Header.Add("Content-Type", "application/json; charset=utf-8")
	} else {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	}

	// Perform the request
	resp, err := utils.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read and validate the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}
	body_s := string(body)
	if !gjson.Valid(body_s) {
		return resp, fmt.Errorf("Invalid JSON: %v", body_s)
	}
	err2 := gjson.Get(body_s, "error")

	if !err2.Exists() {
		return resp, fmt.Errorf("Incompatible request body.")
	}

	if err2.Bool() {
		return resp, fmt.Errorf("Pixiv: Invalid request.")
	}
	return resp, nil
}

// ProxyRequest forwards an HTTP request to the target server and copies the response back
func ProxyRequest(w http.ResponseWriter, req *http.Request) error {
	// Make the request to the target server
	resp, err := utils.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Copy headers from the target server's response to our response
	header := w.Header()
	for k, v := range resp.Header {
		header[k] = v
	}

	// Copy the response body from the target server to our response
	_, err = io.Copy(w, resp.Body)
	return err
}

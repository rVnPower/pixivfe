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

type SimpleHTTPResponse struct {
	StatusCode int
	Body       string
}

func retryRequest(ctx context.Context, reqFunc func(context.Context, string) (SimpleHTTPResponse, *http.Response, error)) (SimpleHTTPResponse, error) {
	var lastErr error
	tokenManager := config.GlobalConfig.TokenManager

	for i := 0; i < tokenManager.GetMaxRetries(); i++ {
		token := tokenManager.GetToken()
		if token == nil {
			return SimpleHTTPResponse{}, errors.New("All tokens are timed out")
		}

		res, resp, err := reqFunc(ctx, token.Value)
		if err == nil && res.StatusCode == http.StatusOK {
			tokenManager.MarkTokenStatus(token, token_manager.Good)
			return res, nil
		}

		lastErr = err
		if err == nil {
			lastErr = fmt.Errorf("HTTP status code: %d", res.StatusCode)
		}

		tokenManager.MarkTokenStatus(token, token_manager.TimedOut)

		backoffDuration := tokenManager.GetBaseTimeout() * time.Duration(1<<uint(i))
		if backoffDuration > tokenManager.GetMaxBackoffTime() {
			backoffDuration = tokenManager.GetMaxBackoffTime()
		}

		select {
		case <-ctx.Done():
			return SimpleHTTPResponse{}, ctx.Err()
		case <-time.After(backoffDuration):
		}

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

// send GET
func API_GET(ctx context.Context, url string, _ string) (SimpleHTTPResponse, error) {
	return retryRequest(ctx, func(ctx context.Context, token string) (SimpleHTTPResponse, *http.Response, error) {
		return _API_GET(ctx, url, token)
	})
}

func _API_GET(ctx context.Context, url string, token string) (SimpleHTTPResponse, *http.Response, error) {
	var res SimpleHTTPResponse

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return res, nil, err
	}

	req.Header.Add("User-Agent", config.GlobalConfig.UserAgent)
	req.Header.Add("Accept-Language", config.GlobalConfig.AcceptLanguage)

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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, resp, err
	}

	res = SimpleHTTPResponse{
		StatusCode: resp.StatusCode,
		Body:       string(body),
	}
	return res, resp, nil
}

func API_GET_UnwrapJson(ctx context.Context, url string, _ string) (string, error) {
	resp, err := API_GET(ctx, url, "")
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

// send POST
func API_POST(ctx context.Context, url, payload, _, csrf string, isJSON bool) error {
	tokenManager := config.GlobalConfig.TokenManager

	var lastErr error
	for i := 0; i < tokenManager.GetMaxRetries(); i++ {
		token := tokenManager.GetToken()
		if token == nil {
			return errors.New("All tokens are timed out")
		}

		start_time := time.Now()
		resp, err := _API_POST(ctx, url, payload, token.Value, csrf, isJSON)
		end_time := time.Now()

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

		if err == nil && resp.StatusCode == http.StatusOK {
			tokenManager.MarkTokenStatus(token, token_manager.Good)
			return nil
		}

		lastErr = err
		if err == nil {
			lastErr = fmt.Errorf("HTTP status code: %d", resp.StatusCode)
		}

		tokenManager.MarkTokenStatus(token, token_manager.TimedOut)

		backoffDuration := tokenManager.GetBaseTimeout() * time.Duration(1<<uint(i))
		if backoffDuration > tokenManager.GetMaxBackoffTime() {
			backoffDuration = tokenManager.GetMaxBackoffTime()
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoffDuration):
		}
	}

	return fmt.Errorf("Max retries reached. Last error: %v", lastErr)
}

func _API_POST(ctx context.Context, url, payload, token, csrf string, isJSON bool) (*http.Response, error) {
	requestBody := []byte(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0")
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

	resp, err := utils.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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

func ProxyRequest(w http.ResponseWriter, req *http.Request) error {
	// Make the request
	resp, err := utils.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// copy headers
	header := w.Header()
	for k, v := range resp.Header {
		header[k] = v
	}

	// copy body
	_, err = io.Copy(w, resp.Body)
	return err
}

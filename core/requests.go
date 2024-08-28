package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"codeberg.org/vnpower/pixivfe/v2/audit"
	config "codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/utils"
	"github.com/tidwall/gjson"
)

type SimpleHTTPResponse struct {
	StatusCode int
	Body       string
}

// send GET
func API_GET(context context.Context, url string, token string) (SimpleHTTPResponse, error) {
	start_time := time.Now()
	res, resp, err := _API_GET(context, url, token)
	end_time := time.Now()
	audit.LogAPIRoundTrip(context, audit.APIPerformance{Response: resp, Error: err, Method: "GET", Url: url, Token: token, Body: res.Body, StartTime: start_time, EndTime: end_time})
	if err != nil {
		return SimpleHTTPResponse{}, fmt.Errorf("While GET %s: %w", url, err)
	}
	return res, nil
}

func _API_GET(context context.Context, url string, token string) (SimpleHTTPResponse, *http.Response, error) {
	var res SimpleHTTPResponse

	req, err := http.NewRequestWithContext(context, "GET", url, nil)
	if err != nil {
		return res, nil, err
	}

	req.Header.Add("User-Agent", config.GlobalConfig.UserAgent)
	req.Header.Add("Accept-Language", config.GlobalConfig.AcceptLanguage)

	if token == "" {
		req.AddCookie(&http.Cookie{
			Name:  "PHPSESSID",
			Value: config.GetRandomDefaultToken(),
		})
	} else {
		req.AddCookie(&http.Cookie{
			Name:  "PHPSESSID",
			Value: token,
		})
	}

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

func API_GET_UnwrapJson(context context.Context, url string, token string) (string, error) {
	resp, err := API_GET(context, url, token)
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
func API_POST(context context.Context, url, payload, token, csrf string, isJSON bool) error {
	start_time := time.Now()
	resp, err := _API_POST(context, url, payload, token, csrf, isJSON)
	end_time := time.Now()
	audit.LogAPIRoundTrip(context, audit.APIPerformance{Response: resp, Error: err, Method: "POST", Url: url, Token: token, Body: "", StartTime: start_time, EndTime: end_time})
	if err != nil {
		return fmt.Errorf("While POST %s: %w", url, err)
	}
	return err
}

func _API_POST(context context.Context, url, payload, token, csrf string, isJSON bool) (*http.Response, error) {
	requestBody := []byte(payload)

	req, err := http.NewRequestWithContext(context, "POST", url, bytes.NewBuffer(requestBody))
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

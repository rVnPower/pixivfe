package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	config "codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/utils"
	"github.com/tidwall/gjson"
)

type SimpleHTTPResponse struct {
	Ok         bool
	StatusCode int
	Body    string
	Message string
}

const DevDir_Response = "/tmp/pixivfe-dev/resp"

func CreateResponseAuditFolder() error {
	// err := os.RemoveAll(DevDir_Response)
	// if err != nil {
	// 	log.Println(err)
	// }

	return os.MkdirAll(DevDir_Response, 0o700)
}

func logResponseBody(body string) (string, error) {
	filename := path.Join(DevDir_Response, time.Now().UTC().Format(time.RFC3339Nano))
	err := os.WriteFile(filename, []byte(body), 0o600)
	if err != nil {
		return "", err
	}
	return filename, nil
}

// send GET
func PixivGetRequest(context context.Context, URL, token string) (SimpleHTTPResponse, error) {
	resp, err := _PixivGETRequest(context, URL, token)
	if err != nil {
		return SimpleHTTPResponse{}, fmt.Errorf("While sending request to %s: %w", URL, err)
	}
	if config.GlobalServerConfig.InDevelopment {
		if resp.Ok {
			filename, err := logResponseBody(resp.Body)
			if err != nil {
				log.Println(err)
			}
			if !(300 > resp.StatusCode && resp.StatusCode >= 200) {
				log.Println("(WARN) non-2xx response from pixiv:")
			}
			log.Println("->", URL, "->", resp.StatusCode, filename)
		} else {
			log.Println("->", URL, "ERR", resp.Message)
		}
	}
	return resp, nil
}

func _PixivGETRequest(context context.Context, URL, token string) (SimpleHTTPResponse, error) {
	req, err := http.NewRequestWithContext(context, "GET", URL, nil)
	if err != nil {
		return SimpleHTTPResponse{}, err
	}

	req.Header.Add("User-Agent", config.GlobalServerConfig.UserAgent)
	req.Header.Add("Accept-Language", config.GlobalServerConfig.AcceptLanguage)

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
		return SimpleHTTPResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SimpleHTTPResponse{}, err
	}

	resp2 := SimpleHTTPResponse{
		Ok:         true,
		StatusCode: resp.StatusCode,
		Body:       string(body),
		Message:    "",
	}

	return resp2, nil
}

func UnwrapWebAPIRequest(context context.Context, URL, token string) (string, error) {
	resp, err := PixivGetRequest(context, URL, token)
	if err != nil {
		return "", err
	}

	if !resp.Ok {
		return "", errors.New(resp.Message)
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
func PixivPostRequest(r *http.Request, url, payload, token, csrf string, isJSON bool) error {
	requestBody := []byte(payload)

	req, err := http.NewRequestWithContext(r.Context(), "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
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
		return errors.New("Failed to do this action.")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Cannot parse the response from Pixiv. Please report this issue.")
	}
	body_s := string(body)
	if !gjson.Valid(body_s) {
		return fmt.Errorf("Invalid JSON: %v", body_s)
	}
	errr := gjson.Get(body_s, "error")

	if !errr.Exists() {
		return errors.New("Incompatible request body.")
	}

	if errr.Bool() {
		return errors.New("Pixiv: Invalid request.")
	}
	return nil
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

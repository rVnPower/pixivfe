package core

import (
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

type HttpResponse struct {
	Ok         bool
	StatusCode int

	// @iacore: this not being []byte might come back to bite us
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

func WebAPIRequest(context context.Context, URL, token string) HttpResponse {
	resp := webAPIRequest(context, URL, token)
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
	return resp
}

func webAPIRequest(context context.Context, URL, token string) HttpResponse {
	req, err := http.NewRequestWithContext(context, "GET", URL, nil)
	if err != nil {
		return HttpResponse{
			Ok:         false,
			StatusCode: 0,
			Body:       "",
			Message:    fmt.Sprintf("Failed to create a request to %s\n.", URL),
		}
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
		return HttpResponse{
			Ok:         false,
			StatusCode: 0,
			Body:       "",
			Message:    fmt.Sprintf("Failed to send a request to %s\n.", URL),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return HttpResponse{
			Ok:         false,
			StatusCode: 0,
			Body:       "",
			Message:    fmt.Sprintln("Failed to parse request data."),
		}
	}

	resp2 := HttpResponse{
		Ok:         true,
		StatusCode: resp.StatusCode,
		Body:       string(body),
		Message:    "",
	}

	return resp2
}

func UnwrapWebAPIRequest(context context.Context, URL, token string) (string, error) {
	resp := WebAPIRequest(context, URL, token)

	if !resp.Ok {
		return "", errors.New(resp.Message)
	}
	if !gjson.Valid(resp.Body) {
		return "", fmt.Errorf("Invalid JSON: %v", resp.Body)
	}

	err := gjson.Get(resp.Body, "error")

	if !err.Exists() {
		return "", errors.New("Incompatible request body")
	}

	if err.Bool() {
		return "", errors.New(gjson.Get(resp.Body, "message").String())
	}

	return gjson.Get(resp.Body, "body").String(), nil
}

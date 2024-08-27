package routes

import (
	"fmt"
	"io"
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/utils"
)

func copyRequest(w http.ResponseWriter, req *http.Request) error {
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

func SPximgProxy(w http.ResponseWriter, r *http.Request) error {
	URL := fmt.Sprintf("https://s.pximg.net/%s", r.URL.Path)
	req, err := http.NewRequestWithContext(r.Context(), "GET", URL, nil)
	if err != nil {
		return err
	}
	return copyRequest(w, req)
}

func IPximgProxy(w http.ResponseWriter, r *http.Request) error {
	URL := fmt.Sprintf("https://i.pximg.net/%s", r.URL.Path)
	req, err := http.NewRequestWithContext(r.Context(), "GET", URL, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Referer", "https://www.pixiv.net/")
	return copyRequest(w, req)
}

func UgoiraProxy(w http.ResponseWriter, r *http.Request) error {
	URL := fmt.Sprintf("https://ugoira.com/api/mp4/%s", r.URL.Path)
	req, err := http.NewRequestWithContext(r.Context(), "GET", URL, nil)
	if err != nil {
		return err
	}
	return copyRequest(w, req)
}

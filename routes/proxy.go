package routes

import (
	"fmt"
	"io"
	"net/http"
)

func makeRequest(w http.ResponseWriter, req *http.Request) error {
	// Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))

	_, err = io.Copy(w, resp.Body)
	return err
}

func SPximgProxy(w http.ResponseWriter, r CompatRequest) error {
	URL := fmt.Sprintf("https://s.pximg.net/%s", r.Params("*"))
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}
	return makeRequest(w, req.WithContext(r.Context()))
}

func IPximgProxy(w http.ResponseWriter, r CompatRequest) error {
	URL := fmt.Sprintf("https://i.pximg.net/%s", r.Params("*"))
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Referer", "https://www.pixiv.net/")
	return makeRequest(w, req.WithContext(r.Context()))
}

func UgoiraProxy(w http.ResponseWriter, r CompatRequest) error {
	URL := fmt.Sprintf("https://ugoira.com/api/mp4/%s", r.Params("*"))
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}
	return makeRequest(w, req.WithContext(r.Context()))
}

package routes

import (
	"fmt"
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/core"
)

func SPximgProxy(w http.ResponseWriter, r *http.Request) error {
	URL := fmt.Sprintf("https://s.pximg.net/%s", r.URL.Path)
	req, err := http.NewRequestWithContext(r.Context(), "GET", URL, nil)
	if err != nil {
		return err
	}
	core.ProxyRequest(w, req)
	return nil
}

func IPximgProxy(w http.ResponseWriter, r *http.Request) error {
	URL := fmt.Sprintf("https://i.pximg.net/%s", r.URL.Path)
	req, err := http.NewRequestWithContext(r.Context(), "GET", URL, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Referer", "https://www.pixiv.net/")
	core.ProxyRequest(w, req)
	return nil
}

func UgoiraProxy(w http.ResponseWriter, r *http.Request) error {
	URL := fmt.Sprintf("https://ugoira.com/api/mp4/%s", r.URL.Path)
	req, err := http.NewRequestWithContext(r.Context(), "GET", URL, nil)
	if err != nil {
		return err
	}
	core.ProxyRequest(w, req)
	return nil
}

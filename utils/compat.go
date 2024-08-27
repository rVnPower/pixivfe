package utils

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

type CompatRequest struct {
	*http.Request
}

func (r CompatRequest) BaseURL() string {
	return (&url.URL{
		Scheme: r.URL.Scheme,
		Opaque: r.URL.Opaque,
		User:   r.URL.User,
		Host:   r.URL.Host,
	}).String()
}
func (r CompatRequest) OriginalURL() string {
	return (&url.URL{
		Path:        r.URL.Path,
		RawPath:     r.URL.RawPath,
		OmitHost:    r.URL.OmitHost,
		ForceQuery:  r.URL.ForceQuery,
		RawQuery:    r.URL.RawQuery,
		Fragment:    r.URL.Fragment,
		RawFragment: r.URL.RawFragment,
	}).String()
}
func (r CompatRequest) PageURL() string {
	return r.URL.String()
}

func (r CompatRequest) Query(name string, defaultValue ...string) string {
	if v := r.URL.Query().Get(name); v != "" {
		return v
	} else {
		if len(defaultValue) == 0 {
			return ""
		} else {
			return defaultValue[0]
		}
	}
}

// get path segment. no idea why it's called "params"
func (r CompatRequest) Params(name string, defaultValue ...string) string {
	if v := mux.Vars(r.Request)[name]; v != "" {
		return v
	} else {
		if len(defaultValue) == 0 {
			return ""
		} else {
			return defaultValue[0]
		}
	}
}

func SendString(w http.ResponseWriter, text string) error {
	w.Header().Set("content-type", "text/plain")
	_, err :=  w.Write([]byte(text))
	return err
}

func RedirectToRoute(w http.ResponseWriter, r CompatRequest, path string, query_params map[string]string, code int) error {
	query := url.Values{}
	for k, v := range query_params {
		query.Add(k, v)
	}
	http.Redirect(w, r.Request, path+query.Encode(), code)
	return nil
}

func WriteJson(w http.ResponseWriter, data any) {
	json.NewEncoder(w).Encode(data)
}
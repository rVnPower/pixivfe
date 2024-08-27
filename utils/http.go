package utils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)


func SendString(w http.ResponseWriter, text string) error {
	w.Header().Set("content-type", "text/plain")
	_, err :=  w.Write([]byte(text))
	return err
}

func SendJson(w http.ResponseWriter, data any) {
	json.NewEncoder(w).Encode(data)
}

func RedirectTo(w http.ResponseWriter, r *http.Request, path string, query_params map[string]string) error {
	query := url.Values{}
	for k, v := range query_params {
		query.Add(k, v)
	}
	http.Redirect(w, r, path+query.Encode(), http.StatusSeeOther)
	return nil
}

func RedirectToWhenceYouCame(w http.ResponseWriter, r *http.Request) {
	referrer := r.Referer()
	if strings.HasPrefix(referrer, Origin(r)) {
		http.Redirect(w, r, referrer, http.StatusSeeOther)
	} else {
		w.WriteHeader(200)
	}
}

func Origin(r *http.Request) string {
	return (&url.URL{
		Scheme: r.URL.Scheme,
		Opaque: r.URL.Opaque,
		User:   r.URL.User,
		Host:   r.URL.Host,
	}).String()
}

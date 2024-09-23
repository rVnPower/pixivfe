package utils

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/goccy/go-json"
)

// SendString writes a plain text response to the provided http.ResponseWriter.
// It sets the content type to "text/plain" and returns any error encountered during writing.
func SendString(w http.ResponseWriter, text string) error {
	w.Header().Set("content-type", "text/plain")
	_, err := w.Write([]byte(text))
	return err
}

// SendJson encodes the provided data as JSON and writes it to the http.ResponseWriter.
// It automatically sets the appropriate content type header.
func SendJson(w http.ResponseWriter, data any) {
	json.NewEncoder(w).Encode(data)
}

// RedirectTo performs a redirect to the specified path with optional query parameters.
// It uses HTTP status 303 (See Other) for the redirect.
func RedirectTo(w http.ResponseWriter, r *http.Request, path string, query_params map[string]string) error {
	query := url.Values{}
	for k, v := range query_params {
		query.Add(k, v)
	}
	http.Redirect(w, r, path+"?"+query.Encode(), http.StatusSeeOther)
	return nil
}

// RedirectToWhenceYouCame redirects the user back to the referring page if it's from the same origin.
// This helps prevent open redirects by checking the referrer against the current origin.
// If the referrer is not from the same origin, it responds with a 200 OK status.
func RedirectToWhenceYouCame(w http.ResponseWriter, r *http.Request) {
	referrer := r.Referer()
	if strings.HasPrefix(referrer, Origin(r)) {
		http.Redirect(w, r, referrer, http.StatusSeeOther)
	} else {
		w.WriteHeader(200)
	}
}

// Origin extracts the origin (scheme and host) from the given request's URL.
// This is useful for comparing against referrers or constructing absolute URLs.
func Origin(r *http.Request) string {
	return (&url.URL{
		Scheme: r.URL.Scheme,
		Opaque: r.URL.Opaque,
		User:   r.URL.User,
		Host:   r.URL.Host,
	}).String()
}

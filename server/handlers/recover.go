package handlers

import (
	"errors"
	"net/http"
)

// RecoverFromPanic wraps an http.Handler and recovers from any panics that occur during its execution.
// If a panic occurs, it sends an HTTP 500 Internal Server Error response with the panic message.
func RecoverFromPanic(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}

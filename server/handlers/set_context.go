package handlers

import (
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/server/request_context"
)

// ProvideUserContext is a middleware that wraps an http.Handler
// to inject a user context into each incoming request.
func ProvideUserContext(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r.WithContext(request_context.ProvideWith(r.Context())))
	})
}

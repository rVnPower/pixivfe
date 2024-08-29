package handlers

import (
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/request_context"
)

func ProvideUserContext(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r.WithContext(request_context.ProvideWith(r.Context())))
	})
}

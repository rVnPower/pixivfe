package handlers

import (
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/handlers/user_context"
)

type UserContext = user_context.UserContext

func GetUserContext(r *http.Request) *UserContext {
	return user_context.GetUserContext(r.Context())
}

func ProvideUserContext(h http.Handler) http.Handler {
	return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r.WithContext(user_context.WithContext(r.Context())))
	})
}

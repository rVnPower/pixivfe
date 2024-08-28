package handlers

import (
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/handlers/user_context"
)

type UserContext = user_context.UserContext

func GetUserContext(r *http.Request) *UserContext {
	return user_context.GetUserContext(r.Context())
}

func SetUserContext(handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r.WithContext(user_context.WithContext(r.Context())))
	}
}

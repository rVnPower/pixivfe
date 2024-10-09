package middleware

import (
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/i18n"
	"codeberg.org/vnpower/pixivfe/v2/server/session"
)

// ProvideUserContext is a middleware that wraps an http.Handler
// to inject a user context into each incoming request.
func SetLocaleFromCookie(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i18n.SetLocale(session.GetCookie(r, session.Cookie_Locale))
		h.ServeHTTP(w, r)
	})
}

package middleware

import (
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/server/session"
)

func SetDefaultCookies(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if it's a first visit
		_, err := r.Cookie("pixivfe-Visited")
		if err == http.ErrNoCookie {
			// It's a first visit, set default cookies
			localCSRF := session.GenerateCSRFToken()
			session.SetCookie(w, session.Cookie_LocalCSRF, localCSRF)

			// Set the visited cookie
			session.SetCookie(w, "pixivfe-Visited", "true")
		}

		next.ServeHTTP(w, r)
	})
}

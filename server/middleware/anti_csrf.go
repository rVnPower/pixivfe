package middleware

import (
	"errors"
    "net/http"
    "codeberg.org/vnpower/pixivfe/v2/server/session"
    "codeberg.org/vnpower/pixivfe/v2/server/routes"
)

func CSRFProtection(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
            csrfToken := r.FormValue("csrf_token")
            expectedToken := session.GetCookie(r, session.Cookie_LocalCSRF)

            if csrfToken == "" || csrfToken != expectedToken {
                routes.ErrorPage(w, r, errors.New("CSRF token validation failed"), http.StatusForbidden)
                return
            }
        }
        next.ServeHTTP(w, r)
    }
}

package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"codeberg.org/vnpower/pixivfe/v2/server/session"
)

// SetPrivacyHeaders is a middleware that adds security headers to HTTP responses.
func SetPrivacyHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Add("Referrer-Policy", "same-origin") // needed for settings redirect
		header.Add("X-Frame-Options", "DENY")
		// use this if need iframe: `X-Frame-Options: SAMEORIGIN`
		header.Add("X-Content-Type-Options", "nosniff")
		header.Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		// use this if need iframe: `frame-ancestors 'self'`
		header.Add("Permissions-Policy", "accelerometer=(), ambient-light-sensor=(), battery=(), camera=(), display-capture=(), document-domain=(), encrypted-media=(), execution-while-not-rendered=(), execution-while-out-of-viewport=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), midi=(), navigation-override=(), payment=(), publickey-credentials-get=(), screen-wake-lock=(), sync-xhr=(), usb=(), web-share=(), xr-spatial-tracking=()")
		if !strings.HasPrefix(r.URL.Path, "/diagnostics") {
			header.Add("Content-Security-Policy", fmt.Sprintf("base-uri 'self'; default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self' %s; media-src 'self' %s; connect-src 'self'; form-action 'self'; frame-ancestors 'none';", session.GetImageProxyOrigin(r), session.GetImageProxyOrigin(r)))
		}

		h.ServeHTTP(w, r)
	})
}

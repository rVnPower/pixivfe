package handlers

import (
	"net/http"
	"strings"
	"time"

	"codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/server/audit"
	"codeberg.org/vnpower/pixivfe/v2/server/request_context"
)

// ResponseWriterInterceptStatus wraps http.ResponseWriter to intercept the status code.
type ResponseWriterInterceptStatus struct {
	statusCode int
	http.ResponseWriter
}

// WriteHeader intercepts the status code before writing it to the response.
func (w *ResponseWriterInterceptStatus) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// CanRequestSkipLogger determines if a request should bypass the logging middleware.
// This is useful for reducing log clutter from static assets and development-specific routes.
func CanRequestSkipLogger(r *http.Request) bool {
	// Uncomment the following line to log all requests
	// return false
	path := r.URL.Path
	return strings.HasPrefix(path, "/img/") ||
		strings.HasPrefix(path, "/css/") ||
		strings.HasPrefix(path, "/js/") ||
		strings.HasPrefix(path, "/diagnostics") ||
		(config.GlobalConfig.InDevelopment &&
			(strings.HasPrefix(path, "/proxy/s.pximg.net/") || strings.HasPrefix(path, "/proxy/i.pximg.net/")))
}

// LogRequest is a middleware that logs incoming HTTP requests and their corresponding responses.
// It wraps the next handler in the chain and provides detailed logging of request/response metrics.
func LogRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w_ http.ResponseWriter, r *http.Request) {
		if CanRequestSkipLogger(r) {
			// If the request should skip logging, pass it directly to the next handler
			h.ServeHTTP(w_, r)
		} else {
			// Wrap the ResponseWriter to intercept the status code
			w := &ResponseWriterInterceptStatus{
				statusCode:     0,
				ResponseWriter: w_,
			}

			// TODO: Set user context here if needed

			// Record the start time of the request
			start_time := time.Now()

			// Call the next handler in the chain
			h.ServeHTTP(w, r)

			// Record the end time of the request
			end_time := time.Now()

			// Log the request details using the audit package
			audit.LogServerRoundTrip(audit.ServedRequestSpan{
				StartTime:  start_time,
				EndTime:    end_time,
				RequestId:  request_context.Get(r).RequestId,
				Method:     r.Method,
				Path:       r.URL.Path,
				Status:     w.statusCode,
				Referer:    r.Referer(),
				RemoteAddr: r.RemoteAddr,
				Error:      request_context.Get(r).CaughtError,
			})
		}
	})
}

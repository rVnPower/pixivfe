package handlers

import (
	"net/http"
	"strings"
	"time"

	"codeberg.org/vnpower/pixivfe/v2/audit"
	"codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/request_context"
)

type ResponseWriterInterceptStatus struct {
	statusCode int
	http.ResponseWriter
}

func (w *ResponseWriterInterceptStatus) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func CanRequestSkipLogger(r *http.Request) bool {
	// return false
	path := r.URL.Path
	return strings.HasPrefix(path, "/img/") ||
		strings.HasPrefix(path, "/css/") ||
		strings.HasPrefix(path, "/js/") ||
		strings.HasPrefix(path, "/diagnostics") ||
		(config.GlobalConfig.InDevelopment &&
			(strings.HasPrefix(path, "/proxy/s.pximg.net/") || strings.HasPrefix(path, "/proxy/i.pximg.net/")))
}

func LogRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w_ http.ResponseWriter, r *http.Request) {
		if CanRequestSkipLogger(r) {
			h.ServeHTTP(w_, r)
		} else {
			w := &ResponseWriterInterceptStatus{
				statusCode:     0,
				ResponseWriter: w_,
			}
			// set user context

			start_time := time.Now()

			h.ServeHTTP(w, r)

			end_time := time.Now()

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

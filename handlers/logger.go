package handlers

import (
	"net/http"
	"strings"
	"time"

	"codeberg.org/vnpower/pixivfe/v2/audit"
	"codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/handlers/user_context"
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
		(config.GlobalConfig.InDevelopment &&
			(strings.HasPrefix(path, "/proxy/s.pximg.net/") || strings.HasPrefix(path, "/proxy/i.pximg.net/")))
}

func LogRequest(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w_ http.ResponseWriter, r *http.Request) {
		if CanRequestSkipLogger(r) {
			f(w_, r)
		} else {
			w := &ResponseWriterInterceptStatus{
				statusCode:     0,
				ResponseWriter: w_,
			}
			// set user context
			r = r.WithContext(user_context.WithContext(r.Context()))

			start_time := time.Now()

			f(w, r)

			end_time := time.Now()

			audit.LogServerRoundTrip(r.Context(), audit.ServerPerformance{
				StartTime:  start_time,
				EndTime:    end_time,
				RemoteAddr: r.RemoteAddr,
				Method:     r.Method,
				Path:       r.URL.Path,
				Status:     w.statusCode,
				Error:      GetUserContext(r).Err,
			})
		}
	}
}

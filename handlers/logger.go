package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"codeberg.org/vnpower/pixivfe/v2/audit"
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
		strings.HasPrefix(path, "/proxy/s.pximg.net/") ||
		strings.HasPrefix(path, "/proxy/i.pximg.net/")
}

func LogRequest(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w_ http.ResponseWriter, r *http.Request) {
		w := &ResponseWriterInterceptStatus{
			statusCode:     0,
			ResponseWriter: w_,
		}
		// set user context
		r = r.WithContext(context.WithValue(r.Context(), UserContextKey, &UserContext{}))

		start_time := time.Now()

		f(w, r)

		end_time := time.Now()

		audit.TraceRoute(audit.RoutePerf{
			StartTime:   start_time,
			EndTime:     end_time,
			RemoteAddr:  r.RemoteAddr,
			Method:      r.Method,
			Path:        r.URL.Path,
			Status:      w.statusCode,
			Err:         GetUserContext(r).Err,
			SkipLogging: CanRequestSkipLogger(r),
		})
	}
}

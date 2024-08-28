package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"
)

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

		if !CanRequestSkipLogger(r) { // logger
			time := start_time
			latency := end_time.Sub(start_time)
			ip := r.RemoteAddr
			method := r.Method
			path := r.URL.Path
			status := w.statusCode
			err := GetUserContext(r).Err

			log.Printf("%v +%v %v %v %v %v %v", time, latency, ip, method, path, status, err)
		}
	}
}

type ResponseWriterInterceptStatus struct {
	statusCode int
	http.ResponseWriter
}

func (w *ResponseWriterInterceptStatus) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

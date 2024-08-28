package handlers

import (
	"bytes"
	"log"
	"maps"
	"net/http"
	"net/http/httptest"
	"slices"

	"codeberg.org/vnpower/pixivfe/v2/routes"
)

func CatchError(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header_backup := http.Header{}
		for k, v := range w.Header() {
			header_backup[k] = slices.Clone(v)
		}
		recorder := httptest.ResponseRecorder{
			HeaderMap: w.Header(),
			Body:      new(bytes.Buffer),
			Code:      200,
		}
		err := handler(&recorder, r)
		if err != nil {
			clear(header_backup)
			maps.Copy(w.Header(), header_backup)
			GetUserContext(r).Err = err
		} else {
			w.WriteHeader(recorder.Code)
			_, _ = recorder.Body.WriteTo(w)
		}
	}
}

func HandleError(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)

		err := GetUserContext(r).Err

		if err != nil {
			code := GetUserContext(r).ErrorStatusCodeOverride
			if code == 0 {
				code = http.StatusInternalServerError
			}
			w.WriteHeader(code)
			// Send custom error page
			err = routes.ErrorPage(w, r, err)
			if err != nil {
				log.Panicf("[fix this ASAP] Error rendering error route: %s", err)
			}
		}
	})
}

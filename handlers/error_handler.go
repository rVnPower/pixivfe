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

type UserContext struct {
	Err error
	ErrorStatusCodeOverride int
}

type userContextKey struct{}

var UserContextKey = userContextKey{}

func GetUserContext(r *http.Request) *UserContext {
	return r.Context().Value(UserContextKey).(*UserContext)
}

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

func ErrorHandler(w http.ResponseWriter, r *http.Request) { // error handler
	err := GetUserContext(r).Err

	if err != nil {
		log.Printf("Internal Server Error: %s", err)
		code := GetUserContext(r).ErrorStatusCodeOverride
		if code == 0 {
			code = http.StatusInternalServerError
		}
		w.WriteHeader(code)
		// Send custom error page
		err = routes.ErrorPage(w, r, err)
		if err != nil {
			log.Printf("Error rendering error route: %s", err)
		}
	}
}

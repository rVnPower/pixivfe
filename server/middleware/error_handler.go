package middleware

import (
	"bytes"
	"maps"
	"net/http"
	"net/http/httptest"
	"slices"

	"codeberg.org/vnpower/pixivfe/v2/server/request_context"
	"codeberg.org/vnpower/pixivfe/v2/server/routes"
)

// CatchError is a middleware that wraps an HTTP handler to catch and manage errors.
// It allows for graceful error handling and response manipulation.
func CatchError(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Backup the original response headers
		header_backup := http.Header{}
		for k, v := range w.Header() {
			header_backup[k] = slices.Clone(v)
		}

		// Create a response recorder to capture the handler's output
		recorder := httptest.ResponseRecorder{
			HeaderMap: w.Header(),
			Body:      new(bytes.Buffer),
			Code:      200,
		}

		// Execute the handler and catch any returned error
		err := handler(&recorder, r)
		if err != nil {
			// If an error occurred, restore the original headers
			clear(header_backup)
			maps.Copy(w.Header(), header_backup)
			// Store the error in the request context for later handling
			request_context.Get(r).CaughtError = err
		} else {
			w.WriteHeader(recorder.Code)
			_, _ = recorder.Body.WriteTo(w)
		}
	}
}

// HandleError is a middleware that checks for errors caught by CatchError and renders an error page if necessary.
func HandleError(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Execute the wrapped handler
		h.ServeHTTP(w, r)

		// Check if an error was caught during the request processing
		err := request_context.Get(r).CaughtError

		if err != nil {
			// If an error was caught, render the error page
			routes.ErrorPage(w, r, err, http.StatusInternalServerError)
		}
	})
}

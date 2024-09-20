// Package audit provides functionality for logging and recording various types of spans
// in the application, including server requests and API calls.
package audit

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"codeberg.org/vnpower/pixivfe/v2/config"
	"github.com/oklog/ulid/v2"
)

// Logger is a custom logger with no timestamp prefix, as we control the timestamps in our log messages.
var Logger = log.New(os.Stderr, "", 0)

// RecordedSpans stores a slice of recorded Span objects for later analysis or debugging.
var RecordedSpans = []Span{}

// LogAndRecord logs the given span and optionally records it in the RecordedSpans slice.
// It manages the RecordedSpans slice to maintain a maximum number of recorded spans.
func LogAndRecord(span Span) {
	// Log the span with a formatted timestamp, duration, and log line
	Logger.Printf("%v +%-5.3f %s", span.GetStartTime().Format("2006-01-02 15:04:05.000"), float64(Duration(span))/float64(time.Second), span.LogLine())

	// If MaxRecordedCount is set, manage the RecordedSpans slice
	if MaxRecordedCount != 0 {
		// Remove the oldest span if we're at capacity
		if len(RecordedSpans)+1 == MaxRecordedCount {
			RecordedSpans = RecordedSpans[1:]
		}
		// Append the new span
		RecordedSpans = append(RecordedSpans, span)
	}
}

// LogServerRoundTrip logs and records a server request span.
// It also logs any internal server errors that occurred during the request.
func LogServerRoundTrip(perf ServedRequestSpan) {
	if perf.Error != nil {
		log.Printf("Internal Server Error: %s", perf.Error)
	}

	LogAndRecord(perf)
}

// LogAPIRoundTrip logs and records an API request span.
// It handles saving the response body to a file if enabled and logs warnings for non-2xx status codes.
func LogAPIRoundTrip(perf APIRequestSpan) {
	if perf.Response != nil {
		// Save response body to file if enabled and body is not empty
		if perf.Body != "" && optionSaveResponse {
			var err error
			perf.ResponseFilename, err = writeResponseBodyToFile(perf.Body)
			if err != nil {
				log.Print("When saving response to file: ", err)
			}
		}
		// Log a warning for non-2xx status codes
		if !(300 > perf.Response.StatusCode && perf.Response.StatusCode >= 200) {
			log.Print("(WARN) non-2xx response from pixiv:")
		}
	}

	LogAndRecord(perf)
}

// writeResponseBodyToFile saves the given response body to a file in the ResponseSaveLocation directory.
// It generates a unique filename using ULID and returns the filename and any error encountered.
func writeResponseBodyToFile(body string) (string, error) {
	// Generate a unique ID using ULID
	id := ulid.Make().String()
	// Create a filename using the last 6 characters of the ID
	filename := path.Join(config.GlobalConfig.ResponseSaveLocation, id[len(id)-6:])
	// Write the body to the file with read/write permissions for the owner only
	err := os.WriteFile(filename, []byte(body), 0o600)
	if err != nil {
		return "", fmt.Errorf("failed to write response body to file %s: %w", filename, err)
	}
	log.Printf("Successfully wrote response body to file: %s", filename)
	return filename, nil
}

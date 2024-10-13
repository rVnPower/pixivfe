package audit

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/oklog/ulid/v2"

	"codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/i18n"
)

// Logger is a custom logger with no timestamp prefix, as we control the timestamps in our log messages.
var Logger = log.New(os.Stderr, "", 0)

// RecordedSpans stores a slice of recorded Span objects for later analysis or debugging.
var RecordedSpans = []Span{}

// standardLog logs messages in a standardized format.
func standardLog(logType string, message string, err error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	errStr := ""
	if err != nil {
		errStr = fmt.Sprintf(" error=%v", err)
	}
	Logger.Printf("%s %s %s locale=%s%s", timestamp, logType, message, i18n.GetLocale(), errStr)
}

// LogAndRecord logs the given span and optionally records it in the RecordedSpans slice.
// It manages the RecordedSpans slice to maintain a maximum number of recorded spans.
func LogAndRecord(span Span) {
	duration := float64(Duration(span)) / float64(time.Second)
	message := fmt.Sprintf("+%-5.3f %s", duration, span.LogLine())
	standardLog("INFO", message, nil)

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
		standardLog("ERROR", "Internal Server Error", perf.Error)
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
				standardLog("ERROR", "Failed to save response to file", err)
			}
		}
		// Log a warning for non-2xx status codes
		if !(300 > perf.Response.StatusCode && perf.Response.StatusCode >= 200) {
			standardLog("WARN", fmt.Sprintf("Non-2xx response from pixiv: %d", perf.Response.StatusCode), nil)
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
		return "", i18n.Errorf("failed to write response body to file %s: %w", filename, err)
	}

	return filename, nil
}

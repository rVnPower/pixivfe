package audit

import (
	"os"
	"path"
	"time"

	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"

	"codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/i18n"
)

// RecordedRequestSpans stores a slice of requestSpan objects.
var RecordedRequestSpans = []Span{}

// LogAndRecord logs the given requestSpan and optionally records it in the RecordedRequestSpans slice.
// It manages the RecordedRequestSpans slice to maintain a maximum number of recorded spans.
func LogAndRecord(requestSpan Span) {
	duration := float64(Duration(requestSpan)) / float64(time.Second)

	logger.Info("Request",
		// Zap already prefixes the log with a timestamp
		// zap.String("timestamp", span.GetStartTime().Format(time.RFC3339)),
		zap.String("component", requestSpan.Component()),
		zap.Any("action", requestSpan.Action()),
		zap.Any("outcome", requestSpan.Outcome()),
		zap.Float64("duration", duration),
	)

	// If MaxRecordedCount is set, manage the RecordedRequestSpans slice
	if MaxRecordedCount != 0 {
		// Remove the oldest span if we're at capacity
		if len(RecordedRequestSpans)+1 == MaxRecordedCount {
			RecordedRequestSpans = RecordedRequestSpans[1:]
		}
		// Append the new span
		RecordedRequestSpans = append(RecordedRequestSpans, requestSpan)
	}
}

// LogServerRoundTrip logs and records a ServerRequestSpan.
// It also logs any internal server errors that occurred during the request.
func LogServerRoundTrip(requestSpan ServerRequestSpan) {
	if requestSpan.Error != nil {
		logger.Error("Internal Server Error",
			zap.Error(requestSpan.Error),
			zap.String("requestId", requestSpan.RequestId),
		)
	}

	LogAndRecord(requestSpan)
}

// LogAPIRoundTrip logs and records an APIRequestSpan.
// It handles saving the response body to a file if enabled and logs warnings for non-2xx status codes.
func LogAPIRoundTrip(requestSpan APIRequestSpan) {
	if requestSpan.Response != nil {
		// Save response body to file if enabled and body is not empty
		if requestSpan.Body != "" && optionSaveResponse {
			var err error
			requestSpan.ResponseFilename, err = writeResponseBodyToFile(requestSpan.Body)
			if err != nil {
				logger.Error("Failed to save response to file",
					zap.Error(err),
					zap.String("requestId", requestSpan.RequestId),
				)
			}
		}
		// Log a warning for non-2xx status codes
		if !(300 > requestSpan.Response.StatusCode && requestSpan.Response.StatusCode >= 200) {
			logger.Warn("Non-2xx response from pixiv",
				zap.Int("status", requestSpan.Response.StatusCode),
				zap.String("requestId", requestSpan.RequestId),
			)
		}
	}

	LogAndRecord(requestSpan)
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

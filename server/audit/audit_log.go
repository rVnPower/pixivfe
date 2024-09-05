package audit

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/oklog/ulid/v2"
)

// logger with no timestamp prefix, because we control the timestamps
var Logger = log.New(os.Stderr, "", 0)

var RecordedSpans = []Span{}

func LogAndRecord(span Span) {
	Logger.Printf("%v +%-5.3f %s", span.GetStartTime().Format("2006-01-02 15:04:05.000"), float64(Duration(span))/float64(time.Second), span.LogLine())

	if MaxRecordedCount != 0 {
		if len(RecordedSpans)+1 == MaxRecordedCount {
			RecordedSpans = RecordedSpans[1:]
		}
		RecordedSpans = append(RecordedSpans, span)
	}
}

func LogServerRoundTrip(perf ServedRequestSpan) {
	if perf.Error != nil {
		log.Printf("Internal Server Error: %s", perf.Error)
	}

	LogAndRecord(perf)
}

func LogAPIRoundTrip(perf APIRequestSpan) {
	if perf.Response != nil {
		if perf.Body != "" && optionSaveResponse {
			var err error
			perf.ResponseFilename, err = writeResponseBodyToFile(perf.Body)
			if err != nil {
				log.Print("When saving response to file: ", err)
			}
		}
		if !(300 > perf.Response.StatusCode && perf.Response.StatusCode >= 200) {
			log.Print("(WARN) non-2xx response from pixiv:")
		}
	}

	LogAndRecord(perf)
}

func writeResponseBodyToFile(body string) (string, error) {
	id := ulid.Make().String()
	filename := path.Join(DevDir_Response, id[len(id)-6:])
	err := os.WriteFile(filename, []byte(body), 0o600)
	if err != nil {
		return "", err
	}
	return filename, nil
}

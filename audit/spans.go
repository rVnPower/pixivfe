package audit

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Span interface {
	GetStartTime() time.Time
	GetEndTime() time.Time
	GetRequestId() string
	LogLine() string
}

func Duration(span Span) time.Duration {
	return span.GetEndTime().Sub(span.GetStartTime())
}

// logger with no timestamp prefix, because we control the timestamps
var Logger = log.New(os.Stderr, "", 0)

var RecordedSpans = []Span{}

// should be configurable. set to 0 to disable recording
var MaxRecordedCount = 128

func LogAndRecord(span Span) {
	Logger.Printf("%v +%-5.3f %s", span.GetStartTime().Format("2006-01-02 15:04:05.000"), float64(Duration(span))/float64(time.Second), span.LogLine())

	if MaxRecordedCount != 0 {
		if len(RecordedSpans)+1 == MaxRecordedCount {
			RecordedSpans = RecordedSpans[1:]
		}
		RecordedSpans = append(RecordedSpans, span)
	}
}

type ServedRequestSpan struct {
	StartTime  time.Time
	EndTime    time.Time
	RequestId  string
	Method     string
	Path       string `json:"Url"`
	Status     int
	Referer    string
	RemoteAddr string
	Error      error
}

func (span ServedRequestSpan) GetStartTime() time.Time {
	return span.StartTime
}
func (span ServedRequestSpan) GetEndTime() time.Time {
	return span.EndTime
}
func (span ServedRequestSpan) GetRequestId() string {
	return span.RequestId
}
func (span ServedRequestSpan) LogLine() string {
	return fmt.Sprintf("%v %v %v %v", span.Method, span.Path, span.Status, span.Error)
}

type APIRequestSpan struct {
	StartTime        time.Time
	EndTime          time.Time
	RequestId        string
	Response         *http.Response `json:"-"`
	Error            error
	Method           string
	Url              string
	Token            string
	Body             string `json:"-"`
	ResponseFilename string
}

func (span APIRequestSpan) GetStartTime() time.Time {
	return span.StartTime
}
func (span APIRequestSpan) GetEndTime() time.Time {
	return span.EndTime
}
func (span APIRequestSpan) GetRequestId() string {
	return span.RequestId
}
func (span APIRequestSpan) LogLine() string {
	return fmt.Sprintf("-> %v %v %v %v", span.Method, span.Url, span.Error, span.ResponseFilename)
}

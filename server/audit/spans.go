package audit

import (
	"fmt"
	"net/http"
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
	return fmt.Sprintf("SERVER method=%s path=%s status=%d error=%v", span.Method, span.Path, span.Status, span.Error)
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
	return fmt.Sprintf("API method=%s url=%s error=%v responseFile=%s", span.Method, span.Url, span.Error, span.ResponseFilename)
}

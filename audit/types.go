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
	LogLine() string
}

func Duration(span Span) time.Duration {
	return span.GetEndTime().Sub(span.GetStartTime())
}

// logger with no timestamp prefix, because we control the timestamps
var Logger = log.New(os.Stderr, "", 0)

func Log(span Span) {
	Logger.Printf("%v +%-5.3f %s", span.GetStartTime().Format("2006-01-02 15:04:05.000"), float64(Duration(span))/float64(time.Second), span.LogLine())
}

type ServedRequestSpan struct {
	StartTime  time.Time
	EndTime    time.Time
	RemoteAddr string
	Method     string
	Path       string
	Status     int
	Error      error
}

func (perf ServedRequestSpan) GetStartTime() time.Time {
	return perf.StartTime
}

func (perf ServedRequestSpan) GetEndTime() time.Time {
	return perf.EndTime
}

func (perf ServedRequestSpan) LogLine() string {
	return fmt.Sprintf("%v %v %v %v", perf.Method, perf.Path, perf.Status, perf.Error)
}

type APIRequestSpan struct {
	StartTime        time.Time
	EndTime          time.Time
	Response         *http.Response
	Error            error
	Method           string
	Url              string
	Token            string
	Body             string
	ResponseFilename string
}

func (perf APIRequestSpan) GetStartTime() time.Time {
	return perf.StartTime
}
func (perf APIRequestSpan) GetEndTime() time.Time {
	return perf.EndTime
}
func (perf APIRequestSpan) LogLine() string {
	return fmt.Sprintf("-> %v %v %v %v", perf.Method, perf.Url, perf.Error, perf.ResponseFilename)
}

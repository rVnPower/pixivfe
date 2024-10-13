package audit

import (
	"codeberg.org/vnpower/pixivfe/v2/i18n"
	"net/http"
	"time"
)

type Span interface {
	GetStartTime() time.Time
	GetEndTime() time.Time
	GetRequestId() string
	Component() string
	Action() map[string]interface{}
	Outcome() map[string]interface{}
}

func Duration(span Span) time.Duration {
	return span.GetEndTime().Sub(span.GetStartTime())
}

type ServerRequestSpan struct {
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

func (span ServerRequestSpan) GetStartTime() time.Time {
	return span.StartTime
}
func (span ServerRequestSpan) GetEndTime() time.Time {
	return span.EndTime
}
func (span ServerRequestSpan) GetRequestId() string {
	return span.RequestId
}
func (span ServerRequestSpan) Component() string {
	return "server"
}
func (span ServerRequestSpan) Action() map[string]interface{} {
	return map[string]interface{}{
		"method": span.Method,
		"path":   span.Path,
	}
}
func (span ServerRequestSpan) Outcome() map[string]interface{} {
	outcome := map[string]interface{}{
		"status": span.Status,
		"error":  "<none>",
		"locale": i18n.GetLocale(),
	}
	if span.Error != nil {
		outcome["error"] = span.Error.Error()
	}
	return outcome
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
func (span APIRequestSpan) Component() string {
	return "API"
}
func (span APIRequestSpan) Action() map[string]interface{} {
	return map[string]interface{}{
		"method":        span.Method,
		"url":           span.Url,
		"response_file": span.ResponseFilename,
	}
}
func (span APIRequestSpan) Outcome() map[string]interface{} {
	outcome := map[string]interface{}{
		"status": "success",
		"error":  "<none>",
		"locale": i18n.GetLocale(),
	}
	if span.Error != nil {
		outcome["status"] = "error"
		outcome["error"] = span.Error.Error()
	}
	if span.Response != nil {
		outcome["status_code"] = span.Response.StatusCode
	}
	return outcome
}

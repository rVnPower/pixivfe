package audit

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"codeberg.org/vnpower/pixivfe/v2/handlers/user_context"
	"codeberg.org/vnpower/pixivfe/v2/utils"
	"github.com/openzipkin/zipkin-go"
)

const DevDir_Response = "/tmp/pixivfe-dev/resp"

var optionSaveResponse bool

func Init(saveResponse bool, tracer *zipkin.Tracer) error {
	optionSaveResponse = saveResponse
	utils.Tracer = tracer
	if optionSaveResponse {
		return os.MkdirAll(DevDir_Response, 0o700)
	} else {
		return nil
	}
}

type ServerPerformance struct {
	StartTime   time.Time
	EndTime     time.Time
	RemoteAddr  string
	Method      string
	Path        string
	Status      int
	Error       error
	SkipLogging bool
}

type APIPerformance struct {
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

func LogServerRoundTrip(context context.Context, perf ServerPerformance) {
	if perf.Error != nil {
		log.Printf("Internal Server Error: %s", perf.Error)
	}

	span, _ := utils.Tracer.StartSpanFromContext(context, fmt.Sprintf("%v %v %v %v", perf.Method, perf.Path, perf.Status, perf.Error), zipkin.StartTime(perf.StartTime), zipkin.Parent(user_context.GetUserContext(context).Parent))
	span.Tag("RemoteAddr", perf.RemoteAddr)
	span.FinishedWithDuration(perf.EndTime.Sub(perf.StartTime))
}

func LogAPIRoundTrip(context context.Context, perf APIPerformance) {
	if perf.Response != nil {
		if perf.Body != "" && optionSaveResponse {
			var err error
			perf.ResponseFilename, err = writeResponseBodyToFile(perf.Body)
			if err != nil {
				log.Println("When saving response to file: ", err)
			} else {
				log.Println(fmt.Sprintf("[API] %v %v saved to %v", perf.Method, perf.Url, perf.ResponseFilename))
			}
		}
		if !(300 > perf.Response.StatusCode && perf.Response.StatusCode >= 200) {
			log.Println("(WARN) non-2xx response from pixiv:")
		}
	}
	span, _ := utils.Tracer.StartSpanFromContext(context, fmt.Sprintf("%v %v %v", perf.Method, perf.Url, perf.Error), zipkin.StartTime(perf.StartTime), zipkin.Parent(user_context.GetUserContext(context).Parent))
	span.Tag("ResponseFilename", perf.ResponseFilename)
	span.FinishedWithDuration(perf.EndTime.Sub(perf.StartTime))
}

func writeResponseBodyToFile(body string) (string, error) {
	filename := path.Join(DevDir_Response, time.Now().UTC().Format(time.RFC3339Nano))
	err := os.WriteFile(filename, []byte(body), 0o600)
	if err != nil {
		return "", err
	}
	return filename, nil
}

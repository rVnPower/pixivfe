package audit

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
	http_reporter "github.com/openzipkin/zipkin-go/reporter/http"

	"codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/handlers/user_context"
	"codeberg.org/vnpower/pixivfe/v2/utils"
)

const DevDir_Response = "/tmp/pixivfe-dev/resp"

var optionSaveResponse bool

func Init(saveResponse bool) error {
	optionSaveResponse = saveResponse
	if optionSaveResponse {
		err := os.MkdirAll(DevDir_Response, 0o700)
		if err != nil {
			return err
		}
	}

	var reporter reporter.Reporter = nil

	_, enableReporting := config.LookupEnv("PIXIVFE_ENABLE_ZIPKIN")
	if enableReporting {
		reporter = http_reporter.NewReporter("http://localhost:9411/api/v2/spans")
		defer func() {
			_ = reporter.Close()
		}()
	} else {
		// comment out this block in logging is too verbose
		reporter = NewLogReporter()
		defer func() {
			_ = reporter.Close()
		}()
	}

	// this is purely theoretical. the port is used for distributed tracing.
	endpoint, err := zipkin.NewEndpoint("pixivfe", "localhost:8282")
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// initialize our tracer
	tracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	utils.Tracer = tracer

	return nil
}

type ServerPerformance struct {
	StartTime   time.Time
	EndTime     time.Time
	RemoteAddr  string
	Method      string
	Path        string
	Status      int
	Error       error
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
				log.Print("When saving response to file: ", err)
			} else {
				log.Printf("[API] %v %v saved to %v", perf.Method, perf.Url, perf.ResponseFilename)
			}
		}
		if !(300 > perf.Response.StatusCode && perf.Response.StatusCode >= 200) {
			log.Print("(WARN) non-2xx response from pixiv:")
		}
	}
	span, _ := utils.Tracer.StartSpanFromContext(context, fmt.Sprintf("API %v %v %v", perf.Method, perf.Url, perf.Error), zipkin.StartTime(perf.StartTime), zipkin.Parent(user_context.GetUserContext(context).Parent))
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

package audit

import (
	"log"
	"net/http"
	"os"
	"path"
	"time"

	config "codeberg.org/vnpower/pixivfe/v2/config"
)

const DevDir_Response = "/tmp/pixivfe-dev/resp"

var optionSaveResponse bool

func Init(saveResponse bool) error {
	optionSaveResponse = saveResponse
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

func LogServerRoundTrip(perf ServerPerformance) {
	if perf.Error != nil {
		log.Printf("Internal Server Error: %s", perf.Error)
	}

	if !perf.SkipLogging {
		// todo: log.Printf("%v +%v %v %v %v %v %v", time, latency, ip, method, path, status, err)
	}
}

func LogAPIRoundTrip(perf APIPerformance) {
	if perf.Response != nil {
		if perf.Body != "" && optionSaveResponse {
			var err error
			perf.ResponseFilename, err = writeResponseBodyToFile(perf.Body)
			if err != nil {
				log.Println("When saving response to file: ", err)
			}
		}
		if !(300 > perf.Response.StatusCode && perf.Response.StatusCode >= 200) {
			log.Println("(WARN) non-2xx response from pixiv:")
		}
	}
	// structured logging
	if config.GlobalServerConfig.InDevelopment {
		// todo
	} else {
		// todo
	}
}

func writeResponseBodyToFile(body string) (string, error) {
	filename := path.Join(DevDir_Response, time.Now().UTC().Format(time.RFC3339Nano))
	err := os.WriteFile(filename, []byte(body), 0o600)
	if err != nil {
		return "", err
	}
	return filename, nil
}

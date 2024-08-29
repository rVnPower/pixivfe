package audit

import (
	"context"
	"log"
	"os"
	"path"

	"github.com/oklog/ulid/v2"
)

const DevDir_Response = "/tmp/pixivfe/r"

var optionSaveResponse bool

func Init(saveResponse bool) error {
	optionSaveResponse = saveResponse
	if optionSaveResponse {
		err := os.MkdirAll(DevDir_Response, 0o700)
		if err != nil {
			return err
		}
	}

	return nil
}

func LogServerRoundTrip(context context.Context, perf ServedRequestSpan) {
	if perf.Error != nil {
		log.Printf("Internal Server Error: %s", perf.Error)
	}

	Log(perf)
}

func LogAPIRoundTrip(context context.Context, perf APIRequestSpan) {
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

	Log(perf)
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

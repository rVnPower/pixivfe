package audit

import (
	"fmt"
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

func LogAPIRoundTrip(resp *http.Response, err error, method, url, token, body string, start_time, end_time time.Time) {
	if config.GlobalServerConfig.InDevelopment {
		errs := ""
		if err != nil {
			errs = fmt.Sprintf("ERR %v", err)
		}
		if resp != nil {
			filename := ""
			if body != "" && optionSaveResponse {
				var err error
				filename, err = writeResponseBodyToFile(body)
				if err != nil {
					log.Println(err)
				}
			}
			if !(300 > resp.StatusCode && resp.StatusCode >= 200) {
				log.Println("(WARN) non-2xx response from pixiv:")
			}
			log.Println("->", method, url, "->", resp.StatusCode, filename, errs)
		} else {
			log.Println("->", method, url, errs)
		}
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

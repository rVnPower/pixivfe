package audit

import "os"

const DevDir_Response = "/tmp/pixivfe/r"

var optionSaveResponse bool

var MaxRecordedCount = 0

func Init(saveResponse bool) error {
	optionSaveResponse = saveResponse
	if optionSaveResponse {
		MaxRecordedCount = 128
		err := os.MkdirAll(DevDir_Response, 0o700)
		if err != nil {
			return err
		}
	}

	return nil
}

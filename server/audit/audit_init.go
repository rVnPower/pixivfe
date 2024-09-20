package audit

import (
	"os"

	"codeberg.org/vnpower/pixivfe/v2/config"
)

var optionSaveResponse bool
var MaxRecordedCount = 0

// Init initializes the audit package and sets up response saving if enabled.
// saveResponse is passed as a boolean from main.go.
// The response save location is taken from the global configuration.
func Init(saveResponse bool) error {
	optionSaveResponse = saveResponse

	if optionSaveResponse {
		MaxRecordedCount = 128
		err := os.MkdirAll(config.GlobalConfig.ResponseSaveLocation, 0o700)
		if err != nil {
			return err
		}
	}

	return nil
}

package audit

import (
	"os"

	"codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/i18n"
)

var optionSaveResponse bool
var MaxRecordedCount = 0

// Init initializes the audit package and sets up response saving if enabled.
// saveResponse is passed as a boolean from main.go.
// The response save location is taken from the global configuration.
func Init(saveResponse bool) error {
	optionSaveResponse = saveResponse

	if !optionSaveResponse {
		return nil
	}

	MaxRecordedCount = 128
	savePath := config.GlobalConfig.ResponseSaveLocation

	if err := os.MkdirAll(savePath, 0o700); err != nil {
		standardLog("ERROR", "Failed to create response save directory", err)
		return i18n.Errorf("failed to create response save directory: %w", err)
	}

	return nil
}

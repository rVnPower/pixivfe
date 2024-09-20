package audit

import (
	"fmt"
	"log"
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

	if !optionSaveResponse {
		return nil
	}

	MaxRecordedCount = 128
	savePath := config.GlobalConfig.ResponseSaveLocation

	if err := os.MkdirAll(savePath, 0o700); err != nil {
		log.Printf("Error creating response save directory: %v", err)
		return fmt.Errorf("failed to create response save directory: %w", err)
	}

	return nil
}

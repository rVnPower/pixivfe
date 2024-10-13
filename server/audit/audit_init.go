package audit

import (
	"os"

	"codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/i18n"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var optionSaveResponse bool
var MaxRecordedCount = 0
var logger *zap.Logger

// Init initializes the audit package and sets up response saving if enabled.
// saveResponse is passed as a boolean from main.go.
// The response save location is taken from the global configuration.
func Init(saveResponse bool) error {
	optionSaveResponse = saveResponse

	// Initialize zap logger
	zapConfig := zap.NewProductionConfig()
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConfig.OutputPaths = []string{"stderr"}
	var err error
	logger, err = zapConfig.Build()
	if err != nil {
		return i18n.Errorf("failed to initialize zap logger: %w", err)
	}

	if !optionSaveResponse {
		return nil
	}

	MaxRecordedCount = 128
	savePath := config.GlobalConfig.ResponseSaveLocation

	if err := os.MkdirAll(savePath, 0o700); err != nil {
		logger.Error("Failed to create response save directory",
			zap.Error(err),
			zap.String("path", savePath),
		)
		return i18n.Errorf("failed to create response save directory: %w", err)
	}

	return nil
}

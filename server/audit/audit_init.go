package audit

import (
	"os"

	"codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/i18n"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	optionSaveResponse bool
	MaxRecordedCount   = 0
	logger             *zap.Logger
)

// Init initializes the audit package and sets up response saving if enabled.
// saveResponse is passed as a boolean from main.go.
// The response save location is taken from the global configuration.
func Init(saveResponse bool) error {
	// Initialize the auditing parameter
	optionSaveResponse = saveResponse

	// Read configuration values from GlobalConfig
	savePath := config.GlobalConfig.ResponseSaveLocation
	logLevel := config.GlobalConfig.LogLevel
	logOutputs := config.GlobalConfig.LogOutputs
	logFormat := config.GlobalConfig.LogFormat

	// Initialize zap logger with custom configuration
	zapConfig := zap.NewProductionConfig()

	// Adjust log encoding format based on config
	switch logFormat {
	case "json":
		zapConfig.Encoding = "json"
	case "console":
		fallthrough
	default:
		zapConfig.Encoding = "console"
	}

	// Adjust log level based on config
	var atom zap.AtomicLevel
	switch logLevel {
	case "debug":
		atom = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		atom = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		atom = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		atom = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		atom = zap.NewAtomicLevelAt(zap.InfoLevel) // Default to info level
	}
	zapConfig.Level = atom

	// Set custom output paths
	if len(logOutputs) > 0 {
		zapConfig.OutputPaths = logOutputs
	} else {
		// Default to standard error if logOutputs is empty
		zapConfig.OutputPaths = []string{"stdout"}
	}

	// Use console-friendly time encoding
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Build and assign the logger
	var err error
	logger, err = zapConfig.Build()
	if err != nil {
		return i18n.Errorf("failed to initialize zap logger: %w", err)
	}

	if !optionSaveResponse {
		return nil
	}

	// Handle saving responses
	MaxRecordedCount = 128
	if err := os.MkdirAll(savePath, 0o700); err != nil {
		logger.Error("Failed to create response save directory",
			zap.Error(err),
			zap.String("path", savePath),
		)
		return i18n.Errorf("failed to create response save directory: %w", err)
	}

	return nil
}

package i18n

import (
	"fmt"
)

// Log returns a localized log message
func Log(text string) string {
	return Lookup(text)
}

// Logf returns a localized formatted log message
func Logf(format string, a ...any) string {
	format = Lookup(format)
	return fmt.Sprintf(format, a...)
}

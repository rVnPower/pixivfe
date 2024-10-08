package i18n

import (
	"errors"
	"fmt"

	"github.com/timandy/routine"
)

var Locale = routine.NewInheritableThreadLocal[string]()

func getLocalizer() Localizer {
	locale := Locale.Get()
	if locale == "" {
		locale = "en"
	}
	return LocalizerOf(locale)
}

func Error(text string) error {
	text = getLocalizer().lookup(text)
	return errors.New(text)
}

func Errorf(format string, a ...any) error {
	format = getLocalizer().lookup(format)
	return fmt.Errorf(format, a...)
}

func Sprintf(format string, a ...any) string {
	format = getLocalizer().lookup(format)
	return fmt.Sprintf(format, a...)
}

// translate string
func Tr(text string) string {
	return getLocalizer().lookup(text)
}

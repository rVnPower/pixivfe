package i18n

import (
	"errors"
	"fmt"
)

var GlobalLocale string = "en"

func Error(text string) error {
	text = LocalizerOf(GlobalLocale).lookup(text)
	return errors.New(text)
}

func Errorf(format string, a ...any) error {
	format = LocalizerOf(GlobalLocale).lookup(format)
	return fmt.Errorf(format, a...)
}

func Sprintf(format string, a ...any) string {
	format = LocalizerOf(GlobalLocale).lookup(format)
	return fmt.Sprintf(format, a...)
}

// translate string
func Tr(text string) string {
	return LocalizerOf(GlobalLocale).lookup(text)
}

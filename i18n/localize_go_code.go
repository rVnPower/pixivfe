package i18n

import (
	"errors"
	"fmt"
)

func Error(text string) error {
	text = Tr(text)
	return errors.New(text)
}

func Errorf(format string, a ...any) error {
	format = Tr(format)
	return fmt.Errorf(format, a...)
}

func Sprintf(format string, a ...any) string {
	format = Tr(format)
	return fmt.Sprintf(format, a...)
}

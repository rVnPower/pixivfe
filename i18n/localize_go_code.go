package i18n

import (
	"errors"
	"fmt"
)

func Error(text string) error {
	text = lookup(text)
	return errors.New(text)
}

func Errorf(format string, a ...any) error {
	format = lookup(format)
	return fmt.Errorf(format, a...)
}

func Sprintf(format string, a ...any) string {
	format = lookup(format)
	return fmt.Sprintf(format, a...)
}

// translate string
func Tr(text string) string {
	return lookup(text)
}


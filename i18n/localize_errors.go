package i18n

import (
	"errors"
	"fmt"
)

func Error(text string) error {
	text = Lookup(text)
	return errors.New(text)
}

func Errorf(format string, a ...any) error {
	format = Lookup(format)
	return fmt.Errorf(format, a...)
}

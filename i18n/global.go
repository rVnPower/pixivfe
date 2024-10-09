package i18n

import (
	"errors"
	"fmt"

	"github.com/timandy/routine"
)

var goroutine_locale = routine.NewInheritableThreadLocal[string]()

const BaseLocale = "en"

func GetLocale() string {
	locale := goroutine_locale.Get()
	if locale == "" {
		locale = BaseLocale
	}
	return locale
}

func SetLocale(locale string) {
	goroutine_locale.Set(locale)
}

func Error(text string) error {
	text = __lookup_skip_stack_2(GetLocale(), text)
	return errors.New(text)
}

func Errorf(format string, a ...any) error {
	format = __lookup_skip_stack_2(GetLocale(), format)
	return fmt.Errorf(format, a...)
}

func Sprintf(format string, a ...any) string {
	format = __lookup_skip_stack_2(GetLocale(), format)
	return fmt.Sprintf(format, a...)
}

// translate string
func Tr(text string) string {
	return __lookup_skip_stack_2(GetLocale(), text)
}

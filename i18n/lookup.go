package i18n

import "runtime"

// translate string
func lookup(text string) string {
	pc, file, line, ok := runtime.Caller(2) // user function -> i18n.XXX("...") -> lookup

	println("recorded stackframe:", pc, file, line, ok, text)

	return text
}

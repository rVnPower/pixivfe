package i18n

import (
	"encoding/base64"
	"encoding/binary"
	"runtime"

	"github.com/zeebo/xxh3"
)

var locales = map[string]map[string]string{}

func init() {
	// todo: load locales into into `locales`
}

type Localizer struct {
	locale string
}

func LocalizerOf(locale string) Localizer {
	return Localizer{locale: locale}
}

// lookup string in the database
func (l Localizer) lookup(text string) string {
	translation_map, exist := locales[l.locale]
	if !exist {
		return text
	}

	_, file, _, ok := runtime.Caller(2) // user function -> i18n.XXX("...") -> lookup
	// file and line is correct
	// println("recorded stackframe:", pc, file, line, ok, text)
	if !ok {
		return text
	}

	translation, exist := translation_map[SuccintId(file, text)]
	if !exist {
		return text
	}

	return translation
}

// pad stack frame by 1
func (l Localizer) Tr(text string) string {
	return l.lookup(text)
}

func SuccintId(file string, text string) string {
	hash := xxh3.HashString(text)
	hash_bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(hash_bytes, uint64(hash))
	digest := base64.RawURLEncoding.EncodeToString(hash_bytes)
	return file + ":" + digest
}

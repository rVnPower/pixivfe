package i18n

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"io/fs"
	"log"
	"maps"
	"os"
	"path"
	"runtime"

	"github.com/zeebo/xxh3"
)

var locales = map[string]map[string]string{}

func init() {
	fs_i18n := os.DirFS("i18n/locale")
	entries, err := fs.ReadDir(fs_i18n, ".")
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if !entry.Type().IsDir() {
			continue
		}
		locales[entry.Name()], err = loadLocale(fs_i18n, entry.Name())
		if err != nil {
			panic(err)
		}
		log.Printf("Loaded locale %s", entry.Name())
	}
}

func loadLocale(fs_i18n fs.FS, locale string) (map[string]string, error) {
	m0, err := loadLocale_helper(fs_i18n, locale, "code.json")
	if err != nil {
		return nil, err
	}
	m1, err := loadLocale_helper(fs_i18n, locale, "template.json")
	if err != nil {
		return nil, err
	}
	maps.Copy(m0, m1)
	return m0, nil
}
func loadLocale_helper(fs_i18n fs.FS, locale string, filename string) (map[string]string, error) {
	file, err := fs_i18n.Open(path.Join(locale, filename))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var translated_strings map[string]string
	err = json.NewDecoder(file).Decode(&translated_strings)
	if err != nil {
		return nil, err
	}
	return translated_strings, nil
}

// (internal, do not use directly) lookup string in the database
// call stack should look like this: caller (in the correct file) -> another function -> __lookup_skip_stack_2
func __lookup_skip_stack_2(locale string, text string) string {
	if locale == BaseLocale {
		return text
	}

	translation_map, exist := locales[locale]
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

func SuccintId(file string, text string) string {
	hash := xxh3.HashString(text)
	hash_bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(hash_bytes, uint64(hash))
	digest := base64.RawURLEncoding.EncodeToString(hash_bytes)
	return file + ":" + digest
}

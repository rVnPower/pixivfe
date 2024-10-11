// rewrite jet template before it went to be parsed by jet

package i18n

import (
	"slices"
	"strings"
)

// func RewriteString(file string, content string) string {
// 	locale := getLocale()
// 	if locale == BaseLocale {
// 		return content
// 	}
// 	replacer := TranslationReplacer(locale, file)
// 	return replacer.Replace(content)
// }

// format: old0, new0, old1, new1, ...
type TrPairs = []string

type cacheKey = struct {
	locale string
	file   string
}

var tm = map[cacheKey]*strings.Replacer{}

// returns nil when nothing need to be replaced
func Replacer(locale string, file string) *strings.Replacer {
	if locale == BaseLocale {
		return nil
	}
	k := cacheKey{locale: locale, file: file}
	v, exist := tm[k]
	if exist {
		return v
	}
	pairs := translationPairs_inner(locale, file)
	if len(pairs) == 0 {
		v = nil
	} else {
		v = strings.NewReplacer(pairs...)
	}
	tm[k] = v
	return v
}

type translation_pair = struct {
	before string
	after  string
}

func translationPairs_inner(locale string, file string) TrPairs {
	from_map := locales[BaseLocale]
	to_map, exist := locales[locale]
	if !exist {
		return TrPairs{}
	}

	staging := []translation_pair{}
	for k, v := range to_map {
		if strings.HasPrefix(k, file+":") && from_map[k] != v {
			staging = append(staging, translation_pair{from_map[k], v})
		}
	}
	// sort by length. longest first. this is to prevent weird stuff when rewriting multiple strings and a short one is substring of a long one.
	slices.SortFunc(staging, func(a translation_pair, b translation_pair) int {
		return len(b.before) - len(a.before)
	})

	result := TrPairs{}
	for _, v := range staging {
		result = append(result, v.before)
		result = append(result, v.after)
	}
	return result
}

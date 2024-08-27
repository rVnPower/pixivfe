package template

import (
	"fmt"
	"strings"
)

type PartialURL struct {
	Path  string
	Query map[string]string
}

func LowercaseFirstChar(s string) string {
	return strings.ToLower(s[0:1]) + s[1:]
}

// Turn `url` into /path?other_key=other_value&`key`=
func UnfinishedQuery(url PartialURL, key string) string {
	result := fmt.Sprintf("/%s", url.Path)
	first_query_pair := true

	for k, v := range url.Query {
		k = LowercaseFirstChar(k)

		if k == key {
			continue
		}

		if v == "" {
			// If the value is empty, ignore to not clutter the URL
			continue
		}

		if first_query_pair {
			result += "?"
			first_query_pair = false
		} else {
			result += "&"
		}
		result += fmt.Sprintf("%s=%s", k, v)
	}

	// This is to move the matched query to the end of the URL
	var t string
	if first_query_pair {
		t = "?"
	} else {
		t = "&"
	}
	result += fmt.Sprintf("%s%s=", t, key)

	return result
}

func ReplaceQuery(url PartialURL, key string, value string) string {
	return UnfinishedQuery(url, key) + value
}

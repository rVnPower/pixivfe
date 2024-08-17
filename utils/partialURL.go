package utils

import "fmt"

type PartialURL struct {
	Path string
	Query map[string]string
}

// Turn `url` into /path?other_key=other_value&`key`=
func unfinishedQuery(url PartialURL, key string) string {
	result := fmt.Sprintf("/%s", url.Path)
	first_query_pair := true
	query_param_exists := false

	for k, v := range url.Query {
		k = lowercaseFirstChar(k)

		if k == key {
			// Reserve this
			query_param_exists = true
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
	if query_param_exists {
		var t string
		if first_query_pair {
			t = "?"
		} else {
			t = "&"
		}
		result += fmt.Sprintf("%s%s=", t, key)
	} else {
		// todo: what now? if it doesn't exist, it's a no-op? that doesn't make sense given how navigation works.
	}

	return result
}

func replaceQuery(url PartialURL, key string, value string) string {
	return unfinishedQuery(url, key) + value
}

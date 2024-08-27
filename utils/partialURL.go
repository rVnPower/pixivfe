package utils

import "fmt"

type PartialURL struct {
	Path  string
	Query map[string]string
}

// Turn `url` into /path?other_key=other_value&`key`=
func unfinishedQuery(url PartialURL, key string) string {
	result := fmt.Sprintf("/%s", url.Path)
	first_query_pair := true

	for k, v := range url.Query {
		k = lowercaseFirstChar(k)

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

func replaceQuery(url PartialURL, key string, value string) string {
	return unfinishedQuery(url, key) + value
}

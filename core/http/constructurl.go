package core

import (
	"fmt"
	"strings"
)

type URLConstructor struct {
	path string
	hash map[string]string
}

func lowercaseFirstChar(s string) string {
	return strings.ToLower(s[0:1]) + s[1:]
}

func (obj *URLConstructor) Replace(name string) string {
	url := fmt.Sprintf("/%s", obj.path)
	first := true

	for k, v := range obj.hash {
		k = lowercaseFirstChar(k)

		if k == name {
			continue
		}
		if first {
			url += "?"
			first = false
		} else {
			url += "&"
		}
		url += fmt.Sprintf("%s=%s", k, v)
	}

	// This is to move the matched query to the end of the URL
	var t string
	if first {
		t = "?"
	} else {
		t = "&"
	}
	url += fmt.Sprintf("%s%s=", t, name)

	return url
}

func NewURLConstruct(path string, obj map[string]string) URLConstructor {
	return URLConstructor{
		path: path,
		hash: obj,
	}
}

package core

import "fmt"

type URLConstructor struct {
	path string
	hash map[string]string
}

func (obj *URLConstructor) Replace(name, value string) string {
	url := fmt.Sprintf("/%s", obj.path)
	first := true

	for k, v := range obj.hash {
		if first {
			url += "?"
			first = false
		} else {
			url += "&"
		}
		if k == name {
			v = value
		}
		url += fmt.Sprintf("%s=%s", k, v)
	}

	return url
}

func NewURLConstruct(path string, obj map[string]string) URLConstructor {
	return URLConstructor{
		path: path,
		hash: obj,
	}
}

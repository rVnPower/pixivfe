package core

type URLConstructor struct {
	Path string
	Hash map[string]string
}

func NewURLConstruct(path string, obj map[string]string) URLConstructor {
	return URLConstructor{
		Path: path,
		Hash: obj,
	}
}

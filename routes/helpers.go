package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)


func GetQueryParam(r *http.Request, name string, defaultValue ...string) string {
	if v := r.URL.Query().Get(name); v != "" {
		return v
	} else {
		if len(defaultValue) == 0 {
			return ""
		} else {
			return defaultValue[0]
		}
	}
}

// get path segment. no idea why it's called "params"
func GetPathVar(r *http.Request, name string, defaultValue ...string) string {
	if v := mux.Vars(r)[name]; v != "" {
		return v
	} else {
		if len(defaultValue) == 0 {
			return ""
		} else {
			return defaultValue[0]
		}
	}
}

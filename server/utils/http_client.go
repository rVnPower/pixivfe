package utils

import "net/http"

var HttpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 20,
	},
}

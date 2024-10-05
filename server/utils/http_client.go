package utils

import "net/http"

// HttpClient is a pre-configured http.Client.
// It serves as a base HTTP client used across different packages.
var HttpClient = &http.Client{
	Transport: &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		MaxIdleConns:        0,
		MaxIdleConnsPerHost: 20,
	},
}

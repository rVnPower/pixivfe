package config

const BuiltinProxyUrl = "/proxy/i.pximg.net" // built-in proxy route

// the list of proxies on /settings
var BuiltinProxyList = []string{
	// !!!! WE ARE NOT AFFILIATED WITH MOST OF THE PROXIES !!!!
	"https://pximg.exozy.me", // except this one. this one we are affiliated with.
	"https://pixiv.ducks.party",
	"https://pximg.cocomi.eu.org",
	// "https://mima.localghost.org/proxy/pximg", // doesn't support HTTP/1.1. only support HTTP/2. need proxy code to use a http client (not http.Client) that supports HTTP/2.
	"https://pixiv.darkness.services",
	"https://i.pixiv.re",
	// "https://pixiv.tatakai.top", // dead due to us :(
	// "https://pximg.chaotic.ninja", // incompatible
}

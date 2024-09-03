package config

const BuiltinProxyUrl = "/proxy/i.pximg.net" // built-in proxy route

// the list of proxies on /settings
var BuiltinProxyList = []string{
	// !!!! WE ARE NOT AFFILIATED WITH MOST OF THE PROXIES !!!!
	"https://pximg.exozy.me", // except this one. this one we are affiliated with.
	"https://pixiv.ducks.party",
	"https://pximg.cocomi.eu.org",
	"https://i.suimoe.com",
	"https://i.yuki.sh",
	"https://pximg.obfs.dev",
	"https://pixiv.darkness.services",
	"https://pixiv.tatakai.top",
	"https://pi.169889.xyz",
	"https://i.pixiv.re",
	// "https://mima.localghost.org/proxy/pximg", // doesn't support HTTP/1.1. only support HTTP/2. need proxy code to use a http client (not http.Client) that supports HTTP/2.
	// "https://pximg.chaotic.ninja", // incompatible

	// VnPower: Please comment non-working sites instead of deleting them.
}

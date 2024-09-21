// this file is used by http handlers. no need to refactor when not adding new features.

package session

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"codeberg.org/vnpower/pixivfe/v2/config"
)

// GetUserToken retrieves the authentication token for the Pixiv API from the 'pixivfe-Token' cookie.
// This token takes precedence over the default one provided by tokenManager.
func GetUserToken(r *http.Request) string {
	return GetCookie(r, Cookie_Token)
}

func GetImageProxy(r *http.Request) url.URL {
	value := GetCookie(r, Cookie_ImageProxy)
	if value == "" {
		// fall through to default case
	} else {
		proxyUrl, err := url.Parse(value)
		if err != nil {
			// fall through to default case
		} else {
			return *proxyUrl
		}
	}
	return config.GlobalConfig.ProxyServer
}

func ProxyImageUrl(r *http.Request, s string) string {
	proxyOrigin := GetImageProxyPrefix(r)
	s = strings.ReplaceAll(s, `https:\/\/i.pximg.net`, proxyOrigin)
	// s = strings.ReplaceAll(s, `https:\/\/i.pximg.net`, "/proxy/i.pximg.net")
	s = strings.ReplaceAll(s, `https:\/\/s.pximg.net`, "/proxy/s.pximg.net")
	return s
}

func ProxyImageUrlNoEscape(r *http.Request, s string) string {
	proxyOrigin := GetImageProxyPrefix(r)
	s = strings.ReplaceAll(s, `https://i.pximg.net`, proxyOrigin)
	// s = strings.ReplaceAll(s, `https:\/\/i.pximg.net`, "/proxy/i.pximg.net")
	s = strings.ReplaceAll(s, `https://s.pximg.net`, "/proxy/s.pximg.net")
	return s
}

func GetImageProxyOrigin(r *http.Request) string {
	url := GetImageProxy(r)
	return urlAuthority(url)
}

func GetImageProxyPrefix(r *http.Request) string {
	url := GetImageProxy(r)
	return urlAuthority(url) + url.Path
	// note: not sure if url.EscapedPath() is useful here. go's standard library is trash at handling URL (:// should be part of the scheme)
}

// note: still cannot believe Go doesn't have this function built-in. if stability is their goal, they really don't have the incentive to add useful, crucial features
func urlAuthority(url url.URL) string {
	r := ""
	if (url.Scheme != "") != (url.Host != "") {
		log.Panicf("url must have both scheme and authority or neither: %s", url.String())
	}
	if url.Scheme != "" {
		r += url.Scheme + "://"
	}
	r += url.Host
	return r
}

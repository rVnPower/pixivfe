// User Settings (Using Browser Cookies)

package session

import (
	"time"

	"net/http"
)

type CookieName string

const ( // the __Host thing force it to be secure and same-origin (no subdomain) >> https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie
	Cookie_Token             CookieName = "pixivfe-Token"
	Cookie_CSRF              CookieName = "pixivfe-CSRF"
	Cookie_ImageProxy        CookieName = "pixivfe-ImageProxy"
	Cookie_NovelFontType     CookieName = "pixivfe-NovelFontType"
	Cookie_NovelViewMode     CookieName = "pixivfe-NovelViewMode"
	Cookie_ThumbnailToNewTab CookieName = "pixivfe-ThumbnailToNewTab"
	Cookie_ArtworkPreview    CookieName = "pixivfe-ArtworkPreview"
	Cookie_ShowArtR18        CookieName = "pixivfe-ShowArtR18"
	Cookie_ShowArtR18G       CookieName = "pixivfe-ShowArtR18G"
	Cookie_ShowArtAI         CookieName = "pixivfe-ShowArtAI"
)

// Go can't make this a const...
var AllCookieNames []CookieName = []CookieName{
	Cookie_Token,
	Cookie_CSRF,
	Cookie_ImageProxy,
	Cookie_NovelFontType,
	Cookie_NovelViewMode,
	Cookie_ThumbnailToNewTab,
	Cookie_ArtworkPreview,
	Cookie_ShowArtR18,
	Cookie_ShowArtR18G,
	Cookie_ShowArtAI,
}

func GetCookie(r *http.Request, name CookieName) string {
	cookie, err := r.Cookie(string(name))
	if err != nil {
		return ""
	}
	return cookie.Value
}

func SetCookie(w http.ResponseWriter, name CookieName, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:  string(name),
		Value: value,
		Path:  "/",
		// expires in 30 days from now
		Expires:  time.Now().Add(30 * (24 * time.Hour)),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode, // bye-bye cross site forgery
	})
}

var CookieExpireDelete = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

func ClearCookie(w http.ResponseWriter, name CookieName) {
	http.SetCookie(w, &http.Cookie{
		Name:  string(name),
		Value: "",
		Path:  "/",
		// expires in 30 days from now
		Expires:  CookieExpireDelete,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func ClearAllCookies(w http.ResponseWriter) {
	for _, name := range AllCookieNames {
		ClearCookie(w, name)
	}
}

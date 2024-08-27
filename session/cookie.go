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

func GetCookie(c *http.Request, name CookieName, defaultValue ...string) string {
	return c.Cookies(string(name), defaultValue...)
}

func SetCookie(c *http.Request, name CookieName, value string) {
	cookie := fiber.Cookie{
		Name:  string(name),
		Value: value,
		Path:  "/",
		// expires in 30 days from now
		Expires:  c.Context().Time().Add(30 * (24 * time.Hour)),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode, // bye-bye cross site forgery
	}
	c.Cookie(&cookie)
}

var CookieExpireDelete = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

func ClearCookie(c *http.Request, name CookieName) {
	cookie := fiber.Cookie{
		Name:  string(name),
		Value: "",
		Path:  "/",
		// expires in 30 days from now
		Expires:  CookieExpireDelete,
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
	}
	c.Cookie(&cookie)
}

func ClearAllCookies(c *http.Request) {
	for _, name := range AllCookieNames {
		ClearCookie(c, name)
	}
}

package middleware

import (
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/server/session"
)

func SetDefaultCookies(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if it's a first visit
		_, err := r.Cookie("pixivfe-visited")
		if err == http.ErrNoCookie {
			// It's a first visit, set default cookies
			defaultCookies := config.GlobalConfig.DefaultCookies
			session.SetCookie(w, session.Cookie_ImageProxy, defaultCookies.ImageProxy)
			session.SetCookie(w, session.Cookie_NovelFontType, defaultCookies.NovelFontType)
			session.SetCookie(w, session.Cookie_NovelViewMode, defaultCookies.NovelViewMode)
			session.SetCookie(w, session.Cookie_ThumbnailToNewTab, defaultCookies.ThumbnailToNewTab)
			session.SetCookie(w, session.Cookie_ArtworkPreview, defaultCookies.ArtworkPreview)
			session.SetCookie(w, session.Cookie_HideArtR18, defaultCookies.HideArtR18)
			session.SetCookie(w, session.Cookie_HideArtR18G, defaultCookies.HideArtR18G)
			session.SetCookie(w, session.Cookie_HideArtAI, defaultCookies.HideArtAI)

			// Set the visited cookie
			session.SetCookie(w, "pixivfe-visited", "true")
		}

		next.ServeHTTP(w, r)
	})
}

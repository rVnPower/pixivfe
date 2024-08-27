package routes

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"codeberg.org/vnpower/pixivfe/v2/config"
	httpc "codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/session"
)

// todo: allow clear proxy
// todo: allow clear all settings

func setToken(w http.ResponseWriter, r CompatRequest) error {
	token := r.FormValue("token")
	if token != "" {
		URL := httpc.GetNewestFromFollowingURL("all", "1")

		_, err := httpc.UnwrapWebAPIRequest(r.Context(), URL, token)
		if err != nil {
			return errors.New("Cannot authorize with supplied token.")
		}

		// Make a test request to verify the token.
		// THE TEST URL IS NSFW!
		req, err := http.NewRequest("GET", "https://www.pixiv.net/en/artworks/115365120", nil)
		if err != nil {
			return err
		}
		req = req.WithContext(r.Context())
		req.Header.Add("User-Agent", "Mozilla/5.0")
		req.AddCookie(&http.Cookie{
			Name:  "PHPSESSID",
			Value: token,
		})

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return errors.New("Cannot authorize with supplied token.")
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.New("Cannot parse the response from Pixiv. Please report this issue.")
		}

		// CSRF token
		r := regexp.MustCompile(`"token":"([0-9a-f]+)"`)
		csrf := r.FindStringSubmatch(string(body))[1]

		if csrf == "" {
			return errors.New("Cannot authorize with supplied token.")
		}

		// Set the token
		session.SetCookie(w, session.Cookie_Token, token)
		session.SetCookie(w, session.Cookie_CSRF, csrf)

		return nil
	}
	return errors.New("You submitted an empty/invalid form.")
}

func setImageServer(w http.ResponseWriter, r CompatRequest) error {
	token := r.FormValue("image-proxy")
	if token != "" {
		session.SetCookie(w, session.Cookie_ImageProxy, token)
	} else {
		session.ClearCookie(r, session.Cookie_ImageProxy)
	}
	return nil
}

func setNovelFontType(w http.ResponseWriter, r CompatRequest) error {
	fontType := r.FormValue("font-type")
	if fontType != "" {
		session.SetCookie(w, session.Cookie_NovelFontType, fontType)
	}

	return nil
}

func setNovelViewMode(w http.ResponseWriter, r CompatRequest) error {
	viewMode := r.FormValue("view-mode")
	if viewMode != "" {
		session.SetCookie(w, session.Cookie_NovelViewMode, viewMode)
	}

	return nil
}

func setThumbnailToNewTab(w http.ResponseWriter, r CompatRequest) error {
	ttnt := r.FormValue("ttnt")
	if ttnt == "_blank" || ttnt == "_self" {
		session.SetCookie(w, session.Cookie_ThumbnailToNewTab, ttnt)
	}

	return nil
}

func setArtworkPreview(w http.ResponseWriter, r CompatRequest) error {
	value := r.FormValue("app")
	if value == "cover" || value == "button" || value == "" {
		session.SetCookie(w, session.Cookie_ArtworkPreview, value)
	}

	return nil
}

func setLogout(w http.ResponseWriter, r CompatRequest) error {
	session.ClearCookie(r, session.Cookie_Token)
	session.ClearCookie(r, session.Cookie_CSRF)
	return nil
}

func setCookie(w http.ResponseWriter, r CompatRequest) error {
	key := r.FormValue("key")
	value := r.FormValue("value")
	for _, cookie_name := range session.AllCookieNames {
		if string(cookie_name) == key {
			session.SetCookie(w, cookie_name, value)
			return nil
		}
	}
	return fmt.Errorf("Invalid Cookie Name: %s", key)
}

func setRawCookie(w http.ResponseWriter, r CompatRequest) error {
	raw := r.FormValue("raw")
	lines := strings.Split(raw, "\n")

	for _, line := range lines {
		sub := strings.Split(line, "=")
		if len(sub) != 2 {
			continue
		}

		name := session.CookieName(sub[0])
		value := sub[1]

		if !slices.Contains(session.AllCookieNames, name) {
			continue
		}

		session.SetCookie(w, name, value)
	}
	return nil
}

func resetAll(w http.ResponseWriter, r CompatRequest) error {
	session.ClearAllCookies(w)
	return nil
}

func SettingsPage(w http.ResponseWriter, r CompatRequest) error {
	return Render(w, r, Data_settings{WorkingProxyList: config.GetWorkingProxies(), ProxyList: config.BuiltinProxyList})
}

func SettingsPost(w http.ResponseWriter, r CompatRequest) error {
	t := r.Params("type")
	noredirect := r.FormValue("noredirect", "") == ""
	var err error

	switch t {
	case "imageServer":
		err = setImageServer(r)
	case "token":
		err = setToken(r)
	case "logout":
		err = setLogout(r)
	case "reset-all":
		err = resetAll(r)
	case "novelFontType":
		err = setNovelFontType(r)
	case "thumbnailToNewTab":
		err = setThumbnailToNewTab(r)
	case "novelViewMode":
		err = setNovelViewMode(r)
	case "artworkPreview":
		err = setArtworkPreview(r)
	case "set-cookie":
		err = setCookie(r)
	case "raw":
		err = setRawCookie(r)
	default:
		err = errors.New("No such setting is available.")
	}

	if err != nil {
		return err
	}

	if !noredirect {
		return nil
	}

	return r.Redirect("/settings", http.StatusSeeOther)
}

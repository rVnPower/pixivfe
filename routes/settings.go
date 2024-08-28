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
	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/session"
	"codeberg.org/vnpower/pixivfe/v2/utils"
)

func setToken(w http.ResponseWriter, r *http.Request) error {
	token := r.FormValue("token")
	if token != "" {
		URL := core.GetNewestFromFollowingURL("all", "1")

		_, err := core.UnwrapWebAPIRequest(r.Context(), URL, token)
		if err != nil {
			return errors.New("Cannot authorize with supplied token.")
		}

		// Make a test request to verify the token.
		// THE TEST URL IS NSFW!
		req, err := http.NewRequestWithContext(r.Context(), "GET", "https://www.pixiv.net/en/artworks/115365120", nil)
		if err != nil {
			return err
		}
		req.Header.Add("User-Agent", "Mozilla/5.0")
		req.AddCookie(&http.Cookie{
			Name:  "PHPSESSID",
			Value: token,
		})

		resp, err := utils.HttpClient.Do(req)
		if err != nil {
			return errors.New("Cannot authorize with supplied token.")
		}
		defer resp.Body.Close()

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

func setImageServer(w http.ResponseWriter, r *http.Request) error {
	token := r.FormValue("image-proxy")
	if token != "" {
		session.SetCookie(w, session.Cookie_ImageProxy, token)
	} else {
		session.ClearCookie(w, session.Cookie_ImageProxy)
	}
	return nil
}

func setNovelFontType(w http.ResponseWriter, r *http.Request) error {
	fontType := r.FormValue("font-type")
	if fontType != "" {
		session.SetCookie(w, session.Cookie_NovelFontType, fontType)
	}

	return nil
}

func setNovelViewMode(w http.ResponseWriter, r *http.Request) error {
	viewMode := r.FormValue("view-mode")
	if viewMode == "1" || viewMode == "2" || viewMode == "" {
		session.SetCookie(w, session.Cookie_NovelViewMode, viewMode)
	}

	return nil
}

func setThumbnailToNewTab(w http.ResponseWriter, r *http.Request) error {
	ttnt := r.FormValue("ttnt")
	if ttnt == "_blank" || ttnt == "_self" {
		session.SetCookie(w, session.Cookie_ThumbnailToNewTab, ttnt)
	}

	return nil
}

func setArtworkPreview(w http.ResponseWriter, r *http.Request) error {
	value := r.FormValue("app")
	if value == "cover" || value == "button" || value == "" {
		session.SetCookie(w, session.Cookie_ArtworkPreview, value)
	}

	return nil
}

func setLogout(w http.ResponseWriter, _ *http.Request) error {
	session.ClearCookie(w, session.Cookie_Token)
	session.ClearCookie(w, session.Cookie_CSRF)
	return nil
}

func setCookie(w http.ResponseWriter, r *http.Request) error {
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

func setRawCookie(w http.ResponseWriter, r *http.Request) error {
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

func resetAll(w http.ResponseWriter, _ *http.Request) error {
	session.ClearAllCookies(w)
	return nil
}

func SettingsPage(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, Data_settings{WorkingProxyList: config.GetWorkingProxies(), ProxyList: config.BuiltinProxyList})
}

func SettingsPost(w http.ResponseWriter, r *http.Request) error {
	t := GetPathVar(r, "type")
	var err error

	switch t {
	case "imageServer":
		err = setImageServer(w, r)
	case "token":
		err = setToken(w, r)
	case "logout":
		err = setLogout(w, r)
	case "reset-all":
		err = resetAll(w, r)
	case "novelFontType":
		err = setNovelFontType(w, r)
	case "thumbnailToNewTab":
		err = setThumbnailToNewTab(w, r)
	case "novelViewMode":
		err = setNovelViewMode(w, r)
	case "artworkPreview":
		err = setArtworkPreview(w, r)
	case "set-cookie":
		err = setCookie(w, r)
	case "raw":
		err = setRawCookie(w, r)
	default:
		err = errors.New("No such setting is available.")
	}

	if err != nil {
		return err
	}

	utils.RedirectToWhenceYouCame(w, r)
	return nil
}

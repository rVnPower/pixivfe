package routes

import (
	"bufio"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/i18n"
	"codeberg.org/vnpower/pixivfe/v2/server/proxy_checker"
	"codeberg.org/vnpower/pixivfe/v2/server/session"
	"codeberg.org/vnpower/pixivfe/v2/server/utils"
)

var r_csrf = regexp.MustCompile(`"token":"([0-9a-f]+)"`)

func setToken(w http.ResponseWriter, r *http.Request) (string, error) {
	token := r.FormValue("token")
	if token != "" {
		URL := core.GetNewestFromFollowingURL("all", "1")

		_, err := core.API_GET_UnwrapJson(r.Context(), URL, token)
		if err != nil {
			return "", i18n.Error("Cannot authorize with supplied token.")
		}

		// Make a test request to verify the token.
		// THE TEST URL IS NSFW!
		resp, err := core.API_GET(r.Context(), "https://www.pixiv.net/en/artworks/115365120", token)
		if err != nil {
			return "", err
		}

		if resp.StatusCode != 200 {
			return "", i18n.Error("Cannot authorize with supplied token.")
		}

		// CSRF token
		csrf := r_csrf.FindStringSubmatch(resp.Body)[1]

		if csrf == "" {
			return "", i18n.Error("Cannot authorize with supplied token.")
		}

		// Set the token
		session.SetCookie(w, session.Cookie_Token, token)
		session.SetCookie(w, session.Cookie_CSRF, csrf)

		return i18n.Sprintf("Successfully logged in."), nil
	}
	return "", i18n.Error("You submitted an empty/invalid form.")
}

func setImageServer(w http.ResponseWriter, r *http.Request) (string, error) {
	customProxy := r.FormValue("custom-image-proxy")
	selectedProxy := r.FormValue("image-proxy")

	if customProxy != "" {
		session.SetCookie(w, session.Cookie_ImageProxy, customProxy)
		return i18n.Sprintf("Custom image proxy server set successfully."), nil
	} else if selectedProxy != "" {
		session.SetCookie(w, session.Cookie_ImageProxy, selectedProxy)
		return i18n.Sprintf("Image proxy server updated successfully."), nil
	} else {
		session.ClearCookie(w, session.Cookie_ImageProxy)
		return i18n.Sprintf("Image proxy server cleared."), nil
	}
}

func setNovelFontType(w http.ResponseWriter, r *http.Request) (string, error) {
	fontType := r.FormValue("font-type")
	if fontType != "" {
		session.SetCookie(w, session.Cookie_NovelFontType, fontType)
		return i18n.Sprintf("Novel font type updated successfully."), nil
	}

	return "", i18n.Error("Invalid font type.")
}

func setNovelViewMode(w http.ResponseWriter, r *http.Request) (string, error) {
	viewMode := r.FormValue("view-mode")
	if viewMode == "1" || viewMode == "2" || viewMode == "" {
		session.SetCookie(w, session.Cookie_NovelViewMode, viewMode)
		return i18n.Sprintf("Novel view mode updated successfully."), nil
	}

	return "", i18n.Error("Invalid view mode.")
}

func setThumbnailToNewTab(w http.ResponseWriter, r *http.Request) (string, error) {
	ttnt := r.FormValue("ttnt")
	if ttnt == "_blank" {
		session.SetCookie(w, session.Cookie_ThumbnailToNewTab, ttnt)
		return i18n.Sprintf("Thumbnails will now open in a new tab."), nil
	} else {
		session.SetCookie(w, session.Cookie_ThumbnailToNewTab, "_self")
		return i18n.Sprintf("Thumbnails will now open in the same tab."), nil
	}
}

func setArtworkPreview(w http.ResponseWriter, r *http.Request) (string, error) {
	value := r.FormValue("app")
	if value == "cover" || value == "button" || value == "" {
		session.SetCookie(w, session.Cookie_ArtworkPreview, value)
		return i18n.Sprintf("Artwork preview setting updated successfully."), nil
	}

	return "", i18n.Error("Invalid artwork preview setting.")
}

func setFilter(w http.ResponseWriter, r *http.Request) (string, error) {
	r18 := r.FormValue("filter-r18")
	r18g := r.FormValue("filter-r18g")
	ai := r.FormValue("filter-ai")

	session.SetCookie(w, session.Cookie_HideArtR18, r18)
	session.SetCookie(w, session.Cookie_HideArtR18G, r18g)
	session.SetCookie(w, session.Cookie_HideArtAI, ai)

	return i18n.Sprintf("Filter settings updated successfully."), nil
}

func setLogout(w http.ResponseWriter, _ *http.Request) (string, error) {
	session.ClearCookie(w, session.Cookie_Token)
	session.ClearCookie(w, session.Cookie_CSRF)
	return i18n.Sprintf("Successfully logged out."), nil
}

func setCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	key := r.FormValue("key")
	value := r.FormValue("value")
	for _, cookie_name := range session.AllCookieNames {
		if string(cookie_name) == key {
			session.SetCookie(w, cookie_name, value)
			return i18n.Sprintf("Cookie %s set successfully.", key), nil
		}
	}
	return "", i18n.Errorf("Invalid Cookie Name: %s", key)
}

func clearCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	key := r.FormValue("key")
	for _, cookie_name := range session.AllCookieNames {
		if string(cookie_name) == key {
			session.ClearCookie(w, cookie_name)
			return i18n.Sprintf("Cookie %s cleared successfully.", key), nil
		}
	}
	return "", i18n.Errorf("Invalid Cookie Name: %s", key)
}

func setRawCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	raw := r.FormValue("raw")
	reader := bufio.NewReader(strings.NewReader(raw))
	for {
		line, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if isPrefix {
			return "", bufio.ErrBufferFull
		}
		if err != nil {
			return "", err
		}

		sub := strings.Split(string(line), "=")
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
	return i18n.Sprintf("Raw settings applied successfully."), nil
}

func resetAll(w http.ResponseWriter, _ *http.Request) (string, error) {
	session.ClearAllCookies(w)
	return i18n.Sprintf("All preferences have been reset to default values."), nil
}

func SettingsPage(w http.ResponseWriter, r *http.Request) error {
	return RenderHTML(w, r, Data_settings{
		WorkingProxyList:   proxy_checker.GetWorkingProxies(),
		ProxyList:          config.BuiltinProxyList,
		ProxyCheckEnabled:  config.GlobalConfig.ProxyCheckEnabled,    // Used to check whether proxy_checker is enabled on the instance
		ProxyCheckInterval: config.GlobalConfig.ProxyCheckInterval,   // Used to display the ProxyCheckInterval configured on the instance
		DefaultProxyServer: config.GlobalConfig.ProxyServer.String(), // Used to display the default image proxy server
	})
}

func handleAjaxResponse(w http.ResponseWriter, message string, err error) {
	w.Header().Set("Content-Type", "text/html")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(i18n.Sprintf(`<div class="alert alert-danger alert-dismissible fade show mt-3" role="alert">%s<button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button></div>`, err.Error())))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(i18n.Sprintf(`<div class="alert alert-success alert-dismissible fade show mt-3" role="alert">%s<button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button></div>`, message)))
	}
}

func SettingsPost(w http.ResponseWriter, r *http.Request) error {
	t := GetPathVar(r, "type")
	var message string
	var err error

	switch t {
	case "imageServer":
		message, err = setImageServer(w, r)
	case "token":
		message, err = setToken(w, r)
	case "logout":
		message, err = setLogout(w, r)
	case "reset-all":
		message, err = resetAll(w, r)
	case "novelFontType":
		message, err = setNovelFontType(w, r)
	case "thumbnailToNewTab":
		message, err = setThumbnailToNewTab(w, r)
	case "novelViewMode":
		message, err = setNovelViewMode(w, r)
	case "artworkPreview":
		message, err = setArtworkPreview(w, r)
	case "filter":
		message, err = setFilter(w, r)
	case "set-cookie":
		message, err = setCookie(w, r)
	case "clear-cookie":
		message, err = clearCookie(w, r)
	case "raw":
		message, err = setRawCookie(w, r)
	default:
		err = i18n.Error("No such setting is available.")
	}

	if r.Header.Get("HX-Request") == "true" {
		handleAjaxResponse(w, message, err)
		return nil
	}

	if err != nil {
		return err
	}

	utils.RedirectToWhenceYouCame(w, r)
	return nil
}

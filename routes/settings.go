package pages

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"codeberg.org/vnpower/pixivfe/v2/config"
	httpc "codeberg.org/vnpower/pixivfe/v2/pixiv"
	"codeberg.org/vnpower/pixivfe/v2/session"
	"github.com/gofiber/fiber/v2"
)

// todo: allow clear proxy
// todo: allow clear all settings

func setToken(c *fiber.Ctx) error {
	// Parse the value from the form
	token := c.FormValue("token")
	if token != "" {
		URL := httpc.GetNewestFromFollowingURL("all", "1")

		_, err := httpc.UnwrapWebAPIRequest(c.Context(), URL, token)
		if err != nil {
			return errors.New("Cannot authorize with supplied token.")
		}

		// Make a test request to verify the token.
		// THE TEST URL IS NSFW!
		req, err := http.NewRequest("GET", "https://www.pixiv.net/en/artworks/115365120", nil)
		if err != nil {
			return err
		}
		req = req.WithContext(c.Context())
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
		session.SetCookie(c, session.Cookie_Token, token)
		session.SetCookie(c, session.Cookie_CSRF, csrf)

		return nil
	}
	return errors.New("You submitted an empty/invalid form.")
}

func setImageServer(c *fiber.Ctx) error {
	// Parse the value from the form
	token := c.FormValue("image-proxy")
	if token != "" {
		session.SetCookie(c, session.Cookie_ImageProxy, token)
	} else {
		session.ClearCookie(c, session.Cookie_ImageProxy)
	}
	return nil
}

func setNovelFontType(c *fiber.Ctx) error {
	fontType := c.FormValue("font-type")
	if fontType != "" {
		session.SetCookie(c, session.Cookie_NovelFontType, fontType)
	}

	return nil
}

func setNovelViewMode(c *fiber.Ctx) error {
	viewMode := c.FormValue("view-mode")
	if viewMode != "" {
		session.SetCookie(c, session.Cookie_NovelViewMode, viewMode)
	}

	return nil
}

func setThumbnailToNewTab(c *fiber.Ctx) error {
	ttnt := c.FormValue("ttnt")
	if ttnt == "_blank" || ttnt == "_self" {
		session.SetCookie(c, session.Cookie_ThumbnailToNewTab, ttnt)
	}

	return nil
}

func setArtworkPreview(c *fiber.Ctx) error {
	value := c.FormValue("app")
	if value == "cover" || value == "button" || value == "" {
		session.SetCookie(c, session.Cookie_ArtworkPreview, value)
	}

	return nil
}

func setLogout(c *fiber.Ctx) error {
	session.ClearCookie(c, session.Cookie_Token)
	session.ClearCookie(c, session.Cookie_CSRF)
	return nil
}

func setCookie(c *fiber.Ctx) error {
	key := c.FormValue("key")
	value := c.FormValue("value")
	for _, cookie_name := range session.AllCookieNames {
		if string(cookie_name) == key {
			session.SetCookie(c, cookie_name, value)
			return nil
		}
	}
	return fmt.Errorf("Invalid Cookie Name: %s", key)
}

func setRawCookie(c *fiber.Ctx) error {
	raw := c.FormValue("raw")
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

		session.SetCookie(c, name, value)
	}
	return nil
}

func resetAll(c *fiber.Ctx) error {
	session.ClearAllCookies(c)
	return nil
}

func SettingsPage(c *fiber.Ctx) error {
	return c.Render("settings", fiber.Map{
		"ProxyList": config.BuiltinProxyList,
	})
}

func SettingsPost(c *fiber.Ctx) error {
	// NOTE: VnPower: Future maintainers should leave this function alone.

	t := c.Params("type")
	noredirect := c.FormValue("noredirect", "") == ""
	var err error

	switch t {
	case "imageServer":
		err = setImageServer(c)
	case "token":
		err = setToken(c)
	case "logout":
		err = setLogout(c)
	case "reset-all":
		err = resetAll(c)
	case "novelFontType":
		err = setNovelFontType(c)
	case "thumbnailToNewTab":
		err = setThumbnailToNewTab(c)
	case "novelViewMode":
		err = setNovelViewMode(c)
	case "artworkPreview":
		err = setArtworkPreview(c)
	case "set-cookie":
		err = setCookie(c)
	case "raw":
		err = setRawCookie(c)
	default:
		err = errors.New("No such setting is available.")
	}

	if err != nil {
		return err
	}

	if !noredirect {
		return nil
	}

	return c.Redirect("/settings", http.StatusSeeOther)
}
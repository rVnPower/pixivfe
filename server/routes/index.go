package routes

import (
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/server/session"
	"codeberg.org/vnpower/pixivfe/v2/server/template"
	"codeberg.org/vnpower/pixivfe/v2/server/utils"
)

func IndexPage(w http.ResponseWriter, r *http.Request) error {
	mode := GetQueryParam(r, "mode", "all")
	isLoggedIn := session.GetUserToken(r) != ""

	works, err := core.GetLanding(r, mode, isLoggedIn)
	if err != nil {
		return err
	}

	urlc := template.PartialURL{Path: "", Query: map[string]string{"mode": mode}}

	return RenderHTML(w, r, Data_index{
		Title:    "Landing",
		Data:     *works,
		LoggedIn: isLoggedIn,
		Queries:  urlc,
	})
}

func Oembed(w http.ResponseWriter, r *http.Request) error {
	pageURL := utils.Origin(r)
	artistName := GetQueryParam(r, "a", "")
	artistURL := GetQueryParam(r, "u", "")

	data := map[string]any{
		"version":       "1.0",
		"embed_type":    "rich",
		"provider_name": "PixivFE",
		"provider_url":  pageURL,
		"author_name":   artistName,
		"author_url":    artistURL,
	}

	utils.SendJson(w, data)
	return nil
}

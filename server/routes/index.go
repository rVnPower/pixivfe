package routes

import (
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/server/session"
	"codeberg.org/vnpower/pixivfe/v2/server/utils"
)

func IndexPage(w http.ResponseWriter, r *http.Request) error {
	// If token is set, do the landing request...
	if token := session.GetUserToken(r); token != "" {
		mode := GetQueryParam(r, "mode", "all")

		works, err := core.GetLanding(r, mode)

		if err != nil {
			return err
		}

		return RenderHTML(w, r, Data_index{
			Title:    "Landing",
			Data:     *works,
			LoggedIn: true,
		})
	}

	// ...otherwise, default to today's illustration ranking
	works, err := core.GetRanking(r, "daily", "illust", "", "1")
	if err != nil {
		return err
	}
	return RenderHTML(w, r, Data_index{
		Title:       "Landing",
		NoTokenData: works,
		LoggedIn:    false,
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

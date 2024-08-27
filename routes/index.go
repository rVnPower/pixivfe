package routes

import (
	"encoding/json"
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/session"
)

func IndexPage(w http.ResponseWriter, r CompatRequest) error {
	// If token is set, do the landing request...
	if token := session.GetPixivToken(r.Request); token != "" {
		mode := r.Query("mode", "all")

		works, err := core.GetLanding(r.Request, mode)

		if err != nil {
			return err
		}

		return Render(w, r, Data_index{
			Title:    "Landing",
			Data:     *works,
			LoggedIn: true,
		})
	}

	// ...otherwise, default to today's illustration ranking
	works, err := core.GetRanking(r.Request, "daily", "illust", "", "1")
	if err != nil {
		return err
	}
	return Render(w, r, Data_index{
		Title:       "Landing",
		NoTokenData: works,
		LoggedIn:    false,
	})
}

func Oembed(w http.ResponseWriter, r CompatRequest) error {
	pageURL := r.BaseURL()
	artistName := r.Query("a", "")
	artistURL := r.Query("u", "")

	data := map[string]any{
		"version":       "1.0",
		"embed_type":    "rich",
		"provider_name": "PixivFE",
		"provider_url":  pageURL,
		"author_name":   artistName,
		"author_url":    artistURL,
	}

	return json.NewEncoder(w).Encode(data)
}

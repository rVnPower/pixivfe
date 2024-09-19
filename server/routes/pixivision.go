package routes

import (
	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/server/session"

	"net/http"
)

func PixivisionHomePage(w http.ResponseWriter, r *http.Request) error {
	data, err := core.PixivisionGetHomepage("1", "en")
	if err != nil {
		return err
	}

	for i := range data {
		data[i].Thumbnail = session.ProxyImageUrlNoEscape(r, data[i].Thumbnail)
	}

	return Render(w, r, Data_pixivisionIndex{Data: data})
}

func PixivisionArticlePage(w http.ResponseWriter, r *http.Request) error {
	id := GetPathVar(r, "id")
	data, err := core.PixivisionGetArticle(id, "en")
	if err != nil {
		return err
	}

	data.Thumbnail = session.ProxyImageUrlNoEscape(r, data.Thumbnail)
	for i := range data.Items {
		data.Items[i].Image = session.ProxyImageUrlNoEscape(r, data.Items[i].Image)
		data.Items[i].Avatar = session.ProxyImageUrlNoEscape(r, data.Items[i].Avatar)
	}

	return Render(w, r, Data_pixivisionArticle{Article: data})
}

func PixivisionCategoryPage(w http.ResponseWriter, r *http.Request) error {
	id := GetPathVar(r, "id")
	data, err := core.PixivisionGetCategory(id, "1", "en")
	if err != nil {
		return err
	}

	for i := range data.Articles {
		data.Articles[i].Thumbnail = session.ProxyImageUrlNoEscape(r, data.Articles[i].Thumbnail)
	}

	return Render(w, r, Data_pixivisionCategory{Category: data})
}

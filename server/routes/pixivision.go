package routes

import (
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/server/session"
)

func PixivisionHomePage(w http.ResponseWriter, r *http.Request) error {
	data, err := core.PixivisionGetHomepage(r, "1", "en")
	if err != nil {
		return err
	}

	for i := range data {
		data[i].Thumbnail = session.ProxyImageUrlNoEscape(r, data[i].Thumbnail)
	}

	return RenderHTML(w, r, Data_pixivisionIndex{Data: data})
}

func PixivisionArticlePage(w http.ResponseWriter, r *http.Request) error {
	id := GetPathVar(r, "id")
	data, err := core.PixivisionGetArticle(r, id, "en")
	if err != nil {
		return err
	}

	data.Thumbnail = session.ProxyImageUrlNoEscape(r, data.Thumbnail)
	for i := range data.Items {
		data.Items[i].Image = session.ProxyImageUrlNoEscape(r, data.Items[i].Image)
		data.Items[i].Avatar = session.ProxyImageUrlNoEscape(r, data.Items[i].Avatar)
	}

	return RenderHTML(w, r, Data_pixivisionArticle{Article: data})
}

func PixivisionCategoryPage(w http.ResponseWriter, r *http.Request) error {
	id := GetPathVar(r, "id")
	data, err := core.PixivisionGetCategory(r, id, "1", "en")
	if err != nil {
		return err
	}

	for i := range data.Articles {
		data.Articles[i].Thumbnail = session.ProxyImageUrlNoEscape(r, data.Articles[i].Thumbnail)
	}

	return RenderHTML(w, r, Data_pixivisionCategory{Category: data})
}

func PixivisionTagPage(w http.ResponseWriter, r *http.Request) error {
	id := GetPathVar(r, "id")
	data, err := core.PixivisionGetTag(r, id, "1", "en")
	if err != nil {
		return err
	}

	for i := range data.Articles {
		data.Articles[i].Thumbnail = session.ProxyImageUrlNoEscape(r, data.Articles[i].Thumbnail)
	}

	return RenderHTML(w, r, Data_pixivisionTag{Tag: data})
}

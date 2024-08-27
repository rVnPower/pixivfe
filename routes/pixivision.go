package routes

import (
	"codeberg.org/vnpower/pixivfe/v2/session"
	"codeberg.org/vnpower/pixivision"

	"net/http"
)

func PixivisionHomePage(w http.ResponseWriter, r *http.Request) error {
	data, err := pixivision.GetHomepage()
	if err != nil {
		return err
	}

	for i := range data {
		data[i].Thumbnail = session.ProxyImageUrlNoEscape(r, data[i].Thumbnail)
	}

	return Render(w, r, Data_pixivision_index{Data: data})
}

func PixivisionArticlePage(w http.ResponseWriter, r *http.Request) error {
	id := GetPathVar(r, "id")
	data, err := pixivision.GetArticle(id)
	if err != nil {
		return err
	}

	data.Thumbnail = session.ProxyImageUrlNoEscape(r, data.Thumbnail)
	for i := range data.Items {
		data.Items[i].Image = session.ProxyImageUrlNoEscape(r, data.Items[i].Image)
		data.Items[i].Avatar = session.ProxyImageUrlNoEscape(r, data.Items[i].Avatar)
	}

	return Render(w, r, Data_pixivision_article{Article: data})
}

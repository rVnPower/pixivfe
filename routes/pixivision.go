package routes

import (
	"codeberg.org/vnpower/pixivfe/v2/session"
	"codeberg.org/vnpower/pixivision"

	"net/http"
)

func PixivisionHomePage(c *http.Request) error {
	// Note: don't process images here?
	data, err := pixivision.GetHomepage()
	if err != nil {
		return err
	}

	for i := range data {
		data[i].Thumbnail = session.ProxyImageUrlNoEscape(c, data[i].Thumbnail)
	}

	return Render(c, Data_pixivision_index{Data: data})
}

func PixivisionArticlePage(c *http.Request) error {
	// Note: don't process images here?
	id := c.Params("id")
	data, err := pixivision.GetArticle(id)
	if err != nil {
		return err
	}

	data.Thumbnail = session.ProxyImageUrlNoEscape(c, data.Thumbnail)
	for i := range data.Items {
		data.Items[i].Image = session.ProxyImageUrlNoEscape(c, data.Items[i].Image)
		data.Items[i].Avatar = session.ProxyImageUrlNoEscape(c, data.Items[i].Avatar)
	}

	return Render(c, Data_pixivision_article{Article: data})
}

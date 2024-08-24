package routes

import (
	"codeberg.org/vnpower/pixivision"
	"github.com/gofiber/fiber/v2"

	"codeberg.org/vnpower/pixivfe/v2/session"
)

func PixivisionHomePage(c *fiber.Ctx) error {
	// Note: don't process images here?
	data, err := pixivision.GetHomepage()
	if err != nil {
		return err
	}

	for i := range data {
		data[i].Thumbnail = session.ProxyImageUrlNoEscape(c, data[i].Thumbnail)
	}

	return c.Render("pixivision", fiber.Map{"Data": data})
}

func PixivisionArticlePage(c *fiber.Ctx) error {
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

	return c.Render("pixivision_article", fiber.Map{"Article": data})
}

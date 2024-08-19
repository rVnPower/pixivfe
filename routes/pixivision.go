package pages

import (
	"codeberg.org/vnpower/pixivision"
	"github.com/gofiber/fiber/v2"

	session "codeberg.org/vnpower/pixivfe/v2/session"
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

	return c.Render("pixivision/index", fiber.Map{"Data": data})
}

func PixivisionArticlePage(c *fiber.Ctx) error {
	// Note: don't process images here?
	id := c.Params("id")
	data, err := pixivision.GetArticle(id)
	if err != nil {
		return err
	}

	return c.Render("pixivision/article", fiber.Map{"Article": data})
}

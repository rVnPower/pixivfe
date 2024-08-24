package routes

import (
	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/session"
	"github.com/gofiber/fiber/v2"
)

func IndexPage(c *fiber.Ctx) error {

	// If token is set, do the landing request...
	if token := session.GetPixivToken(c); token != "" {
		mode := c.Query("mode", "all")

		works, err := core.GetLanding(c, mode)

		if err != nil {
			return err
		}

		return Render(c, Data_index{
			Title:      "Landing",
			Data:       *works,
			LoggedIn: true,
		})
	}

	// ...otherwise, default to today's illustration ranking
	works, err := core.GetRanking(c, "daily", "illust", "", "1")
	if err != nil {
		return err
	}
	return Render(c, Data_index{
		Title:       "Landing",
		NoTokenData: works,
		LoggedIn:  false,
	})
}

func Oembed(c *fiber.Ctx) error {
	pageURL := c.BaseURL()
	artistName := c.Query("a", "")
	artistURL := c.Query("u", "")

	data := fiber.Map{
		"version":       "1.0",
		"embed_type":    "rich",
		"provider_name": "PixivFE",
		"provider_url":  pageURL,
		"author_name":   artistName,
		"author_url":    artistURL,
	}

	return c.JSON(data)
}

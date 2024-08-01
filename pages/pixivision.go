package pages

import (
	"github.com/gofiber/fiber/v2"
	"codeberg.org/vnpower/pixivision"
)

func PixivisionHomePage(c *fiber.Ctx) error {
	data, err := pixivision.PixivisionGetHomepage()
	if err != nil {
		return err
	}
	return c.Render("pixivision", fiber.Map{"Data": data})
}

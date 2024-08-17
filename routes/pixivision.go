package pages

import (
	"codeberg.org/vnpower/pixivision"
	"github.com/gofiber/fiber/v2"
)

func PixivisionHomePage(c *fiber.Ctx) error {
	data, err := pixivision.PixivisionGetHomepage()
	if err != nil {
		return err
	}
	return c.Render("pixivision", fiber.Map{"Data": data})
}

package pages

import (
	core "codeberg.org/vnpower/pixivfe/v2/pixiv"
	"github.com/gofiber/fiber/v2"
)

func NewestPage(c *fiber.Ctx) error {
	worktype := c.Query("type", "illust")

	r18 := c.Query("r18", "false")

	works, err := core.GetNewestArtworks(c, worktype, r18)
	if err != nil {
		return err
	}

	return c.Render("newest", fiber.Map{
		"Items": works,
		"Title": "Newest works",
	})
}

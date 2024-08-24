package routes

import (
	"codeberg.org/vnpower/pixivfe/v2/core"
	"github.com/gofiber/fiber/v2"
)

func NewestPage(c *fiber.Ctx) error {
	worktype := c.Query("type", "illust")

	r18 := c.Query("r18", "false")

	works, err := core.GetNewestArtworks(c, worktype, r18)
	if err != nil {
		return err
	}

	return Render(c, Data_newest{Items: works, Title: "Newest works"})
}

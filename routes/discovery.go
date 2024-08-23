package routes

import (
	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/utils"
	"github.com/gofiber/fiber/v2"
)

func DiscoveryPage(c *fiber.Ctx) error {
	mode := c.Query("mode", "safe")

	works, err := core.GetDiscoveryArtwork(c, mode)
	if err != nil {
		return err
	}

	urlc := utils.PartialURL{Path: "discovery", Query: map[string]string{"mode": mode}}

	return c.Render("discovery", fiber.Map{
		"Artworks": works,
		"Title":    "Discovery",
		"Queries":  urlc,
	})
}

func NovelDiscoveryPage(c *fiber.Ctx) error {
	mode := c.Query("mode", "safe")

	works, err := core.GetDiscoveryNovels(c, mode)
	if err != nil {
		return err
	}

	return c.Render("novelDiscovery", fiber.Map{
		"Novels": works,
		"Title":  "Discovery",
	})
}

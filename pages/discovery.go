package pages

import (
	site "codeberg.org/vnpower/pixivfe/v2/core/http"
	core "codeberg.org/vnpower/pixivfe/v2/core/webapi"
	"github.com/gofiber/fiber/v2"
)

func DiscoveryPage(c *fiber.Ctx) error {
	mode := c.Query("mode", "safe")

	works, err := core.GetDiscoveryArtwork(c, mode)
	if err != nil {
		return err
	}

	urlc := site.NewURLConstruct("discovery", map[string]string{"mode": mode}, "#checkpoint")

	return c.Render("discovery", fiber.Map{
		"Artworks": works,
		"Title":    "Discovery",
		"URLC":     urlc.Replace,
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

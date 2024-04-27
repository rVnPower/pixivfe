package pages

import (
	"net/url"
	"strconv"

	core "codeberg.org/vnpower/pixivfe/v2/core/webapi"
	"github.com/gofiber/fiber/v2"
)

func TagPage(c *fiber.Ctx) error {
	queries := make(map[string]string, 3)
	queries["Mode"] = c.Query("mode", "safe")
	queries["Category"] = c.Query("category", "artworks")
	queries["Order"] = c.Query("order", "date_d")
	queries["Ratio"] = c.Query("ratio", "")

	param := c.Params("name")
	name, err := url.PathUnescape(param)
	if err != nil {
		return err
	}

	page := c.Query("page", "1")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return err
	}

	tag, err := core.GetTagData(c, name)
	if err != nil {
		return err
	}
	result, err := core.GetSearch(c, queries["Category"], name, queries["Order"], queries["Mode"], queries["Ratio"], page)
	if err != nil {
		return err
	}

	return c.Render("pages/tag", fiber.Map{"Title": "Results for " + name, "Tag": tag, "Data": result, "Queries": queries, "TrueTag": param, "Page": pageInt})
}

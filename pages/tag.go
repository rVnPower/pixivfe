package pages

import (
	"net/url"
	"strconv"

	site "codeberg.org/vnpower/pixivfe/v2/core/http"
	core "codeberg.org/vnpower/pixivfe/v2/core/webapi"
	"github.com/gofiber/fiber/v2"
)

func TagPage(c *fiber.Ctx) error {
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

	// Because of the large amount of queries available for this route,
	// I made a struct type just to manage the queries
	queries := core.SearchPageSettings{
		Name:     name,
		Category: c.Query("category", "artworks"),
		Order:    c.Query("order", "date_d"),
		Mode:     c.Query("mode", "safe"),
		Ratio:    c.Query("ratio", ""),
		Page:     page,
	}

	tag, err := core.GetTagData(c, name)
	if err != nil {
		return err
	}
	result, err := core.GetSearch(c, queries)
	if err != nil {
		return err
	}

	urlc := site.NewURLConstruct("tags", queries.ReturnMap())

	return c.Render("tag", fiber.Map{"Title": "Results for " + name, "Tag": tag, "Data": result, "Queries": urlc, "TrueTag": param, "Page": pageInt})
}

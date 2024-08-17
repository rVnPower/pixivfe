package pages

import (
	"net/http"
	"net/url"
	"strconv"

	core "codeberg.org/vnpower/pixivfe/v2/pixiv"
	"codeberg.org/vnpower/pixivfe/v2/utils"
	"github.com/gofiber/fiber/v2"
)

func TagPage(c *fiber.Ctx) error {
	param := c.Params("name", c.Query("name"))
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
		Wlt:      c.Query("wlt", ""),
		Wgt:      c.Query("wgt", ""),
		Hlt:      c.Query("hlt", ""),
		Hgt:      c.Query("hgt", ""),
		Tool:     c.Query("tool", ""),
		Scd:      c.Query("scd", ""),
		Ecd:      c.Query("ecd", ""),
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

	urlc := utils.NewURLConstruct("tags", queries.ReturnMap())

	return c.Render("tag", fiber.Map{"Title": "Results for " + name, "Tag": tag, "Data": result, "QueriesC": urlc, "TrueTag": param, "Page": pageInt})
}

func AdvancedTagPost(c *fiber.Ctx) error {
	return c.RedirectToRoute("/tags", fiber.Map{
		"queries": map[string]string{
			"name":     c.Query("name", c.FormValue("name")),
			"category": c.Query("category", "artworks"),
			"order":    c.Query("order", "date_d"),
			"mode":     c.Query("mode", "safe"),
			"ratio":    c.Query("ratio"),
			"page":     c.Query("page", "1"),
			"wlt":      c.Query("wlt", c.FormValue("wlt")),
			"wgt":      c.Query("wgt", c.FormValue("wgt")),
			"hlt":      c.Query("hlt", c.FormValue("hlt")),
			"hgt":      c.Query("hgt", c.FormValue("hgt")),
			"tool":     c.Query("tool", c.FormValue("tool")),
			"scd":      c.Query("scd", c.FormValue("scd")),
			"ecd":      c.Query("ecd", c.FormValue("ecd")),
		},
	}, http.StatusFound)

}

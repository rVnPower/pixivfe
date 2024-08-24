package routes

import (
	"strconv"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"github.com/gofiber/fiber/v2"
)

func RankingPage(c *fiber.Ctx) error {
	mode := c.Query("mode", "daily")
	content := c.Query("content", "all")
	date := c.Query("date", "")

	page := c.Query("page", "1")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		panic(err)
	}

	works, err := core.GetRanking(c, mode, content, date, page)
	if err != nil {
		return err
	}

	return Render(c, Data_rank{Title: "Ranking", Page: pageInt, PageLimit: 10, Date: date, Data: works})
}

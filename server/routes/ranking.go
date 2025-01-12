package routes

import (
	"net/http"
	"strconv"

	"codeberg.org/vnpower/pixivfe/v2/core"
)

func RankingPage(w http.ResponseWriter, r *http.Request) error {
	mode := GetQueryParam(r, "mode", "daily")
	content := GetQueryParam(r, "content", "all")
	date := GetQueryParam(r, "date", "")

	page := GetQueryParam(r, "page", "1")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return err
	}

	works, err := core.GetRanking(r, mode, content, date, page)
	if err != nil {
		return err
	}

	return RenderHTML(w, r, Data_rank{Title: "Ranking", Page: pageInt, PageLimit: 10, Date: date, Data: works})
}

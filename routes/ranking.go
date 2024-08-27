package routes

import (
	"strconv"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"net/http"
)

func RankingPage(w http.ResponseWriter, r CompatRequest) error {
	mode := r.Query("mode", "daily")
	content := r.Query("content", "all")
	date := r.Query("date", "")

	page := r.Query("page", "1")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return err
	}

	works, err := core.GetRanking(r.Request, mode, content, date, page)
	if err != nil {
		return err
	}

	return Render(w, r, Data_rank{Title: "Ranking", Page: pageInt, PageLimit: 10, Date: date, Data: works})
}

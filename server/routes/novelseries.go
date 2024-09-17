package routes

import (
	"fmt"
	"math"
	"strconv"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"net/http"
)

func NovelSeriesPage(w http.ResponseWriter, r *http.Request) error {
	id := GetPathVar(r, "id")
	if _, err := strconv.Atoi(id); err != nil {
		return fmt.Errorf("Invalid ID: %s", id)
	}

	series, err := core.GetNovelSeriesByID(r, id)
	if err != nil {
		return err
	}

	// Hard coded limit
	perPage := 30
	pageLimit := int(math.Ceil(float64(series.Total) / float64(perPage)))

	page := GetQueryParam(r, "p", "1")
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 || pageNum > pageLimit {
		return fmt.Errorf("Invalid Page")
	}

	// Use token only if R-18/R-18G
	seriesContents, err := core.GetNovelSeriesContentByID(r, id, pageNum, perPage, series.XRestrict > 0)
	if err != nil {
		return err
	}

	user, err := core.GetUserBasicInformation(r, series.UserID)
	if err != nil {
		return err
	}

	title := fmt.Sprintf("%s / %s Series", series.Title, series.UserName)

	return Render(w, r, Data_novelSeries{NovelSeries: series, NovelSeriesContents: seriesContents, Title: title, User: user, Page: pageNum, PageLimit: pageLimit})
}

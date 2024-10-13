package routes

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/i18n"
)

func NovelSeriesPage(w http.ResponseWriter, r *http.Request) error {
	id := GetPathVar(r, "id")
	if _, err := strconv.Atoi(id); err != nil {
		return i18n.Errorf("Invalid ID: %s", id)
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
		return i18n.Errorf("Invalid Page Number: %d", pageNum)
	}

	// TODO should use token only if R-18/R-18G
	seriesContents, err := core.GetNovelSeriesContentByID(r, id, pageNum, perPage)
	if err != nil {
		return err
	}

	user, err := core.GetUserBasicInformation(r, series.UserID)
	if err != nil {
		return err
	}

	title := fmt.Sprintf("%s | %s", series.Title, series.UserName)

	return RenderHTML(w, r, Data_novelSeries{NovelSeries: series, NovelSeriesContents: seriesContents, Title: title, User: user, Page: pageNum, PageLimit: pageLimit})
}

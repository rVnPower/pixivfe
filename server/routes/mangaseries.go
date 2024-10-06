package routes

import (
	"fmt"
	"math"
	"strconv"

	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/i18n"
)

func MangaSeriesPage(w http.ResponseWriter, r *http.Request) error {
	id := GetPathVar(r, "id")
	if _, err := strconv.Atoi(id); err != nil {
		return i18n.Errorf("Invalid ID: %s", id)
	}

	seriesId := GetPathVar(r, "sid")
	if _, err := strconv.Atoi(seriesId); err != nil {
		return i18n.Errorf("Invalid Series ID: %s", seriesId)
	}

	// No way to know total before the GetMangaSeriesContentByID request.
	page := GetQueryParam(r, "p", "1")
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		return i18n.Errorf("Invalid Page")
	}

	user, err := core.GetUserBasicInformation(r, id)
	if err != nil {
		return err
	}

	// Must use token. because R-18 artwork could be included in an All-Age series.
	seriesContent, err := core.GetMangaSeriesContentByID(r, seriesId, pageNum)
	if err != nil {
		return err
	}

	// Pixiv auto redirects to correct uid when requesting the url,
	// But we could only do following logic
	if id != seriesContent.Brief.UserID {
		redirectUrl := fmt.Sprintf("/user/%s/series/%s", seriesContent.Brief.UserID, seriesId)
		if pageNum != 1 {
			redirectUrl += fmt.Sprintf("?p=%d", pageNum)
		}
		http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)
		return nil
	}

	// Hard coded limit
	perPage := 12
	pageLimit := int(math.Ceil(float64(seriesContent.Total) / float64(perPage)))

	// Pixiv display empty (not error page) if page id exceeds the total/12 +1
	if pageNum > pageLimit {
		return i18n.Errorf("Invalid Page")
	}

	title := fmt.Sprintf("%s / %s Series", seriesContent.Brief.Title, user.Name)

	return RenderHTML(w, r, Data_mangaSeries{MangaSeriesContent: seriesContent, Title: title, User: user, Page: pageNum, PageLimit: pageLimit})
}

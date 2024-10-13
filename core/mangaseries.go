package core

import (
	"net/http"
	"strconv"
	"time"

	"codeberg.org/vnpower/pixivfe/v2/server/session"
	"github.com/goccy/go-json"
)

type MangaSeries struct {
	ID             string    `json:"id"`
	UserID         string    `json:"userId"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Caption        string    `json:"caption"`
	Total          int       `json:"total"`
	ContentOrder   any       `json:"content_order"`
	URL            string    `json:"url"`
	CoverImageSl   int       `json:"coverImageSl"`
	FirstIllustID  string    `json:"firstIllustId"`
	LatestIllustID string    `json:"latestIllustId"`
	CreateDate     time.Time `json:"createDate"`
	UpdateDate     time.Time `json:"updateDate"`
	WatchCount     any       `json:"watchCount"`
	IsWatched      bool      `json:"isWatched"`
	IsNotifying    bool      `json:"isNotifying"`
}

type MangaSeriesContent struct {
	Series []struct {
		WorkID string `json:"workId"`
		Order  int    `json:"order"`
		Brief  ArtworkBrief
	} `json:"series"`
	IsSetCover           bool   `json:"isSetCover"`
	SeriesID             int    `json:"seriesId"`
	OtherSeriesID        string `json:"otherSeriesId"`
	RecentUpdatedWorkIds []int  `json:"recentUpdatedWorkIds"`
	Total                int    `json:"total"`
	IsWatched            bool   `json:"isWatched"`
	IsNotifying          bool   `json:"isNotifying"`
	Brief                MangaSeries
}

func GetMangaSeriesContentByID(r *http.Request, id string, page int) (MangaSeriesContent, error) {
	var series_content MangaSeriesContent

	URL := GetMangaSeriesContentURL(id, page)

	response, err := API_GET_UnwrapJson(r.Context(), URL, "")
	if err != nil {
		return series_content, err
	}

	response = session.ProxyImageUrl(r, response)

	var mangas struct {
		Thumbnails struct {
			List []ArtworkBrief `json:"illust"`
		} `json:"thumbnails"`
		IllustSeries []MangaSeries      `json:"illustSeries"`
		Page         MangaSeriesContent `json:"page"`
	}

	err = json.Unmarshal([]byte(response), &mangas)
	if err != nil {
		return series_content, err
	}

	series_content = mangas.Page
	for idx, content := range series_content.Series {
		for _, item := range mangas.Thumbnails.List {
			if item.ID == content.WorkID {
				series_content.Series[idx].Brief = item
				break
			}
		}
	}

	for _, brief := range mangas.IllustSeries {
		if brief.ID == strconv.Itoa(series_content.SeriesID) {
			series_content.Brief = brief
			break
		}
	}

	return series_content, nil
}

package routes

import (
	"fmt"
	"strconv"
	"strings"

	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/server/session"
)

func NovelPage(w http.ResponseWriter, r *http.Request) error {
	id := GetPathVar(r, "id")
	if _, err := strconv.Atoi(id); err != nil {
		return fmt.Errorf("Invalid ID: %s", id)
	}

	novel, err := core.GetNovelByID(r, id)
	if err != nil {
		return err
	}

	related, err := core.GetNovelRelated(r, id)
	if err != nil {
		return err
	}

	var contentTitles []core.NovelSeriesContentTitle
	if novel.SeriesNavData.SeriesID != 0 {
		// Must use token, because we can't determine Series' XRestrict via Novel API here
		// and All-age post could also appears in R-18 series.
		contentTitles, _ = core.GetNovelSeriesContentTitlesByID(r, novel.SeriesNavData.SeriesID)
	}

	if novel.CommentOff == 0 {
		// TODO should use token only if R-18/R-18G
		comments, err := core.GetNovelComments(r, id)
		if err == nil {
			novel.CommentsList = comments
		}
	}

	user, err := core.GetUserBasicInformation(r, novel.UserID)
	if err != nil {
		return err
	}

	fontType := session.GetCookie(r, session.Cookie_NovelFontType)
	if fontType == "" {
		fontType = "gothic"
	}
	viewMode := session.GetCookie(r, session.Cookie_NovelViewMode)
	if viewMode == "" {
		viewMode = strconv.Itoa(novel.Settings.ViewMode)
	}

	// println("fontType", fontType)

	title := novel.Title
	if novel.SeriesNavData.SeriesID != 0 {
		title = fmt.Sprintf("#%d %s | %s", novel.SeriesNavData.Order, novel.Title, novel.SeriesNavData.Title)
	}

	return Render(w, r, Data_novel{Novel: novel, NovelRelated: related, User: user, NovelSeriesContentTitles: contentTitles, Title: title, FontType: fontType, ViewMode: viewMode, Language: strings.ToLower(novel.Language)})
}

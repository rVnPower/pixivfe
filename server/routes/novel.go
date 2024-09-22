package routes

import (
	"fmt"
	"strconv"
	"strings"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/server/session"
	"net/http"
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

	return Render(w, r, Data_novel{Novel: novel, NovelRelated: related, User: user, Title: novel.Title, FontType: fontType, ViewMode: viewMode, Language: strings.ToLower(novel.Language)})
}

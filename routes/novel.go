package routes

import (
	"fmt"
	"strconv"
	"strings"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/session"
	"net/http"
)

func NovelPage(w http.ResponseWriter, r CompatRequest) error {
	id := r.Params("id")
	if _, err := strconv.Atoi(id); err != nil {
		return fmt.Errorf("Invalid ID: %s", id)
	}

	novel, err := core.GetNovelByID(r.Request, id)
	if err != nil {
		return err
	}

	related, err := core.GetNovelRelated(r.Request, id)
	if err != nil {
		return err
	}

	user, err := core.GetUserBasicInformation(r.Request, novel.UserID)
	if err != nil {
		return err
	}

	fontType := session.GetCookie(r.Request, session.Cookie_NovelFontType)
	if fontType == "" {
		fontType = "gothic"
	}
	viewMode := session.GetCookie(r.Request, session.Cookie_NovelViewMode)
	if viewMode == "" {
		viewMode = strconv.Itoa(novel.Settings.ViewMode)
	}

	// println("fontType", fontType)

	return Render(w, r, Data_novel{Novel: novel, NovelRelated: related, User: user, Title: novel.Title, FontType: fontType, ViewMode: viewMode, Language: strings.ToLower(novel.Language)})
}

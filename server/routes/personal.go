package routes

import (
	"net/http"
	"strconv"
	"strings"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/server/request_context"
	"codeberg.org/vnpower/pixivfe/v2/server/session"
)

func PromptUserToLoginPage(w http.ResponseWriter, r *http.Request) error {
	request_context.Get(r).RenderStatusCode = http.StatusUnauthorized
	err := Render(w, r, Data_unauthorized{})
	if err != nil {
		return err
	}
	return nil
}

func LoginUserPage(w http.ResponseWriter, r *http.Request) error {
	token := session.GetPixivToken(r)

	if token == "" {
		return PromptUserToLoginPage(w, r)
	}

	// The left part of the token is the member ID
	userId := strings.Split(token, "_")

	http.Redirect(w, r, "/users/"+userId[0], http.StatusSeeOther)
	return nil
}

func LoginBookmarkPage(w http.ResponseWriter, r *http.Request) error {
	token := session.GetPixivToken(r)
	if token == "" {
		return PromptUserToLoginPage(w, r)
	}

	// The left part of the token is the member ID
	userId := strings.Split(token, "_")

	http.Redirect(w, r, "/users/"+userId[0]+"/bookmarks#checkpoint", http.StatusSeeOther)
	return nil
}

func FollowingWorksPage(w http.ResponseWriter, r *http.Request) error {
	if token := session.GetPixivToken(r); token == "" {
		return PromptUserToLoginPage(w, r)
	}

	mode := GetQueryParam(r, "mode", "all")
	page := GetQueryParam(r, "page", "1")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return err
	}

	works, err := core.GetNewestFromFollowing(r, mode, page)
	if err != nil {
		return err
	}

	return Render(w, r, Data_following{Title: "Following works", Mode: mode, Artworks: works, CurPage: page, Page: pageInt})
}

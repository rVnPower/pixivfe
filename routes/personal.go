package routes

import (
	"net/http"
	"strconv"
	"strings"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/session"
)

func PromptUserToLoginPage(w http.ResponseWriter, r CompatRequest) error {
	r.Response.StatusCode = http.StatusUnauthorized
	return Render(w, r, Data_unauthorized{})
}

func LoginUserPage(w http.ResponseWriter, r CompatRequest) error {
	token := session.GetPixivToken(r.Request)

	if token == "" {
		return PromptUserToLoginPage(w, r)
	}

	// The left part of the token is the member ID
	userId := strings.Split(token, "_")

	http.Redirect(w, r.Request, "/users/" + userId[0], http.StatusSeeOther)
	return nil
}

func LoginBookmarkPage(w http.ResponseWriter, r CompatRequest) error {
	token := session.GetPixivToken(r.Request)
	if token == "" {
		return PromptUserToLoginPage(w, r)
	}

	// The left part of the token is the member ID
	userId := strings.Split(token, "_")

	http.Redirect(w, r.Request, "/users/" + userId[0] + "/bookmarks#checkpoint", http.StatusSeeOther)
	return nil
}

func FollowingWorksPage(w http.ResponseWriter, r CompatRequest) error {
	if token := session.GetPixivToken(r.Request); token == "" {
		return PromptUserToLoginPage(w, r)
	}

	mode := r.Query("mode", "all")
	page := r.Query("page", "1")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return err
	}

	works, err := core.GetNewestFromFollowing(r.Request, mode, page)
	if err != nil {
		return err
	}

	return Render(w, r, Data_following{Title: "Following works", Mode: mode, Artworks: works, CurPage: page, Page: pageInt})
}

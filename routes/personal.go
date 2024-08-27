package routes

import (
	"net/http"
	"strconv"
	"strings"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/session"
)

func PromptUserToLoginPage(w http.ResponseWriter, r CompatRequest) error {
	r.Status(http.StatusUnauthorized)
	return Render(w, r, Data_unauthorized{})
}

func LoginUserPage(w http.ResponseWriter, r CompatRequest) error {
	token := session.GetPixivToken(r)

	if token == "" {
		return PromptUserToLoginPage(r)
	}

	// The left part of the token is the member ID
	userId := strings.Split(token, "_")

	r.Redirect("/users/" + userId[0])
	return nil
}

func LoginBookmarkPage(w http.ResponseWriter, r CompatRequest) error {
	token := session.GetPixivToken(r)
	if token == "" {
		return PromptUserToLoginPage(r)
	}

	// The left part of the token is the member ID
	userId := strings.Split(token, "_")

	r.Redirect("/users/" + userId[0] + "/bookmarks#checkpoint")
	return nil
}

func FollowingWorksPage(w http.ResponseWriter, r CompatRequest) error {
	if token := session.GetPixivToken(r); token == "" {
		return PromptUserToLoginPage(r)
	}

	mode := r.Query("mode", "all")
	page := r.Query("page", "1")

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

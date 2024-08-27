package routes

import (
	"net/http"
	"strconv"
	"strings"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/session"
	"net/http"
)

func PromptUserToLoginPage(c *http.Request) error {
	c.Status(http.StatusUnauthorized)
	return Render(c, Data_unauthorized{})
}

func LoginUserPage(c *http.Request) error {
	token := session.GetPixivToken(c)

	if token == "" {
		return PromptUserToLoginPage(c)
	}

	// The left part of the token is the member ID
	userId := strings.Split(token, "_")

	c.Redirect("/users/" + userId[0])
	return nil
}

func LoginBookmarkPage(c *http.Request) error {
	token := session.GetPixivToken(c)
	if token == "" {
		return PromptUserToLoginPage(c)
	}

	// The left part of the token is the member ID
	userId := strings.Split(token, "_")

	c.Redirect("/users/" + userId[0] + "/bookmarks#checkpoint")
	return nil
}

func FollowingWorksPage(c *http.Request) error {
	if token := session.GetPixivToken(c); token == "" {
		return PromptUserToLoginPage(c)
	}

	mode := c.Query("mode", "all")
	page := c.Query("page", "1")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return err
	}

	works, err := core.GetNewestFromFollowing(c, mode, page)
	if err != nil {
		return err
	}

	return Render(c, Data_following{Title: "Following works", Mode: mode, Artworks: works, CurPage: page, Page: pageInt})
}

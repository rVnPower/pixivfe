package pages

import (
	"net/http"
	"strconv"
	"strings"

	session "codeberg.org/vnpower/pixivfe/v2/core/session"
	core "codeberg.org/vnpower/pixivfe/v2/core/webapi"
	"github.com/gofiber/fiber/v2"
)

func PromptUserToLoginPage(c *fiber.Ctx) error {
	c.Status(http.StatusUnauthorized)
	return c.Render("unauthorized", fiber.Map{})
}

func LoginUserPage(c *fiber.Ctx) error {
	token := session.GetPixivToken(c)

	if token == "" {
		return PromptUserToLoginPage(c)
	}

	// The left part of the token is the member ID
	userId := strings.Split(token, "_")

	c.Redirect("/users/" + userId[0])
	return nil
}

func LoginBookmarkPage(c *fiber.Ctx) error {
	token := session.GetPixivToken(c)
	if token == "" {
		return PromptUserToLoginPage(c)
	}

	// The left part of the token is the member ID
	userId := strings.Split(token, "_")

	c.Redirect("/users/" + userId[0] + "/bookmarks#checkpoint")
	return nil
}

func FollowingWorksPage(c *fiber.Ctx) error {
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

	return c.Render("following", fiber.Map{
		"Title":    "Following works",
		"Mode":     mode,
		"Artworks": works,
		"CurPage":  page,
		"Page":     pageInt,
	})
}

package pages

import (
	"fmt"
	"strconv"
	"strings"

	session "codeberg.org/vnpower/pixivfe/v2/core/session"
	core "codeberg.org/vnpower/pixivfe/v2/core/webapi"
	"github.com/gofiber/fiber/v2"
)

func NovelPage(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := strconv.Atoi(id); err != nil {
		return fmt.Errorf("Invalid ID: %s", id)
	}

	novel, err := core.GetNovelByID(c, id)
	if err != nil {
		return err
	}

	related, err := core.GetNovelRelated(c, id)
	if err != nil {
		return err
	}

	user, err := core.GetUserBasicInformation(c, novel.UserID)
	if err != nil {
		return err
	}

	fontType := session.GetCookie(c, session.Cookie_NovelFontType, "gothic")
	viewMode := session.GetCookie(c, session.Cookie_NovelViewMode, strconv.Itoa(novel.Settings.ViewMode))

	// println("fontType", fontType)

	return c.Render("pages/novel", fiber.Map{
		"Novel":        novel,
		"NovelRelated": related,
		"User":         user,
		"Title":        novel.Title,
		"FontType":     fontType,
		"ViewMode":     viewMode,
		"Language":     strings.ToLower(novel.Language),
	})
}

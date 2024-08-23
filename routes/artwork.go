package routes

import (
	"fmt"
	"strconv"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"github.com/gofiber/fiber/v2"
)

func ArtworkPage(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := strconv.Atoi(id); err != nil {
		return fmt.Errorf("Invalid ID: %s", id)
	}

	illust, err := core.GetArtworkByID(c, id, true)
	if err != nil {
		return err
	}

	metaDescription := ""
	for _, i := range illust.Tags {
		metaDescription += "#" + i.Name + ", "
	}

	// monkey patching. assuming illust.Images[_].Large is used
	for _, img := range illust.Images {
		PreloadImage(c, img.Large)
	}

	return Render(c, Data_artwork{
		Illust:          *illust,
		Title:           illust.Title,
		MetaDescription: metaDescription,
		MetaImage:       illust.Images[0].Original,
		MetaAuthor:      illust.UserName,
		MetaAuthorID:    illust.UserID,
	})
}

func PreloadImage(c *fiber.Ctx, url string) {
	c.Response().Header.Add("Link", fmt.Sprintf("<%s>; rel=preload; as=image", url))
}

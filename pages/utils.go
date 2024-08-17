package pages

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func PreloadImage(c *fiber.Ctx, url string) {
	c.Response().Header.Add("Link", fmt.Sprintf("<%s>; rel=preload; as=image", url))
}

package routes

import (
	"codeberg.org/vnpower/pixivfe/v2/config"
	"github.com/gofiber/fiber/v2"
)

func AboutPage(c *fiber.Ctx) error {
	info := fiber.Map{
		"Time":           config.GlobalServerConfig.StartingTime,
		"Version":        config.GlobalServerConfig.Version,
		"ImageProxy":     config.GlobalServerConfig.ProxyServer.String(),
		"AcceptLanguage": config.GlobalServerConfig.AcceptLanguage,
	}
	return c.Render("about", info)
}

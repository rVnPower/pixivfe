package routes

import (
	"codeberg.org/vnpower/pixivfe/v2/config"
	"github.com/gofiber/fiber/v2"
)

func AboutPage(c *fiber.Ctx) error {
	return Render(c, Data_about{
		Time:           config.GlobalServerConfig.StartingTime,
		Version:        config.GlobalServerConfig.Version,
		ImageProxy:     config.GlobalServerConfig.ProxyServer.String(),
		AcceptLanguage: config.GlobalServerConfig.AcceptLanguage,
	})
}

package routes

import (
	"codeberg.org/vnpower/pixivfe/v2/config"
	"net/http"
)

func AboutPage(c *http.Request) error {
	return Render(c, Data_about{
		Time:           config.GlobalServerConfig.StartingTime,
		Version:        config.GlobalServerConfig.Version,
		ImageProxy:     config.GlobalServerConfig.ProxyServer.String(),
		AcceptLanguage: config.GlobalServerConfig.AcceptLanguage,
	})
}

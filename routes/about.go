package routes

import (
	"codeberg.org/vnpower/pixivfe/v2/config"
	"net/http"
)

func AboutPage(w http.ResponseWriter, r CompatRequest) error {
	return Render(w, r, Data_about{
		Time:           config.GlobalServerConfig.StartingTime,
		Version:        config.GlobalServerConfig.Version,
		ImageProxy:     config.GlobalServerConfig.ProxyServer.String(),
		AcceptLanguage: config.GlobalServerConfig.AcceptLanguage,
	})
}

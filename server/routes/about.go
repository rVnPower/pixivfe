package routes

import (
	"codeberg.org/vnpower/pixivfe/v2/config"
	"net/http"
)

func AboutPage(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, Data_about{
		Time:           config.GlobalConfig.StartingTime,
		Version:        config.GlobalConfig.Version,
		ImageProxy:     config.GlobalConfig.ProxyServer.String(),
		AcceptLanguage: config.GlobalConfig.AcceptLanguage,
	})
}
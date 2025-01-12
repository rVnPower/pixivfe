package routes

import (
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/config"
)

func AboutPage(w http.ResponseWriter, r *http.Request) error {
	return RenderHTML(w, r, Data_about{
		Time:           config.GlobalConfig.StartingTime,
		Version:        config.GlobalConfig.Version,
		DomainName:     r.Host, // Used to template in the actual domain name serving PixivFE
		RepoURL:        config.GlobalConfig.RepoURL,
		Revision:       config.GlobalConfig.Revision,
		RevisionHash:   config.GlobalConfig.RevisionHash, // Used for the link to the source code repo
		ImageProxy:     config.GlobalConfig.ProxyServer.String(),
		AcceptLanguage: config.GlobalConfig.AcceptLanguage,
	})
}

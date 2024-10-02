package routes

import (
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/server/template"
)

func DiscoveryPage(w http.ResponseWriter, r *http.Request) error {
	mode := GetQueryParam(r, "mode", "safe")

	works, err := core.GetDiscoveryArtwork(r, mode)
	if err != nil {
		return err
	}

	urlc := template.PartialURL{Path: "discovery", Query: map[string]string{"mode": mode}}

	return RenderHTML(w, r, Data_discovery{Artworks: works, Title: "Discovery", Queries: urlc})
}

func NovelDiscoveryPage(w http.ResponseWriter, r *http.Request) error {
	mode := GetQueryParam(r, "mode", "safe")

	works, err := core.GetDiscoveryNovels(r, mode)
	if err != nil {
		return err
	}

	return RenderHTML(w, r, Data_novelDiscovery{Novels: works, Title: "Discovery"})
}

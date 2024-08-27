package routes

import (
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/template"
)

func DiscoveryPage(w http.ResponseWriter, r *http.Request) error {
	mode := GetQueryParam(r, "mode", "safe")

	works, err := core.GetDiscoveryArtwork(r, mode)
	if err != nil {
		return err
	}

	urlc := template.PartialURL{Path: "discovery", Query: map[string]string{"mode": mode}}

	return Render(w, r, Data_discovery{Artworks: works, Title: "Discovery", Queries: urlc})
}

func NovelDiscoveryPage(w http.ResponseWriter, r *http.Request) error {
	mode := GetQueryParam(r, "mode", "safe")

	works, err := core.GetDiscoveryNovels(r, mode)
	if err != nil {
		return err
	}

	return Render(w, r, Data_novelDiscovery{Novels: works, Title: "Discovery"})
}

package routes

import (
	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/utils"
	"net/http"
)

func DiscoveryPage(w http.ResponseWriter, r CompatRequest) error {
	mode := r.Query("mode", "safe")

	works, err := core.GetDiscoveryArtwork(r, mode)
	if err != nil {
		return err
	}

	urlc := utils.PartialURL{Path: "discovery", Query: map[string]string{"mode": mode}}

	return Render(w, r, Data_discovery{Artworks: works, Title: "Discovery", Queries: urlc})
}

func NovelDiscoveryPage(w http.ResponseWriter, r CompatRequest) error {
	mode := r.Query("mode", "safe")

	works, err := core.GetDiscoveryNovels(r, mode)
	if err != nil {
		return err
	}

	return Render(w, r, Data_novelDiscovery{Novels: works, Title: "Discovery"})
}

package core

import (
	"codeberg.org/vnpower/pixivfe/v2/server/session"
	"github.com/goccy/go-json"
	"net/http"
)

func GetNewestArtworks(r *http.Request, worktype string, r18 string) ([]ArtworkBrief, error) {
	token := session.GetPixivToken(r)
	URL := GetNewestArtworksURL(worktype, r18, "0")

	var body struct {
		Artworks []ArtworkBrief `json:"illusts"`
		// LastId string
	}

	resp, err := API_GET_UnwrapJson(r.Context(), URL, token)
	if err != nil {
		return nil, err
	}
	resp = session.ProxyImageUrl(r, resp)

	err = json.Unmarshal([]byte(resp), &body)
	if err != nil {
		return nil, err
	}

	return body.Artworks, nil
}

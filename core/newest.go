package core

import (
	"codeberg.org/vnpower/pixivfe/v2/session"
	"github.com/goccy/go-json"
	"net/http"
)

func GetNewestArtworks(c *http.Request, worktype string, r18 string) ([]ArtworkBrief, error) {
	token := session.GetPixivToken(c)
	URL := GetNewestArtworksURL(worktype, r18, "0")

	var body struct {
		Artworks []ArtworkBrief `json:"illusts"`
		// LastId string
	}

	resp, err := UnwrapWebAPIRequest(c.Context(), URL, token)
	if err != nil {
		return nil, err
	}
	resp = session.ProxyImageUrl(c, resp)

	err = json.Unmarshal([]byte(resp), &body)
	if err != nil {
		return nil, err
	}

	return body.Artworks, nil
}

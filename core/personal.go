package core

import (
	"codeberg.org/vnpower/pixivfe/v2/session"

	"github.com/goccy/go-json"
	"net/http"
)

func GetNewestFromFollowing(c *http.Request, mode, page string) ([]ArtworkBrief, error) {
	token := session.GetPixivToken(c)
	URL := GetNewestFromFollowingURL(mode, page)

	var body struct {
		Thumbnails json.RawMessage `json:"thumbnails"`
	}

	var artworks struct {
		Artworks []ArtworkBrief `json:"illust"`
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
	err = json.Unmarshal([]byte(body.Thumbnails), &artworks)
	if err != nil {
		return nil, err
	}

	return artworks.Artworks, nil
}

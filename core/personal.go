package core

import (
	"codeberg.org/vnpower/pixivfe/v2/server/session"

	"github.com/goccy/go-json"
	"net/http"
)

func GetNewestFromFollowing(r *http.Request, mode, page string) ([]ArtworkBrief, error) {
	token := session.GetUserToken(r)
	URL := GetNewestFromFollowingURL(mode, page)

	var body struct {
		Thumbnails json.RawMessage `json:"thumbnails"`
	}

	var artworks struct {
		Artworks []ArtworkBrief `json:"illust"`
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
	err = json.Unmarshal([]byte(body.Thumbnails), &artworks)
	if err != nil {
		return nil, err
	}

	return artworks.Artworks, nil
}

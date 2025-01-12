package core

import (
	"net/http"

	"github.com/goccy/go-json"
	"github.com/tidwall/gjson"

	"codeberg.org/vnpower/pixivfe/v2/i18n"
	"codeberg.org/vnpower/pixivfe/v2/server/session"
)

func GetDiscoveryArtwork(r *http.Request, mode string) ([]ArtworkBrief, error) {
	token := session.GetUserToken(r)

	URL := GetDiscoveryURL(mode, 100)

	var artworks []ArtworkBrief

	resp, err := API_GET_UnwrapJson(r.Context(), URL, token)
	if err != nil {
		return nil, err
	}
	resp = session.ProxyImageUrl(r, resp)
	if !gjson.Valid(resp) {
		return nil, i18n.Errorf("Invalid JSON: %v", resp)
	}
	data := gjson.Get(resp, "thumbnails.illust").String()

	err = json.Unmarshal([]byte(data), &artworks)
	if err != nil {
		return nil, err
	}

	return artworks, nil
}

func GetDiscoveryNovels(r *http.Request, mode string) ([]NovelBrief, error) {
	token := session.GetUserToken(r)

	URL := GetDiscoveryNovelURL(mode, 100)

	var novels []NovelBrief

	resp, err := API_GET_UnwrapJson(r.Context(), URL, token)
	if err != nil {
		return nil, err
	}
	resp = session.ProxyImageUrl(r, resp)
	if !gjson.Valid(resp) {
		return nil, i18n.Errorf("Invalid JSON: %v", resp)
	}
	data := gjson.Get(resp, "thumbnails.novel").String()

	err = json.Unmarshal([]byte(data), &novels)
	if err != nil {
		return nil, err
	}

	return novels, nil
}

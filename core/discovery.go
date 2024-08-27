package core

import (
	"fmt"

	"codeberg.org/vnpower/pixivfe/v2/session"
	"github.com/goccy/go-json"
	"net/http"
	"github.com/tidwall/gjson"
)

func GetDiscoveryArtwork(c *http.Request, mode string) ([]ArtworkBrief, error) {
	token := session.GetPixivToken(c)

	URL := GetDiscoveryURL(mode, 100)

	var artworks []ArtworkBrief

	resp, err := UnwrapWebAPIRequest(c.Context(), URL, token)
	if err != nil {
		return nil, err
	}
	resp = session.ProxyImageUrl(c, resp)
	if !gjson.Valid(resp) {
		return nil, fmt.Errorf("Invalid JSON: %v", resp)
	}
	data := gjson.Get(resp, "thumbnails.illust").String()

	err = json.Unmarshal([]byte(data), &artworks)
	if err != nil {
		return nil, err
	}

	return artworks, nil
}

func GetDiscoveryNovels(c *http.Request, mode string) ([]NovelBrief, error) {
	token := session.GetPixivToken(c)

	URL := GetDiscoveryNovelURL(mode, 100)

	var novels []NovelBrief

	resp, err := UnwrapWebAPIRequest(c.Context(), URL, token)
	if err != nil {
		return nil, err
	}
	resp = session.ProxyImageUrl(c, resp)
	if !gjson.Valid(resp) {
		return nil, fmt.Errorf("Invalid JSON: %v", resp)
	}
	data := gjson.Get(resp, "thumbnails.novel").String()

	err = json.Unmarshal([]byte(data), &novels)
	if err != nil {
		return nil, err
	}

	return novels, nil
}

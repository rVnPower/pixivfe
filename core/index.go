package core

import (
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/tidwall/gjson"

	"codeberg.org/vnpower/pixivfe/v2/i18n"
	"codeberg.org/vnpower/pixivfe/v2/server/session"
)

type Pixivision struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnailUrl"`
	URL       string `json:"url"`
}

type RecommendedTags struct {
	Name     string `json:"tag"`
	Artworks []ArtworkBrief
}
type LandingArtworks struct {
	Commissions     []ArtworkBrief
	Following       []ArtworkBrief
	Recommended     []ArtworkBrief
	Newest          []ArtworkBrief
	Rankings        []ArtworkBrief
	Users           []ArtworkBrief
	Pixivision      []Pixivision
	RecommendByTags []RecommendedTags
}

func GetLanding(r *http.Request, mode string) (*LandingArtworks, error) {
	var pages struct {
		Pixivision  []Pixivision `json:"pixivision"`
		Follow      []int        `json:"follow"`
		Recommended struct {
			IDs []string `json:"ids"`
		} `json:"recommend"`
		// EditorRecommended []any `json:"editorRecommend"`
		// UserRecommended   []any `json:"recommendUser"`
		// Commission        []any `json:"completeRequestIds"`
		RecommendedByTags []struct {
			Name string   `json:"tag"`
			IDs  []string `json:"ids"`
		} `json:"recommendByTag"`
	}

	URL := GetLandingURL(mode)

	var landing LandingArtworks

	resp, err := API_GET_UnwrapJson(r.Context(), URL, "")

	if err != nil {
		return &landing, err
	}
	resp = session.ProxyImageUrl(r, resp)

	if !gjson.Valid(resp) {
		return nil, i18n.Errorf("Invalid JSON: %v", resp)
	}

	artworks := map[string]ArtworkBrief{}

	// Get thumbnails and save it into a map, since they were kept
	// separately and need to the index quickly.
	//
	// Since there are no duplicates in this object, we are unable
	// to rely to ranges (ex. one artwork in two separate sections)
	stuff := gjson.Get(resp, "thumbnails.illust")
	stuff.ForEach(func(key, value gjson.Result) bool {
		var artwork ArtworkBrief
		err = json.Unmarshal([]byte(value.String()), &artwork)

		if err != nil {
			return false
		}

		if artwork.ID != "" {
			artworks[artwork.ID] = artwork
		}

		return true // keep iterating
	})

	pagesStr := gjson.Get(resp, "page").String()
	err = json.Unmarshal([]byte(pagesStr), &pages)

	if err != nil {
		return &landing, err
	}

	// Parse everything
	landing.Pixivision = pages.Pixivision

	landing.Following = make([]ArtworkBrief, len(pages.Follow))
	for _, i := range pages.Follow {
		landing.Following = append(landing.Following, artworks[fmt.Sprint(i)])
	}

	for _, i := range pages.RecommendedByTags {
		temp := make([]ArtworkBrief, 0)
		for _, j := range i.IDs {
			temp = append(temp, artworks[j])
		}
		landing.RecommendByTags = append(landing.RecommendByTags, RecommendedTags{Name: i.Name, Artworks: temp})
	}

	landing.Recommended = make([]ArtworkBrief, 0)
	for _, i := range pages.Recommended.IDs {
		landing.Recommended = append(landing.Recommended, artworks[i])
	}

	return &landing, nil
}

package pixiv

import (
	"strings"

	session "codeberg.org/vnpower/pixivfe/v2/session"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

type TagDetail struct {
	Name            string `json:"tag"`
	AlternativeName string `json:"word"`
	Metadata        struct {
		Detail string      `json:"abstract"`
		Image  string      `json:"image"`
		Name   string      `json:"tag"`
		ID     json.Number `json:"id"`
	} `json:"pixpedia"`
}

type SearchArtworks struct {
	Artworks []ArtworkBrief `json:"data"`
	Total    int            `json:"total"`
}

type SearchResult struct {
	Artworks SearchArtworks
	Popular  struct {
		Permanent []ArtworkBrief `json:"permanent"`
		Recent    []ArtworkBrief `json:"recent"`
	} `json:"popular"`
	RelatedTags []string `json:"relatedTags"`
}

type SearchPageSettings struct {
	Name     string // Tag to search for
	Category string // Filter by type, could be illusts or mangas
	Order    string // Sort by date
	Mode     string // Safe, R18 or both
	Ratio    string // Landscape, portrait, or squared
	Page     string // Page number

	// To implement
	Smode string // Exact match, partial match, or match with title
	Wlt   string // Minimum image width
	Wgt   string // Maximum image width
	Hlt   string // Minimum image height
	Hgt   string // Maximum image height
	Tool  string // Filter by production tools (ex. Photoshop)
	Scd   string // After this date
	Ecd   string // Before this date
}

func (s SearchPageSettings) ReturnMap() map[string]string {
	return map[string]string{
		"Name":     s.Name,
		"Category": s.Category,
		"Order":    s.Order,
		"Mode":     s.Mode,
		"Ratio":    s.Ratio,
		"Smode":    s.Smode,
		"Wlt":      s.Wlt,
		"Wgt":      s.Wgt,
		"Hlt":      s.Hlt,
		"Hgt":      s.Hgt,
		"Scd":      s.Scd,
		"Ecd":      s.Ecd,
		"Tool":     s.Tool,
		"Page":     s.Page,
	}
}

func GetTagData(c *fiber.Ctx, name string) (TagDetail, error) {
	var tag TagDetail

	URL := GetTagDetailURL(name)

	response, err := UnwrapWebAPIRequest(c.Context(), URL, "")
	if err != nil {
		return tag, err
	}

	response = session.ProxyImageUrl(c, response)

	err = json.Unmarshal([]byte(response), &tag)
	if err != nil {
		return tag, err
	}

	return tag, nil
}

func GetSearch(c *fiber.Ctx, settings SearchPageSettings) (*SearchResult, error) {
	URL := GetSearchArtworksURL(settings.ReturnMap())

	response, err := UnwrapWebAPIRequest(c.Context(), URL, "")
	if err != nil {
		return nil, err
	}
	response = session.ProxyImageUrl(c, response)

	// IDK how to do better than this lol
	temp := strings.ReplaceAll(string(response), `"illust"`, `"works"`)
	temp = strings.ReplaceAll(temp, `"manga"`, `"works"`)
	temp = strings.ReplaceAll(temp, `"illustManga"`, `"works"`)

	var resultRaw struct {
		*SearchResult
		ArtworksRaw json.RawMessage `json:"works"`
	}
	var artworks SearchArtworks
	var result *SearchResult

	err = json.Unmarshal([]byte(temp), &resultRaw)
	if err != nil {
		return nil, err
	}

	result = resultRaw.SearchResult

	err = json.Unmarshal([]byte(resultRaw.ArtworksRaw), &artworks)
	if err != nil {
		return nil, err
	}

	result.Artworks = artworks

	return result, nil
}

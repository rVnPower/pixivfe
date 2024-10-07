package core

import (
	"fmt"
	"math"
	"sort"

	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/i18n"
	"codeberg.org/vnpower/pixivfe/v2/server/session"
	"github.com/goccy/go-json"
)

// pixivfe internal data type. not used by pixiv.
type UserArtCategory string

const (
	UserArt_Any          UserArtCategory = ""
	UserArt_AnyAlt       UserArtCategory = "artworks"
	UserArt_Illustration UserArtCategory = "illustrations"
	UserArt_Manga        UserArtCategory = "manga"
	UserArt_Bookmarks    UserArtCategory = "bookmarks" // what this user has bookmarked; not art by this user
	UserArt_Novel        UserArtCategory = "novels"
)

func (s UserArtCategory) Validate() error {
	if s != UserArt_Any &&
		s != UserArt_AnyAlt &&
		s != UserArt_Illustration &&
		s != UserArt_Manga &&
		s != UserArt_Bookmarks &&
		s != UserArt_Novel {
		return i18n.Errorf(`Invalid work category: %#v. Only "%s", "%s", "%s", "%s", "%s" and "%s" are available`, s, UserArt_Any, UserArt_AnyAlt, UserArt_Illustration, UserArt_Manga, UserArt_Bookmarks, UserArt_Novel)
	} else {
		return nil
	}
}

type FrequentTag struct {
	Name           string `json:"tag"`
	TranslatedName string `json:"tag_translation"`
}

type User struct {
	ID              string          `json:"userId"`
	Name            string          `json:"name"`
	Avatar          string          `json:"imageBig"`
	Following       int             `json:"following"`
	MyPixiv         int             `json:"mypixivCount"`
	Comment         HTML            `json:"commentHtml"`
	Webpage         string          `json:"webpage"`
	SocialRaw       json.RawMessage `json:"social"`
	Artworks        []ArtworkBrief  `json:"artworks"`
	Novels          []NovelBrief    `json:"novels"`
	Background      map[string]any  `json:"background"`
	ArtworksCount   int
	FrequentTags    []FrequentTag
	Social          map[string]map[string]string
	BackgroundImage string
	NovelSeries     []NovelSeries
	MangaSeries     []MangaSeries
}

func (s *User) ParseSocial() error {
	if string(s.SocialRaw[:]) == "[]" {
		// Fuck Pixiv
		return nil
	}

	err := json.Unmarshal(s.SocialRaw, &s.Social)
	if err != nil {
		return err
	}
	return nil
}

func GetFrequentTags(r *http.Request, ids string, category UserArtCategory) ([]FrequentTag, error) {
	var tags []FrequentTag
	var URL string

	if category != "novels" {
		URL = GetFrequentArtworkTagsURL(ids)
	} else {
		URL = GetFrequentNovelTagsURL(ids)
	}

	response, err := API_GET_UnwrapJson(r.Context(), URL, "")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(response), &tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

type Work interface {
	GetID() string
}

func GetUserWorks(r *http.Request, id, ids string, category UserArtCategory) ([]Work, error) {
	var works []Work

	var URL string
	if category == UserArt_Novel {
		URL = GetUserFullNovelURL(id, ids)
	} else {
		URL = GetUserFullArtworkURL(id, ids)
	}

	resp, err := API_GET_UnwrapJson(r.Context(), URL, "")
	if err != nil {
		return nil, err
	}
	resp = session.ProxyImageUrl(r, resp)

	var body struct {
		Works map[int]json.RawMessage `json:"works"`
	}

	err = json.Unmarshal([]byte(resp), &body)
	if err != nil {
		return nil, err
	}

	for _, v := range body.Works {
		if category == UserArt_Novel {
			var novel NovelBrief
			err = json.Unmarshal(v, &novel)
			if err != nil {
				return nil, err
			}
			works = append(works, novel)
		} else {
			var illust ArtworkBrief
			err = json.Unmarshal(v, &illust)
			if err != nil {
				return nil, err
			}
			works = append(works, illust)
		}
	}

	return works, nil
}

func GetUserArtworksIDAndSeries(r *http.Request, id string, category UserArtCategory, page int) (string, int, json.RawMessage, error) {
	URL := GetUserArtworksURL(id)

	resp, err := API_GET_UnwrapJson(r.Context(), URL, "")
	if err != nil {
		return "", -1, nil, err
	}

	resp = session.ProxyImageUrl(r, resp)

	var body struct {
		Illusts     json.RawMessage `json:"illusts"`
		Mangas      json.RawMessage `json:"manga"`
		MangaSeries json.RawMessage `json:"mangaSeries"`
		Novels      json.RawMessage `json:"novels"`
		NovelSeries json.RawMessage `json:"novelSeries"`
	}

	err = json.Unmarshal([]byte(resp), &body)
	if err != nil {
		return "", -1, nil, err
	}

	var ids []int
	var idsString string

	// TODO: is this necessary
	err = json.Unmarshal([]byte(resp), &body)
	if err != nil {
		return "", -1, nil, err
	}

	var illusts map[int]string
	var mangas map[int]string
	var novels map[int]string
	var series json.RawMessage
	count := 0

	// Get the keys, because Pixiv only returns IDs (very evil)

	if category == UserArt_Illustration || category == UserArt_Any || category == UserArt_AnyAlt {
		if err = json.Unmarshal(body.Illusts, &illusts); err != nil {
			illusts = make(map[int]string)
		}
		for k := range illusts {
			ids = append(ids, k)
			count++
		}
	}
	if category == UserArt_Manga || category == UserArt_Any {
		if err = json.Unmarshal(body.Mangas, &mangas); err != nil {
			mangas = make(map[int]string)
		}
		for k := range mangas {
			ids = append(ids, k)
			count++
		}
		series = body.MangaSeries
	}
	if category == UserArt_Novel {
		if err = json.Unmarshal(body.Novels, &novels); err != nil {
			novels = make(map[int]string)
		}
		for k := range novels {
			ids = append(ids, k)
			count++
		}
		series = body.NovelSeries
	}

	// Reverse sort the ids
	sort.Sort(sort.Reverse(sort.IntSlice(ids)))

	worksPerPage := 30.0
	start, end, err := computeSliceBounds(page, worksPerPage, len(ids))
	if err != nil {
		return "", -1, nil, err
	}

	for _, k := range ids[start:end] {
		idsString += fmt.Sprintf("&ids[]=%d", k)
	}

	return idsString, count, series, nil
}

func GetUserArtwork(r *http.Request, id string, category UserArtCategory, page int, getTags bool) (User, error) {
	var user User

	token := session.GetUserToken(r)

	URL := GetUserInformationURL(id)

	resp, err := API_GET_UnwrapJson(r.Context(), URL, token)
	if err != nil {
		return user, err
	}

	resp = session.ProxyImageUrl(r, resp)

	err = json.Unmarshal([]byte(resp), &user)
	if err != nil {
		return user, err
	}

	if category == UserArt_Bookmarks {
		// Bookmarks
		works, count, err := GetUserBookmarks(r, id, "show", page)
		if err != nil {
			return user, err
		}

		user.Artworks = works

		// Public bookmarks count
		user.ArtworksCount = count
	} else {
		ids, count, series, err := GetUserArtworksIDAndSeries(r, id, category, page)
		if err != nil {
			return user, err
		}

		if count > 0 {
			// Check if the user has works available or not
			works, err := GetUserWorks(r, id, ids, category)
			if err != nil {
				return user, err
			}

			// IDK but the order got shuffled even though Pixiv sorted the IDs in the response
			sort.Slice(works[:], func(i, j int) bool {
				left := works[i].GetID()
				right := works[j].GetID()
				return numberGreaterThan(left, right)
			})

			if category == UserArt_Novel {
				user.Novels = make([]NovelBrief, len(works))
				for i, w := range works {
					user.Novels[i] = w.(NovelBrief)
				}
			} else {
				user.Artworks = make([]ArtworkBrief, len(works))
				for i, w := range works {
					user.Artworks[i] = w.(ArtworkBrief)
				}
			}

			if getTags {
				user.FrequentTags, err = GetFrequentTags(r, ids, category)
				if err != nil {
					return user, err
				}
			}
		}

		if category == UserArt_Novel {
			var novelSeries []NovelSeries
			if series != nil {
				if err = json.Unmarshal(series, &novelSeries); err == nil {
					user.NovelSeries = novelSeries
				}
			}
		} else {
			var mangaSeries []MangaSeries
			if series != nil {
				if err = json.Unmarshal(series, &mangaSeries); err == nil {
					user.MangaSeries = mangaSeries
				}
			}
		}

		// Artworks count
		user.ArtworksCount = count
	}

	err = user.ParseSocial()
	if err != nil {
		return User{}, err
	}

	if user.Background != nil {
		user.BackgroundImage = user.Background["url"].(string)
	}

	return user, nil
}

func GetUserBookmarks(r *http.Request, id, mode string, page int) ([]ArtworkBrief, int, error) {
	page--

	URL := GetUserBookmarksURL(id, mode, page)

	resp, err := API_GET_UnwrapJson(r.Context(), URL, "")
	if err != nil {
		return nil, -1, err
	}
	resp = session.ProxyImageUrl(r, resp)

	var body struct {
		Artworks []json.RawMessage `json:"works"`
		Total    int               `json:"total"`
	}

	err = json.Unmarshal([]byte(resp), &body)
	if err != nil {
		return nil, -1, err
	}

	artworks := make([]ArtworkBrief, len(body.Artworks))

	for index, value := range body.Artworks {
		var artwork ArtworkBrief

		err = json.Unmarshal([]byte(value), &artwork)
		if err != nil {
			artworks[index] = ArtworkBrief{
				ID:        "#",
				Title:     "Deleted or Private",
				Thumbnail: "https://s.pximg.net/common/images/limit_unknown_360.png",
			}
			continue
		}
		artworks[index] = artwork
	}

	return artworks, body.Total, nil
}

func numberGreaterThan(l, r string) bool {
	if len(l) > len(r) {
		return true
	}
	if len(l) < len(r) {
		return false
	}
	return l > r
}

// Utility function to compute slice bounds safely
func computeSliceBounds(page int, worksPerPage float64, totalItems int) (start, end int, err error) {
	if totalItems == 0 {
		return 0, 0, i18n.Error("No items available.")
	}

	maxPages := int(math.Ceil(float64(totalItems) / worksPerPage))
	if page < 1 || page > maxPages {
		return 0, 0, i18n.Error("Invalid page number.")
	}

	start = (page - 1) * int(worksPerPage)
	end = min(start+int(worksPerPage), totalItems)

	return start, end, nil
}

// Methods to satisfy the Work interface
func (a ArtworkBrief) GetID() string {
	return a.ID
}

func (n NovelBrief) GetID() string {
	return n.ID
}

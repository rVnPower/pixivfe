package core

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/goccy/go-json"

	"codeberg.org/vnpower/pixivfe/v2/i18n"
	"codeberg.org/vnpower/pixivfe/v2/server/session"
)

// Pixiv returns 0, 1, 2 to filter SFW and/or NSFW artworks.
// Those values are saved in `XRestrict`
// 0: Safe
// 1: R18
// 2: R18G
type XRestrict int

const (
	Safe XRestrict = 0
	R18  XRestrict = 1
	R18G XRestrict = 2
)

func (x XRestrict) String() string {
	switch x {
	case Safe:
		return i18n.Tr("Safe")
	case R18:
		return i18n.Tr("R18")
	case R18G:
		return i18n.Tr("R18G")
	}
	log.Panicf("invalid value: %#v", int(x))
	return ""
}

// Pixiv returns 0, 1, 2 to filter SFW and/or NSFW artworks.
// Those values are saved in `aiType`
// 0: Not rated / Unknown
// 1: Not AI-generated
// 2: AI-generated

type AiType int

const (
	Unrated AiType = 0
	NotAI   AiType = 1
	AI      AiType = 2
)

func (x AiType) String() string {
	switch x {
	case Unrated:
		return i18n.Tr("Unrated")
	case NotAI:
		return i18n.Tr("Not AI")
	case AI:
		return i18n.Tr("AI")
	}
	log.Panicf("invalid value: %#v", int(x))
	return ""
}

type ImageResponse struct {
	Width  int               `json:"width"`
	Height int               `json:"height"`
	Urls   map[string]string `json:"urls"`
}

type Image struct {
	Width      int
	Height     int
	Small      string
	Medium     string
	Large      string
	Original   string
	IllustType int
}

type Tag struct {
	Name           string `json:"tag"`
	TranslatedName string `json:"translation"`
}

type Comment struct {
	AuthorID   string `json:"userId"`
	AuthorName string `json:"userName"`
	Avatar     string `json:"img"`
	Context    string `json:"comment"`
	Stamp      string `json:"stampId"`
	Date       string `json:"commentDate"`
}

type UserBrief struct {
	ID     string `json:"userId"`
	Name   string `json:"name"`
	Avatar string `json:"imageBig"`
}

type ArtworkBrief struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	ArtistID     string `json:"userId"`
	ArtistName   string `json:"userName"`
	ArtistAvatar string `json:"profileImageUrl"`
	Thumbnail    string `json:"url"`
	Pages        int    `json:"pageCount"`
	XRestrict    int    `json:"xRestrict"`
	AiType       int    `json:"aiType"`
	Bookmarked   any    `json:"bookmarkData"`
	IllustType   int    `json:"illustType"`
}

type Illust struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Description     HTML      `json:"description"`
	UserID          string    `json:"userId"`
	UserName        string    `json:"userName"`
	UserAccount     string    `json:"userAccount"`
	Date            time.Time `json:"uploadDate"`
	Images          []Image
	Tags            []Tag     `json:"tags"`
	Pages           int       `json:"pageCount"`
	Bookmarks       int       `json:"bookmarkCount"`
	Likes           int       `json:"likeCount"`
	Comments        int       `json:"commentCount"`
	Views           int       `json:"viewCount"`
	CommentDisabled int       `json:"commentOff"`
	SanityLevel     int       `json:"sl"`
	XRestrict       XRestrict `json:"xRestrict"`
	AiType          AiType    `json:"aiType"`
	BookmarkData    any       `json:"bookmarkData"`
	Liked           bool      `json:"likeData"`
	SeriesNavData   struct {
		SeriesType  string `json:"seriesType"`
		SeriesID    string `json:"seriesId"`
		Title       string `json:"title"`
		IsWatched   bool   `json:"isWatched"`
		IsNotifying bool   `json:"isNotifying"`
		Order       int    `json:"order"`
		Next        struct {
			Title string `json:"title"`
			Order int    `json:"order"`
			ID    string `json:"id"`
		} `json:"next"`
		Prev struct {
			Title string `json:"title"`
			Order int    `json:"order"`
			ID    string `json:"id"`
		} `json:"prev"`
	} `json:"seriesNavData"`
	User         UserBrief
	RecentWorks  []ArtworkBrief
	RelatedWorks []ArtworkBrief
	CommentsList []Comment
	IsUgoira     bool
	BookmarkID   string
	IllustType   int `json:"illustType"`
}

func GetUserBasicInformation(r *http.Request, id string) (UserBrief, error) {
	var user UserBrief

	URL := GetUserInformationURL(id)

	response, err := API_GET_UnwrapJson(r.Context(), URL, "")
	if err != nil {
		return user, err
	}
	response = session.ProxyImageUrl(r, response)

	err = json.Unmarshal([]byte(response), &user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func GetArtworkImages(r *http.Request, id string, illustType int) ([]Image, error) {
	var resp []ImageResponse
	var images []Image

	URL := GetArtworkImagesURL(id)

	response, err := API_GET_UnwrapJson(r.Context(), URL, "")
	if err != nil {
		return nil, err
	}
	response = session.ProxyImageUrl(r, response)

	err = json.Unmarshal([]byte(response), &resp)
	if err != nil {
		return images, err
	}

	// Extract and proxy every images
	for _, imageRaw := range resp {
		var image Image

		// this is the original art dimention, not the "regular" art dimension
		// the image ratio of "regular" is close to Width/Height
		// maybe not useful
		image.Width = imageRaw.Width
		image.Height = imageRaw.Height

		image.Small = imageRaw.Urls["thumb_mini"]
		image.Medium = imageRaw.Urls["small"]
		image.Large = imageRaw.Urls["regular"]
		image.Original = imageRaw.Urls["original"]

		// Required for logic to display manga differently
		image.IllustType = illustType

		// Debug statement
		// log.Printf("Artwork ID: %s, IllustType set to %d", id, image.IllustType)

		images = append(images, image)
	}

	return images, nil
}

func GetArtworkComments(r *http.Request, id string) ([]Comment, error) {
	var body struct {
		Comments []Comment `json:"comments"`
	}

	URL := GetArtworkCommentsURL(id)

	response, err := API_GET_UnwrapJson(r.Context(), URL, "")
	if err != nil {
		return nil, err
	}
	response = session.ProxyImageUrl(r, response)

	err = json.Unmarshal([]byte(response), &body)
	if err != nil {
		return nil, err
	}

	return body.Comments, nil
}

func GetRelatedArtworks(r *http.Request, id string) ([]ArtworkBrief, error) {
	var body struct {
		Illusts []ArtworkBrief `json:"illusts"`
	}

	// TODO: keep the hard-coded limit?
	URL := GetArtworkRelatedURL(id, 96)

	response, err := API_GET_UnwrapJson(r.Context(), URL, "")
	if err != nil {
		return nil, err
	}

	response = session.ProxyImageUrl(r, response)

	err = json.Unmarshal([]byte(response), &body)
	if err != nil {
		return nil, err
	}

	return body.Illusts, nil
}

func GetArtworkByID(r *http.Request, id string, full bool) (*Illust, error) {
	token := session.GetUserToken(r)

	var illust struct {
		Illust
		UserIllusts map[int]any     `json:"userIllusts"`
		RawTags     json.RawMessage `json:"tags"`

		// In this object:
		// "userId": "54395645",
		// "userName": "ムロマキ",
		// "userAccount": "doberman6969",
		// missing avatar. if we have avatar, we can skip the user request below
	}
	var illust2 struct {
		Images       []Image
		RelatedWorks []ArtworkBrief
		CommentsList []Comment
	}

	wg := sync.WaitGroup{}
	cerr := make(chan error, 7)

	// Get basic illust information
	wg.Add(1)
	go func() {
		defer wg.Done()

		urlArtInfo := GetArtworkInformationURL(id)
		response, err := API_GET_UnwrapJson(r.Context(), urlArtInfo, token)
		if err != nil {
			cerr <- err
			return
		}

		err = json.Unmarshal([]byte(response), &illust)
		if err != nil {
			cerr <- err
			return
		}

		if illust.BookmarkData != nil {
			t := illust.BookmarkData.(map[string]any)
			illust.BookmarkID = t["id"].(string)
		}

		// Get basic user information (the URL above does not contain avatars)
		wg.Add(1)
		go func() {
			defer wg.Done()

			userInfo, err := GetUserBasicInformation(r, illust.UserID)
			if err != nil {
				cerr <- err
				return
			}
			illust.User = userInfo
		}()

		// Get illust images
		// Done after basic information is retrieved in order to properly populate IllustType
		wg.Add(1)
		go func() {
			defer wg.Done()

			images, err := GetArtworkImages(r, id, illust.IllustType)
			if err != nil {
				cerr <- err
				return
			}
			illust2.Images = images
		}()

		if full {
			// Get related artworks
			wg.Add(1)
			go func() {
				defer wg.Done()

				var err error
				related, err := GetRelatedArtworks(r, id)
				if err != nil {
					cerr <- err
					return
				}
				illust2.RelatedWorks = related
			}()
		}

		// translate tags
		wg.Add(1)
		go func() {
			defer wg.Done()

			var tags struct {
				Tags []struct {
					Tag         string            `json:"tag"`
					Translation map[string]string `json:"translation"`
				} `json:"tags"`
			}
			err := json.Unmarshal(illust.RawTags, &tags)
			if err != nil {
				cerr <- err
				return
			}

			var tagsList []Tag
			for _, tag := range tags.Tags {
				var newTag Tag
				newTag.Name = tag.Tag
				newTag.TranslatedName = tag.Translation["en"]

				tagsList = append(tagsList, newTag)
			}
			illust.Tags = tagsList
		}()

		if full {
			// Get recent artworks
			wg.Add(1)
			go func() {
				defer wg.Done()

				var err error
				ids := make([]int, 0)

				for k := range illust.UserIllusts {
					ids = append(ids, k)
				}

				sort.Sort(sort.Reverse(sort.IntSlice(ids)))

				idsString := ""
				count := min(len(ids), 20)

				for i := 0; i < count; i++ {
					idsString += fmt.Sprintf("&ids[]=%d", ids[i])
				}

				recent, err := GetUserArtworkList(r, illust.UserID, idsString)
				if err != nil {
					cerr <- err
					return
				}
				sort.Slice(recent[:], func(i, j int) bool {
					left := recent[i].ID
					right := recent[j].ID
					return numberGreaterThan(left, right)
				})
				illust.RecentWorks = recent
			}()
		}

		// Get reader comments
		//
		// Only fetch the comments if 'full' is requested and comments are not disabled (illust.CommentDisabled != 1).
		// This check needs to happen *after* fetching the basic artwork information, since that's when
		// 'CommentDisabled' is populated. If we check it too early, it would default to 0 (enabled),
		// leading to an invalid API call to fetch comments even when they are disabled, and causing an HTTP 500 error
		// on our end when we receive HTTP 400 from the Pixiv API as a result.
		if full && illust.CommentDisabled != 1 {
			wg.Add(1)
			go func() {
				defer wg.Done()

				comments, err := GetArtworkComments(r, id)
				if err != nil {
					cerr <- err
					return
				}
				illust2.CommentsList = comments
			}()
		}
	}()

	wg.Wait()
	close(cerr)

	illust.Images = illust2.Images
	illust.RelatedWorks = illust2.RelatedWorks
	illust.CommentsList = illust2.CommentsList

	all_errors := []error{}
	for suberr := range cerr {
		all_errors = append(all_errors, suberr)
	}
	err_summary := errors.Join(all_errors...)
	if err_summary != nil {
		return nil, err_summary
	}

	// If this artwork is an ugoira
	illust.IsUgoira = strings.Contains(illust.Images[0].Original, "ugoira")

	return &illust.Illust, nil
}

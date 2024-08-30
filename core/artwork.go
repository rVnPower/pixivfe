package core

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"codeberg.org/vnpower/pixivfe/v2/session"
	"github.com/goccy/go-json"
	"net/http"
)

// Pixiv returns 0, 1, 2 to filter SFW and/or NSFW artworks.
// Those values are saved in `xRestrict`
// 0: Safe
// 1: R18
// 2: R18G
type xRestrict int

const (
	Safe xRestrict = 0
	R18  xRestrict = 1
	R18G xRestrict = 2
)

var xRestrictModel = map[xRestrict]string{
	Safe: "",
	R18:  "R18",
	R18G: "R18G",
}

// Pixiv returns 0, 1, 2 to filter SFW and/or NSFW artworks.
// Those values are saved in `aiType`
// 0: Not rated / Unknown
// 1: Not AI-generated
// 2: AI-generated

type aiType int

const (
	Unrated aiType = 0
	NotAI   aiType = 1
	AI      aiType = 2
)

var aiTypeModel = map[aiType]string{
	Unrated: "Unrated",
	NotAI:   "Not AI",
	AI:      "AI",
}

type ImageResponse struct {
	Width  int               `json:"width"`
	Height int               `json:"height"`
	Urls   map[string]string `json:"urls"`
}

type Image struct {
	Width    int
	Height   int
	Small    string
	Medium   string
	Large    string
	Original string
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

type BookmarkData struct {
	Id string `json:"id"`
}

type ArtworkBrief struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	ArtistID     string        `json:"userId"`
	ArtistName   string        `json:"userName"`
	ArtistAvatar string        `json:"profileImageUrl"`
	Thumbnail    string        `json:"url"`
	Pages        int           `json:"pageCount"`
	XRestrict    int           `json:"xRestrict"`
	AiType       int           `json:"aiType"`
	BookmarkData *BookmarkData `json:"bookmarkData"`
	IllustType   int           `json:"illustType"`
}

type Illust struct {
	ID              string        `json:"id"`
	Title           string        `json:"title"`
	Description     HTML          `json:"description"`
	UserID          string        `json:"userId"`
	UserName        string        `json:"userName"`
	UserAccount     string        `json:"userAccount"`
	Date            time.Time     `json:"uploadDate"`
	Tags            []Tag         `json:"tags"`
	Pages           int           `json:"pageCount"`
	Bookmarks       int           `json:"bookmarkCount"`
	Likes           int           `json:"likeCount"`
	Comments        int           `json:"commentCount"`
	Views           int           `json:"viewCount"`
	CommentDisabled int           `json:"commentOff"`
	SanityLevel     int           `json:"sl"`
	XRestrict       xRestrict     `json:"xRestrict"`
	AiType          aiType        `json:"aiType"`
	BookmarkData    *BookmarkData `json:"bookmarkData"`
	Liked           bool          `json:"likeData"`
	Images          []Image
	User            UserBrief
	RecentWorks     []ArtworkBrief
	RelatedWorks    []ArtworkBrief
	CommentsList    []Comment
	IsUgoira        bool
	BookmarkID      string
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

func GetArtworkImages(r *http.Request, id string) ([]Image, error) {
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
	token := session.GetPixivToken(r)

	var illust struct {
		Illust
		UserIllusts map[int]*struct{} `json:"userIllusts"`
		RawTags     json.RawMessage   `json:"tags"`

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

	// Get illust images
	wg.Add(1)
	go func() {
		defer wg.Done()

		images, err := GetArtworkImages(r, id)
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

		// Get reader comments
		wg.Add(1)
		go func() {
			defer wg.Done()

			// if illust.CommentDisabled == 1 {
			// 	return
			// }
			comments, err := GetArtworkComments(r, id)
			if err != nil {
				cerr <- err
				return
			}
			illust2.CommentsList = comments
		}()
	}

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
			illust.BookmarkID = illust.BookmarkData.Id
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

				recent, err := GetUserArtworks(r, illust.UserID, idsString)
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

package core

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/server/session"
	"github.com/goccy/go-json"
)

type Novel struct {
	Bookmarks      int       `json:"bookmarkCount"`
	CommentCount   int       `json:"commentCount"`
	MarkerCount    int       `json:"markerCount"`
	CreateDate     time.Time `json:"createDate"`
	UploadDate     time.Time `json:"uploadDate"`
	Description    string    `json:"description"`
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Likes          int       `json:"likeCount"`
	Pages          int       `json:"pageCount"`
	UserID         string    `json:"userId"`
	UserName       string    `json:"userName"`
	Views          int       `json:"viewCount"`
	IsOriginal     bool      `json:"isOriginal"`
	IsBungei       bool      `json:"isBungei"`
	XRestrict      int       `json:"xRestrict"`
	Restrict       int       `json:"restrict"`
	Content        string    `json:"content"`
	CoverURL       string    `json:"coverUrl"`
	IsBookmarkable bool      `json:"isBookmarkable"`
	BookmarkData   any       `json:"bookmarkData"`
	LikeData       bool      `json:"likeData"`
	PollData       any       `json:"pollData"`
	Marker         any       `json:"marker"`
	Tags           struct {
		AuthorID string `json:"authorId"`
		IsLocked bool   `json:"isLocked"`
		Tags     []struct {
			Name string `json:"tag"`
		} `json:"tags"`
		Writable bool `json:"writable"`
	} `json:"tags"`
	SeriesNavData any  `json:"seriesNavData"`
	HasGlossary   bool `json:"hasGlossary"`
	IsUnlisted    bool `json:"isUnlisted"`
	// seen values: zh-cn, ja
	Language       string `json:"language"`
	CommentOff     int    `json:"commentOff"`
	CharacterCount int    `json:"characterCount"`
	WordCount      int    `json:"wordCount"`
	UseWordCount   bool   `json:"useWordCount"`
	ReadingTime    int    `json:"readingTime"`
	AiType         int    `json:"aiType"`
	Genre          string `json:"genre"`
	Settings       struct {
		ViewMode int `json:"viewMode"`
		// ...
	} `json:"suggestedSettings"`
	TextEmbeddedImages map[string]struct {
		NovelImageId string `json:"novelImageId"`
		SanityLevel  string `json:"sl"`
		Urls         struct {
			Two40Mw     string `json:"240mw"`
			Four80Mw    string `json:"480mw"`
			One200X1200 string `json:"1200x1200"`
			One28X128   string `json:"128x128"`
			Original    string `json:"original"`
		} `json:"urls"`
	} `json:"textEmbeddedImages"`
	CommentsList []Comment
}

type NovelBrief struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	XRestrict      int       `json:"xRestrict"`
	Restrict       int       `json:"restrict"`
	CoverURL       string    `json:"url"`
	Tags           []string  `json:"tags"`
	UserID         string    `json:"userId"`
	UserName       string    `json:"userName"`
	UserAvatar     string    `json:"profileImageUrl"`
	TextCount      int       `json:"textCount"`
	WordCount      int       `json:"wordCount"`
	ReadingTime    int       `json:"readingTime"`
	Description    string    `json:"description"`
	IsBookmarkable bool      `json:"isBookmarkable"`
	BookmarkData   any       `json:"bookmarkData"`
	Bookmarks      int       `json:"bookmarkCount"`
	IsOriginal     bool      `json:"isOriginal"`
	CreateDate     time.Time `json:"createDate"`
	UpdateDate     time.Time `json:"updateDate"`
	IsMasked       bool      `json:"isMasked"`
	SeriesID       string    `json:"seriesId"`
	SeriesTitle    string    `json:"seriesTitle"`
	IsUnlisted     bool      `json:"isUnlisted"`
	AiType         int       `json:"aiType"`
	Genre          string    `json:"genre"`
}

func GetNovelByID(r *http.Request, id string) (Novel, error) {
	var novel Novel

	URL := GetNovelURL(id)

	response, err := API_GET_UnwrapJson(r.Context(), URL, "")
	if err != nil {
		return novel, err
	}
	response = session.ProxyImageUrl(r, response)

	err = json.Unmarshal([]byte(response), &novel)
	if err != nil {
		return novel, err
	}

	// Novel embedded illusts
	re_r := regexp.MustCompile(`\[pixivimage:(\d+)\]`)
	re_d := regexp.MustCompile(`\d+`)
	re_t := regexp.MustCompile(`\"original\":\"(.+?)\"`)

	novel.Content = re_r.ReplaceAllStringFunc(novel.Content, func(s string) string {
		illustid := re_d.FindString(s)

		URL := GetInsertIllustURL(novel.ID, illustid)
		response, err := API_GET_UnwrapJson(r.Context(), URL, "")
		if err != nil {
			return "Cannot insert illust" + illustid
		}

		url := re_t.FindString(response)
		url = session.ProxyImageUrl(r, url[11:]) // truncate the "original":

		// make [pixivimage:illustid-index] jump to anchor
		link := fmt.Sprintf("/artworks/%s", strings.ReplaceAll(illustid, "-", "#"))
		return fmt.Sprintf(`<a href="%s" target="_blank"><img src=%s alt="%s"/></a>`, link, url, s)
	})

	re_u := regexp.MustCompile(`\[uploadedimage:(\d+)\]`)
	re_id := regexp.MustCompile(`\d+`)
	novel.Content = re_u.ReplaceAllStringFunc(novel.Content, func(s string) string {
		imageId := re_id.FindString(s)
		if val, ok := novel.TextEmbeddedImages[imageId]; ok {
			return fmt.Sprintf(`<img src=%s alt="%s"/>`, val.Urls.Original, s)
		}
		return s
	})

	return novel, nil
}

func GetNovelRelated(r *http.Request, id string) ([]NovelBrief, error) {
	var novels struct {
		List []NovelBrief `json:"novels"`
	}

	// hard-coded value, may change
	URL := GetNovelRelatedURL(id, 50)

	response, err := API_GET_UnwrapJson(r.Context(), URL, "")
	if err != nil {
		return novels.List, err
	}
	response = session.ProxyImageUrl(r, response)

	err = json.Unmarshal([]byte(response), &novels)
	if err != nil {
		return novels.List, err
	}

	return novels.List, nil
}

func GetNovelComments(r *http.Request, id string) ([]Comment, error) {
	var body struct {
		Comments []Comment `json:"comments"`
	}

	URL := GetNovelCommentsURL(id)

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

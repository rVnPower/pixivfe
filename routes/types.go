package routes

import (
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/template"
	"codeberg.org/vnpower/pixivision"
)

func Render[T any](w http.ResponseWriter, r *http.Request, data T) error {
	return template.Render(w, r, data)
}

// Tutorial: adding new types in this file
// Whenever you add new types, update `TestTemplates` in render_test.go to include the type in the test
// Do not use pointer in Data_* struct. faker will insert nil.
// Do not name template file a.b.jet.html or it won't be able to be used here, since Data_a.b is not a valid identifier.

type Data_error struct {
	Title string
	Error error
}
type Data_about struct {
	Time           string
	Version        string
	ImageProxy     string
	AcceptLanguage string
}
type Data_artwork struct {
	Illust          core.Illust
	Title           string
	MetaDescription string
	MetaImage       string
	MetaAuthor      string
	MetaAuthorID    string
}
type Data_artworkMulti struct {
	Artworks []core.Illust
	Title    string
}
type Data_userAtom struct {
	URL       string
	Title     string
	User      core.User
	Category  core.UserArtCategory
	Updated   string
	PageLimit int
	Page      int
	// MetaImage string
}
type Data_index struct {
	Title       string
	LoggedIn    bool
	Data        core.LandingArtworks
	NoTokenData core.Ranking
}
type Data_unauthorized struct{}
type Data_discovery struct {
	Artworks []core.ArtworkBrief
	Title    string
	Queries  template.PartialURL
}
type Data_novelDiscovery struct {
	Novels []core.NovelBrief
	Title  string
}
type Data_newest struct {
	Items []core.ArtworkBrief
	Title string
}
type Data_novel struct {
	Novel        core.Novel
	NovelRelated []core.NovelBrief
	User         core.UserBrief
	Title        string
	FontType     string
	ViewMode     string
	Language     string
}
type Data_following struct {
	Title    string
	Mode     string
	Artworks []core.ArtworkBrief
	CurPage  string
	Page     int
}
type Data_pixivision_index struct {
	Data []pixivision.Article
}

type Data_pixivision_article struct {
	Article pixivision.Article
}
type Data_rank struct {
	Title     string
	Page      int
	PageLimit int
	Date      string
	Data      core.Ranking
}
type Data_rankingCalendar struct {
	Title       string
	Render      core.HTML
	Mode        string
	Year        int
	MonthBefore DateWrap
	MonthAfter  DateWrap
	ThisMonth   DateWrap
}
type Data_settings struct {
	ProxyList        []string
	WorkingProxyList []string
}
type Data_tag struct {
	Title    string
	Tag      core.TagDetail
	Data     core.SearchResult
	QueriesC template.PartialURL
	TrueTag  string
	Page     int
}
type Data_user struct {
	Title     string
	User      core.User
	Category  core.UserArtCategory
	PageLimit int
	Page      int
	MetaImage string
}

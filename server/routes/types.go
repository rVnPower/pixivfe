package routes

import (
	"net/http"
	"net/url"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/server/request_context"
	"codeberg.org/vnpower/pixivfe/v2/server/template"
	"time"
)

func RenderHTML[T any](w http.ResponseWriter, r *http.Request, data T) error {
	return RenderWithContentType(w, r, "text/html; charset=utf-8", data)
}

func RenderWithContentType[T any](w http.ResponseWriter, r *http.Request, contentType string, data T) error {
	w.Header().Set("content-type", contentType)
	w.WriteHeader(request_context.Get(r).RenderStatusCode)
	return template.Render(w, template.GetTemplatingVariables(r), data)
}

// Tutorial: adding new types in this file
// Whenever you add new types, update `TestTemplates` in render_test.go to include the type in the test
//
// Warnings:
// - Do not use pointer in Data_* struct. faker will insert nil.
// - Do not name template file a.b.jet.html or it won't be able to be used here, since Data_a.b is not a valid identifier.

type Data_about struct {
	Time           string
	Version        string
	DomainName     string
	RepoURL        string
	Revision       string
	RevisionHash   string
	ImageProxy     url.URL
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
type Data_diagnostics struct{}
type Data_discovery struct {
	Artworks []core.ArtworkBrief
	Title    string
	Queries  template.PartialURL
}
type Data_error struct {
	Title string
	Error error
}
type Data_following struct {
	Title    string
	Mode     string
	Artworks []core.ArtworkBrief
	CurPage  string
	Page     int
}
type Data_index struct {
	Title       string
	LoggedIn    bool
	Data        core.LandingArtworks
	NoTokenData core.Ranking
	Queries     template.PartialURL
}
type Data_newest struct {
	Items []core.ArtworkBrief
	Title string
}
type Data_novel struct {
	Novel                    core.Novel
	NovelRelated             []core.NovelBrief
	NovelSeriesContentTitles []core.NovelSeriesContentTitle
	NovelSeriesIDs           []string
	NovelSeriesTitles        []string
	User                     core.UserBrief
	Title                    string
	FontType                 string
	ViewMode                 string
	Language                 string
}
type Data_novelSeries struct {
	NovelSeries         core.NovelSeries
	NovelSeriesContents []core.NovelSeriesContent
	User                core.UserBrief
	Title               string
	Page                int
	PageLimit           int
}
type Data_novelDiscovery struct {
	Novels  []core.NovelBrief
	Title   string
	Queries template.PartialURL
}
type Data_pixivisionIndex struct {
	Data []core.PixivisionArticle
}

type Data_pixivisionArticle struct {
	Article core.PixivisionArticle
}

type Data_pixivisionCategory struct {
	Category core.PixivisionCategory
}

type Data_pixivisionTag struct {
	Tag core.PixivisionTag
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
	Calendar    []core.DayCalendar
	Mode        string
	Year        int
	MonthBefore DateWrap
	MonthAfter  DateWrap
	ThisMonth   DateWrap
}
type Data_settings struct {
	ProxyList          []string
	WorkingProxyList   []string
	ProxyCheckEnabled  bool
	ProxyCheckInterval time.Duration
	DefaultProxyServer url.URL
}
type Data_tag struct {
	Title            string
	Tag              core.TagDetail
	Data             core.SearchResult
	QueriesC         template.PartialURL
	TrueTag          string
	Page             int
	ActiveCategory   string
	ActiveOrder      string
	ActiveMode       string
	ActiveRatio      string
	ActiveSearchMode string
}
type Data_unauthorized struct{}
type Data_user struct {
	Title     string
	User      core.User
	Category  core.UserArtCategory
	PageLimit int
	Page      int
	MetaImage string
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
type Data_mangaSeries struct {
	MangaSeriesContent core.MangaSeriesContent
	Title              string
	User               core.UserBrief
	Page               int
	PageLimit          int
}

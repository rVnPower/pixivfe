package routes

import "codeberg.org/vnpower/pixivfe/v2/core"

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
	Illust          core.Illust // faker can't fill this
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
	Title string
	IsLoggedIn bool
	Data  core.LandingArtworks
	NoTokenData core.Ranking
}

// below are unconverted. types may be wrong.
type Data_discovery struct {
	Artworks string
	Title    string
	Queries  string
}
type Data_novelDiscovery struct {
	Novels string
	Title  string
}
type Data_newest struct {
	Items string
	Title string
}
type Data_novel struct {
	Novel        string
	NovelRelated string
	User         string
	Title        string
	FontType     string
	ViewMode     string
	Language     string
}
type Data_unauthorized struct{}
type Data_following struct {
	Title    string
	Mode     string
	Artworks string
	CurPage  string
	Page     string
}

//	type Data_pixivisionindex struct {
//		Data string
//	}
//
//	type Data_pixivisionarticle struct {
//		Article string
//	}
type Data_rank struct {
	Title     string
	Page      string
	PageLimit int
	Date      string
	Data      string
}
type Data_rankingCalendar struct {
	Title       string
	Render      string
	Mode        string
	Year        string
	MonthBefore string
	MonthAfter  string
	ThisMonth   string
}

//	type Data_settings struct {
//		ProxyList string
//	}
//
//	type Data_tag struct {
//		Title string
//	}
type Data_user struct {
	Title     string
	User      string
	Category  string
	PageLimit int
	Page      string
	MetaImage string
}

// add new types above this line
// whenever you add new types, update `TestTemplates` in render_test.go to include the type in the test
// caution: do not use pointer in Data_* struct. faker will insert nil.
// caution: do not name template file a.b.jet.html or it won't be able to be used here, since Data_a.b is not a valid identifier.

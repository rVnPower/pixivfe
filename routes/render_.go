package routes

import (
	"log"
	"reflect"
	"strings"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"github.com/gofiber/fiber/v2"
)

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
type Data_index struct {
	Title string
	Data  string
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
// type Data_pixivisionindex struct {
// 	Data string
// }
// type Data_pixivisionarticle struct {
// 	Article string
// }
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
// type Data_settings struct {
// 	ProxyList string
// }
// type Data_tag struct {
// 	Title string
// }
type Data_user struct {
	Title     string
	User      string
	Category  string
	PageLimit int
	Page      string
	MetaImage string
}
// add new types above

// the migration plan
//
// 1. find and replace every occurance of `c.Render("abc" struct {...})` with `Render(c, Data_abc{...})` (except in this file)
// 2. create type `Data_abc` in this file (see `Data_error` above)
// 3. update `TestTemplates` in render_test.go to include `Data_abc`

func Render[T interface{}](c *fiber.Ctx, data T) error {
	template_name, found := strings.CutPrefix(reflect.TypeFor[T]().Name(), "Data_")
	if !found {
		log.Panicf("struct name does not start with 'Data_': %s", template_name)
	}
	return c.Render(template_name, StructToMap(data))
}

func StructToMap[T interface{}](data T) map[string]interface{} {
	result := map[string]interface{}{}
	Type := reflect.TypeFor[T]()
	for i := 0; i < Type.NumField(); i += 1 {
		field := Type.Field(i)
		result[field.Name] = reflect.ValueOf(data).FieldByName(field.Name).Interface()
	}
	return result
}

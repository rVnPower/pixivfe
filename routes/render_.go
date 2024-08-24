package routes

import (
	"log"
	"reflect"
	"strings"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/session"
	"codeberg.org/vnpower/pixivfe/v2/utils"

	"github.com/CloudyKit/jet/v6"
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

// global variable, yes.
var views *jet.Set

func InitTemplatingEngine(InDevelopment bool) {
	if InDevelopment {
		views = jet.NewSet(
			jet.NewOSFileSystemLoader("assets/views"),
			jet.InDevelopmentMode(), // disable cache
		)
	} else {
		views = jet.NewSet(
			jet.NewOSFileSystemLoader("assets/views"),
		)
	}
	for fn_name, fn := range utils.GetTemplateFunctions() {
		views.AddGlobal(fn_name, fn)
	}
}

func Render[T any](c *fiber.Ctx, data T) error {
	template_name, found := strings.CutPrefix(reflect.TypeFor[T]().Name(), "Data_")
	if !found {
		log.Panicf("struct name does not start with 'Data_': %s", template_name)
	}

	// Pass in values that we want to be available to all pages here
	token := session.GetPixivToken(c)
	pageURL := c.BaseURL() + c.OriginalURL()

	cookies := map[string]string{}
	for _, name := range session.AllCookieNames {
		value := session.GetCookie(c, name)
		cookies[string(name)] = value
	}

	template, err := views.GetTemplate(template_name + ".jet.html")
	if err != nil {
		return err
	}

	views.Parse(template_name + ".jet.html", template.String())

	variables := jet.VarMap{}

	// The middleware at line 99 in `main.go` cannot bind these values below if we use this function.
	variables.Set("BaseURL", c.BaseURL())
	variables.Set("OriginalURL", c.OriginalURL())
	variables.Set("PageURL", pageURL)
	variables.Set("LoggedIn", token != "")
	variables.Set("Queries", c.Queries())
	variables.Set("CookieList", cookies)

	// Type := reflect.TypeFor[T]()
	// for _, special_varname := range []string{"Title", "MetaAuthor", "MetaDescription", "MetaImage"} {
	// 	_, has_field := Type.FieldByName(special_varname)
	// 	if has_field {
	// 		variables.Set(special_varname, FieldName(data, special_varname))
	// 	}
	// }

	c.Context().SetContentType("text/html; charset=utf-8")
	return template.Execute(c.Response().BodyWriter(), variables, data)
}

func StructToMap[T any](data T) map[string]any {
	result := map[string]any{}
	Type := reflect.TypeFor[T]()
	for i := 0; i < Type.NumField(); i += 1 {
		field := Type.Field(i)
		result[field.Name] = FieldName(data, field.Name)
	}
	return result
}

// assumes that the field `field_name` exists, panics otherwise
func FieldName[T any](data T, field_name string) any {	
	return reflect.ValueOf(data).FieldByName(field_name).Interface()
}

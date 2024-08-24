package routes

import (
	"log"
	"reflect"
	"strings"

	"codeberg.org/vnpower/pixivfe/v2/session"
	"codeberg.org/vnpower/pixivfe/v2/utils"

	"github.com/CloudyKit/jet/v6"
	"github.com/gofiber/fiber/v2"
)

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

	c.Context().SetContentType("text/html; charset=utf-8")
	return template.Execute(c.Response().BodyWriter(), variables, data)
}

// func structToMap[T any](data T) map[string]any {
// 	result := map[string]any{}
// 	Type := reflect.TypeFor[T]()
// 	for i := 0; i < Type.NumField(); i += 1 {
// 		field := Type.Field(i)
// 		result[field.Name] = fieldName(data, field.Name)
// 	}
// 	return result
// }

// // assumes that the field `field_name` exists, panics otherwise
// func fieldName[T any](data T, field_name string) any {	
// 	return reflect.ValueOf(data).FieldByName(field_name).Interface()
// }

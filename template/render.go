package template

import (
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"

	"codeberg.org/vnpower/pixivfe/v2/session"
	"codeberg.org/vnpower/pixivfe/v2/utils"

	"github.com/CloudyKit/jet/v6"
)

// global variable, yes.
var views *jet.Set

func InitTemplatingEngine(DisableCache bool) {
	if DisableCache {
		views = jet.NewSet(
			jet.NewOSFileSystemLoader("assets/views"),
			jet.InDevelopmentMode(), // disable cache
		)
	} else {
		views = jet.NewSet(
			jet.NewOSFileSystemLoader("assets/views"),
		)
	}
	for fn_name, fn := range GetTemplateFunctions() {
		views.AddGlobal(fn_name, fn)
	}
}

// render the template selected based on the name of type `T`
func Render[T any](w http.ResponseWriter, r *http.Request, data T) error {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	// todo: think about caching a bit more
	// w.Header().Set("expires", time.Now().Add(config.ExpiresIn).Format(time.RFC1123))
	SetHTMLPrivacyHeaders(w, r)
	return RenderInner(w, GetTemplatingVariables(r), data)
}

func RenderInner[T any](w io.Writer, variables jet.VarMap, data T) error {
	template_name, found := strings.CutPrefix(reflect.TypeFor[T]().Name(), "Data_")
	if !found {
		log.Panicf("struct name does not start with 'Data_': %s", template_name)
	}

	template, err := views.GetTemplate(template_name + ".jet.html")
	if err != nil {
		return err
	}

	views.Parse(template_name+".jet.html", template.String())

	return template.Execute(w, variables, data)
}

func GetTemplatingVariables(r *http.Request) jet.VarMap {
	// Pass in values that we want to be available to all pages here
	token := session.GetPixivToken(r)
	baseURL := utils.Origin(r)
	pageURL := r.URL.String()

	cookies := map[string]string{}
	for _, name := range session.AllCookieNames {
		value := session.GetCookie(r, name)
		cookies[string(name)] = value
	}

	return jet.VarMap{}.
		Set("BaseURL", baseURL).
		Set("PageURL", pageURL).
		Set("LoggedIn", token != "").
		Set("Queries", r.URL.Query().Encode()).
		Set("CookieList", cookies)
}

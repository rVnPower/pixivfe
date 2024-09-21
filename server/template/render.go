package template

import (
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"

	"codeberg.org/vnpower/pixivfe/v2/server/request_context"
	"codeberg.org/vnpower/pixivfe/v2/server/session"
	"codeberg.org/vnpower/pixivfe/v2/server/utils"

	"github.com/CloudyKit/jet/v6"
)

// global variable, yes.
var views *jet.Set

func Init(DisableCache bool, assetsLocation string) {
	if DisableCache {
		views = jet.NewSet(
			jet.NewOSFileSystemLoader(assetsLocation),
			jet.InDevelopmentMode(), // disable cache
		)
	} else {
		views = jet.NewSet(
			jet.NewOSFileSystemLoader(assetsLocation),
		)
	}
	for fn_name, fn := range GetTemplateFunctions() {
		views.AddGlobal(fn_name, fn)
	}
}

// render the template selected based on the name of type `T`
func Render[T any](w http.ResponseWriter, r *http.Request, data T) error {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	// todo: think about caching a bit more. see doc/dev/features/caching.md
	// w.Header().Set("expires", time.Now().Add(config.ExpiresIn).Format(time.RFC1123))
	w.WriteHeader(request_context.Get(r).RenderStatusCode)
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
	token := session.GetUserToken(r)
	baseURL := utils.Origin(r)
	pageURL := r.URL.String()

	cookies := map[string]string{}
	for _, name := range session.AllCookieNames {
		value := session.GetCookie(r, name)
		cookies[string(name)] = value
	}

	queries := make(map[string]string)

	for k, v := range r.URL.Query() {
		queries[k] = v[0]
	}

	return jet.VarMap{}.
		Set("BaseURL", baseURL).
		Set("PageURL", pageURL).
		Set("LoggedIn", token != "").
		Set("Queries", queries).
		Set("CookieList", cookies)
}

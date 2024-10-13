package template

import (
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"

	"codeberg.org/vnpower/pixivfe/v2/server/session"
	"codeberg.org/vnpower/pixivfe/v2/server/utils"

	"github.com/CloudyKit/jet/v6"
)

// global variable, yes.
var views *jet.Set

func Init(DisableCache bool, assetsLocation string) {
	if DisableCache {
		views = jet.NewSet(
			NewLocalizedFSLoader(assetsLocation),
			jet.InDevelopmentMode(), // disable cache
		)
	} else {
		views = jet.NewSet(
			NewLocalizedFSLoader(assetsLocation),
		)
	}
	for fn_name, fn := range GetTemplateFunctions() {
		views.AddGlobal(fn_name, fn)
	}
}

func Render[T any](w io.Writer, variables jet.VarMap, data T) error {
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
	cookies_ordered := []struct {
		k string
		v string
	}{}
	for _, name := range session.AllCookieNames {
		value := session.GetCookie(r, name)
		cookies[string(name)] = value
		cookies_ordered = append(cookies_ordered, struct {
			k string
			v string
		}{string(name), value})
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
		Set("CookieList", cookies).
		Set("CookieListOrdered", cookies_ordered)
}

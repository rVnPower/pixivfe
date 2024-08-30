package template_test

import (
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/CloudyKit/jet/v6"
	"pgregory.net/rapid"

	. "codeberg.org/vnpower/pixivfe/v2/routes"
	"codeberg.org/vnpower/pixivfe/v2/template"
)

func TestMain(m *testing.M) {
	err := os.Chdir("..")
	if err != nil {
		panic(err)
	}
	template.Init(false)

	m.Run()
}

func allowAll[T any](T) bool {
	return true
}

func TestAutoRender(t *testing.T) {
	test[Data_about](t, allowAll)
	test[Data_artwork](t, allowAll)
	test[Data_artworkMulti](t, allowAll)
	test[Data_diagnostics](t, allowAll)
	test[Data_discovery](t, allowAll)
	test[Data_error](t, func(d Data_error) bool { return d.Error != nil })
	test[Data_following](t, allowAll)
	test[Data_index](t, allowAll)
	test[Data_newest](t, allowAll)
	test[Data_novel](t, allowAll)
	test[Data_novelDiscovery](t, allowAll)
	test[Data_pixivision_article](t, allowAll)
	test[Data_pixivision_index](t, allowAll)
	test[Data_rank](t, allowAll)
	test[Data_rankingCalendar](t, allowAll)
	test[Data_settings](t, allowAll)
	test[Data_tag](t, allowAll)
	test[Data_unauthorized](t, allowAll)
	test[Data_user](t, allowAll)
	test[Data_userAtom](t, allowAll)
}

func fakeData[T any](t *rapid.T, label string) T {
	return rapid.Make[T]().Draw(t, label)
}

// test template with fake data
func test[T any](t *testing.T, filter func(T) bool) {
	t.Run(
		reflect.TypeFor[T]().Name(),
		rapid.MakeCheck(func(t *rapid.T) {
			gen := rapid.Make[T]().Filter(filter)
			sample := gen.Draw(t, "sample")
			testWith(t, sample)
		}),
	)
}

func testWith[T any](t *rapid.T, data T) {
	route_name, found := strings.CutPrefix(reflect.TypeFor[T]().Name(), "Data_")
	if !found {
		log.Panicf("struct name does not start with 'Data_': %s", route_name)
	}

	// log.Print("Testing " + route_name)

	variables := jet.VarMap{}

	for k, v := range map[string]any{
		"BaseURL":    fakeData[string](t, "BaseURL"),
		"PageURL":    fakeData[string](t, "PageURL"),
		"LoggedIn":   fakeData[bool](t, "LoggedIn"),
		"Queries":    fakeData[map[string]string](t, "Queries"),
		"CookieList": fakeData[map[string]string](t, "CookieList"),
	} {
		variables.Set(k, v)
	}

	err := template.RenderInner(io.Discard, variables, data)

	if err != nil {
		template_name, _ := strings.CutPrefix(reflect.TypeFor[T]().Name(), "Data_")
		t.Errorf("while rendering template %s: %v", template_name, err)
	}
}

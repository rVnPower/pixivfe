// fly check templates
// by no means comprehensive

package template_test

import (
	"io"
	"log"
	"reflect"
	"strings"
	"testing"

	. "codeberg.org/vnpower/pixivfe/v2/server/routes"
	"codeberg.org/vnpower/pixivfe/v2/server/template"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-faker/faker/v4"
)

func TestMain(m *testing.M) {
	template.Init(false, "../../assets/views")
	m.Run()
}

func TestAutoRender(t *testing.T) {
	test[Data_about](t)
	test[Data_artwork](t)
	test[Data_artworkMulti](t)
	test[Data_diagnostics](t)
	test[Data_discovery](t)
	test[Data_error](t, Data_error{Title: fakeData[string](), Error: io.EOF})
	test[Data_following](t)
	test[Data_index](t)
	test[Data_newest](t)
	test[Data_novel](t)
	test[Data_novelDiscovery](t)
	test[Data_pixivisionArticle](t)
	test[Data_pixivisionIndex](t)
	test[Data_rank](t)
	test[Data_rankingCalendar](t)
	test[Data_settings](t)
	test[Data_tag](t)
	test[Data_unauthorized](t)
	test[Data_user](t)
	test[Data_userAtom](t)
	test[Data_novelSeries](t)
	test[Data_mangaSeries](t)
}

func fakeData[T any]() T {
	var data T
	faker.FakeData(&data)
	return data
}

// test template with fake data
func test[T any](t *testing.T, data ...T) {
	if len(data) == 0 {
		testWith(t, fakeData[T]())
	} else {
		testWith(t, data[0])
	}
}

func testWith[T any](t *testing.T, data T) {
	route_name, found := strings.CutPrefix(reflect.TypeFor[T]().Name(), "Data_")
	if !found {
		log.Panicf("struct name does not start with 'Data_': %s", route_name)
	}

	// log.Print("Testing " + route_name)

	variables := jet.VarMap{}

	for k, v := range map[string]any{
		"BaseURL":    fakeData[string](),
		"PageURL":    fakeData[string](),
		"LoggedIn":   fakeData[bool](),
		"Queries":    fakeData[map[string]string](),
		"CookieList": fakeData[map[string]string](),
	} {
		variables.Set(k, v)
	}

	err := template.Render(io.Discard, variables, data)

	if err != nil {
		template_name, _ := strings.CutPrefix(reflect.TypeFor[T]().Name(), "Data_")
		t.Errorf("while rendering template %s: %v", template_name, err)
	}
}

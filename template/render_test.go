package template_test

import (
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	. "codeberg.org/vnpower/pixivfe/v2/routes"
	template "codeberg.org/vnpower/pixivfe/v2/template"
	"github.com/CloudyKit/jet/v6"
	"github.com/go-faker/faker/v4"
)

func TestTemplates(t *testing.T) {
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
	test[Data_pixivision_article](t)
	test[Data_pixivision_index](t)
	test[Data_rank](t)
	test[Data_rankingCalendar](t)
	test[Data_settings](t)
	test[Data_tag](t)
	test[Data_unauthorized](t)
	test[Data_user](t)
	test[Data_userAtom](t)
}

func TestMain(m *testing.M) {
	err := os.Chdir("..")
	if err != nil {
		panic(err)
	}
	template.Init(false)

	m.Run()
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

	err := template.RenderInner(io.Discard, variables, data)

	if err != nil {
		template_name, _ := strings.CutPrefix(reflect.TypeFor[T]().Name(), "Data_")
		t.Errorf("while rendering template %s: %v", template_name, err)
	}
}

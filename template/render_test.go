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
	{
		data := fakeData[Data_error]()
		data.Error = io.EOF
		manualTest(t, data)
	}
	autoTest[Data_about](t)
	autoTest[Data_artwork](t)
	autoTest[Data_artworkMulti](t)
	autoTest[Data_discovery](t)
	autoTest[Data_following](t)
	autoTest[Data_index](t)
	autoTest[Data_newest](t)
	autoTest[Data_novel](t)
	autoTest[Data_novelDiscovery](t)
	autoTest[Data_pixivision_article](t)
	autoTest[Data_pixivision_index](t)
	autoTest[Data_rank](t)
	autoTest[Data_rankingCalendar](t)
	autoTest[Data_settings](t)
	autoTest[Data_tag](t)
	autoTest[Data_unauthorized](t)
	autoTest[Data_user](t)
	autoTest[Data_userAtom](t)
}

func TestMain(m *testing.M) {
	err := os.Chdir("..")
	if err != nil {
		panic(err)
	}
	template.InitTemplatingEngine(false)

	m.Run()
}

func fakeData[T any]() T {
	var data T
	faker.FakeData(&data)
	return data
}

// autoTest template with fake data
func autoTest[T any](t *testing.T) {
	manualTest(t, fakeData[T]())
}

func manualTest[T any](t *testing.T, data T) {
	route_name, found := strings.CutPrefix(reflect.TypeFor[T]().Name(), "Data_")
	if !found {
		log.Panicf("struct name does not start with 'Data_': %s", route_name)
	}

	// log.Println("Testing " + route_name)

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

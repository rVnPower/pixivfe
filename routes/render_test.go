package routes

import (
	"io"
	"log"
	"reflect"
	"strings"
	"testing"

	"codeberg.org/vnpower/pixivfe/v2/utils"
	"github.com/go-faker/faker/v4"
	"github.com/gofiber/template/jet/v2"
)

func TestTemplates(t *testing.T) {
	autoTest[Data_error](t)
	autoTest[Data_about](t)
	autoTest[Data_artwork](t)
	autoTest[Data_artworkMulti](t)
	autoTest[Data_discovery](t)
	autoTest[Data_novelDiscovery](t)
	autoTest[Data_index](t)
	autoTest[Data_newest](t)
	autoTest[Data_novel](t)
	autoTest[Data_unauthorized](t)
	autoTest[Data_following](t)
	autoTest[Data_pixivision_index](t)
	autoTest[Data_pixivision_article](t)
	autoTest[Data_rank](t)
	autoTest[Data_rankingCalendar](t)
	autoTest[Data_settings](t)
	autoTest[Data_tag](t)
	autoTest[Data_user](t)
	autoTest[Data_userAtom](t)
}

var engine *jet.Engine

func TestMain(m *testing.M) {
	engine = jet.New("../assets/views", ".jet.html")
	engine.AddFuncMap(utils.GetTemplateFunctions())

	// gofiber bug: no error even if the templates are invalid??? https://github.com/gofiber/template/issues/341
	err := engine.Load()
	if err != nil {
		panic(err)
	}

	m.Run()
}

// autoTest template with fake data
func autoTest[T any](t *testing.T) {
	var data T
	faker.FakeData(&data)
	manualTest(t, data)
}

func manualTest[T any](t *testing.T, data T) {
	route_name, found := strings.CutPrefix(reflect.TypeFor[T]().Name(), "Data_")
	if !found {
		log.Panicf("struct name does not start with 'Data_': %s", route_name)
	}
	bindings := structToMap(data)

	for k, v := range map[string]any{
		"BaseURL":     "",
		"OriginalURL": "",
		"PageURL":     "",
		"LoggedIn":    false,
		"Queries":     map[string]string{},
		"CookieList":  map[string]string{},
	} {
		bindings[k] = v
	}

	err := engine.Render(io.Discard, route_name, bindings)

	if err != nil {
		template_name, _ := strings.CutPrefix(reflect.TypeFor[T]().Name(), "Data_")
		t.Errorf("while rendering template %s: %v", template_name, err)
	}
}

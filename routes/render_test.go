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
	test[Data_error](t)
	test[Data_about](t)
	test[Data_artwork](t)
	test[Data_artworkMulti](t)
	test[Data_discovery](t)
	test[Data_novelDiscovery](t)
	test[Data_index](t)
	test[Data_newest](t)
	test[Data_novel](t)
	test[Data_unauthorized](t)
	test[Data_following](t)
	test[Data_pixivisionindex](t)
	test[Data_pixivisionarticle](t)
	test[Data_rank](t)
	test[Data_rankingCalendar](t)
	test[Data_settings](t)
	test[Data_tag](t)
	test[Data_user](t)
	test[Data_userAtom](t)
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

// test template
func test[T any](t *testing.T) {	
	var data T
	faker.FakeData(&data)

	route_name, found := strings.CutPrefix(reflect.TypeFor[T]().Name(), "Data_")
	if !found {
		log.Panicf("struct name does not start with 'Data_': %s", route_name)
	}
	bindings := StructToMap(data)

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


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
}

var engine *jet.Engine

func TestMain(m *testing.M) {
	engine = jet.New("../assets/layout", ".jet.html")
	engine.AddFuncMap(utils.GetTemplateFunctions())

	// gofiber bug: no error even if the templates are invalid??? https://github.com/gofiber/template/issues/341
	err := engine.Load()
	if err != nil {
		panic(err)
	}

	m.Run()
}

// test template
func test[T interface{}](t *testing.T) {
	var data T
	faker.FakeData(&data)

	route_name, found := strings.CutPrefix(reflect.TypeFor[T]().Name(), "Data_")
	if !found {
		log.Panicf("struct name does not start with 'Data_': %s", route_name)
	}
	bindings := StructToMap(data)

	err := engine.Render(io.Discard, route_name, bindings)

	if err != nil {
		template_name, _ := strings.CutPrefix(reflect.TypeFor[T]().Name(), "Data_")
		t.Errorf("while rendering template %s: %v", template_name, err)
	}
}


package routes

import (
	"log"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Data_error struct {
	Title string
	Error error
}

// add new types above

// the migration plan
//
// 1. find and replace every occurance of `c.Render("abc", fiber.Map{...})` with `Render(c, Data_abc{...})` (except in this file)
// 2. create type `Data_abc` in this file (see `Data_error` above)
// 3. update `TestTemplates` in render_test.go to include `Data_abc`

func Render[T interface{}](c *fiber.Ctx, data T) error {
	template_name, found := strings.CutPrefix(reflect.TypeFor[T]().Name(), "Data_")
	if !found {
		log.Panicf("struct name does not start with 'Data_': %s", template_name)
	}
	return c.Render(template_name, StructToMap(data))
}

func StructToMap[T interface{}](data T) map[string]interface{} {
	result := map[string]interface{}{}
	Type := reflect.TypeFor[T]()
	for i := 0; i < Type.NumField(); i += 1 {
		field := Type.Field(i)
		result[field.Name] = reflect.ValueOf(data).FieldByName(field.Name).Interface()
	}
	return result
}

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

func Render[T interface{}](c *fiber.Ctx, data T) error {
	route_name, found := strings.CutPrefix(reflect.TypeFor[T]().Name(), "Data_")
	if !found { log.Panicf("struct name does not start with 'Data_': %s", route_name) }
	return c.Render(route_name, structToMap(data))
}

func structToMap[T interface{}](data T) map[string]interface{} {
	result := map[string]interface{}{}
	Type := reflect.TypeFor[T]()
	for i := 0; i < Type.NumField(); i += 1 {
		field := Type.Field(i)
		result[field.Name] = reflect.ValueOf(data).FieldByName(field.Name).Interface()
	}
	return result
}

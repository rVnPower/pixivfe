package routes

import (
	"log"
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/server/request_context"
	"codeberg.org/vnpower/pixivfe/v2/server/template"
)

func ErrorPage(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
	request_context.Get(r).RenderStatusCode = statusCode
	err = template.Render(w, r, Data_error{Title: "Error", Error: err})
	if err != nil {
		log.Printf("Error rendering error route: %s", err)
	}
}
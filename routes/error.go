package routes

import (
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/template"
)

func ErrorPage(w http.ResponseWriter, r *http.Request, err error) error {
	return template.Render(w, r, Data_error{Title: "Error", Error: err})
}

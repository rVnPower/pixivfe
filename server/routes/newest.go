package routes

import (
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/core"
)

func NewestPage(w http.ResponseWriter, r *http.Request) error {
	worktype := GetQueryParam(r, "type", "illust")

	r18 := GetQueryParam(r, "r18", "false")

	works, err := core.GetNewestArtworks(r, worktype, r18)
	if err != nil {
		return err
	}

	return RenderHTML(w, r, Data_newest{Items: works, Title: "Newest works"})
}

package routes

import (
	"codeberg.org/vnpower/pixivfe/v2/core"
	"net/http"
)

func NewestPage(w http.ResponseWriter, r CompatRequest) error {
	worktype := r.Query("type", "illust")

	r18 := r.Query("r18", "false")

	works, err := core.GetNewestArtworks(r.Request, worktype, r18)
	if err != nil {
		return err
	}

	return Render(w, r, Data_newest{Items: works, Title: "Newest works"})
}

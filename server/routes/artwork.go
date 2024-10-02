package routes

import (
	"fmt"
	"strconv"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"net/http"
)

func ArtworkPage(w http.ResponseWriter, r *http.Request) error {
	id := GetPathVar(r, "id")
	if _, err := strconv.Atoi(id); err != nil {
		return fmt.Errorf("Invalid ID: %s", id)
	}

	illust, err := core.GetArtworkByID(r, id, true)
	if err != nil {
		return err
	}

	metaDescription := ""
	for _, i := range illust.Tags {
		metaDescription += "#" + i.Name + ", "
	}

	// monkey patching. assuming illust.Images[_].Large is used
	for _, img := range illust.Images {
		PreloadImage(w, img.Large)
	}

	return RenderHTML(w, r, Data_artwork{
		Illust:          *illust,
		Title:           illust.Title,
		MetaDescription: metaDescription,
		MetaImage:       illust.Images[0].Original,
		MetaAuthor:      illust.UserName,
		MetaAuthorID:    illust.UserID,
	})
}

func PreloadImage(w http.ResponseWriter, url string) {
	w.Header().Add("Link", fmt.Sprintf("<%s>; rel=preload; as=image", url))
}

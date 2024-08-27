package routes

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"net/http"
)

func ArtworkMultiPage(w http.ResponseWriter, r *http.Request) error {
	ids_ := GetPathVar(r, "ids")
	ids := strings.Split(ids_, ",")

	artworks := make([]core.Illust, len(ids))

	wg := sync.WaitGroup{}
	// // gofiber/fasthttp's API is trash
	// // i can't replace r.Context() with this
	// // so i guess we will have to wait for network traffic to finish on error
	// ctx, cancel := context.WithCancel(r.Context())
	// defer cancel()
	// r.SetUserContext(ctx)
	var err_global error = nil
	for i, id := range ids {
		if _, err := strconv.Atoi(id); err != nil {
			err_global = fmt.Errorf("Invalid ID: %s", id)
			break
		}

		wg.Add(1)
		go func(i int, id string) {
			defer wg.Done()

			illust, err := core.GetArtworkByID(r, id, false)
			if err != nil {
				artworks[i] = core.Illust{
					Title: err.Error(), // this might be flaky
				}
				return
			}

			metaDescription := ""
			for _, i := range illust.Tags {
				metaDescription += "#" + i.Name + ", "
			}

			artworks[i] = *illust
		}(i, id)
	}
	// if err_global != nil {
	// 	cancel()
	// }
	wg.Wait()
	if err_global != nil {
		return err_global
	}

	return Render(w, r, Data_artworkMulti{
		Artworks: artworks,
		Title:    fmt.Sprintf("(%d images)", len(artworks)),
	})
}

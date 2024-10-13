package routes

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"codeberg.org/vnpower/pixivfe/v2/core"
)

type userPageData struct {
	user      core.User
	category  core.UserArtCategory
	pageLimit int
	page      int
}

func fetchData(r *http.Request, getTags bool) (userPageData, error) {
	id := GetPathVar(r, "id")
	if _, err := strconv.Atoi(id); err != nil {
		return userPageData{}, err
	}
	category := core.UserArtCategory(GetPathVar(r, "category", string(core.UserArt_Any)))
	err := category.Validate()
	if err != nil {
		return userPageData{}, err
	}

	page_param := GetQueryParam(r, "page", "1")
	page, err := strconv.Atoi(page_param)
	if err != nil {
		return userPageData{}, err
	}

	user, err := core.GetUserProfile(r, id, category, page, getTags)
	if err != nil {
		return userPageData{}, err
	}

	var worksCount int
	var worksPerPage float64

	if category == core.UserArt_Bookmarks {
		worksPerPage = 48.0
	} else {
		worksPerPage = 30.0
	}

	worksCount = user.CategoryItemCount
	pageLimit := int(math.Ceil(float64(worksCount) / worksPerPage))

	return userPageData{
		user:      user,
		category:  category,
		pageLimit: pageLimit,
		page:      page,
	}, nil
}

func UserPage(w http.ResponseWriter, r *http.Request) error {
	data, err := fetchData(r, true)
	if err != nil {
		return err
	}

	return RenderHTML(w, r, Data_user{
		Title:     data.user.Name,
		User:      data.user,
		Category:  data.category,
		PageLimit: data.pageLimit,
		Page:      data.page,
		MetaImage: data.user.BackgroundImage,
	})
}

func UserAtomFeed(w http.ResponseWriter, r *http.Request) error {
	data, err := fetchData(r, false)
	if err != nil {
		return err
	}

	return RenderWithContentType(w, r, "application/atom+xml", Data_userAtom{
		URL:       r.RequestURI,
		Title:     data.user.Name,
		User:      data.user,
		Category:  data.category,
		Updated:   time.Now().Format(time.RFC3339),
		PageLimit: data.pageLimit,
		Page:      data.page,
		// MetaImage: data.user.BackgroundImage,
	})
}

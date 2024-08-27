package routes

import (
	"net/http"
	"net/url"
	"strconv"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/utils"
)

func TagPage(w http.ResponseWriter, r CompatRequest) error {
	param := r.Params("name", r.Query("name"))
	name, err := url.PathUnescape(param)
	if err != nil {
		return err
	}

	page := r.Query("page", "1")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return err
	}

	// Because of the large amount of queries available for this route,
	// I made a struct type just to manage the queries
	queries := core.SearchPageSettings{
		Name:     name,
		Category: r.Query("category", "artworks"),
		Order:    r.Query("order", "date_d"),
		Mode:     r.Query("mode", "safe"),
		Ratio:    r.Query("ratio", ""),
		Wlt:      r.Query("wlt", ""),
		Wgt:      r.Query("wgt", ""),
		Hlt:      r.Query("hlt", ""),
		Hgt:      r.Query("hgt", ""),
		Tool:     r.Query("tool", ""),
		Scd:      r.Query("scd", ""),
		Ecd:      r.Query("ecd", ""),
		Page:     page,
	}

	tag, err := core.GetTagData(r.Request, name)
	if err != nil {
		return err
	}
	result, err := core.GetSearch(r.Request, queries)
	if err != nil {
		return err
	}

	urlc := utils.PartialURL{Path: "tags", Query: queries.ReturnMap()}

	return Render(w, r, Data_tag{Title: "Results for " + name, Tag: tag, Data: *result, QueriesC: urlc, TrueTag: param, Page: pageInt})
}

func AdvancedTagPost(w http.ResponseWriter, r CompatRequest) error {
	return utils.RedirectToRoute(w, r,"/tags", map[string]string{
		"name":     r.Query("name", r.FormValue("name")),
		"category": r.Query("category", "artworks"),
		"order":    r.Query("order", "date_d"),
		"mode":     r.Query("mode", "safe"),
		"ratio":    r.Query("ratio"),
		"page":     r.Query("page", "1"),
		"wlt":      r.Query("wlt", r.FormValue("wlt")),
		"wgt":      r.Query("wgt", r.FormValue("wgt")),
		"hlt":      r.Query("hlt", r.FormValue("hlt")),
		"hgt":      r.Query("hgt", r.FormValue("hgt")),
		"tool":     r.Query("tool", r.FormValue("tool")),
		"scd":      r.Query("scd", r.FormValue("scd")),
		"ecd":      r.Query("ecd", r.FormValue("ecd")),
	}, http.StatusFound)
}

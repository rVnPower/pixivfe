package routes

import (
	"net/http"
	"net/url"
	"strconv"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/utils"
)

func TagPage(w http.ResponseWriter, r *http.Request) error {
	param := GetPathVar(r, "name", GetQueryParam(r, "name"))
	name, err := url.PathUnescape(param)
	if err != nil {
		return err
	}

	page := GetQueryParam(r, "page", "1")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return err
	}

	// Because of the large amount of queries available for this route,
	// I made a struct type just to manage the queries
	queries := core.SearchPageSettings{
		Name:     name,
		Category: GetQueryParam(r, "category", "artworks"),
		Order:    GetQueryParam(r, "order", "date_d"),
		Mode:     GetQueryParam(r, "mode", "safe"),
		Ratio:    GetQueryParam(r, "ratio", ""),
		Wlt:      GetQueryParam(r, "wlt", ""),
		Wgt:      GetQueryParam(r, "wgt", ""),
		Hlt:      GetQueryParam(r, "hlt", ""),
		Hgt:      GetQueryParam(r, "hgt", ""),
		Tool:     GetQueryParam(r, "tool", ""),
		Scd:      GetQueryParam(r, "scd", ""),
		Ecd:      GetQueryParam(r, "ecd", ""),
		Page:     page,
	}

	tag, err := core.GetTagData(r, name)
	if err != nil {
		return err
	}
	result, err := core.GetSearch(r, queries)
	if err != nil {
		return err
	}

	urlc := utils.PartialURL{Path: "tags", Query: queries.ReturnMap()}

	return Render(w, r, Data_tag{Title: "Results for " + name, Tag: tag, Data: *result, QueriesC: urlc, TrueTag: param, Page: pageInt})
}

func AdvancedTagPost(w http.ResponseWriter, r *http.Request) error {
	return utils.RedirectTo(w, r,"/tags", map[string]string{
		"name":     GetQueryParam(r, "name", r.FormValue("name")),
		"category": GetQueryParam(r, "category", "artworks"),
		"order":    GetQueryParam(r, "order", "date_d"),
		"mode":     GetQueryParam(r, "mode", "safe"),
		"ratio":    GetQueryParam(r, "ratio"),
		"page":     GetQueryParam(r, "page", "1"),
		"wlt":      GetQueryParam(r, "wlt", r.FormValue("wlt")),
		"wgt":      GetQueryParam(r, "wgt", r.FormValue("wgt")),
		"hlt":      GetQueryParam(r, "hlt", r.FormValue("hlt")),
		"hgt":      GetQueryParam(r, "hgt", r.FormValue("hgt")),
		"tool":     GetQueryParam(r, "tool", r.FormValue("tool")),
		"scd":      GetQueryParam(r, "scd", r.FormValue("scd")),
		"ecd":      GetQueryParam(r, "ecd", r.FormValue("ecd")),
	})
}

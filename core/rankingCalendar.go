package core

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"

	"codeberg.org/vnpower/pixivfe/v2/session"
)

func get_weekday(n time.Weekday) int {
	switch n {
	case time.Sunday:
		return 1
	case time.Monday:
		return 2
	case time.Tuesday:
		return 3
	case time.Wednesday:
		return 4
	case time.Thursday:
		return 5
	case time.Friday:
		return 6
	case time.Saturday:
		return 7
	}
	return 0
}

var selector_img = cascadia.MustCompile("img")

// note(@iacore):
// so the funny thing about Pixiv is that they will return this month's data for a request of a future date
// is it a bug or a feature?
func GetRankingCalendar(r *http.Request, mode string, year, month int) (HTML, error) {
	token := session.GetPixivToken(r)
	URL := GetRankingCalendarURL(mode, year, month)

	resp, err := PixivGetRequest(r.Context(), URL, token)
	if err != nil {
		return "", err
	}

	// Use the html package to parse the response body from the request
	doc, err := html.Parse(strings.NewReader(resp.Body))
	if err != nil {
		return "", err
	}
	
	// Find and print all links on the web page
	var links []string
	for _, node := range cascadia.QueryAll(doc, selector_img) {
		for _, attr := range node.Attr {
			if attr.Key == "data-src" {
				// adds a new link entry when the attribute matches
				links = append(links, session.ProxyImageUrlNoEscape(r, attr.Val))
			}
		}
	}

	// now := r.Context().Time()
	// yearNow := now.Year()
	// monthNow := now.Month()
	lastMonth := time.Date(year, time.Month(month), 0, 0, 0, 0, 0, time.UTC)
	thisMonth := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC)

	renderString := ""
	for i := 0; i < get_weekday(lastMonth.Weekday()); i++ {
		renderString += "<div class=\"calendar-node calendar-node-empty\"></div>"
	}
	for i := 0; i < thisMonth.Day(); i++ {
		date := fmt.Sprintf("%d%02d%02d", year, month, i+1)
		if len(links) > i {
			renderString += fmt.Sprintf(`<a href="/ranking?mode=%s&date=%s"><div class="calendar-node"><img src="%s" alt="Day %d" /><span>%d</span></div></a>`, mode, date, links[i], i+1, i+1)
		} else {
			renderString += fmt.Sprintf(`<div class="calendar-node"><span>%d</span></div>`, i+1)
		}
	}
	return HTML(renderString), nil
}

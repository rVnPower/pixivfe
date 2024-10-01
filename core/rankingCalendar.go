package core

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"

	"codeberg.org/vnpower/pixivfe/v2/server/session"
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
	token := session.GetUserToken(r)
	URL := GetRankingCalendarURL(mode, year, month)

	resp, err := API_GET(r.Context(), URL, token)
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

	renderString := "<tr>"
	dayCount := 0
	for i := 0; i < get_weekday(lastMonth.Weekday()); i++ {
		renderString += `<td class="calendar-node calendar-node-empty"></td>`
		dayCount++
	}
	for i := 0; i < thisMonth.Day(); i++ {
		if dayCount == 7 {
			renderString += "</tr><tr>"
			dayCount = 0
		}
		date := fmt.Sprintf("%d%02d%02d", year, month, i+1)
		if len(links) > i {
			renderString += fmt.Sprintf(`<td class="calendar-node"><a href="/ranking?mode=%s&date=%s" class="d-block position-relative"><img src="%s" alt="Day %d" class="img-fluid" /><span class="position-absolute bottom-0 end-0 bg-white px-2 rounded-pill">%d</span></a></td>`, mode, date, links[i], i+1, i+1)
		} else {
			renderString += fmt.Sprintf(`<td class="calendar-node"><span class="d-block text-center">%d</span></td>`, i+1)
		}
		dayCount++
	}
	for dayCount < 7 {
		renderString += `<td class="calendar-node calendar-node-empty"></td>`
		dayCount++
	}
	renderString += "</tr>"
	return HTML(renderString), nil
}

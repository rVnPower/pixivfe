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

// get_weekday converts a time.Weekday to an integer representation.
// Sunday is 1, Monday is 2, and so on. This is used for calendar calculations.
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

// selector_img is a pre-compiled CSS selector for finding <img> tags in HTML.
var selector_img = cascadia.MustCompile("img")

// GetRankingCalendar retrieves and processes the ranking calendar data from Pixiv.
// It returns an HTML string representation of the calendar and any error encountered.
//
// iacore: so the funny thing about Pixiv is that they will return this month's data for a request of a future date. is it a bug or a feature?
func GetRankingCalendar(r *http.Request, mode string, year, month int) (HTML, error) {
	// Retrieve the user token from the session
	token := session.GetUserToken(r)
	URL := GetRankingCalendarURL(mode, year, month)

	// Make an API request to Pixiv
	resp, err := API_GET(r.Context(), URL, token)
	if err != nil {
		return "", err
	}

	// Parse the HTML response
	doc, err := html.Parse(strings.NewReader(resp.Body))
	if err != nil {
		return "", err
	}

	// Extract image links from the parsed HTML
	var links []string
	for _, node := range cascadia.QueryAll(doc, selector_img) {
		for _, attr := range node.Attr {
			if attr.Key == "data-src" {
				// Proxy the image URL to avoid direct requests to Pixiv
				links = append(links, session.ProxyImageUrlNoEscape(r, attr.Val))
			}
		}
	}

	// Calculate the last day of the previous month and the current month
	lastMonth := time.Date(year, time.Month(month), 0, 0, 0, 0, 0, time.UTC)
	thisMonth := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC)

	// Generate the HTML for the calendar
	renderString := "<tr>"
	dayCount := 0

	// Add empty cells for days before the 1st of the month
	for i := 0; i < get_weekday(lastMonth.Weekday()); i++ {
		renderString += `<td class="calendar-node calendar-node-empty"></td>`
		dayCount++
	}

	// Add cells for each day of the month
	for i := 0; i < thisMonth.Day(); i++ {
		// Start a new row if necessary
		if dayCount == 7 {
			renderString += "</tr><tr>"
			dayCount = 0
		}

		// Format the date string
		date := fmt.Sprintf("%d%02d%02d", year, month, i+1)

		// Add a cell with an image link if available, otherwise just the day number
		if len(links) > i {
			renderString += fmt.Sprintf(`<td class="calendar-node"><a href="/ranking?mode=%s&date=%s" class="d-block position-relative"><img src="%s" alt="Day %d" class="img-fluid" /><span class="position-absolute bottom-0 end-0 bg-body-tertiary px-2 rounded-pill">%d</span></a></td>`, mode, date, links[i], i+1, i+1)
		} else {
			renderString += fmt.Sprintf(`<td class="calendar-node"><span class="d-block text-center">%d</span></td>`, i+1)
		}
		dayCount++
	}

	// Add empty cells to complete the last row if necessary
	for dayCount < 7 {
		renderString += `<td class="calendar-node calendar-node-empty"></td>`
		dayCount++
	}

	renderString += "</tr>"
	return HTML(renderString), nil
}

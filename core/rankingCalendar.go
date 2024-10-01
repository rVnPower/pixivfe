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
// It returns an HTML string representation of the calendar using Bootstrap cards and any error encountered.
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

	// Generate the HTML for the calendar using Bootstrap cards
	renderString := ""
	dayCount := 0

	// Add empty cards for days before the 1st of the month
	for i := 0; i < get_weekday(lastMonth.Weekday()); i++ {
		renderString += `<div class="col">
			<div class="card h-100 border-0 ratio ratio-1x1"></div>
		</div>`
		dayCount++
	}

	// Add cards for each day of the month
	for i := 0; i < thisMonth.Day(); i++ {
		// Start a new row if necessary
		if dayCount == 7 {
			renderString += `</div><div class="row g-3 mb-3">`
			dayCount = 0
		}

		// Format the date string
		date := fmt.Sprintf("%d%02d%02d", year, month, i+1)

		// Add a card with an image link if available, otherwise just the day number
		if len(links) > i {
			renderString += fmt.Sprintf(`
				<div class="col">
					<div class="card h-100 ratio ratio-1x1">
						<a href="/ranking?mode=%s&date=%s" class="text-decoration-none">
							<img src="%s" alt="Day %d" class="card-img-top img-fluid object-fit-cover h-100" />
							<div class="card-img-overlay d-flex align-items-end">
								<p class="card-text text-white bg-dark bg-opacity-50 rounded px-2 py-1 mb-0">%d</p>
							</div>
						</a>
					</div>
				</div>`, mode, date, links[i], i+1, i+1)
		} else {
			renderString += fmt.Sprintf(`
				<div class="col">
					<div class="card h-100 ratio ratio-1x1">
						<div class="card-body d-flex align-items-center justify-content-center">
							<p class="card-text text-center mb-0">%d</p>
						</div>
					</div>
				</div>`, i+1)
		}
		dayCount++
	}

	// Add empty cards to complete the last row if necessary
	for dayCount < 7 {
		renderString += `<div class="col">
			<div class="card h-100 border-0 ratio ratio-1x1"></div>
		</div>`
		dayCount++
	}

	return HTML(renderString), nil
}

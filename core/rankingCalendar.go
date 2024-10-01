package core

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"

	"codeberg.org/vnpower/pixivfe/v2/server/session"
)

// DayCalendar represents the data for a single day in the ranking calendar
type DayCalendar struct {
	DayNumber   int    // The day of the month
	ImageURL    string // Proxy URL to the image (optional, can be empty when no image is available)
	ArtworkLink string // The link to the artwork page for this day
}

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

// extractArtworkID extracts the artwork ID from the image URL
func extractArtworkID(imageURL string) string {
	re := regexp.MustCompile(`/(\d+)_p0_(custom|square)1200\.jpg`)
	matches := re.FindStringSubmatch(imageURL)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// GetRankingCalendar retrieves and processes the ranking calendar data from Pixiv.
// It returns a slice of DayCalendar structs and any error encountered.
//
// iacore: so the funny thing about Pixiv is that they will return this month's data for a request of a future date. is it a bug or a feature?
func GetRankingCalendar(r *http.Request, mode string, year, month int) ([]DayCalendar, error) {
	// Retrieve the user token from the session
	token := session.GetUserToken(r)
	URL := GetRankingCalendarURL(mode, year, month)

	// Make an API request to Pixiv
	resp, err := API_GET(r.Context(), URL, token)
	if err != nil {
		return nil, err
	}

	// Parse the HTML response
	doc, err := html.Parse(strings.NewReader(resp.Body))
	if err != nil {
		return nil, err
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

	// Generate the calendar data
	var calendar []DayCalendar
	dayCount := 0

	// Add empty days for days before the 1st of the month
	for i := 0; i < get_weekday(lastMonth.Weekday()); i++ {
		calendar = append(calendar, DayCalendar{DayNumber: 0})
		dayCount++
	}

	// Add data for each day of the month
	for i := 0; i < thisMonth.Day(); i++ {
		day := DayCalendar{
			DayNumber: i + 1,
		}
		if len(links) > i {
			day.ImageURL = links[i]
			artworkID := extractArtworkID(links[i])
			if artworkID != "" {
				day.ArtworkLink = fmt.Sprintf("/artworks/%s", artworkID)
			}
		}
		calendar = append(calendar, day)
		dayCount++
	}

	// Add empty days to complete the last week if necessary
	for dayCount%7 != 0 {
		calendar = append(calendar, DayCalendar{DayNumber: 0})
		dayCount++
	}

	return calendar, nil
}

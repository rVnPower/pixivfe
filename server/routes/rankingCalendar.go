package routes

import (
	"fmt"
	"strconv"
	"time"

	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/server/utils"
)

// DateWrap is a struct that encapsulates date-related information for easier handling in templates.
type DateWrap struct {
	Link         string // URL-friendly date string
	Year         int
	Month        int
	MonthPadded  string // Two-digit representation of the month
	MonthLiteral string // Full name of the month
}

// parseDate converts a time.Time value into a DateWrap struct.
// This function is used to prepare date information for display and navigation.
func parseDate(t time.Time) DateWrap {
	var d DateWrap

	year := t.Year()
	month := t.Month()
	monthPadded := fmt.Sprintf("%02d", month)

	d.Link = fmt.Sprintf("%d-%s-01", year, monthPadded)
	d.Year = year
	d.Month = int(month)
	d.MonthPadded = monthPadded
	d.MonthLiteral = month.String()

	return d
}

// RankingCalendarPicker handles the form submission for selecting a ranking calendar.
// It redirects to the RankingCalendarPage with the appropriate query parameters.
func RankingCalendarPicker(w http.ResponseWriter, r *http.Request) error {
	mode := r.FormValue("mode")
	if mode == "" {
		mode = "daily" // Default to daily mode if not specified
	}
	date := r.FormValue("date")

	return utils.RedirectTo(w, r, "/rankingCalendar", map[string]string{
		"mode": mode,
		"date": date,
	})
}

// RankingCalendarPage generates and renders the ranking calendar page.
// It handles date parsing, retrieves the calendar data, and prepares the context for template rendering.
func RankingCalendarPage(w http.ResponseWriter, r *http.Request) error {
	mode := GetQueryParam(r, "mode", "daily")
	date := GetQueryParam(r, "date", "")

	var year int
	var month int

	// Parse the date from the query parameter if provided
	if len(date) == 10 {
		var err error
		year, err = strconv.Atoi(date[:4])
		if err != nil {
			return err
		}
		month, err = strconv.Atoi(date[5:7])
		if err != nil {
			return err
		}
	} else {
		// Use current date if no date is provided
		now := time.Now()
		year = now.Year()
		month = int(now.Month())
	}

	// Calculate dates for navigation
	realDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	monthBefore := realDate.AddDate(0, -1, 0)
	monthAfter := realDate.AddDate(0, 1, 0)

	// Retrieve the ranking calendar HTML
	render, err := core.GetRankingCalendar(r, mode, year, month)
	if err != nil {
		return err
	}

	// Prepare and render the template with the calendar data
	return RenderHTML(w, r, Data_rankingCalendar{
		Title:       "Ranking calendar",
		Render:      render,
		Mode:        mode,
		Year:        year,
		MonthBefore: parseDate(monthBefore),
		MonthAfter:  parseDate(monthAfter),
		ThisMonth:   parseDate(realDate),
	})
}

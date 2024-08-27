package routes

import (
	"fmt"
	"strconv"
	"time"

	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/utils"
)

type DateWrap struct {
	Link         string
	Year         int
	Month        int
	MonthPadded  string
	MonthLiteral string
}

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

func RankingCalendarPicker(w http.ResponseWriter, r CompatRequest) error {
	mode := r.FormValue("mode")
	if mode == "" {
		mode = "daily"
	}
	date := r.FormValue("date")

	return utils.RedirectToRoute(w, r,"/rankingCalendar",  map[string]string{
			"mode": mode,
			"date": date,
		
	}, http.StatusFound)
}

func RankingCalendarPage(w http.ResponseWriter, r CompatRequest) error {
	mode := r.Query("mode", "daily")
	date := r.Query("date", "")

	var year int
	var month int

	// If the user supplied a date
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
		now := time.Now()
		year = now.Year()
		month = int(now.Month())
	}

	realDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	monthBefore := realDate.AddDate(0, -1, 0)
	monthAfter := realDate.AddDate(0, 1, 0)

	render, err := core.GetRankingCalendar(r.Request, mode, year, month)
	if err != nil {
		return err
	}

	return Render(w, r, Data_rankingCalendar{Title: "Ranking calendar", Render: render, Mode: mode, Year: year, MonthBefore: parseDate(monthBefore), MonthAfter: parseDate(monthAfter), ThisMonth: parseDate(realDate)})
}

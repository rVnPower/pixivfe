package main

import (
	"io"
	"log"
	"net/http"
	"testing"
)

func setup() {
}

func teardown() {
}

// TestMain can be used for global setup and teardown
func TestMain(m *testing.M) {
	setup()
	_ = m.Run()
	teardown()
}

func getBaseURL() string {
	return "http://0.0.0.0:8282"
	// return "https://pixivfe.exozy.me"
}

func generateRequest(link, method string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, link, body)
	if err != nil {
		log.Fatalf("Failed to generate a request: %s", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:122.0) Gecko/20100101 Firefox/122.0")

	return req
}

func executeRequest(req *http.Request) *http.Response {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to execute a request: %s", err)
	}

	return resp
}

func TestBasicAllRoutes(t *testing.T) {
	testPositiveURLs := []string{
		// Discovery pages
		"/discovery",
		"/discovery?mode=r18",
		"/discovery/novel",
		"/discovery/novel?mode=r18",

		// Ranking pages
		"/ranking",
		"/ranking?content=all&date=20230212&page=1&mode=male",
		"/ranking?content=manga&page=2&mode=weekly_r18",
		"/ranking?content=ugoira&mode=daily_r18",
		"/rankingCalendar?mode=daily_r18&date=2018-08-01",

		// Artwork page
		"/artworks/121247335",
		"/artworks/120131626", // NSFW
		"/artworks-multi/121289276,121247331,121200724",

		// Users page
		"/users/810305",
		"/users/810305/novels",
		"/users/810305/bookmarks",
	}

	for _, path := range testPositiveURLs {
		URL := getBaseURL() + path
		t.Logf("GETting: %s", URL)
		req := generateRequest(URL, "GET", nil)
		resp := executeRequest(req)

		if resp.StatusCode != 200 {
			t.Errorf("Request route response NOT OK: %d", resp.StatusCode)
		}
	}
}

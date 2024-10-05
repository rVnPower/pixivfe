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

type HTTPTestCase struct {
	URL     string
	Method  string
	Queries map[string]string
	Expect  string
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
	testCases := []HTTPTestCase{
		{
			URL:    "/newest",
			Method: "GET",
		},
		{
			URL:    "/discovery",
			Method: "GET",
		},
		{
			URL:    "/discovery?mode=r18",
			Method: "GET",
		},
		{
			URL:    "/discovery/novel",
			Method: "GET",
		},
		{
			URL:    "/discovery/novel?mode=r18",
			Method: "GET",
		},

		// Ranking pages
		{
			URL:    "/ranking",
			Method: "GET",
		},
		{
			URL:    "/ranking?content=all&date=20230212&page=1&mode=male",
			Method: "GET",
		},
		{
			URL:    "/ranking?content=manga&page=2&mode=weekly_r18",
			Method: "GET",
		},
		{
			URL:    "/ranking?content=ugoira&mode=daily_r18",
			Method: "GET",
		},
		{
			URL:    "/rankingCalendar?mode=daily_r18&date=2018-08-01",
			Method: "GET",
		},

		// Artwork page
		{
			URL:    "/artworks/121247335",
			Method: "GET",
		},
		{
			URL:    "/artworks/120131626",
			Method: "GET",
		},
		{
			URL:    "/artworks-multi/121289276,121247331,121200724",
			Method: "GET",
		},
		// User page
		{
			URL:    "/users/810305",
			Method: "GET",
		},
		{
			URL:    "/users/810305.atom.xml",
			Method: "GET",
		},
		{
			URL:    "/users/810305/manga.atom.xml",
			Method: "GET",
		},
		{
			URL:    "/users/810305/novels",
			Method: "GET",
		},
		{
			URL:    "/users/810305/bookmarks",
			Method: "GET",
		},
		// Pixivision page
		{
			URL:    "/pixivision/",
			Method: "GET",
		},
		{
			URL:    "/pixivision/a/10128",
			Method: "GET",
		},
		{
			URL:    "/pixivision/t/27",
			Method: "GET",
		},
		{
			URL:    "/pixivision/c/manga",
			Method: "GET",
		},

		// Tag page
		{
			URL:    "/tags/original",
			Method: "GET",
		},
		{
			URL:    "/tags?category=manga&ecd=&hgt=1000&hlt=&mode=r18&name=original&order=date&page=1&ratio=0&scd=&tool=&wgt=&wlt=",
			Method: "GET",
		},
	}

	for _, testCase := range testCases {
		URL := getBaseURL() + testCase.URL
		t.Logf("%s: %s", testCase.Method, testCase.URL)
		req := generateRequest(URL, testCase.Method, nil)
		resp := executeRequest(req)

		if resp.StatusCode != 200 {
			t.Errorf("Request route response NOT OK: %d", resp.StatusCode)
		}
	}
}

package main_test

import (
	"log"
	"testing"

	"github.com/playwright-community/playwright-go"
)

var pw *playwright.Playwright
var browser playwright.Browser

func setup() {
	var err error

	runOption := &playwright.RunOptions{
		SkipInstallBrowsers: true,
	}
	if err = playwright.Install(runOption); err != nil {
		log.Fatalf("could not install playwright dependencies. %s", err)
	}

	pw, err = playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright")
	}

	option := playwright.BrowserTypeLaunchOptions{
		Channel: playwright.String("firefox"),
	}

	browser, err = pw.Firefox.Launch(option)
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}

	log.Print("Setup is complete")
}

func teardown() {
	if err := browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}

	if err := pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
	log.Print("Teardown is complete")
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

func checkIfPageHasError(page playwright.Page) bool {
	hasErr, _ := page.Locator(".error").IsVisible()

	return hasErr
}

func TestBasicGetHomepage(t *testing.T) {
	page, err := browser.NewPage()
	if err != nil {
		t.Errorf("could not create page: %v", err)
	}
	if _, err = page.Goto(getBaseURL() + "/"); err != nil {
		t.Errorf("could not goto: %v", err)
	}
	artworks, err := page.Locator(".artwork-small").All()
	if err != nil {
		t.Errorf("could not get entries: %v", err)
	}

	if len(artworks) != 50 {
		t.Errorf("number of daily ranking artworks is %d. expected: 50", len(artworks))
	}
}

func TestBasicAllRoutes(t *testing.T) {
	page, err := browser.NewPage()
	if err != nil {
		t.Errorf("could not create page: %v", err)
	}

	testPositiveURLs := []string{
		// Discovery pages
		"/discovery",
		"/discovery?mode=r18",
		"/discovery/novel",
		"/discovery/novel?mode=r18",

		// Ranking pages
		"/ranking",
		"/ranking?content=all&date=20230212&page=1&mode=male",
		"/ranking?content=mangas&date=20221022&page=3&mode=weekly_r18",
		"/ranking?content=ugoira&date=&page=1&mode=daily_r18",
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

	for _, url := range testPositiveURLs {
		if _, err = page.Goto(getBaseURL() + url); err != nil {
			t.Errorf("could not goto: %v", err)
		}

		if checkIfPageHasError(page) {
			log.Fatalf("Path's response NOT OK: %s", url)
		}
	}
}

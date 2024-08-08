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

	if err = playwright.Install(); err != nil {
		log.Fatalf("could not install playwright dependencies")
	}

	pw, err = playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright")
	}

	browser, err = pw.Chromium.Launch()
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}

	log.Println("Setup is complete")
}

func teardown() {
	if err := browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}

	if err := pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
	log.Println("Teardown is complete")
}

// TestMain can be used for global setup and teardown
func TestMain(m *testing.M) {
	setup()
	_ = m.Run()
	teardown()
}

func getBaseURL() string {
	// return "http://0.0.0.0:8282"
	return "https://pixivfe.exozy.me"
}

func TestGetHomepage(t *testing.T) {
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

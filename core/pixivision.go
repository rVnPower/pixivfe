package core

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"codeberg.org/vnpower/pixivfe/v2/i18n"
)

// PixivDatetimeLayout defines the date format used by Pixiv
const PixivDatetimeLayout = "2006-01-02"

// re_lang is a regular expression to extract the language code from a URL
var re_lang = regexp.MustCompile(`.*\/\/.*?\/(.*?)\/`)

// generateRequest creates a new HTTP request with appropriate headers for Pixiv
func generateRequest(r *http.Request, link, method string, body io.Reader) *http.Request {
	req, err := http.NewRequestWithContext(r.Context(), method, link, body)
	if err != nil {
		panic(err)
	}

	// Extract language from the URL
	lang := re_lang.FindStringSubmatch(link)[1]

	// Set headers to mimic a browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:122.0) Gecko/20100101 Firefox/122.0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cookie", "user_lang="+lang)
	req.Header.Set("Referer", link)

	return req
}

// executeRequest sends an HTTP request and handles common error cases
func executeRequest(req *http.Request) (*http.Response, error) {
	client := http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		return nil, i18n.Error("Pixivision: Page not found")
	}
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// PixivisionArticle represents a single article on Pixivision
type PixivisionArticle struct {
	ID          string
	Title       string
	Description []string
	Category    string
	Thumbnail   string
	Date        time.Time
	Items       []PixivisionArticleItem
	Tags        []PixivisionEmbedTag
}

// PixivisionEmbedTag represents a tag associated with a Pixivision article
type PixivisionEmbedTag struct {
	ID   string
	Name string
}

// PixivisionArticleItem represents an item (artwork) within a Pixivision article
type PixivisionArticleItem struct {
	Username string
	UserID   string
	Title    string
	ID       string
	Avatar   string
	Image    string
}

// PixivisionTag represents a tag page on Pixivision
type PixivisionTag struct {
	Title       string
	Description string
	Articles    []PixivisionArticle
	Total       int // The total number of articles
}

// PixivisionCategory represents a category page on Pixivision
type PixivisionCategory struct {
	Articles    []PixivisionArticle
	Title       string
	Description string
}

// generatePixivisionURL creates a URL for Pixivision based on the route and language
func generatePixivisionURL(route string, lang []string) string {
	template := "https://www.pixivision.net/%s/%s"
	language := "en" // Default
	availableLangs := []string{"en", "zh", "ja", "zh-tw", "ko"}

	// Validate and set the language if provided
	if len(lang) > 0 {
		t := lang[0]

		for _, i := range availableLangs {
			if t == i {
				language = t
			}
		}
	}

	return fmt.Sprintf(template, language, route)
}

// re_findid is a regular expression to extract the ID from a Pixiv link
var re_findid = regexp.MustCompile(`.*\/(\d+)`)

// parseIDFromPixivLink extracts the numeric ID from a Pixiv URL
func parseIDFromPixivLink(link string) string {
	return re_findid.FindStringSubmatch(link)[1]
}

// r_img is a regular expression to extract the image URL from a CSS background-image property
var r_img = regexp.MustCompile(`.*\((.*)\)`)

// parseBackgroundImage extracts the image URL from a CSS background-image property
func parseBackgroundImage(link string) string {
	return r_img.FindStringSubmatch(link)[1]
}

// PixivisionGetHomepage fetches and parses the Pixivision homepage
func PixivisionGetHomepage(r *http.Request, page string, lang ...string) ([]PixivisionArticle, error) {
	var articles []PixivisionArticle

	URL := generatePixivisionURL(fmt.Sprintf("?p=%s", page), lang)
	req := generateRequest(r, URL, "GET", nil)
	resp, err := executeRequest(req)
	if err != nil {
		return articles, err
	}

	if resp.StatusCode == 404 {
		return articles, i18n.Error("We couldn't find the page you're looking for")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse each article on the homepage
	doc.Find("article.spotlight").Each(func(i int, s *goquery.Selection) {
		var article PixivisionArticle

		date := s.Find("time._date").AttrOr("datetime", "")

		article.ID = s.Find(`a[data-gtm-action=ClickTitle]`).AttrOr("data-gtm-label", "")
		article.Title = s.Find(`a[data-gtm-action=ClickTitle]`).Text()
		article.Category = s.Find("._category-label").Text()
		article.Thumbnail = parseBackgroundImage(s.Find("._thumbnail").AttrOr("style", ""))
		article.Date, _ = time.Parse(PixivDatetimeLayout, date)

		// Parse tags associated with the article
		s.Find("._tag-list a").Each(func(i int, t *goquery.Selection) {
			var tag PixivisionEmbedTag
			tag.ID = parseIDFromPixivLink(t.AttrOr("href", ""))
			tag.Name = t.AttrOr("data-gtm-label", "")

			article.Tags = append(article.Tags, tag)
		})

		articles = append(articles, article)
	})

	return articles, nil
}

// PixivisionGetTag fetches and parses a tag page on Pixivision
func PixivisionGetTag(r *http.Request, id string, page string, lang ...string) (PixivisionTag, error) {
	var tag PixivisionTag

	URL := generatePixivisionURL(fmt.Sprintf("t/%s/?p=%s", id, page), lang)
	req := generateRequest(r, URL, "GET", nil)
	resp, err := executeRequest(req)
	if err != nil {
		return tag, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return tag, err
	}

	tag.Title = doc.Find(".tdc__header h1").Text()

	// Extract and process the description
	fullDescription := doc.Find(".tdc__description").Text()
	parts := strings.Split(fullDescription, "pixivision") // split once the boilerplate about "pixivision currently has ..." starts
	tag.Description = strings.TrimSpace(parts[0])

	// Extract total number of articles if available
	if len(parts) > 1 {
		re := regexp.MustCompile(`(\d+)\s+article\(s\)`)
		matches := re.FindStringSubmatch(parts[1])
		if len(matches) > 1 {
			tag.Total, _ = strconv.Atoi(matches[1])
		}
	}

	// Parse each article in the tag page
	doc.Find("._article-card").Each(func(i int, s *goquery.Selection) {
		var article PixivisionArticle

		// article.ID = s.Find(".arc__title a").AttrOr("data-gtm-label", "")
		// article.Title = s.Find(".arc__title a").Text()

		article.ID = s.Find(`a[data-gtm-action="ClickTitle"]`).AttrOr("data-gtm-label", "")
		article.Title = s.Find(`a[data-gtm-action="ClickTitle"]`).Text()
		article.Category = s.Find(".arc__thumbnail-label").Text()
		article.Thumbnail = parseBackgroundImage(s.Find("._thumbnail").AttrOr("style", ""))

		date := s.Find("time._date").AttrOr("datetime", "")
		article.Date, _ = time.Parse(PixivDatetimeLayout, date)

		// Parse tags associated with the article
		s.Find("._tag-list a").Each(func(i int, t *goquery.Selection) {
			var tag PixivisionEmbedTag
			tag.ID = parseIDFromPixivLink(t.AttrOr("href", ""))
			tag.Name = t.AttrOr("data-gtm-label", "")

			article.Tags = append(article.Tags, tag)
		})

		tag.Articles = append(tag.Articles, article)
	})

	return tag, nil
}

// PixivisionGetArticle fetches and parses a single article on Pixivision
func PixivisionGetArticle(r *http.Request, id string, lang ...string) (PixivisionArticle, error) {
	var article PixivisionArticle

	URL := generatePixivisionURL(fmt.Sprintf("a/%s", id), lang)
	req := generateRequest(r, URL, "GET", nil)
	resp, err := executeRequest(req)
	if err != nil {
		return article, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return article, err
	}

	// Parse article metadata
	date := doc.Find("time._date").AttrOr("datetime", "")
	article.Title = doc.Find("h1.am__title").Text()
	article.Category = doc.Find(".am__categoty-pr ._category-label").Text()
	article.Thumbnail = doc.Find(".aie__image").AttrOr("src", "")
	article.Date, _ = time.Parse(PixivDatetimeLayout, date)

	// Parse article description
	doc.Find(".fab__paragraph p").Each(func(i int, s *goquery.Selection) {
		desc, err := s.Html()
		if err != nil {
			return
		}

		if desc == "<br/>" {
			return
		}

		desc = html.EscapeString(desc)

		article.Description = append(article.Description, desc)
	})

	// Parse artworks featured in the article
	doc.Find("._feature-article-body__pixiv_illust").Each(func(i int, s *goquery.Selection) {
		var item PixivisionArticleItem

		item.Title = s.Find(".am__work__title a.inner-link").Text()
		item.Username = s.Find(".am__work__user-name a.inner-link").Text()
		item.ID = parseIDFromPixivLink(s.Find(".am__work__title a.inner-link").AttrOr("href", ""))
		item.UserID = parseIDFromPixivLink(s.Find(".am__work__user-name a.inner-link").AttrOr("href", ""))
		item.Avatar = s.Find(".am__work__user-icon-container img.am__work__uesr-icon").AttrOr("src", "")
		item.Image = s.Find("img.am__work__illust").AttrOr("src", "")

		article.Items = append(article.Items, item)
	})

	// Parse tags associated with the article
	doc.Find("._tag-list a").Each(func(i int, s *goquery.Selection) {
		var tag PixivisionEmbedTag
		tag.ID = parseIDFromPixivLink(s.AttrOr("href", ""))
		tag.Name = s.AttrOr("data-gtm-label", "")

		article.Tags = append(article.Tags, tag)
	})

	return article, nil
}

// PixivisionGetCategory fetches and parses a category page on Pixivision
func PixivisionGetCategory(r *http.Request, id string, page string, lang ...string) (PixivisionCategory, error) {
	var category PixivisionCategory

	URL := generatePixivisionURL(fmt.Sprintf("c/%s/?p=%s", id, page), lang)
	req := generateRequest(r, URL, "GET", nil)
	resp, err := executeRequest(req)
	if err != nil {
		return category, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return category, err
	}

	category.Title = doc.Find(".ssc__name").Text()
	category.Description = doc.Find(".ssc__descriotion").Text() // NOTE: This is a typo in the original HTML

	// Parse each article in the category page
	doc.Find("._article-card").Each(func(i int, s *goquery.Selection) {
		var article PixivisionArticle

		// article.ID = s.Find(".arc__title a").AttrOr("data-gtm-label", "")
		// article.Title = s.Find(".arc__title a").Text()

		article.ID = s.Find(`a[data-gtm-action="ClickTitle"]`).AttrOr("data-gtm-label", "")
		article.Title = s.Find(`a[data-gtm-action="ClickTitle"]`).Text()
		article.Category = s.Find(".arc__thumbnail-label").Text()
		article.Thumbnail = parseBackgroundImage(s.Find("._thumbnail").AttrOr("style", ""))

		date := s.Find("time._date").AttrOr("datetime", "")
		article.Date, _ = time.Parse(PixivDatetimeLayout, date)

		category.Articles = append(category.Articles, article)
	})

	return category, nil
}

package core

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"

	"codeberg.org/vnpower/pixivfe/v2/i18n"
)

const PixivDatetimeLayout = "2006-01-02"

var re_lang = regexp.MustCompile(`.*\/\/.*?\/(.*?)\/`)

func generateRequest(r *http.Request, link, method string, body io.Reader) *http.Request {
	req, err := http.NewRequestWithContext(r.Context(), method, link, body)
	if err != nil {
		panic(err)
	}

	lang := re_lang.FindStringSubmatch(link)[1]

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:122.0) Gecko/20100101 Firefox/122.0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cookie", "user_lang="+lang)
	req.Header.Set("Referer", link)

	return req
}

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

type PixivisionEmbedTag struct {
	ID   string
	Name string
}

type PixivisionArticleItem struct {
	Username string
	UserID   string
	Title    string
	ID       string
	Avatar   string
	Image    string
}

type PixivisionTag struct {
	Title       string
	Description string
	Articles    []PixivisionArticle
}

type PixivisionCategory struct {
	Articles    []PixivisionArticle
	Title       string
	Description string
}

func generatePixivisionURL(route string, lang []string) string {
	template := "https://www.pixivision.net/%s/%s"
	language := "en" // Default
	availableLangs := []string{"en", "zh", "ja", "zh-tw", "ko"}

	// Validation
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

var re_findid = regexp.MustCompile(`.*\/(\d+)`)

func parseIDFromPixivLink(link string) string {
	return re_findid.FindStringSubmatch(link)[1]
}

var r_img = regexp.MustCompile(`.*\((.*)\)`)

func parseBackgroundImage(link string) string {
	return r_img.FindStringSubmatch(link)[1]
}

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

	// TODO: Re-use this function
	doc.Find("article.spotlight").Each(func(i int, s *goquery.Selection) {
		var article PixivisionArticle

		date := s.Find("time._date").AttrOr("datetime", "")

		article.ID = s.Find(`a[data-gtm-action=ClickTitle]`).AttrOr("data-gtm-label", "")
		article.Title = s.Find(`a[data-gtm-action=ClickTitle]`).Text()
		article.Category = s.Find("._category-label").Text()
		article.Thumbnail = parseBackgroundImage(s.Find("._thumbnail").AttrOr("style", ""))
		article.Date, _ = time.Parse(PixivDatetimeLayout, date)

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
	tag.Description = doc.Find(".tdc__description").Text()

	doc.Find("._article-card").Each(func(i int, s *goquery.Selection) {
		var article PixivisionArticle

		//article.ID = s.Find(".arc__title a").AttrOr("data-gtm-label", "")
		//article.Title = s.Find(".arc__title a").Text()

		article.ID = s.Find(`a[data-gtm-action="ClickTitle"]`).AttrOr("data-gtm-label", "")
		article.Title = s.Find(`a[data-gtm-action="ClickTitle"]`).Text()
		article.Category = s.Find(".arc__thumbnail-label").Text()
		article.Thumbnail = parseBackgroundImage(s.Find("._thumbnail").AttrOr("style", ""))

		date := s.Find("time._date").AttrOr("datetime", "")
		article.Date, _ = time.Parse(PixivDatetimeLayout, date)

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

	// Title
	date := doc.Find("time._date").AttrOr("datetime", "")
	article.Title = doc.Find("h1.am__title").Text()
	article.Category = doc.Find(".am__categoty-pr ._category-label").Text()
	article.Thumbnail = doc.Find(".aie__image").AttrOr("src", "")
	article.Date, _ = time.Parse(PixivDatetimeLayout, date)

	// Description
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

	doc.Find("._tag-list a").Each(func(i int, s *goquery.Selection) {
		var tag PixivisionEmbedTag
		tag.ID = parseIDFromPixivLink(s.AttrOr("href", ""))
		tag.Name = s.AttrOr("data-gtm-label", "")

		article.Tags = append(article.Tags, tag)
	})

	return article, nil
}

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
	category.Description = doc.Find(".ssc__descriotion").Text() // Not a typo

	doc.Find("._article-card").Each(func(i int, s *goquery.Selection) {
		var article PixivisionArticle

		//article.ID = s.Find(".arc__title a").AttrOr("data-gtm-label", "")
		//article.Title = s.Find(".arc__title a").Text()

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

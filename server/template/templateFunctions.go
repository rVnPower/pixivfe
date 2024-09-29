package template

import (
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"regexp"
	"strings"
	"time"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/i18n"
)

type HTML = core.HTML

func GetRandomColor() string {
	// Some color shade I stole
	colors := []string{
		// Green
		"#C8847E",
		"#C8A87E",
		"#C8B87E",
		"#C8C67E",
		"#C7C87E",
		"#C2C87E",
		"#BDC87E",
		"#82C87E",
		"#82C87E",
		"#7EC8AF",
		"#7EAEC8",
		"#7EA6C8",
		"#7E99C8",
		"#7E87C8",
		"#897EC8",
		"#967EC8",
		"#AE7EC8",
		"#B57EC8",
		"#C87EA5",
	}

	// Randomly choose one and return
	return colors[rand.Intn(len(colors))]
}

var re_emoji = regexp.MustCompile(`\(([^)]+)\)`)

func ParseEmojis(s string) HTML {
	emojiList := map[string]string{
		"normal":        "101",
		"surprise":      "102",
		"serious":       "103",
		"heaven":        "104",
		"happy":         "105",
		"excited":       "106",
		"sing":          "107",
		"cry":           "108",
		"normal2":       "201",
		"shame2":        "202",
		"love2":         "203",
		"interesting2":  "204",
		"blush2":        "205",
		"fire2":         "206",
		"angry2":        "207",
		"shine2":        "208",
		"panic2":        "209",
		"normal3":       "301",
		"satisfaction3": "302",
		"surprise3":     "303",
		"smile3":        "304",
		"shock3":        "305",
		"gaze3":         "306",
		"wink3":         "307",
		"happy3":        "308",
		"excited3":      "309",
		"love3":         "310",
		"normal4":       "401",
		"surprise4":     "402",
		"serious4":      "403",
		"love4":         "404",
		"shine4":        "405",
		"sweat4":        "406",
		"shame4":        "407",
		"sleep4":        "408",
		"heart":         "501",
		"teardrop":      "502",
		"star":          "503",
	}

	parsedString := re_emoji.ReplaceAllStringFunc(s, func(s string) string {
		s = s[1 : len(s)-1] // Get the string inside
		id := emojiList[s]

		return fmt.Sprintf(`<img src="/proxy/s.pximg.net/common/images/emoji/%s.png" alt="(%s)" class="emoji" />`, id, s)
	})
	return HTML(parsedString)
}

type PageInfo struct {
    Number int
    URL    string
}

type PaginationData struct {
    CurrentPage int
    MaxPage     int
    Pages       []PageInfo
    HasPrevious bool
    HasNext     bool
    PreviousURL string
    NextURL     string
    FirstURL    string
    LastURL     string
    HasMaxPage  bool
    LastPage    int
}

func ParsePixivRedirect(s string) HTML {
	regex := regexp.MustCompile(`\/jump\.php\?(http[^"]+)`)

func ParsePixivRedirect(s string) HTML {
	parsedString := re_jump.ReplaceAllStringFunc(s, func(s string) string {
		s = s[10:]
		return s
	})
	escaped, err := url.QueryUnescape(parsedString)
	if err != nil {
		return HTML(s)
	}
	return HTML(escaped)
}

func EscapeString(s string) string {
	escaped := url.QueryEscape(s)
	return escaped
}

func ParseTime(date time.Time) string {
	return date.Format("2006-01-02 15:04")
}

func ParseTimeCustomFormat(date time.Time, format string) string {
	return date.Format(format)
}
func CreatePaginator(base, ending string, current_page, max_page int) PaginationData {
    pageUrl := func(page int) string {
        return fmt.Sprintf(`%s%d%s`, base, page, ending)
    }

    const peek = 2 // Number of pages to show on each side of the current page
    hasMaxPage := max_page != -1

    start := max(1, current_page-peek)
    end := current_page + peek
    if hasMaxPage {
        end = min(max_page, end)
    }

    pages := make([]PageInfo, 0, end-start+1)
    for i := start; i <= end; i++ {
        pages = append(pages, PageInfo{Number: i, URL: pageUrl(i)})
    }

    lastPage := pages[len(pages)-1].Number

    return PaginationData{
        CurrentPage: current_page,
        MaxPage:     max_page,
        Pages:       pages,
        HasPrevious: current_page > 1,
        HasNext:     !hasMaxPage || current_page < max_page,
        PreviousURL: pageUrl(current_page - 1),
        NextURL:     pageUrl(current_page + 1),
        FirstURL:    pageUrl(1),
        LastURL:     pageUrl(max_page),
        HasMaxPage:  hasMaxPage,
        LastPage:    lastPage,
    }
}

func GetNovelGenre(s string) string {
	switch s {
	case "1":
		return i18n.Tr("Romance")
	case "2":
		return i18n.Tr("Isekai fantasy")
	case "3":
		return i18n.Tr("Contemporary fantasy")
	case "4":
		return i18n.Tr("Mystery")
	case "5":
		return i18n.Tr("Horror")
	case "6":
		return i18n.Tr("Sci-fi")
	case "7":
		return i18n.Tr("Literature")
	case "8":
		return i18n.Tr("Drama")
	case "9":
		return i18n.Tr("Historical pieces")
	case "10":
		return i18n.Tr("BL (yaoi)")
	case "11":
		return i18n.Tr("Yuri")
	case "12":
		return i18n.Tr("For kids")
	case "13":
		return i18n.Tr("Poetry")
	case "14":
		return i18n.Tr("Essays/non-fiction")
	case "15":
		return i18n.Tr("Screenplays/scripts")
	case "16":
		return i18n.Tr("Reviews/opinion pieces")
	case "17":
		return i18n.Tr("Other")
	}

	return i18n.Sprintf("(Unknown Genre: %s)", s)
}

func SwitchButtonAttributes(baseURL, selection, currentSelection string) string {
	var cur string = "false"
	if selection == currentSelection {
		cur = "true"
	}

	return fmt.Sprintf(`href=%s%s class=switch-button selected=%s`, baseURL, selection, cur)
}

var furiganaPattern = regexp.MustCompile(`\[\[rb:\s*(.+?)\s*>\s*(.+?)\s*\]\]`)
var chapterPattern = regexp.MustCompile(`\[chapter:\s*(.+?)\s*\]`)
var jumpUriPattern = regexp.MustCompile(`\[\[jumpuri:\s*(.+?)\s*>\s*(.+?)\s*\]\]`)
var jumpPagePattern = regexp.MustCompile(`\[jump:\s*(\d+?)\s*\]`)
var newPagePattern = regexp.MustCompile(`\s*\[newpage\]\s*`)

func GetTemplateFunctions() map[string]any {
	return map[string]any{
		"parseEmojis": func(s string) HTML {
			return ParseEmojis(s)
		},

		"parsePixivRedirect": func(s string) HTML {
			return ParsePixivRedirect(s)
		},
		"escapeString": func(s string) string {
			return EscapeString(s)
		},

		"randomColor": func() string {
			return GetRandomColor()
		},

		"isEmpty": func(s string) bool {
			return len(s) < 1
		},

		"isEmphasize": func(s string) bool {
			switch s {
			case
				"R-18",
				"R-18G":
				return true
			}
			return false
		},
		"reformatDate": func(s string) string {
			if len(s) != 8 {
				return s
			}
			return fmt.Sprintf("%s-%s-%s", s[4:], s[2:4], s[:2])
		},
		"parseTime": func(date time.Time) string {
			return ParseTime(date)
		},
		"parseTimeCustomFormat": func(date time.Time, format string) string {
			return ParseTimeCustomFormat(date, format)
		},
		"createPaginator": func(base, ending string, current_page, max_page int) PaginationData {
			return CreatePaginator(base, ending, current_page, max_page)
		},
		"joinArtworkIds": func(artworks []core.ArtworkBrief) string {
			ids := []string{}
			for _, art := range artworks {
				ids = append(ids, art.ID)
			}
			return strings.Join(ids, ",")
		},
		"stripEmbed": func(s string) string {
			// this is stupid
			return s[:len(s)-6]
		},
		"renderNovel": func(s string) HTML {
			furiganaTemplate := `<ruby>$1<rp>(</rp><rt>$2</rt><rp>)</rp></ruby>`
			s = furiganaPattern.ReplaceAllString(s, furiganaTemplate)

			chapterTemplate := `<h2>$1</h2>`
			s = chapterPattern.ReplaceAllString(s, chapterTemplate)

			jumpUriTemplate := `<a href="$2" target="_blank">$1</a>`
			s = jumpUriPattern.ReplaceAllString(s, jumpUriTemplate)

			jumpPageTemplate := `<a href="#$1">To page $1</a>`
			s = jumpPagePattern.ReplaceAllString(s, jumpPageTemplate)

			if strings.Contains(s, "[newpage]") {
				// if [newpage] in content , then prepend <hr id="1"/> to the page
				s = `<hr id="1"/>` + s
				pageIdx := 1

				// Should run before replace `\n` -> `<br />`
				s = newPagePattern.ReplaceAllStringFunc(s, func(_ string) string {
					pageIdx += 1
					return fmt.Sprintf(`<br /><hr id="%d"/>`, pageIdx)
				})
			}
			s = strings.ReplaceAll(s, "\n", "<br />")
			return HTML(s)
		},
		"novelGenre": GetNovelGenre,
		"floor": func(i float64) int {
			return int(math.Floor(i))
		},
		"unfinishedQuery": UnfinishedQuery,
		"replaceQuery":    ReplaceQuery,
		// "AttrGen": SwitchButtonAttributes,
	}
}

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

// HTML is an alias for the HTML type from the core package
type HTML = core.HTML

// GetRandomColor returns a randomly selected color from a predefined list of color shades
func GetRandomColor() string {
	// VnPower: Some color shade I stole
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

	// Randomly choose one color and return it
	return colors[rand.Intn(len(colors))]
}

// ParseEmojis replaces emoji shortcodes in a string with corresponding image tags
func ParseEmojis(s string) HTML {
	// Map of emoji shortcodes to their corresponding image IDs
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

	// Regular expression to match emoji shortcodes
	regex := regexp.MustCompile(`\(([^)]+)\)`)

	// Replace shortcodes with corresponding image tags
	parsedString := regex.ReplaceAllStringFunc(s, func(s string) string {
		s = s[1 : len(s)-1]       // Get the string inside parentheses
		id, found := emojiList[s] // Check if the string has a corresponding emoji ID
		if !found {
			return fmt.Sprintf("(%s)", s) // No replacement, return the original text
		}

		return fmt.Sprintf(`<img src="/proxy/s.pximg.net/common/images/emoji/%s.png" alt="(%s)" class="emoji" />`, id, s)
	})
	return HTML(parsedString)
}

// PageInfo represents information about a single page in pagination
type PageInfo struct {
	Number int
	URL    string
}

// PaginationData contains all necessary information for rendering pagination controls
type PaginationData struct {
	CurrentPage   int
	MaxPage       int
	Pages         []PageInfo
	HasPrevious   bool
	HasNext       bool
	PreviousURL   string
	NextURL       string
	FirstURL      string
	LastURL       string
	HasMaxPage    bool
	LastPage      int
	DropdownPages []PageInfo
}

// ParsePixivRedirect extracts and unescapes URLs from Pixiv's redirect links
func ParsePixivRedirect(s string) HTML {
	// Regular expression to match Pixiv's redirect URLs
	regex := regexp.MustCompile(`\/jump\.php\?(http[^"]+)`)

	// Extract the actual URL from the redirect link
	parsedString := regex.ReplaceAllStringFunc(s, func(s string) string {
		s = s[10:]
		return s
	})

	// Unescape the URL
	escaped, err := url.QueryUnescape(parsedString)
	if err != nil {
		return HTML(s)
	}
	return HTML(escaped)
}

// EscapeString escapes a string for use in a URL query
func EscapeString(s string) string {
	escaped := url.QueryEscape(s)
	return escaped
}

// ParseTime formats a time.Time value as a string in the format "2006-01-02 15:04"
func ParseTime(date time.Time) string {
	return date.Format("2006-01-02 15:04")
}

// ParseTimeCustomFormat formats a time.Time value as a string using a custom format
func ParseTimeCustomFormat(date time.Time, format string) string {
	return date.Format(format)
}

// CreatePaginator generates pagination data based on the current page and maximum number of pages.
// It returns a PaginationData struct containing all necessary information for rendering pagination controls.
//
// Parameters:
// - base: The base part of the pagination URL
// - ending: The ending part of the pagination URL
// - current_page: The current page being displayed
// - max_page: Maximum number of pages (-1 if unknown)
// - page_margin: Number of pages to display before and after the current page
// - dropdown_offset: Number of pages to include in dropdown before and after current page
func CreatePaginator(base, ending string, current_page, max_page, page_margin, dropdown_offset int) (PaginationData, error) {
	// Validate input parameters
	if current_page < 1 {
		return PaginationData{}, fmt.Errorf("current_page must be greater than or equal to 1, got %d", current_page)
	}
	if page_margin < 0 {
		return PaginationData{}, fmt.Errorf("page_margin must be non-negative, got %d", page_margin)
	}
	if dropdown_offset < 0 {
		return PaginationData{}, fmt.Errorf("dropdown_offset must be non-negative, got %d", dropdown_offset)
	}

	// Validation for users that don't have any artworks
	// NOTE: the following breaks the current max_page implementation, commenting it out for now
	// if max_page < 1 {
	// 	max_page = 1
	// }

	// Validation for max_page in relation to current_page
	// if max_page < current_page {
	// 	return PaginationData{}, fmt.Errorf("max_page (%d) must be greater than or equal to current_page (%d) when specified", max_page, current_page)
	// }

	hasMaxPage := max_page != -1

	pageUrl := func(page int) string {
		return fmt.Sprintf(`%s%d%s`, base, page, ending)
	}

	// Helper function to generate a range of pages
	generatePageRange := func(start, end int) []PageInfo {
		if start > end {
			return []PageInfo{}
		}
		pages := make([]PageInfo, 0, end-start+1)
		for i := start; i <= end; i++ {
			pages = append(pages, PageInfo{Number: i, URL: pageUrl(i)})
		}
		return pages
	}

	// Calculate the range of pages to display
	start := max(1, current_page-page_margin)
	end := current_page + page_margin
	if hasMaxPage {
		end = min(max_page, end)
	}

	// Generate page information for the range
	pages := generatePageRange(start, end)

	var lastPage int
	if len(pages) > 0 {
		lastPage = pages[len(pages)-1].Number
	} else {
		lastPage = current_page
	}

	// Generate dropdown pages
	dropdownStart := max(1, current_page-dropdown_offset)
	dropdownEnd := current_page + dropdown_offset
	if hasMaxPage {
		dropdownEnd = min(max_page, dropdownEnd)
	}
	dropdownPages := generatePageRange(dropdownStart, dropdownEnd)

	// Calculate previous and next URLs
	var previousURL, nextURL string
	if current_page > 1 {
		previousURL = pageUrl(current_page - 1)
	}
	if !hasMaxPage || current_page < max_page {
		nextURL = pageUrl(current_page + 1)
	}

	// Create and return the PaginationData struct
	return PaginationData{
		CurrentPage:   current_page,
		MaxPage:       max_page,
		Pages:         pages,
		HasPrevious:   current_page > 1,
		HasNext:       !hasMaxPage || current_page < max_page,
		PreviousURL:   previousURL,
		NextURL:       nextURL,
		FirstURL:      pageUrl(1),
		LastURL:       pageUrl(max_page),
		HasMaxPage:    hasMaxPage,
		LastPage:      lastPage,
		DropdownPages: dropdownPages,
	}, nil
}

// GetNovelGenre returns the genre name for a given genre ID
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

// SwitchButtonAttributes generates HTML attributes for a switch button based on the current selection
func SwitchButtonAttributes(baseURL, selection, currentSelection string) string {
	var cur string = "false"
	if selection == currentSelection {
		cur = "true"
	}

	return fmt.Sprintf(`href=%s%s class=switch-button selected=%s`, baseURL, selection, cur)
}

var (
	furiganaPattern = regexp.MustCompile(`\[\[rb:\s*(.+?)\s*>\s*(.+?)\s*\]\]`)
	chapterPattern  = regexp.MustCompile(`\[chapter:\s*(.+?)\s*\]`)
	jumpUriPattern  = regexp.MustCompile(`\[\[jumpuri:\s*(.+?)\s*>\s*(.+?)\s*\]\]`)
	jumpPagePattern = regexp.MustCompile(`\[jump:\s*(\d+?)\s*\]`)
	newPagePattern  = regexp.MustCompile(`\s*\[newpage\]\s*`)
)

// GetTemplateFunctions returns a map of custom template functions for use in HTML templates
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
			case "R-18", "R-18G":
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
		"createPaginator": func(base, ending string, current_page, max_page, page_margin, dropdown_offset int) PaginationData {
			paginationData, err := CreatePaginator(base, ending, current_page, max_page, page_margin, dropdown_offset)
			if err != nil {
				fmt.Printf("Error creating paginator: %v", err)
				// Return an empty PaginationData in case of error
				return PaginationData{}
			}
			return paginationData
		},
		"joinArtworkIds": func(artworks []core.ArtworkBrief) string {
			ids := []string{}
			for _, art := range artworks {
				ids = append(ids, art.ID)
			}
			return strings.Join(ids, ",")
		},
		"stripEmbed": func(s string) string {
			// Remove the last 6 characters from the string (assumes "_embed" suffix)
			return s[:len(s)-6]
		},
		"renderNovel": func(s string) HTML {
			// Replace furigana markup with HTML ruby tags
			furiganaTemplate := `<ruby>$1<rp>(</rp><rt>$2</rt><rp>)</rp></ruby>`
			s = furiganaPattern.ReplaceAllString(s, furiganaTemplate)

			// Replace chapter markup with HTML h2 tags
			chapterTemplate := `<h2>$1</h2>`
			s = chapterPattern.ReplaceAllString(s, chapterTemplate)

			// Replace jump URI markup with HTML anchor tags
			jumpUriTemplate := `<a href="$2" target="_blank">$1</a>`
			s = jumpUriPattern.ReplaceAllString(s, jumpUriTemplate)

			// Replace jump page markup with HTML anchor tags
			jumpPageTemplate := `<a href="#$1">To page $1</a>`
			s = jumpPagePattern.ReplaceAllString(s, jumpPageTemplate)

			// Handle newpage markup
			if strings.Contains(s, "[newpage]") {
				// Prepend <hr id="1"/> to the page if [newpage] is present
				s = `<hr id="1"/>` + s
				pageIdx := 1

				// Should run before replace `\n` -> `<br />`
				s = newPagePattern.ReplaceAllStringFunc(s, func(_ string) string {
					pageIdx += 1
					return fmt.Sprintf(`<br /><hr id="%d"/>`, pageIdx)
				})
			}

			// Replace newlines with HTML line breaks
			s = strings.ReplaceAll(s, "\n", "<br />")
			return HTML(s)
		},
		"novelGenre": GetNovelGenre,
		"floor": func(i float64) int {
			return int(math.Floor(i))
		},
		"unfinishedQuery": UnfinishedQuery,
		"replaceQuery":    ReplaceQuery,
		// TODO: what is AttrGen for
		// "AttrGen": SwitchButtonAttributes,
	}
}

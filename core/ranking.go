package core

import (
	"log"
	"strings"
	"time"

	"codeberg.org/vnpower/pixivfe/v2/server/session"
	"codeberg.org/vnpower/pixivfe/v2/i18n"

	"github.com/goccy/go-json"
	"net/http"
)

/*
TODO: occasionally, Pixiv feels like returning an HTML response instead of the expected JSON,
			causing an unmarshalling error on our end and forcing a refresh by the user

			as a (super secret and hardcoded) workaround, we implement automatic retries to
			eventually get the JSON response we need (usually only takes a single retry)

			from qualitative testing, it seems like the interval between requests doesn't actually matter,
			just the fact that we send a second request is enough
*/

const (
	maxRetries = 5
	retryDelay = 10 * time.Millisecond
)

type Ranking struct {
	Contents []struct {
		Title        string `json:"title"`
		Thumbnail    string `json:"url"`
		Pages        int    `json:"illust_page_count,string"`
		ArtistName   string `json:"user_name"`
		ArtistAvatar string `json:"profile_img"`
		ID           int    `json:"illust_id"`
		ArtistID     int    `json:"user_id"`
		Rank         int    `json:"rank"`
		IllustType   int    `json:"illust_type,string"`
		XRestrict    int    // zero field for template
	} `json:"contents"`

	Mode        string          `json:"mode"`
	Content     string          `json:"content"`
	Page        int             `json:"page"`
	RankTotal   int             `json:"rank_total"`
	CurrentDate string          `json:"date"`
	PrevDateRaw json.RawMessage `json:"prev_date"`
	NextDateRaw json.RawMessage `json:"next_date"`
	PrevDate    string
	NextDate    string
}

func GetRanking(r *http.Request, mode, content, date, page string) (Ranking, error) {
	// log.Printf("GetRanking called with mode: %s, content: %s, date: %s, page: %s", mode, content, date, page)

	URL := GetRankingURL(mode, content, date, page)
	// log.Printf("Ranking URL: %s", URL)

	var ranking Ranking
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		i18n.Sprintf("Attempt %d of %d", attempt+1, maxRetries)

		resp, err := API_GET(r.Context(), URL, "")
		if err != nil {
			i18n.Errorf("API_GET error: %v", err)
			return ranking, err
		}

		proxiedResp := session.ProxyImageUrl(r, resp.Body)

		err = json.Unmarshal([]byte(proxiedResp), &ranking)
		if err == nil {
			i18n.Sprintf("JSON unmarshalling successful on attempt %d", attempt+1)
			lastErr = nil // Reset lastErr on successful unmarshalling
			break
		}

		lastErr = err
		i18n.Errorf("JSON unmarshalling error on attempt %d: %v", attempt+1, err)

		if attempt < maxRetries-1 {
			i18n.Sprintf("Retrying in %v", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	if lastErr != nil {
		// NOTE: This shouldn't happen; indicative of a different issue
		i18n.Errorf("All attempts failed. Last error: %v", lastErr)
		return ranking, lastErr
	}

	ranking.PrevDate = strings.ReplaceAll(string(ranking.PrevDateRaw[:]), "\"", "")
	ranking.NextDate = strings.ReplaceAll(string(ranking.NextDateRaw[:]), "\"", "")

	// log.Printf("GetRanking completed successfully. PrevDate: %s, NextDate: %s", ranking.PrevDate, ranking.NextDate)

	return ranking, nil
}

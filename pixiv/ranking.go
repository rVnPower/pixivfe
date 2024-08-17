package pixiv

import (
	"errors"
	"strings"

	session "codeberg.org/vnpower/pixivfe/v2/session"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

type Ranking struct {
	Contents []struct {
		Title        string `json:"title"`
		Image        string `json:"url"`
		Pages        int    `json:"illust_page_count,string"`
		ArtistName   string `json:"user_name"`
		ArtistAvatar string `json:"profile_img"`
		ID           int    `json:"illust_id"`
		ArtistID     int    `json:"user_id"`
		Rank         int    `json:"rank"`
		IllustType   int    `json:"illust_type,string"`
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

func GetRanking(c *fiber.Ctx, mode, content, date, page string) (Ranking, error) {
	URL := GetRankingURL(mode, content, date, page)

	var ranking Ranking

	resp := WebAPIRequest(c.Context(), URL, "")
	if !resp.Ok {
		return ranking, errors.New(resp.Message)
	}
	proxiedResp := session.ProxyImageUrl(c, resp.Body)

	err := json.Unmarshal([]byte(proxiedResp), &ranking)
	if err != nil {
		return ranking, err
	}

	ranking.PrevDate = strings.ReplaceAll(string(ranking.PrevDateRaw[:]), "\"", "")
	ranking.NextDate = strings.ReplaceAll(string(ranking.NextDateRaw[:]), "\"", "")

	return ranking, nil
}

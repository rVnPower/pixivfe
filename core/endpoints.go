package core

import (
	"fmt"
	"net/url"
)

func GetNewestArtworksURL(worktype, r18, lastID string) string {
	base := "https://www.pixiv.net/ajax/illust/new?limit=30&type=%s&r18=%s&lastId=%s"
	return fmt.Sprintf(base, worktype, r18, lastID)
}

func GetDiscoveryURL(mode string, limit int) string {
	base := "https://www.pixiv.net/ajax/discovery/artworks?mode=%s&limit=%d"
	return fmt.Sprintf(base, mode, limit)
}

func GetDiscoveryNovelURL(mode string, limit int) string {
	base := "https://www.pixiv.net/ajax/discovery/novels?mode=%s&limit=%d"
	return fmt.Sprintf(base, mode, limit)
}

func GetRankingURL(mode, content, date, page string) string {
	base := "https://www.pixiv.net/ranking.php?format=json&mode=%s&content=%s&date=%s&p=%s"
	baseNoDate := "https://www.pixiv.net/ranking.php?format=json&mode=%s&content=%s&p=%s"

	if date != "" {
		return fmt.Sprintf(base, mode, content, date, page)
	}

	return fmt.Sprintf(baseNoDate, mode, content, page)
}

func GetRankingCalendarURL(mode string, year, month int) string {
	base := "https://www.pixiv.net/ranking_log.php?mode=%s&date=%d%02d"

	return fmt.Sprintf(base, mode, year, month)
}

func GetUserInformationURL(id string) string {
	base := "https://www.pixiv.net/ajax/user/%s?full=1"

	return fmt.Sprintf(base, id)
}

func GetUserArtworksURL(id string) string {
	base := "https://www.pixiv.net/ajax/user/%s/profile/all"

	return fmt.Sprintf(base, id)
}

func GetUserFullArtworkURL(id, ids string) string {
	base := "https://www.pixiv.net/ajax/user/%s/profile/illusts?work_category=illustManga&is_first_page=0&lang=en%s"

	return fmt.Sprintf(base, id, ids)
}

func GetUserFullNovelURL(id, ids string) string {
	base := "https://www.pixiv.net/ajax/user/%s/profile/novels?is_first_page=0&lang=en%s"

	return fmt.Sprintf(base, id, ids)
}

func GetUserBookmarksURL(id, mode string, page int) string {
	base := "https://www.pixiv.net/ajax/user/%s/illusts/bookmarks?tag=&offset=%d&limit=48&rest=%s"

	return fmt.Sprintf(base, id, page*48, mode)
}

func GetFrequentArtworkTagsURL(ids string) string {
	base := "https://www.pixiv.net/ajax/tags/frequent/illust?%s"

	return fmt.Sprintf(base, ids)
}

func GetFrequentNovelTagsURL(ids string) string {
	base := "https://www.pixiv.net/ajax/tags/frequent/novel?%s"

	return fmt.Sprintf(base, ids)
}

func GetNewestFromFollowingURL(mode, page string) string {
	base := "https://www.pixiv.net/ajax/follow_latest/%s?mode=%s&p=%s"

	// TODO: Recheck this URL
	return fmt.Sprintf(base, "illust", mode, page)
}

func GetArtworkInformationURL(id string) string {
	base := "https://www.pixiv.net/ajax/illust/%s"

	return fmt.Sprintf(base, id)
}

func GetArtworkImagesURL(id string) string {
	base := "https://www.pixiv.net/ajax/illust/%s/pages"

	return fmt.Sprintf(base, id)
}

func GetArtworkRelatedURL(id string, limit int) string {
	base := "https://www.pixiv.net/ajax/illust/%s/recommend/init?limit=%d"

	return fmt.Sprintf(base, id, limit)
}

func GetArtworkCommentsURL(id string) string {
	base := "https://www.pixiv.net/ajax/illusts/comments/roots?illust_id=%s&limit=100"

	return fmt.Sprintf(base, id)
}

func GetTagDetailURL(unescapedTag string) string {
	base := "https://www.pixiv.net/ajax/search/tags/%s"

	return fmt.Sprintf(base, url.PathEscape(unescapedTag))
}

func GetSearchArtworksURL(s map[string]string) string {
	// Long.
	base := "https://www.pixiv.net/ajax/search/%s/%s?order=%s&mode=%s&ratio=%s&s_mode=%s&wlt=%s&wgt=%s&hlt=%s&hgt=%s&tool=%s&scd=%s&ecd=%s&p=%s"

	return fmt.Sprintf(base, s["Category"], s["Name"], s["Order"], s["Mode"], s["Ratio"], s["Smode"], s["Wlt"], s["Wgt"], s["Hlt"], s["Hgt"], s["Tool"], s["Scd"], s["Ecd"], s["Page"])
}

func GetLandingURL(mode string) string {
	base := "https://www.pixiv.net/ajax/top/illust?mode=%s"

	return fmt.Sprintf(base, mode)
}

func GetNovelURL(id string) string {
	base := "https://www.pixiv.net/ajax/novel/%s"

	return fmt.Sprintf(base, id)
}

func GetNovelRelatedURL(id string, limit int) string {
	base := "https://www.pixiv.net/ajax/novel/%s/recommend/init?limit=%d"

	return fmt.Sprintf(base, id, limit)
}

func GetNovelCommentsURL(id string) string {
	base := "https://www.pixiv.net/ajax/novels/comments/roots?novel_id=%s&offset=0&limit=100"

	return fmt.Sprintf(base, id)
}

func GetNovelSeriesURL(id string) string {
	base := "https://www.pixiv.net/ajax/novel/series/%s"

	return fmt.Sprintf(base, id)
}

func GetNovelSeriesContentURL(id string, page int, perPage int) string {
	base := "https://www.pixiv.net/ajax/novel/series_content/%s?limit=%d&last_order=%d&order_by=asc"

	return fmt.Sprintf(base, id, perPage, perPage*(page-1))
}

func GetNovelSeriesContentTitlesURL(id int) string {
	base := "https://www.pixiv.net/ajax/novel/series/%d/content_titles"

	return fmt.Sprintf(base, id)
}

func GetInsertIllustURL(novelid, id string) string {
	base := "https://www.pixiv.net/ajax/novel/%s/insert_illusts?id[]=%s"

	return fmt.Sprintf(base, novelid, id)
}

func GetMangaSeriesContentURL(id string, page int) string {
	base := "https://www.pixiv.net/ajax/series/%s?p=%d"

	return fmt.Sprintf(base, id, page)
}

package views

import (
	"errors"
	"math"
	"net/http"
	"pixivfe/configs"
	"pixivfe/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func get_session_value(c *fiber.Ctx, key string) *string {
	sess, err := configs.Store.Get(c)
	if err != nil {
		panic(err)
	}
	value := sess.Get(key)
	if value != nil {
		placeholder := value.(string)
		return &placeholder
	}
	return nil
}

func artwork_page(c *fiber.Ctx) error {
	image_proxy := get_session_value(c, "image-proxy")
	if image_proxy == nil {
		image_proxy = &configs.ProxyServer
	}

	id := c.Params("id")
	if _, err := strconv.Atoi(id); err != nil {
		return errors.New("Bad id")
	}

	illust, err := PC.GetArtworkByID(id)
	if err != nil {
		return err
	}

	illust.ProxyImages(*image_proxy)

	related, err := PC.GetRelatedArtworks(id)
	if err != nil {
		return err
	}
	related = models.ProxyShortArtworkSlice(related, *image_proxy)

	comments, _ := PC.GetArtworkComments(id)
	comments = models.ProxyCommentsSlice(comments, *image_proxy)

	// Optimize this
	return c.Render("artwork", fiber.Map{
		"Illust":   illust,
		"Related":  related,
		"Comments": comments,
		"Title":    illust.Title,
	})
}

func index_page(c *fiber.Ctx) error {
	// recommended, _ := handler.GetRecommendedIllust(c)
	// ranking, _ := handler.GetRankingIllust(c, "day")
	// spotlight := handler.GetSpotlightArticle(c)
	// newest, _ := handler.GetNewestIllust(c)
	// return c.Render(http.StatusOK, "index.html", fiber.Map{
	// 	"Recommended": recommended,
	// 	"Rankings":    ranking,
	// 	"Spotlights":  spotlight,
	// 	"Newest":      newest,
	// })
	sess, err := configs.Store.Get(c)
	if err != nil {
		panic(err)
	}
	token := sess.Get("token")

	if token != nil {
		println(token.(string))
	}

	return c.Render("temp", fiber.Map{"Title": "Landing"})
}

func user_page(c *fiber.Ctx) error {
	image_proxy := get_session_value(c, "image-proxy")
	if image_proxy == nil {
		image_proxy = &configs.ProxyServer
	}

	id := c.Params("id")
	if _, err := strconv.Atoi(id); err != nil {
		return err
	}
	category := c.Params("category", "artworks")
	if !(category == "artworks" || category == "illustrations" || category == "manga" || category == "bookmarks") {
		return errors.New("Invalid work category: only illustrations, manga, artworks and bookmarks are available")
	}

	page := c.Query("page", "1")
	pageInt, _ := strconv.Atoi(page)

	user, err := PC.GetUserInformation(id, category, pageInt)
	if err != nil {
		return err
	}

	user.ProxyImages(*image_proxy)

	var worksCount int

	worksCount = user.ArtworksCount
	pageLimit := math.Ceil(float64(worksCount)/30.0) + 1.0

	return c.Render("user", fiber.Map{"Title": user.Name, "User": user, "Category": category, "PageLimit": int(pageLimit), "Page": pageInt})
}

func ranking_page(c *fiber.Ctx) error {
	image_proxy := get_session_value(c, "image-proxy")
	if image_proxy == nil {
		image_proxy = &configs.ProxyServer
	}
	mode := c.Query("mode", "daily")

	content := c.Query("content", "all")

	page := c.Query("page", "1")

	pageInt, _ := strconv.Atoi(page)

	response, err := PC.GetRanking(mode, content, page)
	if err != nil {
		return err
	}

	artworks := response.Artworks

	for i := range artworks {
		artworks[i].Image = models.ProxyImage(artworks[i].Image, *image_proxy)
		artworks[i].ArtistAvatar = models.ProxyImage(artworks[i].ArtistAvatar, *image_proxy)
	}

	return c.Render("rank", fiber.Map{
		"Title":   "Ranking",
		"Items":   artworks,
		"Mode":    mode,
		"Content": content,
		"Page":    pageInt})
}

func newest_artworks_page(c *fiber.Ctx) error {
	image_proxy := get_session_value(c, "image-proxy")
	if image_proxy == nil {
		image_proxy = &configs.ProxyServer
	}

	worktype := c.Query("type", "illust")

	r18 := c.Query("r18", "false")

	works, err := PC.GetNewestArtworks(worktype, r18)
	if err != nil {
		return err
	}

	works = models.ProxyShortArtworkSlice(works, *image_proxy)

	return c.Render("newest", fiber.Map{
		"Items": works,
		"Title": "Newest works",
	})
}

func search_page(c *fiber.Ctx) error {
	image_proxy := get_session_value(c, "image-proxy")
	if image_proxy == nil {
		image_proxy = &configs.ProxyServer
	}

	name := c.Params("name")

	page := c.Query("page", "1")

	order := c.Query("order", "date_d")

	mode := c.Query("mode", "safe")

	category := c.Query("category", "artworks")

	tag, err := PC.GetTagData(name)
	if err != nil {
		return err
	}
	result, err := PC.GetSearch(category, name, order, mode, page)
	if err != nil {
		return err
	}

	result.ProxyImages(*image_proxy)

	queries := map[string]string{
		"Page":     page,
		"Order":    order,
		"Mode":     mode,
		"Category": category,
	}
	return c.Render("tag", fiber.Map{"Title": "Results for " + tag.Name, "Tag": tag, "Data": result, "Queries": queries})
}

func search(c *fiber.Ctx) error {
	name := c.FormValue("name")

	return c.Redirect("/tags/"+name, http.StatusFound)
}

func discovery_page(c *fiber.Ctx) error {
	image_proxy := get_session_value(c, "image-proxy")
	if image_proxy == nil {
		image_proxy = &configs.ProxyServer
	}

	mode := c.Query("mode", "safe")

	artworks, err := PC.GetDiscoveryArtwork(mode, 100)
	if err != nil {
		return err

	}
	artworks = models.ProxyShortArtworkSlice(artworks, *image_proxy)

	return c.Render("discovery", fiber.Map{"Title": "Discovery", "Artworks": artworks})
}

func settings_page(c *fiber.Ctx) error {
	return c.Render("settings", fiber.Map{})
}

func settings_post(c *fiber.Ctx) error {
	t := c.Params("type")
	error := ""

	if t == "image_server" {
		error = set_image_server(c)
	}

	if error != "" {
		return errors.New(error)
	}
	c.Redirect("/settings")
	return nil
}

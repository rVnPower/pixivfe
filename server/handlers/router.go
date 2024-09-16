package handlers

import (
	"errors"
	"net/http"
	"strings"

	"codeberg.org/vnpower/pixivfe/v2/server/routes"
	"github.com/gorilla/mux"
)

func handleStripPrefix(router *mux.Router, pathPrefix string, handler http.Handler) *mux.Route {
	return router.PathPrefix(pathPrefix).Handler(http.StripPrefix(pathPrefix, handler))
}

func serveFile(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	}
}

func DefineRoutes() *mux.Router {
	router := mux.NewRouter()

	// redirect handler: strip trailing / to make router behave
	router.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return r.URL.Path != "/" && strings.HasSuffix(r.URL.Path, "/")
	}).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL
		url.Path = url.Path[0 : len(url.Path)-1]
		http.Redirect(w, r, url.String(), http.StatusPermanentRedirect)
	})

	router.HandleFunc("/robots.txt", serveFile("./assets/robots.txt"))
	handleStripPrefix(router, "/img/", http.FileServer(http.Dir("./assets/img")))
	handleStripPrefix(router, "/css/", http.FileServer(http.Dir("./assets/css")))
	handleStripPrefix(router, "/js/", http.FileServer(http.Dir("./assets/js")))

	// Proxy routes. cache headers set by upstream servers.
	handleStripPrefix(router, "/proxy/i.pximg.net/", CatchError(routes.IPximgProxy)).Methods("GET")
	handleStripPrefix(router, "/proxy/s.pximg.net/", CatchError(routes.SPximgProxy)).Methods("GET")
	handleStripPrefix(router, "/proxy/ugoira.com/", CatchError(routes.UgoiraProxy)).Methods("GET")

	router.HandleFunc("/", CatchError(routes.IndexPage)).Methods("GET")
	router.HandleFunc("/about", CatchError(routes.AboutPage)).Methods("GET")
	router.HandleFunc("/newest", CatchError(routes.NewestPage)).Methods("GET")
	router.HandleFunc("/discovery", CatchError(routes.DiscoveryPage)).Methods("GET")
	router.HandleFunc("/discovery/novel", CatchError(routes.NovelDiscoveryPage)).Methods("GET")

	router.HandleFunc("/ranking", CatchError(routes.RankingPage)).Methods("GET")
	router.HandleFunc("/rankingCalendar", CatchError(routes.RankingCalendarPage)).Methods("GET")
	router.HandleFunc("/rankingCalendar", CatchError(routes.RankingCalendarPicker)).Methods("POST")

	router.HandleFunc("/users/{id}.atom.xml", CatchError(routes.UserAtomFeed)).Methods("GET")
	router.HandleFunc("/users/{id}/{category}.atom.xml", CatchError(routes.UserAtomFeed)).Methods("GET")
	router.HandleFunc("/users/{id}", CatchError(routes.UserPage)).Methods("GET")
	router.HandleFunc("/users/{id}/{category}", CatchError(routes.UserPage)).Methods("GET")

	router.HandleFunc("/artworks/{id}", CatchError(routes.ArtworkPage)).Methods("GET")
	router.HandleFunc("/artworks-multi/{ids}", CatchError(routes.ArtworkMultiPage)).Methods("GET")
	// Legacy illust URL
	router.HandleFunc("/member_illust.php", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/artworks/"+routes.GetQueryParam(r, "illust_id"), http.StatusPermanentRedirect)
	}).Methods("GET")

	router.HandleFunc("/novel/show.php", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/novel/"+routes.GetQueryParam(r, "id"), http.StatusPermanentRedirect)
	}).Methods("GET")
	router.HandleFunc("/novel/{id}", CatchError(routes.NovelPage)).Methods("GET")

	router.HandleFunc("/pixivision", CatchError(routes.PixivisionHomePage)).Methods("GET")
	router.HandleFunc("/pixivision/a/{id}", CatchError(routes.PixivisionArticlePage)).Methods("GET")

	router.HandleFunc("/settings", CatchError(routes.SettingsPage)).Methods("GET")
	router.HandleFunc("/settings/{type}", CatchError(routes.SettingsPost)).Methods("POST")

	router.HandleFunc("/self", CatchError(routes.LoginUserPage)).Methods("GET")
	router.HandleFunc("/self/followingWorks", CatchError(routes.FollowingWorksPage)).Methods("GET")
	router.HandleFunc("/self/bookmarks", CatchError(routes.LoginBookmarkPage)).Methods("GET")
	router.HandleFunc("/self/addBookmark/{id}", CatchError(routes.AddBookmarkRoute)).Methods("GET")
	router.HandleFunc("/self/deleteBookmark/{id}", CatchError(routes.DeleteBookmarkRoute)).Methods("GET")
	router.HandleFunc("/self/like/{id}", CatchError(routes.LikeRoute)).Methods("GET")

	router.HandleFunc("/oembed", CatchError(routes.Oembed)).Methods("GET")

	router.HandleFunc("/tags/{name}", CatchError(routes.TagPage)).Methods("GET")
	router.HandleFunc("/tags/{name}", CatchError(routes.TagPage)).Methods("POST")
	router.HandleFunc("/tags", CatchError(routes.TagPage)).Methods("GET")
	router.HandleFunc("/tags", CatchError(routes.AdvancedTagPost)).Methods("POST")

	router.HandleFunc("/diagnostics", CatchError(routes.Diagnostics)).Methods("GET")
	router.HandleFunc("/diagnostics/spans.json", CatchError(routes.DiagnosticsData)).Methods("GET")
	router.HandleFunc("/diagnostics/reset", routes.ResetDiagnosticsData)

	// fallback route (if nothing else matches)
	router.NewRoute().HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routes.ErrorPage(w, r, errors.New("Route not found"), http.StatusNotFound)
	})

	return router
}

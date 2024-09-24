package handlers

import (
	"errors"
	"net/http"
	"strings"

	"codeberg.org/vnpower/pixivfe/v2/server/routes"
	"github.com/gorilla/mux"
)

// handleStripPrefix is a utility function that combines path prefix matching with
// stripping the prefix from the request URL before passing it to the handler.
func handleStripPrefix(router *mux.Router, pathPrefix string, handler http.Handler) *mux.Route {
	return router.PathPrefix(pathPrefix).Handler(http.StripPrefix(pathPrefix, handler))
}

// serveFile returns an http.HandlerFunc that serves a specific file.
// This is useful for serving static files like robots.txt.
func serveFile(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	}
}

// DefineRoutes sets up all the routes for the application.
// It returns a configured mux.Router with all paths and their corresponding handlers.
func DefineRoutes() *mux.Router {
	router := mux.NewRouter()

	// Tutorial: Adding new routes
	// 1. Use router.HandleFunc to define the path and handler
	// 2. Wrap the handler function, defined in package routes, with CatchError for error handling
	// 3. Specify the HTTP method(s) using .Methods()
	// 4. For URL parameters, use curly braces in the path, e.g., "/users/{id}"
	// 5. Group similar routes together for better organization
	//
	// Example:
	// router.HandleFunc("/new/route/{param}", CatchError(routes.NewHandler)).Methods("GET", "POST")

	// Redirect handler: strip trailing / to make router behave consistently
	// This ensures that URLs with and without trailing slashes are treated the same
	router.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return r.URL.Path != "/" && strings.HasSuffix(r.URL.Path, "/")
	}).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL
		url.Path = url.Path[0 : len(url.Path)-1]
		http.Redirect(w, r, url.String(), http.StatusPermanentRedirect)
	})

	// Serve static files
	router.HandleFunc("/robots.txt", serveFile("./assets/robots.txt"))
	handleStripPrefix(router, "/img/", http.FileServer(http.Dir("./assets/img")))
	handleStripPrefix(router, "/css/", http.FileServer(http.Dir("./assets/css")))
	handleStripPrefix(router, "/js/", http.FileServer(http.Dir("./assets/js")))

	// Proxy routes for handling image requests
	// These routes maintain cache headers set by upstream servers
	handleStripPrefix(router, "/proxy/i.pximg.net/", CatchError(routes.IPximgProxy)).Methods("GET")
	handleStripPrefix(router, "/proxy/s.pximg.net/", CatchError(routes.SPximgProxy)).Methods("GET")
	handleStripPrefix(router, "/proxy/ugoira.com/", CatchError(routes.UgoiraProxy)).Methods("GET")

	// Main application routes
	router.HandleFunc("/", CatchError(routes.IndexPage)).Methods("GET")
	router.HandleFunc("/about", CatchError(routes.AboutPage)).Methods("GET")
	router.HandleFunc("/newest", CatchError(routes.NewestPage)).Methods("GET")
	router.HandleFunc("/discovery", CatchError(routes.DiscoveryPage)).Methods("GET")
	router.HandleFunc("/discovery/novel", CatchError(routes.NovelDiscoveryPage)).Methods("GET")

	// Ranking related routes
	router.HandleFunc("/ranking", CatchError(routes.RankingPage)).Methods("GET")
	router.HandleFunc("/rankingCalendar", CatchError(routes.RankingCalendarPage)).Methods("GET")
	router.HandleFunc("/rankingCalendar", CatchError(routes.RankingCalendarPicker)).Methods("POST")

	// User related routes, including Atom feeds
	router.HandleFunc("/users/{id}.atom.xml", CatchError(routes.UserAtomFeed)).Methods("GET")
	router.HandleFunc("/users/{id}/{category}.atom.xml", CatchError(routes.UserAtomFeed)).Methods("GET")
	router.HandleFunc("/users/{id}", CatchError(routes.UserPage)).Methods("GET")
	router.HandleFunc("/users/{id}/{category}", CatchError(routes.UserPage)).Methods("GET")

	// Artwork related routes
	router.HandleFunc("/artworks/{id}", CatchError(routes.ArtworkPage)).Methods("GET")
	router.HandleFunc("/artworks-multi/{ids}", CatchError(routes.ArtworkMultiPage)).Methods("GET")
	// Legacy illust URL redirect
	router.HandleFunc("/member_illust.php", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/artworks/"+routes.GetQueryParam(r, "illust_id"), http.StatusPermanentRedirect)
	}).Methods("GET")

	// Novel related routes
	router.HandleFunc("/novel/show.php", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/novel/"+routes.GetQueryParam(r, "id"), http.StatusPermanentRedirect)
	}).Methods("GET")
	router.HandleFunc("/novel/{id}", CatchError(routes.NovelPage)).Methods("GET")
	router.HandleFunc("/novel/series/{id}", CatchError(routes.NovelSeriesPage)).Methods("GET")

	// Pixivision related routes
	router.HandleFunc("/pixivision", CatchError(routes.PixivisionHomePage)).Methods("GET")
	router.HandleFunc("/pixivision/a/{id}", CatchError(routes.PixivisionArticlePage)).Methods("GET")
	router.HandleFunc("/pixivision/c/{id}", CatchError(routes.PixivisionCategoryPage)).Methods("GET")
	router.HandleFunc("/pixivision/t/{id}", CatchError(routes.PixivisionTagPage)).Methods("GET")

	// Settings related routes
	router.HandleFunc("/settings", CatchError(routes.SettingsPage)).Methods("GET")
	router.HandleFunc("/settings/{type}", CatchError(routes.SettingsPost)).Methods("POST")

	// User action routes (login, bookmarks, likes, etc.)
	router.HandleFunc("/self", CatchError(routes.LoginUserPage)).Methods("GET")
	router.HandleFunc("/self/followingWorks", CatchError(routes.FollowingWorksPage)).Methods("GET")
	router.HandleFunc("/self/bookmarks", CatchError(routes.LoginBookmarkPage)).Methods("GET")
	router.HandleFunc("/self/addBookmark/{id}", CatchError(routes.AddBookmarkRoute)).Methods("GET")
	router.HandleFunc("/self/deleteBookmark/{id}", CatchError(routes.DeleteBookmarkRoute)).Methods("GET")
	router.HandleFunc("/self/like/{id}", CatchError(routes.LikeRoute)).Methods("GET")

	// oEmbed endpoint for embedding Pixiv content
	router.HandleFunc("/oembed", CatchError(routes.Oembed)).Methods("GET")

	// Tag related routes
	router.HandleFunc("/tags/{name}", CatchError(routes.TagPage)).Methods("GET")
	router.HandleFunc("/tags/{name}", CatchError(routes.TagPage)).Methods("POST")
	router.HandleFunc("/tags", CatchError(routes.TagPage)).Methods("GET")
	router.HandleFunc("/tags", CatchError(routes.AdvancedTagPost)).Methods("POST")

	// Diagnostic routes for monitoring and debugging
	router.HandleFunc("/diagnostics", CatchError(routes.Diagnostics)).Methods("GET")
	router.HandleFunc("/diagnostics/spans.json", CatchError(routes.DiagnosticsData)).Methods("GET")
	router.HandleFunc("/diagnostics/reset", routes.ResetDiagnosticsData)

	// Fallback route (if nothing else matches)
	// This ensures that a proper HTTP 404 error is returned for undefined routes
	router.NewRoute().HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routes.ErrorPage(w, r, errors.New("Route not found"), http.StatusNotFound)
	})

	return router
}

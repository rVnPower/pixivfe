package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/routes"
	"codeberg.org/vnpower/pixivfe/v2/session"
	"codeberg.org/vnpower/pixivfe/v2/utils"
)

func CanRequestSkipLimiter(r *http.Request) bool {
	path := r.URL.Path
	return strings.HasPrefix(path, "/img/") ||
		strings.HasPrefix(path, "/css/") ||
		strings.HasPrefix(path, "/js/") ||
		strings.HasPrefix(path, "/proxy/s.pximg.net/")
}

func CanRequestSkipLogger(r *http.Request) bool {
	// return false
	path := r.URL.Path
	return CanRequestSkipLimiter(r) ||
		strings.HasPrefix(path, "/proxy/i.pximg.net/")
}

type UserContext struct {
	err error
}

type userContextKey struct{}

var UserContextKey = userContextKey{}

func GetUserContext(r *http.Request) *UserContext {
	return r.Context().Value(UserContextKey).(*UserContext)
}

func main() {
	config.GlobalServerConfig.InitializeConfig()
	if config.GlobalServerConfig.InDevelopment {
		core.CreateResponseAuditFolder()
	}
	routes.InitTemplatingEngine(config.GlobalServerConfig.InDevelopment)

	router := defineRoutes()

	main_handler := func(w_ http.ResponseWriter, r *http.Request) {
		w := &ResponseWriterInterceptStatus{
			statusCode:     0,
			ResponseWriter: w_,
		}
		// set user context
		r = r.WithContext(context.WithValue(r.Context(), UserContextKey, &UserContext{}))

		start_time := time.Now()

		setGlobalHeaders(w, r)

		if r.URL.Path != "/" && strings.HasSuffix(r.URL.Path, "/") {
			// strip trailing / to make router behave
			url := r.URL
			url.Path, _ = strings.CutSuffix(url.Path, "/")
			http.Redirect(w, r, url.String(), http.StatusPermanentRedirect)
		} else {
			// all the routes are listed here
			router.ServeHTTP(w, r)
		}


		CatchError(func(w http.ResponseWriter, r *http.Request) error {
			err := GetUserContext(r).err
			if err != nil { // error handler
				log.Println("Within handler: ", err)
				code := http.StatusInternalServerError
				w.WriteHeader(code)
				// Send custom error page
				err = routes.Render(w, r, routes.Data_error{Title: "Error", Error: err})
				if err != nil {
					err = utils.SendString(w, (fmt.Sprintf("Internal Server Error: %s", err)))
					if err != nil {
						return err
					}
				}
			}
			return nil
		})(w, r)

		end_time := time.Now()

		if !CanRequestSkipLogger(r) { // logger
			time := start_time
			latency := end_time.Sub(start_time)
			ip := r.RemoteAddr
			method := r.Method
			path := r.URL.Path
			status := w.statusCode
			err := GetUserContext(r).err

			log.Printf("%v +%v %v %v %v %v %v", time, latency, ip, method, path, status, err)
		}
	}

	// Initialize and start the proxy checker
	ctx_timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	config.InitializeProxyChecker(ctx_timeout)

	// run sass when in development mode
	if config.GlobalServerConfig.InDevelopment {
		go func() {
			cmd := exec.Command("sass", "--watch", "assets/css")
			cmd.Stdout = os.Stderr // Sass quirk
			cmd.Stderr = os.Stderr
			cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pdeathsig: syscall.SIGHUP}
			runtime.LockOSThread() // Go quirk https://github.com/golang/go/issues/27505
			err := cmd.Run()
			if err != nil {
				log.Println(fmt.Errorf("when running sass: %w", err))
			}
		}()
	}

	// Listen
	var l net.Listener
	if config.GlobalServerConfig.UnixSocket != "" {
		ln, err := net.Listen("unix", config.GlobalServerConfig.UnixSocket)
		if err != nil {
			panic(err)
		}
		l = ln
		log.Printf("Listening on domain socket %v\n", config.GlobalServerConfig.UnixSocket)
	} else {
		addr := config.GlobalServerConfig.Host + ":" + config.GlobalServerConfig.Port
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			log.Panicf("failed to listen: %v", err)
		}
		l = ln
		addr = ln.Addr().String()
		log.Printf("Listening on http://%v/\n", addr)
	}
	http.Serve(l, http.HandlerFunc(main_handler))
}

func setGlobalHeaders(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Add("X-Frame-Options", "DENY")
	// use this if need iframe: `X-Frame-Options: SAMEORIGIN`
	header.Add("X-Content-Type-Options", "nosniff")
	header.Add("Referrer-Policy", "same-origin") // needed for settings redirect
	header.Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	header.Add("Content-Security-Policy", fmt.Sprintf("base-uri 'self'; default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self' %s; media-src 'self' %s; connect-src 'self'; form-action 'self'; frame-ancestors 'none';", session.GetImageProxyOrigin(r), session.GetImageProxyOrigin(r)))
	// use this if need iframe: `frame-ancestors 'self'`
	header.Add("Permissions-Policy", "accelerometer=(), ambient-light-sensor=(), battery=(), camera=(), display-capture=(), document-domain=(), encrypted-media=(), execution-while-not-rendered=(), execution-while-out-of-viewport=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), midi=(), navigation-override=(), payment=(), publickey-credentials-get=(), screen-wake-lock=(), sync-xhr=(), usb=(), web-share=(), xr-spatial-tracking=()")
}

func serveFile(filename string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, filename) }
}

func handlePrefix(router *mux.Router, pathPrefix string, handler http.Handler) *mux.Route {
	return router.PathPrefix(pathPrefix).Handler(http.StripPrefix(pathPrefix, handler))
}
func defineRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/favicon.ico", serveFile("./assets/img/favicon.ico"))
	router.HandleFunc("/robots.txt", serveFile("./assets/robots.txt"))
	handlePrefix(router, "/img/", http.FileServer(http.Dir("./assets/img")))
	handlePrefix(router, "/css/", http.FileServer(http.Dir("./assets/css")))
	handlePrefix(router, "/js/", http.FileServer(http.Dir("./assets/js")))

	// Routes
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
	router.HandleFunc("/novel/{id}", CatchError(routes.NovelPage)).Methods("GET")
	router.HandleFunc("/pixivision", CatchError(routes.PixivisionHomePage)).Methods("GET")
	router.HandleFunc("/pixivision/a/{id}", CatchError(routes.PixivisionArticlePage)).Methods("GET")

	// Settings group
	router.HandleFunc("/settings", CatchError(routes.SettingsPage)).Methods("GET")
	router.HandleFunc("/settings/{type}", CatchError(routes.SettingsPost)).Methods("POST")

	// Personal group
	router.HandleFunc("/self", CatchError(routes.LoginUserPage)).Methods("GET")
	router.HandleFunc("/self/followingWorks", CatchError(routes.FollowingWorksPage)).Methods("GET")
	router.HandleFunc("/self/bookmarks", CatchError(routes.LoginBookmarkPage)).Methods("GET")
	router.HandleFunc("/self/addBookmark/{id}", CatchError(routes.AddBookmarkRoute)).Methods("GET")
	router.HandleFunc("/self/deleteBookmark/{id}", CatchError(routes.DeleteBookmarkRoute)).Methods("GET")
	router.HandleFunc("/self/like/{id}", CatchError(routes.LikeRoute)).Methods("GET")

	// Oembed group
	router.HandleFunc("/oembed", CatchError(routes.Oembed)).Methods("GET")

	router.HandleFunc("/tags/{name}", CatchError(routes.TagPage)).Methods("GET")
	router.HandleFunc("/tags/{name}", CatchError(routes.TagPage)).Methods("POST")
	router.HandleFunc("/tags", CatchError(routes.TagPage)).Methods("GET")
	router.HandleFunc("/tags", CatchError(routes.AdvancedTagPost)).Methods("POST")

	// Legacy illust URL
	router.HandleFunc("/member_illust.php", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/artworks/"+ routes.GetQueryParam(r, "illust_id"), http.StatusPermanentRedirect)
	}).Methods("GET")

	// Proxy routes
	handlePrefix(router, "/proxy/i.pximg.net/", CatchError(routes.IPximgProxy)).Methods("GET")
	handlePrefix(router, "/proxy/s.pximg.net/", CatchError(routes.SPximgProxy)).Methods("GET")
	handlePrefix(router, "/proxy/ugoira.com/", CatchError(routes.UgoiraProxy)).Methods("GET")

	router.NewRoute().HandlerFunc(CatchError(func(w http.ResponseWriter, r *http.Request) error {
		return errors.New("Route not found")
	}))

	return router
}

func CatchError(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		GetUserContext(r).err = handler(w, r)
	}
}

type ResponseWriterInterceptStatus struct {
	statusCode int
	http.ResponseWriter
}

func (w *ResponseWriterInterceptStatus) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"maps"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"

	"codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/routes"
	"codeberg.org/vnpower/pixivfe/v2/template"
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
	return strings.HasPrefix(path, "/img/") ||
		strings.HasPrefix(path, "/css/") ||
		strings.HasPrefix(path, "/js/") ||
		strings.HasPrefix(path, "/proxy/s.pximg.net/") ||
		strings.HasPrefix(path, "/proxy/i.pximg.net/")
}

type UserContext struct {
	err        error
	statusCode int
}

type userContextKey struct{}

var UserContextKey = userContextKey{}

func GetUserContext(r *http.Request) *UserContext {
	return r.Context().Value(UserContextKey).(*UserContext)
}

// Todo: Should we put middlewares in a separate file?
// IPRateLimiter represents an IP rate limiter.
type IPRateLimiter struct {
	ips     map[string]*rate.Limiter
	mu      *sync.RWMutex
	limiter *rate.Limiter
}

// NewIPRateLimiter creates a new instance of IPRateLimiter with the given rate limit.
func NewIPRateLimiter(r rate.Limit, burst int) *IPRateLimiter {
	return &IPRateLimiter{
		ips:     make(map[string]*rate.Limiter),
		mu:      &sync.RWMutex{},
		limiter: rate.NewLimiter(r, burst),
	}
}

// Allow checks if the request from the given IP is allowed.
func (lim *IPRateLimiter) Allow(ip string) bool {
	lim.mu.RLock()
	rl, exists := lim.ips[ip]
	lim.mu.RUnlock()

	if !exists {
		lim.mu.Lock()
		rl, exists = lim.ips[ip]
		if !exists {
			rl = rate.NewLimiter(lim.limiter.Limit(), lim.limiter.Burst())
			lim.ips[ip] = rl
		}
		lim.mu.Unlock()
	}

	return rl.Allow()
}

var limiter *IPRateLimiter = NewIPRateLimiter(0, 0)

func MiddlewareChain(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		if CanRequestSkipLogger(r) {
			handler.ServeHTTP(w, r)
			return
		}

		if !limiter.Allow(ip) {
			CatchError(func(w http.ResponseWriter, r *http.Request) error {
				GetUserContext(r).statusCode = http.StatusTooManyRequests
				return errors.New("Too many requests")
			})(w, r)
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}

func main() {
	config.GlobalServerConfig.InitializeConfig()
	if config.GlobalServerConfig.InDevelopment {
		core.CreateResponseAuditFolder()
	}
	template.InitTemplatingEngine(config.GlobalServerConfig.InDevelopment)

	// Initialize and start the proxy checker
	ctx_timeout, cancel := context.WithTimeout(context.Background(), config.ProxyCheckerTimeout)
	defer cancel()
	config.InitializeProxyChecker(ctx_timeout)

	router := defineRoutes()

	main_handler := func(w_ http.ResponseWriter, r *http.Request) {
		println("main handler")
		w := &ResponseWriterInterceptStatus{
			statusCode:     0,
			ResponseWriter: w_,
		}
		// set user context
		r = r.WithContext(context.WithValue(r.Context(), UserContextKey, &UserContext{}))

		start_time := time.Now()

		if r.URL.Path != "/" && strings.HasSuffix(r.URL.Path, "/") {
			// redirect handler: strip trailing / to make router behave
			url := r.URL
			url.Path, _ = strings.CutSuffix(url.Path, "/")
			http.Redirect(w, r, url.String(), http.StatusPermanentRedirect)
		} else {
			// handler of all other routes
			router.ServeHTTP(w, r)
		}

		{ // error handler
			err := GetUserContext(r).err

			if err != nil {
				log.Printf("Internal Server Error: %s", err)
				code := GetUserContext(r).statusCode
				if code == 0 {
					code = http.StatusInternalServerError
				}
				w.WriteHeader(code)
				// Send custom error page
				err = routes.ErrorPage(w, r, err)
				if err != nil {
					log.Printf("Error rendering error route: %s", err)
				}
			}
		}

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

func serveFile(filename string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, filename) }
}

func handlePrefix(router *mux.Router, pathPrefix string, handler http.Handler) *mux.Route {
	return router.PathPrefix(pathPrefix).Handler(http.StripPrefix(pathPrefix, handler))
}
func defineRoutes() *mux.Router {
	router := mux.NewRouter()

	//router.Use(MiddlewareChain)

	router.HandleFunc("/favicon.ico", serveFile("./assets/img/favicon.ico"))
	router.HandleFunc("/robots.txt", serveFile("./assets/robots.txt"))
	handlePrefix(router, "/img/", http.FileServer(http.Dir("./assets/img")))
	handlePrefix(router, "/css/", http.FileServer(http.Dir("./assets/css")))
	handlePrefix(router, "/js/", http.FileServer(http.Dir("./assets/js")))

	// Proxy routes. cache headers set by upstream servers.
	handlePrefix(router, "/proxy/i.pximg.net/", CatchError(routes.IPximgProxy)).Methods("GET")
	handlePrefix(router, "/proxy/s.pximg.net/", CatchError(routes.SPximgProxy)).Methods("GET")
	handlePrefix(router, "/proxy/ugoira.com/", CatchError(routes.UgoiraProxy)).Methods("GET")

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

	// Legacy illust URL
	router.HandleFunc("/member_illust.php", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/artworks/"+routes.GetQueryParam(r, "illust_id"), http.StatusPermanentRedirect)
	}).Methods("GET")

	router.NewRoute().HandlerFunc(CatchError(func(w http.ResponseWriter, r *http.Request) error {
		return errors.New("Route not found")
	}))

	return router
}

func CatchError(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header_backup := http.Header{}
		for k, v := range w.Header() {
			header_backup[k] = slices.Clone(v)
		}
		recorder := httptest.ResponseRecorder{
			HeaderMap: w.Header(),
			Body:      new(bytes.Buffer),
			Code:      200,
		}
		err := handler(&recorder, r)
		if err != nil {
			clear(header_backup)
			maps.Copy(w.Header(), header_backup)
			GetUserContext(r).err = err
		} else {
			_, _ = recorder.Body.WriteTo(w)
			w.WriteHeader(recorder.Code)
		}
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

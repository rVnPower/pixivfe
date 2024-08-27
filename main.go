package main

import (
	"context"
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

	// server := fiber.New(fiber.Config{
	// 	AppName:                 "PixivFE",
	// 	DisableStartupMessage:   true,
	// 	Prefork:                 false,
	// 	JSONEncoder:             json.Marshal,
	// 	JSONDecoder:             json.Unmarshal,
	// 	EnableTrustedProxyCheck: true,
	// 	TrustedProxies:          []string{"0.0.0.0/0"},
	// 	ProxyHeader:             fiber.HeaderXForwardedFor,
	// 	ErrorHandler: func(r *http.Request, err error) error {
	// 		log.Println(err)

	// 		// Status code defaults to 500
	// 		code := fiber.StatusInternalServerError

	// 		// // Retrieve the custom status code if it's a *fiber.Error
	// 		// var e *fiber.Error
	// 		// if errors.As(err, &e) {
	// 		// 	code = e.Code
	// 		// }

	// 		// Send custom error page
	// 		r.Status(code)
	// 		err = routes.Render(w, r, routes.Data_error{Title: "Error", Error: err})
	// 		if err != nil {
	// 			return r.Status(code).SendString(fmt.Sprintf("Internal Server Error: %s", err))
	// 		}

	// 		return nil
	// 	},
	// })

	// todo: limiter
	// if config.GlobalServerConfig.RequestLimit > 0 {
	// 	keyedSleepingSpot := kmutex.New()
	// 	server.Use(limiter.New(limiter.Config{
	// 		Next:              CanRequestSkipLimiter,
	// 		Expiration:        30 * time.Second,
	// 		Max:               config.GlobalServerConfig.RequestLimit,
	// 		LimiterMiddleware: limiter.SlidingWindow{},
	// 		LimitReached: func(r *http.Request) error {
	// 			// limit response throughput by pacing, since not every bot reads X-RateLimit-*
	// 			// on limit reached, they just have to wait
	// 			// the design of this means that if they send multiple requests when reaching rate limit, they will wait even longer (since `retryAfter` is calculated before anything has slept)
	// 			retryAfter_s := r.GetRespHeader(fiber.HeaderRetryAfter)
	// 			retryAfter, err := strconv.ParseUint(retryAfter_s, 10, 64)
	// 			if err != nil {
	// 				log.Panicf("response header 'RetryAfter' should be a number: %v", err)
	// 			}
	// 			requestIP := r.IP()
	// 			refcount := keyedSleepingSpot.Lock(requestIP)
	// 			defer keyedSleepingSpot.Unlock(requestIP)
	// 			if refcount >= 4 { // on too much concurrent requests
	// 				// todo: maybe blackhole `requestIP` here
	// 				log.Println("Limit Reached (Hard)!", requestIP)
	// 				// close the connection immediately
	// 				_ = r.Context().Conn().Close()
	// 				return nil
	// 			}

	// 			// sleeping
	// 			// here, sleeping is not the best solution.
	// 			// todo: close this connection when this IP reaches hard limit
	// 			dur := time.Duration(retryAfter) * time.Second
	// 			log.Println("Limit Reached (Soft)! Sleeping for ", dur)
	// 			ctx, cancel := context.WithTimeout(r.Context(), dur)
	// 			defer cancel()
	// 			<-ctx.Done()

	// 			return r.Next()
	// 		},
	// 	}))
	// }

	// todo: caching
	// if !config.GlobalServerConfig.InDevelopment {
	// 	server.Use(cache.New(
	// 		cache.Config{
	// 			Next: func(r *http.Request) bool {
	// 				resp_code := r.Response().StatusCode()
	// 				if resp_code < 200 || resp_code >= 300 {
	// 					return true
	// 				}

	// 				// Disable cache for settings page
	// 				return strings.Contains(r.Path(), "/settings") || r.Path() == "/"
	// 			},
	// 			Expiration:           5 * time.Minute,
	// 			CacheControl:         true,
	// 			StoreResponseHeaders: true,

	// 			KeyGenerator: func(r *http.Request) string {
	// 				key := fiber_utils.CopyString(r.OriginalURL())
	// 				for _, cookieName := range session.AllCookieNames {
	// 					cookieValue := session.GetCookie(r.Request, cookieName)
	// 					if cookieValue != "" {
	// 						key += "\x00\x00"
	// 						key += string(cookieName)
	// 						key += "\x00"
	// 						key += cookieValue
	// 					}
	// 				}
	// 				return key
	// 			},
	// 		},
	// 	))
	// }

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

		router.ServeHTTP(w, r)

		// todo: test this
		// redirect any request with ?r=url
		ret := r.URL.Query().Get("r")
		if ret != "" {
			// could this be unsafe since this redirects to any website?
			http.Redirect(w, r, ret, http.StatusTemporaryRedirect)
		}

		end_time := time.Now()

		if !CanRequestSkipLogger(r) {
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
	header.Add("Referrer-Policy", "no-referrer")
	header.Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	header.Add("Content-Security-Policy", fmt.Sprintf("base-uri 'self'; default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self' %s; media-src 'self' %s; connect-src 'self'; form-action 'self'; frame-ancestors 'none';", session.GetImageProxyOrigin(r), session.GetImageProxyOrigin(r)))
	// use this if need iframe: `frame-ancestors 'self'`
	header.Add("Permissions-Policy", "accelerometer=(), ambient-light-sensor=(), battery=(), camera=(), display-capture=(), document-domain=(), encrypted-media=(), execution-while-not-rendered=(), execution-while-out-of-viewport=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), midi=(), navigation-override=(), payment=(), publickey-credentials-get=(), screen-wake-lock=(), sync-xhr=(), usb=(), web-share=(), xr-spatial-tracking=()")
}

func serveFile(filename string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, filename) }
}
func defineRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/favicon.ico", serveFile("./assets/img/favicon.ico"))
	router.HandleFunc("/robots.txt", serveFile("./assets/robots.txt"))
	router.Handle("/img/", http.FileServer(http.Dir("./assets/img")))
	router.Handle("/css/", http.FileServer(http.Dir("./assets/css")))
	router.Handle("/js/", http.FileServer(http.Dir("./assets/js")))

	// Routes
	router.HandleFunc("/", CatchError(routes.IndexPage)).Methods("GET")
	router.HandleFunc("/about", CatchError(routes.AboutPage)).Methods("GET")
	router.HandleFunc("/newest", CatchError(routes.NewestPage)).Methods("GET")
	router.HandleFunc("/discovery", CatchError(routes.DiscoveryPage)).Methods("GET")
	router.HandleFunc("/discovery/novel", CatchError(routes.NovelDiscoveryPage)).Methods("GET")
	router.HandleFunc("/ranking", CatchError(routes.RankingPage)).Methods("GET")
	router.HandleFunc("/rankingCalendar", CatchError(routes.RankingCalendarPage)).Methods("GET")
	router.HandleFunc("/rankingCalendar", CatchError(routes.RankingCalendarPicker)).Methods("POST")
	router.HandleFunc("/users/:id.atom.xml", CatchError(routes.UserAtomFeed)).Methods("GET")
	router.HandleFunc("/users/:id/:category.atom.xml", CatchError(routes.UserAtomFeed)).Methods("GET")
	router.HandleFunc("/users/:id/:category?", CatchError(routes.UserPage)).Methods("GET")
	router.HandleFunc("/artworks/:id/", CatchError(routes.ArtworkPage)).Methods("GET")
	router.HandleFunc("/artworks-multi/:ids/", CatchError(routes.ArtworkMultiPage)).Methods("GET")
	router.HandleFunc("/novel/:id/", CatchError(routes.NovelPage)).Methods("GET")
	router.HandleFunc("/pixivision", CatchError(routes.PixivisionHomePage)).Methods("GET")
	router.HandleFunc("/pixivision/a/:id", CatchError(routes.PixivisionArticlePage)).Methods("GET")

	// Settings group
	settings := router.NewRoute().Path("/settings").Subrouter()
	settings.HandleFunc("/", CatchError(routes.SettingsPage)).Methods("GET")
	settings.HandleFunc("/:type/:noredirect?", CatchError(routes.SettingsPost)).Methods("POST")

	// Personal group
	self := router.NewRoute().Path("/self").Subrouter()
	self.HandleFunc("/", CatchError(routes.LoginUserPage)).Methods("GET")
	self.HandleFunc("/followingWorks", CatchError(routes.FollowingWorksPage)).Methods("GET")
	self.HandleFunc("/bookmarks", CatchError(routes.LoginBookmarkPage)).Methods("GET")
	self.HandleFunc("/addBookmark/:id", CatchError(routes.AddBookmarkRoute)).Methods("GET")
	self.HandleFunc("/deleteBookmark/:id", CatchError(routes.DeleteBookmarkRoute)).Methods("GET")
	self.HandleFunc("/like/:id", CatchError(routes.LikeRoute)).Methods("GET")

	// Oembed group
	router.HandleFunc("/oembed", CatchError(routes.Oembed)).Methods("GET")

	router.HandleFunc("/tags/:name", CatchError(routes.TagPage)).Methods("GET")
	router.HandleFunc("/tags/:name", CatchError(routes.TagPage)).Methods("POST")
	router.HandleFunc("/tags", CatchError(routes.TagPage)).Methods("GET")
	router.HandleFunc("/tags", CatchError(routes.AdvancedTagPost)).Methods("POST")

	// Legacy illust URL
	router.HandleFunc("/member_illust.php", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/artworks/"+routes.CompatRequest{Request: r}.Query("illust_id"), http.StatusFound)
	}).Methods("GET")

	// Proxy routes
	proxy := router.NewRoute().Path("/proxy").Subrouter()
	proxy.HandleFunc("/i.pximg.net/*", CatchError(routes.IPximgProxy)).Methods("GET")
	proxy.HandleFunc("/s.pximg.net/*", CatchError(routes.SPximgProxy)).Methods("GET")
	proxy.HandleFunc("/ugoira.com/*", CatchError(routes.UgoiraProxy)).Methods("GET")

	return router
}

func CatchError(handler func(w http.ResponseWriter, r routes.CompatRequest) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		GetUserContext(r).err = handler(w, routes.CompatRequest{Request: r})
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

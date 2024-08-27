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

	"codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/routes"
	"codeberg.org/vnpower/pixivfe/v2/session"
	// "codeberg.org/vnpower/pixivfe/v2/utils/kmutex"
	// "github.com/CloudyKit/jet/v6"
	// "github.com/goccy/go-json"
	// "net/http"
	// "net/http/middleware/cache"
	// "net/http/middleware/compress"
	// "net/http/middleware/limiter"
	// "net/http/middleware/logger"
	// "net/http/middleware/recover"
	// fiber_utils "net/http/utils"
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
	// 	ErrorHandler: func(c *http.Request, err error) error {
	// 		log.Println(err)

	// 		// Status code defaults to 500
	// 		code := fiber.StatusInternalServerError

	// 		// // Retrieve the custom status code if it's a *fiber.Error
	// 		// var e *fiber.Error
	// 		// if errors.As(err, &e) {
	// 		// 	code = e.Code
	// 		// }

	// 		// Send custom error page
	// 		c.Status(code)
	// 		err = routes.Render(c, routes.Data_error{Title: "Error", Error: err})
	// 		if err != nil {
	// 			return c.Status(code).SendString(fmt.Sprintf("Internal Server Error: %s", err))
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
	// 		LimitReached: func(c *http.Request) error {
	// 			// limit response throughput by pacing, since not every bot reads X-RateLimit-*
	// 			// on limit reached, they just have to wait
	// 			// the design of this means that if they send multiple requests when reaching rate limit, they will wait even longer (since `retryAfter` is calculated before anything has slept)
	// 			retryAfter_s := c.GetRespHeader(fiber.HeaderRetryAfter)
	// 			retryAfter, err := strconv.ParseUint(retryAfter_s, 10, 64)
	// 			if err != nil {
	// 				log.Panicf("response header 'RetryAfter' should be a number: %v", err)
	// 			}
	// 			requestIP := c.IP()
	// 			refcount := keyedSleepingSpot.Lock(requestIP)
	// 			defer keyedSleepingSpot.Unlock(requestIP)
	// 			if refcount >= 4 { // on too much concurrent requests
	// 				// todo: maybe blackhole `requestIP` here
	// 				log.Println("Limit Reached (Hard)!", requestIP)
	// 				// close the connection immediately
	// 				_ = c.Context().Conn().Close()
	// 				return nil
	// 			}

	// 			// sleeping
	// 			// here, sleeping is not the best solution.
	// 			// todo: close this connection when this IP reaches hard limit
	// 			dur := time.Duration(retryAfter) * time.Second
	// 			log.Println("Limit Reached (Soft)! Sleeping for ", dur)
	// 			ctx, cancel := context.WithTimeout(c.Context(), dur)
	// 			defer cancel()
	// 			<-ctx.Done()

	// 			return c.Next()
	// 		},
	// 	}))
	// }

	// todo: caching
	// if !config.GlobalServerConfig.InDevelopment {
	// 	server.Use(cache.New(
	// 		cache.Config{
	// 			Next: func(c *http.Request) bool {
	// 				resp_code := c.Response().StatusCode()
	// 				if resp_code < 200 || resp_code >= 300 {
	// 					return true
	// 				}

	// 				// Disable cache for settings page
	// 				return strings.Contains(c.Path(), "/settings") || c.Path() == "/"
	// 			},
	// 			Expiration:           5 * time.Minute,
	// 			CacheControl:         true,
	// 			StoreResponseHeaders: true,

	// 			KeyGenerator: func(c *http.Request) string {
	// 				key := fiber_utils.CopyString(c.OriginalURL())
	// 				for _, cookieName := range session.AllCookieNames {
	// 					cookieValue := session.GetCookie(c, cookieName)
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

	main_handler := func(w http.ResponseWriter, r *http.Request) {
		start_time := time.Now()

		setGlobalHeaders(r)

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
			status := r.Response.Status

			log.Printf("%v +%v %v %v %v %v {todo: print error (where is error?)}", time, latency, ip, method, path, status)
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

// todo: if this doesn't work, need to use `w.Header()`
func setGlobalHeaders(r *http.Request) {
	// Respond with global HTTP headers
	header := r.Response.Header
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
func defineRoutes() http.ServeMux {
	router := http.ServeMux{}

	router.HandleFunc("/favicon.ico", serveFile("./assets/img/favicon.ico"))
	router.HandleFunc("/robots.txt", serveFile("./assets/robots.txt"))
	router.Handle("/img/", http.FileServer(http.Dir("./assets/img")))
	router.Handle("/css/", http.FileServer(http.Dir("./assets/css")))
	router.Handle("/js/", http.FileServer(http.Dir("./assets/js")))

	// server.Use(recover.New(recover.Config{EnableStackTrace: config.GlobalServerConfig.InDevelopment}))

	// // Routes

	// server.Get("/", routes.IndexPage)
	// server.Get("/about", routes.AboutPage)
	// server.Get("/newest", routes.NewestPage)
	// server.Get("/discovery", routes.DiscoveryPage)
	// server.Get("/discovery/novel", routes.NovelDiscoveryPage)
	// server.Get("/ranking", routes.RankingPage)
	// server.Get("/rankingCalendar", routes.RankingCalendarPage)
	// server.Post("/rankingCalendar", routes.RankingCalendarPicker)
	// server.Get("/users/:id.atom.xml", routes.UserAtomFeed)
	// server.Get("/users/:id/:category.atom.xml", routes.UserAtomFeed)
	// server.Get("/users/:id/:category?", routes.UserPage)
	// server.Get("/artworks/:id/", routes.ArtworkPage).Name("artworks")
	// server.Get("/artworks-multi/:ids/", routes.ArtworkMultiPage)
	// server.Get("/novel/:id/", routes.NovelPage)
	// server.Get("/pixivision", routes.PixivisionHomePage)
	// server.Get("/pixivision/a/:id", routes.PixivisionArticlePage)

	// // Settings group
	// settings := server.Group("/settings")
	// settings.Get("/", routes.SettingsPage)
	// settings.Post("/:type/:noredirect?", routes.SettingsPost)

	// // Personal group
	// self := server.Group("/self")
	// self.Get("/", routes.LoginUserPage)
	// self.Get("/followingWorks", routes.FollowingWorksPage)
	// self.Get("/bookmarks", routes.LoginBookmarkPage)
	// self.Get("/addBookmark/:id", routes.AddBookmarkRoute)
	// self.Get("/deleteBookmark/:id", routes.DeleteBookmarkRoute)
	// self.Get("/like/:id", routes.LikeRoute)

	// // Oembed group
	// server.Get("/oembed", routes.Oembed)

	// server.Get("/tags/:name", routes.TagPage)
	// server.Post("/tags/:name", routes.TagPage)
	// server.Get("/tags", routes.TagPage)
	// server.Post("/tags", routes.AdvancedTagPost)

	// // Legacy illust URL
	// server.Get("/member_illust.php", func(c *http.Request) error {
	// 	return c.Redirect("/artworks/" + c.Query("illust_id"))
	// })

	// // Proxy routes
	// proxy := server.Group("/proxy")
	// proxy.Get("/i.pximg.net/*", routes.IPximgProxy)
	// proxy.Get("/s.pximg.net/*", routes.SPximgProxy)
	// proxy.Get("/ugoira.com/*", routes.UgoiraProxy)

	return router
}

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/routes"
	"codeberg.org/vnpower/pixivfe/v2/session"
	"codeberg.org/vnpower/pixivfe/v2/utils/kmutex"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiber_utils "github.com/gofiber/fiber/v2/utils"
)

func CanRequestSkipLimiter(c *fiber.Ctx) bool {
	path := c.Path()
	return strings.HasPrefix(path, "/img/") ||
		strings.HasPrefix(path, "/css/") ||
		strings.HasPrefix(path, "/js/") ||
		strings.HasPrefix(path, "/proxy/s.pximg.net/")
}

func CanRequestSkipLogger(c *fiber.Ctx) bool {
	// return false
	path := c.Path()
	return CanRequestSkipLimiter(c) ||
		strings.HasPrefix(path, "/proxy/i.pximg.net/")
}

func main() {
	config.GlobalServerConfig.InitializeConfig()
	core.CreateResponseAuditFolder()

	routes.InitTemplatingEngine(config.GlobalServerConfig.InDevelopment)

	server := fiber.New(fiber.Config{
		AppName:                 "PixivFE",
		DisableStartupMessage:   true,
		Prefork:                 false,
		JSONEncoder:             json.Marshal,
		JSONDecoder:             json.Unmarshal,
		ViewsLayout:             "_layout",
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"0.0.0.0/0"},
		ProxyHeader:             fiber.HeaderXForwardedFor,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Println(err)

			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// // Retrieve the custom status code if it's a *fiber.Error
			// var e *fiber.Error
			// if errors.As(err, &e) {
			// 	code = e.Code
			// }

			// Send custom error page
			c.Status(code)
			err = routes.Render(c, routes.Data_error{Title: "Error", Error: err})
			if err != nil {
				return c.Status(code).SendString(fmt.Sprintf("Internal Server Error: %s", err))
			}

			return nil
		},
	})

	server.Use(func(c *fiber.Ctx) error {
		// Pass in values that we want to be available to all pages here
		token := session.GetPixivToken(c)
		pageURL := c.BaseURL() + c.OriginalURL()

		cookies := map[string]string{}
		for _, name := range session.AllCookieNames {
			value := session.GetCookie(c, name)
			cookies[string(name)] = value
		}

		c.Bind(fiber.Map{
			"BaseURL":     c.BaseURL(),
			"OriginalURL": c.OriginalURL(),
			"PageURL":     pageURL,
			"LoggedIn":    token != "",
			"Queries":     c.Queries(),
			"CookieList":  cookies,
		})
		return c.Next()
	})

	if config.GlobalServerConfig.RequestLimit > 0 {
		keyedSleepingSpot := kmutex.New()
		server.Use(limiter.New(limiter.Config{
			Next:              CanRequestSkipLimiter,
			Expiration:        30 * time.Second,
			Max:               config.GlobalServerConfig.RequestLimit,
			LimiterMiddleware: limiter.SlidingWindow{},
			LimitReached: func(c *fiber.Ctx) error {
				// limit response throughput by pacing, since not every bot reads X-RateLimit-*
				// on limit reached, they just have to wait
				// the design of this means that if they send multiple requests when reaching rate limit, they will wait even longer (since `retryAfter` is calculated before anything has slept)
				retryAfter_s := c.GetRespHeader(fiber.HeaderRetryAfter)
				retryAfter, err := strconv.ParseUint(retryAfter_s, 10, 64)
				if err != nil {
					log.Panicf("response header 'RetryAfter' should be a number: %v", err)
				}
				requestIP := c.IP()
				refcount := keyedSleepingSpot.Lock(requestIP)
				defer keyedSleepingSpot.Unlock(requestIP)
				if refcount >= 4 { // on too much concurrent requests
					// todo: maybe blackhole `requestIP` here
					log.Println("Limit Reached (Hard)!", requestIP)
					// close the connection immediately
					_ = c.Context().Conn().Close()
					return nil
				}

				// sleeping
				// here, sleeping is not the best solution.
				// todo: close this connection when this IP reaches hard limit
				dur := time.Duration(retryAfter) * time.Second
				log.Println("Limit Reached (Soft)! Sleeping for ", dur)
				ctx, cancel := context.WithTimeout(c.Context(), dur)
				defer cancel()
				<-ctx.Done()

				return c.Next()
			},
		}))
	}

	server.Use(logger.New(
		logger.Config{
			Format: "${time} +${latency} ${ip} ${method} ${path} ${status} ${error} \n",
			Next:   CanRequestSkipLogger,
			CustomTags: map[string]logger.LogFunc{
				// make latency always print in seconds
				logger.TagLatency: func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
					latency := data.Stop.Sub(data.Start).Seconds()
					return output.WriteString(fmt.Sprintf("%.6f", latency))
				},
			},
		},
	))

	server.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

	if !config.GlobalServerConfig.InDevelopment {
		server.Use(cache.New(
			cache.Config{
				Next: func(c *fiber.Ctx) bool {
					resp_code := c.Response().StatusCode()
					if resp_code < 200 || resp_code >= 300 {
						return true
					}

					// Disable cache for settings page
					return strings.Contains(c.Path(), "/settings") || c.Path() == "/"
				},
				Expiration:           5 * time.Minute,
				CacheControl:         true,
				StoreResponseHeaders: true,

				KeyGenerator: func(c *fiber.Ctx) string {
					key := fiber_utils.CopyString(c.OriginalURL())
					for _, cookieName := range session.AllCookieNames {
						cookieValue := session.GetCookie(c, cookieName)
						if cookieValue != "" {
							key += "\x00\x00"
							key += string(cookieName)
							key += "\x00"
							key += cookieValue
						}
					}
					return key
				},
			},
		))
	}

	// redirect any round with ?r=url
	// could this be unsafe with cross-site scripting?
	server.Use(func(c *fiber.Ctx) error {
		ret := c.Query("r")
		if ret != "" {
			c.Redirect(ret)
		}
		return c.Next()
	})

	// Global HTTP headers
	server.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			return err
		}
		if strings.HasPrefix(string(c.Response().Header.ContentType()), "text/html") {
			c.Set("X-Frame-Options", "DENY")
			// use this if need iframe: `X-Frame-Options: SAMEORIGIN`
			c.Set("X-Content-Type-Options", "nosniff")
			c.Set("Referrer-Policy", "no-referrer")
			c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
			c.Set("Content-Security-Policy", fmt.Sprintf("base-uri 'self'; default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self' %s; media-src 'self' %s; connect-src 'self'; form-action 'self'; frame-ancestors 'none';", session.GetImageProxyOrigin(c), session.GetImageProxyOrigin(c)))
			// use this if need iframe: `frame-ancestors 'self'`
			c.Set("Permissions-Policy", "accelerometer=(), ambient-light-sensor=(), battery=(), camera=(), display-capture=(), document-domain=(), encrypted-media=(), execution-while-not-rendered=(), execution-while-out-of-viewport=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), midi=(), navigation-override=(), payment=(), publickey-credentials-get=(), screen-wake-lock=(), sync-xhr=(), usb=(), web-share=(), xr-spatial-tracking=()")
		}
		return nil
	})

	server.Static("/favicon.ico", "./assets/img/favicon.ico")
	server.Static("/robots.txt", "./assets/robots.txt")
	server.Static("/img/", "./assets/img")
	server.Static("/css/", "./assets/css")
	server.Static("/js/", "./assets/js")

	server.Use(recover.New(recover.Config{EnableStackTrace: config.GlobalServerConfig.InDevelopment}))

	// Routes

	server.Get("/", routes.IndexPage)
	server.Get("/about", routes.AboutPage)
	server.Get("/newest", routes.NewestPage)
	server.Get("/discovery", routes.DiscoveryPage)
	server.Get("/discovery/novel", routes.NovelDiscoveryPage)
	server.Get("/ranking", routes.RankingPage)
	server.Get("/rankingCalendar", routes.RankingCalendarPage)
	server.Post("/rankingCalendar", routes.RankingCalendarPicker)
	server.Get("/users/:id.atom.xml", routes.UserAtomFeed)
	server.Get("/users/:id/:category?.atom.xml", routes.UserAtomFeed)
	server.Get("/users/:id/:category?", routes.UserPage)
	server.Get("/artworks/:id/", routes.ArtworkPage).Name("artworks")
	server.Get("/artworks-multi/:ids/", routes.ArtworkMultiPage)
	server.Get("/novel/:id/", routes.NovelPage)
	server.Get("/pixivision", routes.PixivisionHomePage)
	server.Get("/pixivision/a/:id", routes.PixivisionArticlePage)

	// Settings group
	settings := server.Group("/settings")
	settings.Get("/", routes.SettingsPage)
	settings.Post("/:type/:noredirect?", routes.SettingsPost)

	// Personal group
	self := server.Group("/self")
	self.Get("/", routes.LoginUserPage)
	self.Get("/followingWorks", routes.FollowingWorksPage)
	self.Get("/bookmarks", routes.LoginBookmarkPage)
	self.Get("/addBookmark/:id", routes.AddBookmarkRoute)
	self.Get("/deleteBookmark/:id", routes.DeleteBookmarkRoute)
	self.Get("/like/:id", routes.LikeRoute)

	// Oembed group
	server.Get("/oembed", routes.Oembed)

	server.Get("/tags/:name", routes.TagPage)
	server.Post("/tags/:name", routes.TagPage)
	server.Get("/tags", routes.TagPage)
	server.Post("/tags", routes.AdvancedTagPost)

	// Legacy illust URL
	server.Get("/member_illust.php", func(c *fiber.Ctx) error {
		return c.Redirect("/artworks/" + c.Query("illust_id"))
	})

	// Proxy routes
	proxy := server.Group("/proxy")
	proxy.Get("/i.pximg.net/*", routes.IPximgProxy)
	proxy.Get("/s.pximg.net/*", routes.SPximgProxy)
	proxy.Get("/ugoira.com/*", routes.UgoiraProxy)

	// Initialize and start the proxy checker
	ctx_timeout, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
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
	if config.GlobalServerConfig.UnixSocket != "" {
		ln, err := net.Listen("unix", config.GlobalServerConfig.UnixSocket)
		if err != nil {
			panic(err)
		}
		log.Printf("Listening on domain socket %v\n", config.GlobalServerConfig.UnixSocket)
		err = server.Listener(ln)
		if err != nil {
			panic(err)
		}
	} else {
		addr := config.GlobalServerConfig.Host + ":" + config.GlobalServerConfig.Port
		ln, err := net.Listen(server.Config().Network, addr)
		if err != nil {
			log.Panicf("failed to listen: %v", err)
		}
		addr = ln.Addr().String()
		log.Printf("Listening on http://%v/\n", addr)
		err = server.Listener(ln)
		if err != nil {
			panic(err)
		}
	}
}

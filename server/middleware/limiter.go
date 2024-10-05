package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/sethvargo/go-limiter/memorystore"

	"codeberg.org/vnpower/pixivfe/v2/config"
)

// CanRequestSkipLimiter determines if a request should bypass the rate limiter.
// It exempts static assets and proxied image requests from rate limiting.
func CanRequestSkipLimiter(r *http.Request) bool {
	path := r.URL.Path
	return strings.HasPrefix(path, "/img/") ||
		strings.HasPrefix(path, "/css/") ||
		strings.HasPrefix(path, "/js/") ||
		strings.HasPrefix(path, "/proxy/s.pximg.net/")
}

// NewIPRateLimiter creates a new instance of IPRateLimiter with the specified rate limit and burst.
//
// ## Arguments
//
// Tokens: Number of tokens allowed per interval.
// Interval: Interval until tokens reset.
func NewIPRateLimiter(Tokens uint64, Interval time.Duration) (*httplimit.Middleware, error) {
	store, err := memorystore.New(&memorystore.Config{
		Tokens:   Tokens,
		Interval: Interval,
	})
	if err != nil {
		return nil, err
	}
	return httplimit.NewMiddleware(store, httplimit.IPKeyFunc("X-Forwarded-For"))
}

// Global rate limiter instance
var limiter *httplimit.Middleware

// InitializeRateLimiter sets up the global rate limiter based on the application's configuration.
// If the request limit is less than 1, it sets an infinite rate limit.
//
// Returns the rate limit middleware
func InitializeRateLimiter() func(http.Handler) http.Handler {
	if config.GlobalConfig.RequestLimit < 1 {
		limiter = nil
	} else {
		var err error
		limiter, err = NewIPRateLimiter(config.GlobalConfig.RequestLimit, 30*time.Second)
		if err != nil {
			log.Panic(err)
		}
	}
	return rateLimitRequest
}

// RateLimitRequest is a middleware that applies rate limiting to incoming HTTP requests.
// It exempts certain requests (as defined by CanRequestSkipLimiter) from rate limiting.
func rateLimitRequest(h http.Handler) http.Handler {
	if limiter == nil {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if CanRequestSkipLimiter(r) {
			h.ServeHTTP(w, r)
			return
		}
		limiter.Handle(h).ServeHTTP(w, r)
	})
}

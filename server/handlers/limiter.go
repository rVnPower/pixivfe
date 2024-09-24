package handlers

import (
	"errors"
	"math"
	"net"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/time/rate"

	"codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/server/routes"
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

// IPRateLimiter manages rate limiting on a per-IP basis.
//
// TODO: Should we put middlewares in a separate file?
type IPRateLimiter struct {
	ips     map[string]*rate.Limiter // Maps IP addresses to their respective rate limiters
	mu      *sync.RWMutex            // Ensures thread-safe access to the map
	limiter *rate.Limiter            // Global rate limiter used as a template for per-IP limiters
}

// NewIPRateLimiter creates a new instance of IPRateLimiter with the specified rate limit and burst.
func NewIPRateLimiter(r rate.Limit, burst int) *IPRateLimiter {
	return &IPRateLimiter{
		ips:     make(map[string]*rate.Limiter),
		mu:      &sync.RWMutex{},
		limiter: rate.NewLimiter(r, burst),
	}
}

// Allow checks if a request from the given IP is allowed based on the rate limit.
// If the IP doesn't have a limiter, a new one is created.
func (lim *IPRateLimiter) Allow(ip string) bool {
	lim.mu.RLock()
	rl, exists := lim.ips[ip]
	lim.mu.RUnlock()

	if !exists {
		lim.mu.Lock()
		rl, exists = lim.ips[ip]
		if !exists {
			// Create a new limiter for this IP using the global limiter's settings
			rl = rate.NewLimiter(lim.limiter.Limit(), lim.limiter.Burst())
			lim.ips[ip] = rl
		}
		lim.mu.Unlock()
	}

	return rl.Allow()
}

// Global rate limiter instance
var limiter *IPRateLimiter

// InitializeRateLimiter sets up the global rate limiter based on the application's configuration.
// If the request limit is less than 1, it sets an infinite rate limit.
func InitializeRateLimiter() {
	r := float64(config.GlobalConfig.RequestLimit) / 30.0
	if config.GlobalConfig.RequestLimit < 1 {
		r = math.Inf(1)
	}
	limiter = NewIPRateLimiter(rate.Limit(r), 3)
}

// RateLimitRequest is a middleware that applies rate limiting to incoming HTTP requests.
// It exempts certain requests (as defined by CanRequestSkipLimiter) from rate limiting.
func RateLimitRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		if CanRequestSkipLimiter(r) {
			h.ServeHTTP(w, r)
			return
		}

		if !limiter.Allow(ip) {
			// If the request exceeds the rate limit, return an HTTP 429 Too Many Requests error
			routes.ErrorPage(w, r, errors.New("Too many requests"), http.StatusTooManyRequests)
		} else {
			// If the request is within the rate limit, proceed to the next handler
			h.ServeHTTP(w, r)
		}
	})
}

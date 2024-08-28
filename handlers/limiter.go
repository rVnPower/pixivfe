package handlers

import (
	"errors"
	"net"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/time/rate"
)

func CanRequestSkipLimiter(r *http.Request) bool {
	path := r.URL.Path
	return strings.HasPrefix(path, "/img/") ||
		strings.HasPrefix(path, "/css/") ||
		strings.HasPrefix(path, "/js/") ||
		strings.HasPrefix(path, "/proxy/s.pximg.net/")
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

func RateLimitRequest(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		if CanRequestSkipLimiter(r) {
			handler(w, r)
			return
		}

		if !limiter.Allow(ip) {
			CatchError(func(w http.ResponseWriter, r *http.Request) error {
				GetUserContext(r).StatusCode = http.StatusTooManyRequests
				return errors.New("Too many requests")
			})(w, r)
		} else {
			handler(w, r)
		}
	}
}

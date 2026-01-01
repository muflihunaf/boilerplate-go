package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/muflihunaf/boilerplate-go/pkg/response"
)

// RateLimiter implements a simple in-memory rate limiter.
type RateLimiter struct {
	requests map[string]*clientRequests
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

type clientRequests struct {
	count    int
	lastSeen time.Time
}

// NewRateLimiter creates a rate limiter with the given limit per window.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*clientRequests),
		limit:    limit,
		window:   window,
	}

	// Cleanup old entries periodically
	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	for range ticker.C {
		rl.mu.Lock()
		for ip, client := range rl.requests {
			if time.Since(client.lastSeen) > rl.window {
				delete(rl.requests, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Limit returns a middleware that limits requests per IP.
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		rl.mu.Lock()
		client, exists := rl.requests[ip]
		if !exists {
			rl.requests[ip] = &clientRequests{count: 1, lastSeen: time.Now()}
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		// Reset if window has passed
		if time.Since(client.lastSeen) > rl.window {
			client.count = 1
			client.lastSeen = time.Now()
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		// Check limit
		if client.count >= rl.limit {
			rl.mu.Unlock()
			w.Header().Set("Retry-After", rl.window.String())
			response.Error(w, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "Too many requests, please try again later")
			return
		}

		client.count++
		client.lastSeen = time.Now()
		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}


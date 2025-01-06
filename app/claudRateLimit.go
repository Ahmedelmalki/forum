package forum

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
	}
}

func (rl *RateLimiter) Limit(requests int, duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			rl.mu.Lock()
			now := time.Now()

			// Clean old requests
			if times, exists := rl.requests[ip]; exists {
				var valid []time.Time
				for _, t := range times {
					if now.Sub(t) <= duration {
						valid = append(valid, t)
					}
				}
				rl.requests[ip] = valid
			}

			// Check rate limit
			if len(rl.requests[ip]) >= requests {
				rl.mu.Unlock()
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			// Add new request
			rl.requests[ip] = append(rl.requests[ip], now)
			rl.mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}

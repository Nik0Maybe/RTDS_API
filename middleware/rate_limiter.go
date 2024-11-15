package middleware

import (
	"net/http"
	"sync"

	"RTDS_API/config"

	"golang.org/x/time/rate"
)

type Limiter struct {
	clients map[string]*rate.Limiter
	mu      sync.Mutex
}

func NewLimiter() *Limiter {
	return &Limiter{clients: make(map[string]*rate.Limiter)}
}

func (l *Limiter) GetLimiter(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	limiter, exists := l.clients[ip]
	if !exists {
		limiter = rate.NewLimiter(rate.Limit(config.RateLimitRequestsPerSecond), config.RateLimitBurst)
		l.clients[ip] = limiter
	}
	return limiter
}

func RateLimitMiddleware(limiter *Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if !limiter.GetLimiter(ip).Allow() {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

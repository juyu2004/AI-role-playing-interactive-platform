package middleware

import (
	"net/http"
	"time"
)

// Simple fixed-window rate limiter: N requests per window across the process.
func RateLimit(maxPerSecond int) func(http.Handler) http.Handler {
	if maxPerSecond <= 0 { maxPerSecond = 50 }
	// Token ticker
	tokens := make(chan struct{}, maxPerSecond)
	for i := 0; i < maxPerSecond; i++ { tokens <- struct{}{} }
	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			// refill
			for i := 0; i < maxPerSecond-len(tokens); i++ {
				select { case tokens <- struct{}{}: default: }
			}
		}
	}()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			select {
			case <-tokens:
				next.ServeHTTP(w, r)
			default:
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte("rate limit exceeded"))
			}
		})
	}
}

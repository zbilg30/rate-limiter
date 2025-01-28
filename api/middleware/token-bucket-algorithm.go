package middleware

import (
	"net/http"
	"sync"
	"time"
)

type TokenBucket struct {
	capacity   int
	tokens     int
	rate       time.Duration
	lastRefill time.Time
	mu         sync.Mutex
}

func (tb *TokenBucket) allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)

	if elapsed >= tb.rate {
		tb.tokens = tb.capacity
		tb.lastRefill = now
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

func newTokenBucket(capacity int, rate time.Duration) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		rate:       rate,
		lastRefill: time.Now(),
	}
}

func TokenBucketAlgorithmMiddleware(next http.Handler) http.Handler {
	tb := newTokenBucket(5, time.Minute)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !tb.allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

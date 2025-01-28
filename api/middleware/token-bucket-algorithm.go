/*
	Token Bucket Algorithm

	1. A token bucket is a container that has pre-defined capacity. Tokens are put in the bucket at preset rates periodically. Once the bucket is full, no more tokens are added. As shown in Figure 4, the token bucket capacity is 4. The refiller puts 2 tokens into the bucket every second. Once the bucket is full, extra tokens will overflow.

	2. Each request consumes one token. When a request arrives, we check if there are enough tokens in the bucket. Figure 5 explains how it works.

	3. If there are enough tokens, we take one token out for each request, and the request goes through.

	4. If there are not enough tokens, the request is dropped.
*/

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
	tb := newTokenBucket(5, 5*time.Second)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !tb.allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

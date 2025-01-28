package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Request struct {
	Path   string
	Method string
}

type LeakyBucket struct {
	// this represents our bucket - any new request will be sent in this queue
	queue []Request

	// this is the capacity of our bucket.
	// We will throttle requests if the bucket is full
	capacity int

	// this represents how often we will take requests
	// from our bucket and send to servers
	emptyRate time.Duration

	//we have this variable to signal in case we want to shut down
	stopRefiller chan struct{}

	//this is to handle data race conditions
	mu sync.Mutex
}

func NewLeakyBucket(capacity int, emptyRate time.Duration) *LeakyBucket {
	lb := &LeakyBucket{
		capacity:     capacity,
		emptyRate:    emptyRate,
		stopRefiller: make(chan struct{}),
	}

	go lb.removeRequestsFromQueue()
	return lb
}

func (lb *LeakyBucket) removeRequestsFromQueue() {
	ticker := time.NewTicker(lb.emptyRate)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			lb.mu.Lock()
			if len(lb.queue) > 0 {
				fmt.Printf("Serving %d requests\n", len(lb.queue))
				lb.queue = nil
			}
			lb.mu.Unlock()
		case <-lb.stopRefiller:
			return
		}
	}
}

func (lb *LeakyBucket) StopRefiller() {
	close(lb.stopRefiller)
}

func (lb *LeakyBucket) AddRequest(req Request) bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if len(lb.queue) < lb.capacity {
		lb.queue = append(lb.queue, req)
		return true
	}
	return false
}

func LeakingBucketAlgorithmMiddleware(next http.Handler) http.Handler {
	lb := NewLeakyBucket(3, 1*time.Second)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received", r.URL.Path, r.Method)
		if !lb.AddRequest(Request{Path: r.URL.Path, Method: r.Method}) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Request accepted and will be processed shortly."))

		next.ServeHTTP(w, r)
	})
}

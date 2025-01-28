package main

import (
	"net/http"
	"rate-limiter/api"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/v1/", api.RegisterV1Routes())

	http.ListenAndServe(":8080", mux)

	// We create a bucket with capacity 3 and requests
	// get sent to the server every 1 second.
	// lb := middleware.NewLeakyBucket(3, time.Second)

	// for i := 1; i <= 20; i++ {
	// 	// make a new request
	// 	req := middleware.Request{ID: i}
	// 	if lb.AddRequest(req) {
	// 		// this request was added to the bucket
	// 		fmt.Printf("Request %d added to the queue\n", i)
	// 	} else {
	// 		// this request couldn't be added
	// 		fmt.Printf("Request %d throttled\n", i)
	// 	}
	// 	// some time between requests
	// 	time.Sleep(300 * time.Millisecond)
	// }
	// lb.StopRefiller()
}

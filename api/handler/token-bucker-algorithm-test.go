package handler

import "net/http"

func TokenBucketAlgorithmTestHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World! Token Bucket Algorithm"))
}

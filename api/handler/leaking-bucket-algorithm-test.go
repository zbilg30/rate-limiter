package handler

import (
	"fmt"
	"net/http"
)

func LeakingBucketAlgorithmTestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request received", r.URL.Path, r.Method)
}

package api

import (
	"net/http"
	"rate-limiter/api/handler"
	"rate-limiter/api/middleware"
)

func RegisterV1Routes() http.Handler {
	v1Mux := http.NewServeMux()

	v1Mux.Handle("/token-bucket-algorithm-route", middleware.TokenBucketAlgorithmMiddleware(http.HandlerFunc(handler.TokenBucketAlgorithmTestHandler)))

	v1Mux.Handle("/leaking-bucket-algorithm-route", middleware.LeakingBucketAlgorithmMiddleware(http.HandlerFunc(handler.LeakingBucketAlgorithmTestHandler)))

	return http.StripPrefix("/v1", v1Mux)
}

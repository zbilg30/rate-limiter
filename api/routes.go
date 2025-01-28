package api

import (
	"net/http"
	"rate-limiter/api/handler"
	"rate-limiter/api/middleware"
)

func RegisterV1Routes() http.Handler {
	v1Mux := http.NewServeMux()

	v1Mux.Handle("/token-bucket-algorithm-route", middleware.TokenBucketAlgorithmMiddleware(http.HandlerFunc(handler.TestHandler)))

	return http.StripPrefix("/v1", v1Mux)
}

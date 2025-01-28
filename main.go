package main

import (
	"net/http"
	"rate-limiter/api"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/v1/", api.RegisterV1Routes())

	http.ListenAndServe(":8080", mux)
}

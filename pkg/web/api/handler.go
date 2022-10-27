package api

import (
	"fmt"
	"log"
	"net/http"
)

type Middleware interface {
	func(next http.Handler) http.Handler
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "[NOT FOUND] called")
}

func Options(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "[OPTIONS] called")
}

func BasicLogger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Printf("method=%q, path=%q, params=%v\n", r.Method, r.RequestURI, r.URL.Query())
	}
	return http.HandlerFunc(fn)
}

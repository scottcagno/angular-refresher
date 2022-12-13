package middleware

import (
	"fmt"
	"log"
	"net/http"
)

type Middleware = func(next http.Handler) http.Handler

func NotFound(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "[NOT FOUND] called")
}

func Options(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "[OPTIONS] called")
}

// func _Secure(conf *BasicAuthConfig, next http.Handler) http.Handler {
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		_, _, isAuthenticated := r.BasicAuth()
// 		if isAuthenticated {
// 			next.ServeHTTP(w, r)
// 			return
// 		}
// 		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
// 		return
// 	}
// 	return http.HandlerFunc(fn)
// }

func _WithLogging2(logger *log.Logger, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		logger.Printf("HttpReq <- [method=%q, path=%q, params=%v]\n", r.Method, r.RequestURI, r.URL.Query())
		logger.Printf("HttpRes -> [method=%q, path=%q, params=%v]\n", r.Method, r.RequestURI, r.URL.Query())
	}
	return http.HandlerFunc(fn)
}

func _HandleCORS(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// func AddContextToRequestIn()
//
// func CheckAuth(token string) Middleware {
// 	return func(h http.Handler) http.Handler {
// 		return http.HandlerFunc(
// 			func(w http.ResponseWriter, r *http.Request) {
// 				if r.Header.Get("Auth") != authToken {
// 					util.SendError(w, "...", http.StatusForbidden, false)
// 					return
// 				}
// 				h.ServeHTTP(w, r)
// 			},
// 		)
// 	}
// }
//
// func YourMiddleware(h http.HandlerFunc) http.HandlerFunc {
// 	// do the thing you need
// 	// get the value
//
// 	// store it in the context, below
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		context.Set(r, mwKey, "whatyouwanttostore")
//
// 		h.ServeHTTP(w, r)
// 	}
//
// 	return fn
// }
//
// func YourHandler(w http.ResponseWriter, r *http.Request) {
// 	val, ok := context.Get(r, mwKey)
// 	// Type assert - we'll assume it's a string we want
// 	if !val.(string) {
// 		// Handle the case where the "wrong" type has
// 		// been stored in our context against the
// 		// expected key. HTTP 500 recommended here.
// 	}
// }
//
// func YourMiddleware(h http.HandlerFunc) http.HandlerFunc {
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		// do the thing you need
// 		// get the value
// 		// doing this inside fn ensures this is one every time a request comes in for the handler `h`
//
// 		// store it in the context, below
// 		context.Set(r, mwKey, "whatyouwanttostore")
//
// 		h.ServeHTTP(w, r)
// 	}
//
// 	return fn
// }
//
//
// On a side note, `context.Get` might have a different signature. If I recall right, it only returns one value (an interface{}):
//
// val := context.Get(r, mwKey)

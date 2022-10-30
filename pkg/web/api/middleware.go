package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Middleware = func(next http.Handler) http.Handler

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

type CORSConfig struct {
	// AllowOrigin defines a list of origins that may access the resource.
	//
	// Optional. Default value "*"
	// Header definition: Access-Control-Allow-Origin: <origin> | *
	AllowOrigins string

	// AllowMethods defines a list methods allowed when accessing the resource.
	// This is used in response to a preflight request.
	//
	// Optional. Default value "GET,POST,HEAD,PUT,DELETE,PATCH"
	// Header definition: Access-Control-Allow-Methods: <method>[, <method>]*
	AllowMethods string

	// AllowHeaders defines a list of request headers that can be used when
	// making the actual request. This is in response to a preflight request.
	//
	// Optional. Default value "".
	// Header definition: Access-Control-Allow-Headers: <header-name>[, <header-name>]*
	AllowHeaders string

	// AllowCredentials indicates whether the response to the request can
	// be exposed when the credentials flag is true. When used as part of
	// a response to a preflight request, this indicates whether the actual
	// request can be made using credentials.
	//
	// Optional. Default value false.
	// Header definition: Access-Control-Allow-Credentials: true|false
	AllowCredentials bool

	// ExposeHeaders defines a whitelist headers that clients are allowed to
	// access.
	//
	// Optional. Default value "".
	// Header definition: Access-Control-Expose-Headers: <header-name>[, <header-name>]*
	ExposeHeaders string

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached.
	//
	// Optional. Default value 0.
	// Header definition: Access-Control-Max-Age: <delta-seconds>
	MaxAge int
}

var defaultCORSConfig = &CORSConfig{
	AllowOrigins:     "*",
	AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
	AllowHeaders:     "",
	AllowCredentials: false,
	ExposeHeaders:    "",
	MaxAge:           0,
}

func CORSHandler(c *CORSConfig) http.Handler {
	if c == nil {
		c = defaultCORSConfig
		goto skipSetDefaults
	}
	if c.AllowOrigins != "" {
		c.AllowOrigins = defaultCORSConfig.AllowOrigins
	}
	if c.AllowMethods != "" {
		c.AllowMethods = defaultCORSConfig.AllowMethods
	}
	if c.AllowHeaders != "" {
		c.AllowHeaders = defaultCORSConfig.AllowHeaders
	}
	if !c.AllowCredentials {
		c.AllowCredentials = defaultCORSConfig.AllowCredentials
	}
	if c.ExposeHeaders != "" {
		c.ExposeHeaders = defaultCORSConfig.ExposeHeaders
	}
	if c.MaxAge == 0 {
		c.ExposeHeaders = defaultCORSConfig.ExposeHeaders
	}
skipSetDefaults:
	fn := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		// handle a pre-flight request
		case http.MethodOptions:
			w.Header().Add(HeaderVary, HeaderOrigin)
			w.Header().Add(HeaderVary, HeaderAccessControlRequestMethod)
			w.Header().Add(HeaderVary, HeaderAccessControlRequestHeaders)
			w.Header().Set(HeaderAccessControlAllowOrigin, c.AllowOrigins)
			w.Header().Set(HeaderAccessControlAllowMethods, c.AllowMethods)
			if c.AllowCredentials {
				w.Header().Set(HeaderAccessControlAllowCredentials, "true")
			}
			if c.AllowHeaders != "" {
				w.Header().Set(HeaderAccessControlAllowHeaders, c.AllowHeaders)
			} else {
				if rh := r.Header.Get(HeaderAccessControlRequestHeaders); rh != "" {
					w.Header().Set(HeaderAccessControlAllowHeaders, rh)
				}
			}
			if c.MaxAge > 0 {
				w.Header().Set(HeaderAccessControlMaxAge, strconv.Itoa(c.MaxAge))
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			// simple, non pre-flight request
			w.Header().Add(HeaderVary, HeaderOrigin)
			w.Header().Set(HeaderAccessControlAllowOrigin, c.AllowOrigins)
			w.Header().Set(HeaderAccessControlAllowMethods, c.AllowMethods)
			if c.AllowCredentials {
				w.Header().Set(HeaderAccessControlAllowCredentials, "true")
			}
			if c.ExposeHeaders != "" {
				w.Header().Set(HeaderAccessControlExposeHeaders, c.ExposeHeaders)
			}
		}
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

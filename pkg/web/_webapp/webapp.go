package webapp

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// A RouteHandler responds to an HTTP request. ServeRoute should write reply headers
type RouteHandler interface {
	ServeRoute(c Context) error
}

// The RouteHandlerFunc type is an adapter to allow the use of ordinary functions
// as HTTP handlers. If f is a function with the appropriate signature,
// then using RouteHandlerFunc(f) acts as a RouteHandler that calls f.
type RouteHandlerFunc func(c Context) error

// ServeRoute calls f(c)
func (f RouteHandlerFunc) ServeRoute(c Context) error {
	return f(c)
}

// RouteHandlerToHandler takes a RouteHandler and returns an http.Handler
func RouteHandlerToHandler(handler RouteHandler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		err := handler.ServeRoute(newContext(w, r))
		if err != nil {
			http.Error(w, err.Error(), StatusCode(err.Error()))
			return
		}
	}
	return http.HandlerFunc(fn)
}

func CodeFromError(err error) int {
	return StatusCode(err.Error())
}

func ErrorFromCode(code int) error {
	return errors.New(http.StatusText(code))
}

// Router is an HTTP request multiplexer. It matches the URL of each incoming
// request against a list of registered patterns and calls the handler for the
// pattern that most closely matches the URL.
type Router struct {
	lock     sync.RWMutex
	routeMap map[string]routeEntry
	routeSet []routeEntry
}

// routeEntry acts as a single route entry in the router.
type routeEntry struct {
	method  string
	pattern string
	handler RouteHandler
}

// NewRouter allocates and returns a new Router
func NewRouter() *Router {
	return new(Router)
}

// DefaultRouter is the default Router used by Serve
var DefaultRouter = &defaultRouter
var defaultRouter Router

func newContext(w http.ResponseWriter, r *http.Request) *context {
	return &context{
		res:   newResponse(w),
		req:   r,
		path:  r.URL.Path,
		query: r.URL.Query(),
	}
}

func (rm *Router) newContext(w http.ResponseWriter, r *http.Request) *context {
	return &context{
		res:   newResponse(w),
		req:   r,
		path:  r.URL.Path,
		query: r.URL.Query(),
	}
}

func (rm *Router) Handle(method string, pattern string, handler RouteHandler) {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	if pattern == "" {
		panic("http: invalid pattern")
	}
	if handler == nil {
		panic("http: nil handler")
	}
	if _, exists := rm.routeMap[pattern]; exists {
		panic("http: multiple registrations for " + pattern)
	}
	if rm.routeMap == nil {
		rm.routeMap = make(map[string]routeEntry)
	}
	entry := routeEntry{
		method:  method,
		pattern: pattern,
		handler: handler,
	}
	rm.routeMap[pattern] = entry
	if pattern[len(pattern)-1] == '/' {
		rm.routeSet = append(rm.routeSet, entry)
	}
}

func (rm *Router) HandleFunc(method string, pattern string, handler RouteHandlerFunc) {
	if handler == nil {
		panic("http: nil handler")
	}
	rm.Handle(method, pattern, handler)
}

func (rm *Router) Forward(oldpattern string, newpattern string) {
	rm.Handle(http.MethodGet, oldpattern, RedirectHandler(http.StatusTemporaryRedirect, newpattern))
}

func (rm *Router) Get(pattern string, handler RouteHandler) {
	rm.Handle(http.MethodGet, pattern, handler)
}

func (rm *Router) Post(pattern string, handler RouteHandler) {
	rm.Handle(http.MethodPost, pattern, handler)
}

func (rm *Router) Put(pattern string, handler RouteHandler) {
	rm.Handle(http.MethodPut, pattern, handler)
}

func (rm *Router) Delete(pattern string, handler RouteHandler) {
	rm.Handle(http.MethodDelete, pattern, handler)
}

// func (rm *Router) Static(pattern string, path string) {
// 	staticHandler := func(c Context) error {
// 		http.StripPrefix(pattern, http.FileServer(http.Dir(path)))
// 		return nil
// 	}
// 	rm.Handle(http.MethodGet, pattern, staticHandler)
// }

func (rm *Router) entries() []string {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	var entries []string
	for _, entry := range rm.routeMap {
		entries = append(entries, fmt.Sprintf("%s %s\n", entry.method, entry.pattern))
	}
	return entries
}

// match attempts to locate a handler on a handler map given a
// path string; most-specific (longest) pattern wins
func (rm *Router) match(path string) (string, string, RouteHandler) {
	// first, check for exact match
	e, ok := rm.routeMap[path]
	if ok {
		return e.method, e.pattern, e.handler
	}
	// then, check for longest valid match. mux.es
	// contains all patterns that end in "/" sorted
	// from longest to shortest
	for _, e = range rm.routeSet {
		if strings.HasPrefix(path, e.pattern) {
			return e.method, e.pattern, e.handler
		}
	}
	return "", "", nil
}

func (rm *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var err error
	var c *context
	c.Reset(w, r)
	meth, _, h := rm.match(r.URL.Path)
	if meth != r.Method && meth != "*" {
		h = RouteHandlerFunc(MethodNotAllowed)
	}
	if h == nil {
		h = RouteHandlerFunc(NotFound)
	}
	c.handler = h
	err = h.ServeRoute(c)
	if err != nil {
		panic(err)
	}
}

func (rm *Router) foo2() {}

func (rm *Router) foo3() {}

func (rm *Router) foo4() {}

func (rm *Router) foo5() {}

func (rm *Router) foo6() {}

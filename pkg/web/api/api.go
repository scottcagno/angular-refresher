package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/scottcagno/angular-refresher/pkg/web/api/middleware"
)

type M = map[string]any

type API struct {
	base     string
	cors     http.Handler
	logger   *log.Logger
	mux      *http.ServeMux
	handlers []handler
}

func NewAPI(base string, cors http.Handler, logger *log.Logger, mux *http.ServeMux) *API {
	if logger == nil {
		logger = log.New(os.Stderr, "[DEFAULT] ", log.LstdFlags)
	}
	if mux == nil {
		mux = http.NewServeMux()
	}
	api := &API{
		base:     base,
		cors:     cors,
		logger:   logger,
		mux:      mux,
		handlers: make([]handler, 0),
	}
	api.mux.Handle("/", http.RedirectHandler(api.base, http.StatusSeeOther))
	api.mux.Handle(filepath.ToSlash(filepath.Join(api.base, "stats")), api.StatsHandler())
	return api
}

func (api *API) StatsHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, map[string]any{"routes": api.handlers})
	}
	return http.HandlerFunc(fn)
}

func (api *API) Register(name string, re Resource) {
	h := &handler{
		name: name,
		path: filepath.ToSlash(filepath.Join(api.base, name)),
		reso: re,
	}
	api.handlers = append(api.handlers, *h)
	api.mux.Handle(h.path, middleware.WithLogging(api.logger, h))
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// apply cors handler if we have one
	if api.cors != nil {
		api.cors.ServeHTTP(w, r)
	}
	// lookup resource handler
	rh, pat := api.mux.Handler(r)
	// do something with the pattern if we need to
	_ = pat
	// call the resource handler
	rh.ServeHTTP(w, r)
}

type handler struct {
	name string
	path string
	reso Resource
}

func (h handler) String() string {
	return fmt.Sprintf("name=%q, path=%q", h.name, h.path)
}

func writeReqCtx(r *http.Request, key, val any) *http.Request {
	// create a new context from the parent context in the incoming request
	ctx := context.WithValue(r.Context(), key, val)
	// create and return a new request using the new context
	return r.WithContext(ctx)
}

func readReqCtx(r *http.Request, key any) any {
	// get the value from the request context
	return r.Context().Value(key)
}

func GetParam(r *http.Request, key string) (string, bool) {
	params := r.URL.Query()
	if params.Has(key) {
		return params.Get(key), true
	}
	return "", false
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.reso.Get(w, r)
	case http.MethodPost:
		h.reso.Add(w, r)
	case http.MethodPut:
		h.reso.Set(w, r)
	case http.MethodDelete:
		h.reso.Del(w, r)
	case http.MethodOptions:
		middleware.Options(w, r)
	default:
		middleware.NotFound(w, r)
	}
}

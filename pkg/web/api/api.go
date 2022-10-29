package api

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
)

type M = map[string]any

type API struct {
	base     string
	mux      *http.ServeMux
	handlers []handler
}

func NewAPI(base string, mux *http.ServeMux) *API {
	if mux == nil {
		mux = http.NewServeMux()
	}
	api := &API{
		base:     base,
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
	api.mux.Handle(h.path, BasicLogger(h))
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		Options(w, r)
	default:
		NotFound(w, r)
	}
}

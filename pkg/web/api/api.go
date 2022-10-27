package api

import (
	"net/http"
	"path/filepath"
)

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
	return api
}

func (api *API) Register(name string, re Resource) {
	h := &handler{
		name: name,
		path: filepath.ToSlash(filepath.Join(api.base, name)),
		reso: re,
	}
	h.reso.Init()
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

// func (rh *handler) _ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	hasID := len(r.URL.Query()) > 0
// 	switch r.Method {
// 	case http.MethodGet:
// 		if hasID {
// 			rh.reso.GetOne(w, r)
// 			return
// 		}
// 		rh.reso.GetAll(w, r)
// 		return
// 	case http.MethodPost:
// 		rh.reso.AddOne(w, r)
// 		return
// 	case http.MethodPut:
// 		if hasID {
// 			rh.reso.SetOne(w, r)
// 			return
// 		}
// 	case http.MethodDelete:
// 		if hasID {
// 			rh.reso.DelOne(w, r)
// 		}
// 	default:
// 		http.NotFoundHandler().ServeHTTP(w, r)
// 		return
// 	}
// }

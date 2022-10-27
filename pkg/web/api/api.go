package api

import (
	"net/http"
)

type API struct {
	base string
	mux  *http.ServeMux
}

func NewAPI(base string, mux *http.ServeMux) *API {
	if mux == nil {
		mux = http.NewServeMux()
	}
	api := &API{
		base: base,
		mux:  mux,
	}
	api.mux.Handle("/", http.RedirectHandler(api.base, http.StatusSeeOther))
	return api
}

func (api *API) Register(name string, re Resource) {
	r := &resourceHandler{
		name: name,
		path: api.base + name,
		re:   re,
	}
	api.mux.Handle(r.path, r)
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// lookup resource handler
	rh, pat := api.mux.Handler(r)
	// do something with the pattern if we need to
	_ = pat
	// call the resource handler
	rh.ServeHTTP(w, r)
}

type resourceHandler struct {
	name string
	path string
	re   Resource
}

func (rh *resourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hasID := len(r.URL.Query()) > 0
	switch r.Method {
	case http.MethodGet:
		if hasID {
			rh.re.GetOne(w, r)
			return
		}
		rh.re.GetAll(w, r)
		return
	case http.MethodPost:
		rh.re.AddOne(w, r)
		return
	case http.MethodPut:
		if hasID {
			rh.re.SetOne(w, r)
			return
		}
	case http.MethodDelete:
		if hasID {
			rh.re.DelOne(w, r)
		}
	default:
		http.NotFoundHandler().ServeHTTP(w, r)
		return
	}
}

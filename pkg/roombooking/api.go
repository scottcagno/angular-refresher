package roombooking

import (
	"fmt"
	"net/http"
	"path"
)

const version = "v1"

type RestAPI struct {
	addr string
	base string
}

func NewRestAPI(base, addr string) *RestAPI {
	return &RestAPI{
		addr: addr,
		base: path.Join(base, version),
	}
}

func (api *RestAPI) Addr() string {
	return api.addr
}

func (api *RestAPI) BasePath() string {
	return api.base
}

func (api *RestAPI) Version() string {
	return version
}

func (api *RestAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Path == api.base {
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", api.base)
		return
	}
}

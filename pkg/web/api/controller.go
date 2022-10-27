package api

import (
	"net/http"
	"path/filepath"
)

type Controller struct {
	*API
	*http.ServeMux
}

func NewController(route string, resource Resource) *Controller {
	c := &Controller{
		API: NewAPI(filepath.Dir(route), nil),
	}
	c.ServeMux = c.mux
	c.API.Register(filepath.Base(route), resource)
	return c
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.API.ServeHTTP(w, r)
}

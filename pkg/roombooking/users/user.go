package users

import (
	"net/http"

	"github.com/scottcagno/angular-refresher/pkg/web/api"
)

type User struct {
	ID   string
	Name string
}

type Controller struct {
	users []User
}

func (c *Controller) Inject(s api.Service) {
	c.users = s.Get("UserRepo").([]User)
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	id, found := api.GetParam(r, "id")
	if !found {
		id = "none"
	}
	api.WriteJSON(w, http.StatusOK, api.M{"user controller": "get", "id": api.M{"id": id}})
}

func (c *Controller) Add(w http.ResponseWriter, r *http.Request) {
	api.WriteJSON(w, http.StatusOK, api.M{"user controller": "add"})
}

func (c *Controller) Set(w http.ResponseWriter, r *http.Request) {
	id, found := api.GetParam(r, "id")
	if !found {
		id = "none"
	}
	api.WriteJSON(w, http.StatusOK, api.M{"user controller": "set", "id": api.M{"id": id}})
}

func (c Controller) Del(w http.ResponseWriter, r *http.Request) {
	id, found := api.GetParam(r, "id")
	if !found {
		id = "none"
	}
	api.WriteJSON(w, http.StatusOK, api.M{"user controller": "del", "id": api.M{"id": id}})
}

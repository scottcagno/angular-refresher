package booking

import (
	"net/http"

	"github.com/scottcagno/angular-refresher/pkg/web/api"
)

type Controller struct {
	*BookingRepository
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	id, found := api.GetParam(r, "id")
	if !found {
		id = "none"
	}
	api.WriteJSON(w, http.StatusOK, api.M{"booking controller": "get", "id": api.M{"id": id}})
}

func (c *Controller) Add(w http.ResponseWriter, r *http.Request) {

	api.WriteJSON(w, http.StatusOK, api.M{"booking controller": "add"})
}

func (c *Controller) Set(w http.ResponseWriter, r *http.Request) {
	id, found := api.GetParam(r, "id")
	if !found {
		id = "none"
	}
	api.WriteJSON(w, http.StatusOK, api.M{"booking controller": "set", "id": api.M{"id": id}})
}

func (c *Controller) Del(w http.ResponseWriter, r *http.Request) {
	id, found := api.GetParam(r, "id")
	if !found {
		id = "none"
	}
	api.WriteJSON(w, http.StatusOK, api.M{"booking controller": "del", "id": api.M{"id": id}})
}

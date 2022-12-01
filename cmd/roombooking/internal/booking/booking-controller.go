package booking

import (
	"net/http"
	"strconv"

	"github.com/scottcagno/angular-refresher/pkg/web/api"
)

type Controller struct {
	*BookingRepository
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	param, found := api.GetParam(r, "date")
	if !found {
		// handle, get all
		users, err := c.Repository.Find(func(b *Booking) bool { return b != nil })
		if err != nil {
			api.WriteJSON(w, http.StatusExpectationFailed, err)
			return
		}
		api.WriteJSON(w, http.StatusOK, users)
		return
	}
	// handle get all by date
	// date, err := time.Parse("2006-01-02", param)
	// if err != nil {
	// 	api.WriteJSON(w, http.StatusExpectationFailed, err)
	// 	return
	// }
	user, err := c.Repository.Find(func(b *Booking) bool { return b.Date == param })
	if err != nil || len(user) != 1 {
		api.WriteJSON(w, http.StatusExpectationFailed, err)
		return
	}
	api.WriteJSON(w, http.StatusOK, user)
	return
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
	// get id we need to update
	param, found := api.GetParam(r, "id")
	if !found {
		api.WriteJSON(w, http.StatusNotFound, "error: id required, but not found")
		return
	}
	// locate booking using id
	id, err := strconv.Atoi(param)
	if err != nil {
		api.WriteJSON(w, http.StatusExpectationFailed, err)
		return
	}
	err = c.Repository.Delete(id)
	if err != nil {
		api.WriteJSON(w, http.StatusExpectationFailed, err)
		return
	}
	api.WriteJSON(w, http.StatusOK, nil)
}

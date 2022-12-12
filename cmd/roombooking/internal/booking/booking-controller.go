package booking

import (
	"log"
	"net/http"
	"strconv"

	"github.com/scottcagno/angular-refresher/pkg/web/api"
)

type Controller struct {
	*BookingRepository
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	if api.HasParam(r, "date") {
		log.Println("BookingsController with DATE called...")
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
		bookings, err := c.Repository.Find(func(b *Booking) bool { return b.Date == param })
		if err != nil || len(bookings) < 1 {
			api.WriteJSON(w, http.StatusExpectationFailed, err)
			return
		}
		api.WriteJSON(w, http.StatusOK, bookings)
		return
	}
	if api.HasParam(r, "id") {
		log.Println("BookingsController with ID called...")
		param, _ := api.GetParam(r, "id")
		id, err := strconv.Atoi(param)
		if err != nil {
			api.WriteJSON(w, http.StatusExpectationFailed, err)
			return
		}
		booking, err := c.Repository.FindOne(func(b *Booking) bool { return b.ID == id })
		if err != nil {
			api.WriteJSON(w, http.StatusExpectationFailed, err)
			return
		}
		// WE ONLY WANT TO RETURN ONE ITEM
		api.WriteJSON(w, http.StatusOK, booking)
		return
	}
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

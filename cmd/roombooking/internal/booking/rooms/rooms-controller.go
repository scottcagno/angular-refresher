package rooms

import (
	"net/http"

	"github.com/scottcagno/angular-refresher/pkg/web/api"
)

type Controller struct {
	*RoomRepository
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	id, found := api.GetParam(r, "id")
	if !found {
		// handle, get all
		rooms, err := c.Find(func(r *Room) bool { return r != nil })
		if err != nil {
			api.WriteJSON(w, http.StatusExpectationFailed, err)
			return
		}
		api.WriteJSON(w, http.StatusOK, rooms)
		return
	}
	// handle get one
	room, err := c.Find(func(r *Room) bool { return r.ID == id })
	if err != nil || len(room) != 1 {
		api.WriteJSON(w, http.StatusExpectationFailed, err)
		return
	}
	api.WriteJSON(w, http.StatusOK, room)
	return
}

func (c *Controller) Add(w http.ResponseWriter, r *http.Request) {
	var newRoom Room
	err := api.ReadJSON(r, &newRoom)
	if err != nil {
		api.WriteJSON(w, http.StatusExpectationFailed, err)
		return
	}
	err = c.Insert(newRoom.ID, &newRoom)
	if err != nil {
		api.WriteJSON(w, http.StatusExpectationFailed, err)
		return
	}
	api.WriteJSON(w, http.StatusCreated, newRoom)
}

func (c *Controller) Set(w http.ResponseWriter, r *http.Request) {
	// get id we need to update
	id, found := api.GetParam(r, "id")
	if !found {
		api.WriteJSON(w, http.StatusNotFound, "error: id required, but not found")
		return
	}
	// get the updated room
	var updatedRoom Room
	err := api.ReadJSON(r, &updatedRoom)
	if err != nil {
		api.WriteJSON(w, http.StatusExpectationFailed, err)
		return
	}
	err = c.Update(id, &updatedRoom)
	if err != nil {
		api.WriteJSON(w, http.StatusExpectationFailed, err)
		return
	}
	api.WriteJSON(w, http.StatusOK, updatedRoom)
}

func (c *Controller) Del(w http.ResponseWriter, r *http.Request) {
	// get id we need to update
	id, found := api.GetParam(r, "id")
	if !found {
		api.WriteJSON(w, http.StatusNotFound, "error: id required, but not found")
		return
	}
	// locate room using id
	err := c.Delete(id)
	if err != nil {
		api.WriteJSON(w, http.StatusExpectationFailed, err)
		return
	}
	api.WriteJSON(w, http.StatusOK, nil)
}

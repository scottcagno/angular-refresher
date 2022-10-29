package rooms

import (
	"net/http"

	"github.com/scottcagno/angular-refresher/pkg/web/api"
)

type Controller struct {
	repo RoomRepository
}

func (c *Controller) Inject(s api.Service) {
	c.repo = s.GetRepository("RoomRepo").(RoomRepository)
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	id, found := api.GetParam(r, "id")
	if !found {
		// handle, get all
		rooms, err := c.repo.FindAll()
		if err != nil {
			api.WriteJSON(w, http.StatusExpectationFailed, err)
			return
		}
		api.WriteJSON(w, http.StatusOK, rooms)
		return
	}
	// handle get one
	room, err := c.repo.FindOne(id)
	if err != nil {
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
	err = c.repo.Insert(newRoom)
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
	err = c.repo.Update(id, updatedRoom)
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
	err := c.repo.Delete(id)
	if err != nil {
		api.WriteJSON(w, http.StatusExpectationFailed, err)
		return
	}
	api.WriteJSON(w, http.StatusOK, nil)
}

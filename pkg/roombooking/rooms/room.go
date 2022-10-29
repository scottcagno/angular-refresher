package rooms

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/scottcagno/angular-refresher/pkg/web/api"
)

const (
	roomController = "RoomHandler"
)

type Room struct {
	ID    string
	Title string
	Time  time.Time
}

type Controller struct {
	rooms []Room
}

func (c *Controller) Inject(s api.Service) {
	c.rooms = s.Get("RoomRepo").([]Room)
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	id, found := api.GetParam(r, "id")
	if !found {
		api.WriteJSON(w, http.StatusOK, c.rooms)
		return
	}
	for _, room := range c.rooms {
		if room.ID == id {
			api.WriteJSON(w, http.StatusOK, room)
			return
		}
	}
}

func (c *Controller) Add(w http.ResponseWriter, r *http.Request) {
	var newRoom Room
	err := api.ReadJSON(r, &newRoom)
	if err != nil {
		api.WriteJSON(w, http.StatusExpectationFailed, err.Error())
		return
	}
	c.rooms = append(c.rooms, newRoom)
	api.WriteJSON(w, http.StatusCreated, c.rooms)
}

func (c *Controller) Set(w http.ResponseWriter, r *http.Request) {
	// get id we need to update
	id, found := api.GetParam(r, "id")
	if !found {
		api.WriteJSON(w, http.StatusNotFound, "error: id required, but not found")
		return
	}
	// locate room using id
	at := -1
	for i, room := range c.rooms {
		if room.ID == id {
			at = i
			break
		}
	}
	if at == -1 {
		api.WriteJSON(w, http.StatusNotFound, "error: matching room id not found")
		return
	}
	// get the updated room
	var updatedRoom Room
	err := json.NewDecoder(r.Body).Decode(&updatedRoom)
	if err != nil {
		api.WriteJSON(w, http.StatusInternalServerError, "error: decoding json from body")
		return
	}
	// update the room
	c.rooms[at] = updatedRoom
	api.WriteJSON(w, http.StatusOK, c.rooms)
}

func (c *Controller) Del(w http.ResponseWriter, r *http.Request) {
	// get id we need to update
	id, found := api.GetParam(r, "id")
	if !found {
		api.WriteJSON(w, http.StatusNotFound, "error: id required, but not found")
		return
	}
	// locate room using id
	at := -1
	for i, room := range c.rooms {
		if room.ID == id {
			at = i
			break
		}
	}
	if at == -1 {
		api.WriteJSON(w, http.StatusNotFound, "error: matching room id not found")
		return
	}
	// delete the found room
	if at < len(c.rooms)-1 {
		copy(c.rooms[at:], c.rooms[at+1:])
	}
	c.rooms[len(c.rooms)-1] = Room{} // or the zero value of T
	c.rooms = c.rooms[:len(c.rooms)-1]
	// return some response
	api.WriteJSON(w, http.StatusOK, c.rooms)
}

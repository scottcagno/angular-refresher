package rooms

import (
	"encoding/json"
	"fmt"
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

type Rooms struct {
	rooms []Room
}

func (c *Rooms) Init() {
	c.InitFakeData()
}

func (c *Rooms) Get(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("id") {
		api.AsJSON(w, c.rooms)
		// fmt.Fprintf(w, "[%s] Get()", roomController)
		return
	}
	id := r.URL.Query().Get("id")
	for _, room := range c.rooms {
		if room.ID == id {
			api.AsJSON(w, room)
			return
		}
	}
	fmt.Fprintf(w, "[%s] Get(id: %s)", roomController, id)
}

func (c *Rooms) Add(w http.ResponseWriter, r *http.Request) {
	var newRoom Room
	err := json.NewDecoder(r.Body).Decode(&newRoom)
	if err != nil {
		api.AsJSON(w, map[string]string{"error": "error decoding json from body"})
		return
	}
	c.rooms = append(c.rooms, newRoom)
	fmt.Fprintf(w, "[%s] Add()", roomController)
}

func (c *Rooms) Set(w http.ResponseWriter, r *http.Request) {
	// get id we need to update
	if !r.URL.Query().Has("id") {
		api.AsJSON(w, map[string]string{"error": "id required, but not found"})
		return
	}
	id := r.URL.Query().Get("id")
	// locate room using id
	at := -1
	for i, room := range c.rooms {
		if room.ID == id {
			at = i
			break
		}
	}
	if at == -1 {
		api.AsJSON(w, map[string]string{"error": "error matching room id not found"})
		return
	}
	// get the updated room
	var updatedRoom Room
	err := json.NewDecoder(r.Body).Decode(&updatedRoom)
	if err != nil {
		api.AsJSON(w, map[string]string{"error": "error decoding json from body"})
		return
	}
	// update the room
	c.rooms[at] = updatedRoom
	fmt.Fprintf(w, "[%s] Set(id: %s)", roomController, r.URL.Query().Get("id"))
}

func (c *Rooms) Del(w http.ResponseWriter, r *http.Request) {
	// get id we need to update
	if !r.URL.Query().Has("id") {
		api.AsJSON(w, map[string]string{"error": "id required, but not found"})
		return
	}
	id := r.URL.Query().Get("id")
	// locate room using id
	at := -1
	for i, room := range c.rooms {
		if room.ID == id {
			at = i
			break
		}
	}
	if at == -1 {
		api.AsJSON(w, map[string]string{"error": "error matching room id not found"})
		return
	}
	// delete the found room
	if at < len(c.rooms)-1 {
		copy(c.rooms[at:], c.rooms[at+1:])
	}
	c.rooms[len(c.rooms)-1] = Room{} // or the zero value of T
	c.rooms = c.rooms[:len(c.rooms)-1]

	fmt.Fprintf(w, "[%s] Del(id: %s)", roomController, r.URL.Query().Get("id"))
}

func (c *Rooms) InitFakeData() {
	room1 := Room{
		ID:    "1",
		Title: "Room number one",
		Time:  time.Now(),
	}
	room2 := Room{
		ID:    "2",
		Title: "Room number two",
		Time:  time.Now().Add(5 * time.Hour),
	}
	room3 := Room{
		ID:    "3",
		Title: "Room number three",
		Time:  time.Now().Add(2 * time.Hour),
	}
	c.rooms = append(c.rooms, room1)
	c.rooms = append(c.rooms, room2)
	c.rooms = append(c.rooms, room3)
}

type Params = map[string]string

func Route(method string, path string, r *http.Request, params ...string) bool {
	if len(params) > 0 {
		for _, param := range params {
			if !r.URL.Query().Has(param) {
				return false
			}
		}
	}
	return r.Method == method && r.URL.Path == path
}

func RestRoomsController() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if Route(http.MethodGet, "/api/rooms", r) {
			// handle
		}
		if Route(http.MethodGet, "/api/rooms", r, "id") {
			// handle
		}
		if Route(http.MethodGet, "/api/rooms", r) {
			// handle
		}
		if Route(http.MethodGet, "/api/rooms", r) {
			// handle
		}
		if Route(http.MethodGet, "/api/rooms", r) {
			// handle
		}
	}
	return http.HandlerFunc(fn)
}

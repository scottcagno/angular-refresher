package rooms

import (
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

type RoomController struct {
	rooms []Room
}

func (c *RoomController) GetAll(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "[%s] GetAll()", roomController)
}

func (c *RoomController) GetOne(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	for _, room := range c.rooms {
		if room.ID == id {
			api.AsJSON(w, room)
			return
		}
	}
	fmt.Fprintf(w, "[%s] GetOne(id: %s)", roomController, id)
}

func (c *RoomController) AddOne(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "[%s] AddOne()", roomController)
}

func (c *RoomController) SetOne(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "[%s] SetOne(id: %s)", roomController, r.URL.Query().Get("id"))
}

func (c *RoomController) DelOne(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "[%s] DelOne(id: %s)", roomController, r.URL.Query().Get("id"))
}

func (c *RoomController) InitFakeData() {
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

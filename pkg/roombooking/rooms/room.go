package rooms

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/scottcagno/angular-refresher/pkg/roombooking/api"
)

type Room struct {
	ID    int
	Title string
	Time  time.Time
}

func (r *Room) GetID() string {
	return strconv.Itoa(r.ID)
}

func (r *Room) SetID(id string) {
	myID, err := strconv.Atoi(id)
	if err != nil {
		log.Panicf("Room.SetID(): %s\n", err)
	}
	r.ID = myID
}

type RoomController struct {
	*api.Controller[*Room]
	*api.Service[*Room]
}

func NewRoomController(mux *http.ServeMux) *RoomController {
	rc := &RoomController{
		Controller: api.NewController[*Room](mux),
		Service:    api.NewService[*Room](),
	}
	rc.Controller.Service = rc.Service
	return rc
}

func (rc *RoomController) InitFakeData() {
	room1 := Room{
		ID:    1,
		Title: "Room number one",
		Time:  time.Now(),
	}
	room2 := Room{
		ID:    2,
		Title: "Room number two",
		Time:  time.Now().Add(5 * time.Hour),
	}
	room3 := Room{
		ID:    3,
		Title: "Room number three",
		Time:  time.Now().Add(2 * time.Hour),
	}
	rc.Service.Insert(&room1)
	rc.Service.Insert(&room2)
	rc.Service.Insert(&room3)
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

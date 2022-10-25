package rooms

import (
	"log"
	"strconv"

	"github.com/scottcagno/angular-refresher/pkg/roombooking/api"
)

type Room struct {
	ID int
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

var RoomService = api.NewService[*Room]()

type RestRoomController struct {
	baseURI string
}

func (rc *RestRoomController) GetAllRooms() {

}

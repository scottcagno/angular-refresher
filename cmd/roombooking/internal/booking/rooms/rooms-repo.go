package rooms

import (
	"time"

	"github.com/scottcagno/angular-refresher/pkg/web/api"
)

type RoomRepository struct {
	api.Repository[*Room, string]
}

func NewRoomRepository() *RoomRepository {
	u := &RoomRepository{
		Repository: api.NewMemoryRepository[*Room, string](),
	}
	u.init()
	return u
}

func (repo *RoomRepository) init() {
	room1 := &Room{
		ID:    "1",
		Title: "Room number one",
		Time:  time.Now(),
	}
	room2 := &Room{
		ID:    "2",
		Title: "Room number two",
		Time:  time.Now().Add(5 * time.Hour),
	}
	room3 := &Room{
		ID:    "3",
		Title: "Room number three",
		Time:  time.Now().Add(2 * time.Hour),
	}
	var err error
	err = repo.Insert(room1.ID, room1)
	if err != nil {
		panic(err)
	}
	err = repo.Insert(room2.ID, room2)
	if err != nil {
		panic(err)
	}
	err = repo.Insert(room3.ID, room3)
	if err != nil {
		panic(err)
	}
}

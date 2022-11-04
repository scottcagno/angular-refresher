package rooms

import (
	"github.com/scottcagno/angular-refresher/pkg/web/api"
)

type RoomRepository struct {
	api.Repository[*Room, int]
}

func NewRoomRepository() *RoomRepository {
	u := &RoomRepository{
		Repository: api.NewMemoryRepository[*Room, int](),
	}
	u.init()
	return u
}

func (repo *RoomRepository) init() {
	room1 := NewRoom(1, "Blue meeting room", "First floor")
	room1.AddLayoutCapacity(NewLayoutCapacity(Layout_USHAPE, 5))

	room2 := NewRoom(2, "Red meeting room", "Second floor")
	room2.AddLayoutCapacity(NewLayoutCapacity(Layout_THEATER, 9))

	room3 := NewRoom(3, "Main conference room", "Third floor")
	room3.AddLayoutCapacity(NewLayoutCapacity(Layout_BOARD, 12))

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

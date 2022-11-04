package booking

import (
	"time"

	"github.com/scottcagno/angular-refresher/cmd/roombooking/internal/booking/rooms"
	"github.com/scottcagno/angular-refresher/cmd/roombooking/internal/booking/users"
	"github.com/scottcagno/angular-refresher/pkg/web/api"
)

type BookingRepository struct {
	nextID         int
	UserRepository *users.UserRepository
	RoomRepository *rooms.RoomRepository
	Repository     api.Repository[*Booking, string]
}

func NewBookingRepository(userRepo *users.UserRepository, roomRepo *rooms.RoomRepository) *BookingRepository {
	b := &BookingRepository{
		UserRepository: userRepo,
		RoomRepository: roomRepo,
		Repository:     api.NewMemoryRepository[*Booking, string](),
	}
	b.init()
	return b
}

func (b *BookingRepository) init() {
	userData, err := b.UserRepository.Find(func(u *users.User) bool { return u != nil })
	if err != nil {
		panic(err)
	}
	roomData, err := b.RoomRepository.Find(func(r *rooms.Room) bool { return r != nil })
	if err != nil {
		panic(err)
	}
	booking1 := &Booking{
		ID:           "1",
		Title:        "Conference call with CEO",
		User:         *userData[1],
		Room:         *roomData[0],
		Date:         time.Now().String(),
		StartTime:    time.Now().Add(30 * time.Minute),
		EndTime:      time.Now().Add(3 * time.Hour),
		Participants: 4,
	}
	booking2 := &Booking{
		ID:           "2",
		Title:        "Some important meeting",
		User:         *userData[0],
		Room:         *roomData[1],
		Date:         time.Now().Add(29 * time.Hour).String(),
		StartTime:    time.Now().Add(30 * time.Hour),
		EndTime:      time.Now().Add(31 * time.Hour),
		Participants: 7,
	}
	b.nextID = 3
	err = b.Repository.Insert(booking1.ID, booking1)
	if err != nil {
		panic(err)
	}
	err = b.Repository.Insert(booking2.ID, booking2)
	if err != nil {
		panic(err)
	}
}

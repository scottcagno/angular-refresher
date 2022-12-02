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
	Repository     api.Repository[*Booking, int]
}

func NewBookingRepository(userRepo *users.UserRepository, roomRepo *rooms.RoomRepository) *BookingRepository {
	b := &BookingRepository{
		UserRepository: userRepo,
		RoomRepository: roomRepo,
		Repository:     api.NewMemoryRepository[*Booking, int](),
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
		ID:           1,
		Title:        "Conference call with CEO",
		User:         *userData[1],
		Room:         *roomData[0],
		Date:         GetDate(Today, 0),
		StartTime:    time.Now().Add(0 * time.Minute),
		EndTime:      time.Now().Add(30 * time.Minute),
		Participants: 4,
	}
	booking2 := &Booking{
		ID:           2,
		Title:        "Some important meeting",
		User:         *userData[0],
		Room:         *roomData[1],
		Date:         GetDate(Today, 0),
		StartTime:    time.Now().Add(45 * time.Minute),
		EndTime:      time.Now().Add(90 * time.Minute),
		Participants: 7,
	}
	booking3 := &Booking{
		ID:           3,
		Title:        "The most important meeting",
		User:         *userData[0],
		Room:         *roomData[1],
		Date:         GetDate(Add, 1),
		StartTime:    time.Now().Add(24 * time.Hour),
		EndTime:      time.Now().Add(25 * time.Hour),
		Participants: 2,
	}
	b.nextID = 4
	err = b.Repository.Insert(booking1.ID, booking1)
	if err != nil {
		panic(err)
	}
	err = b.Repository.Insert(booking2.ID, booking2)
	if err != nil {
		panic(err)
	}
	err = b.Repository.Insert(booking3.ID, booking3)
	if err != nil {
		panic(err)
	}
}

package services

import (
	"sync"
	"time"

	"github.com/scottcagno/angular-refresher/pkg/roombooking/booking"
	"github.com/scottcagno/angular-refresher/pkg/roombooking/rooms"
	"github.com/scottcagno/angular-refresher/pkg/roombooking/users"
)

var once sync.Once

type DataService struct {
	RoomRepo    []rooms.Room
	UserRepo    []users.User
	BookingRepo []booking.Booking
}

var DataServiceInstance *DataService

func NewDataService() *DataService {
	once.Do(
		func() {
			DataServiceInstance = initDataServiceInstance()
		},
	)
	return DataServiceInstance
}

func initDataServiceInstance() *DataService {
	ds := &DataService{
		RoomRepo:    make([]rooms.Room, 0),
		UserRepo:    make([]users.User, 0),
		BookingRepo: make([]booking.Booking, 0),
	}
	ds.InitService()
	return ds
}

func (ds *DataService) InitService() {

	// add initial data to room repo
	room1 := rooms.Room{
		ID:    "1",
		Title: "Room number one",
		Time:  time.Now(),
	}
	room2 := rooms.Room{
		ID:    "2",
		Title: "Room number two",
		Time:  time.Now().Add(5 * time.Hour),
	}
	room3 := rooms.Room{
		ID:    "3",
		Title: "Room number three",
		Time:  time.Now().Add(2 * time.Hour),
	}
	ds.RoomRepo = append(ds.RoomRepo, room1)
	ds.RoomRepo = append(ds.RoomRepo, room2)
	ds.RoomRepo = append(ds.RoomRepo, room3)

	// add initial data to user repo
	user1 := users.User{
		ID:   "1",
		Name: "Dick Chesterwood",
	}
	user2 := users.User{
		ID:   "2",
		Name: "Matt Greencroft",
	}
	ds.UserRepo = append(ds.UserRepo, user1)
	ds.UserRepo = append(ds.UserRepo, user2)

	// add initial data to booking repo
	booking1 := booking.Booking{
		ID:           "1",
		Title:        "Conference call with CEO",
		User:         user2,
		Room:         room1,
		Date:         time.Now().String(),
		StartTime:    time.Now().Add(30 * time.Minute),
		EndTime:      time.Now().Add(3 * time.Hour),
		Participants: 4,
	}
	booking2 := booking.Booking{
		ID:           "2",
		Title:        "Some important meeting",
		User:         user1,
		Room:         room2,
		Date:         time.Now().Add(29 * time.Hour).String(),
		StartTime:    time.Now().Add(30 * time.Hour),
		EndTime:      time.Now().Add(31 * time.Hour),
		Participants: 7,
	}
	ds.BookingRepo = append(ds.BookingRepo, booking1)
	ds.BookingRepo = append(ds.BookingRepo, booking2)
}

func (ds *DataService) Get(key string) any {
	if key == "RoomRepo" {
		return ds.RoomRepo
	}
	if key == "UserRepo" {
		return ds.UserRepo
	}
	if key == "BookingRepo" {
		return ds.BookingRepo
	}
	return nil
}

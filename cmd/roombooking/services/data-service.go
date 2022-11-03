package services

import (
	"sync"

	"github.com/scottcagno/angular-refresher/cmd/roombooking/internal/booking"
	"github.com/scottcagno/angular-refresher/cmd/roombooking/internal/booking/rooms"
	"github.com/scottcagno/angular-refresher/cmd/roombooking/internal/booking/users"
)

var once sync.Once

type DataService struct {
	RoomRepo    *rooms.RoomRepository
	UserRepo    *users.UserRepository
	BookingRepo *booking.BookingRepository
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
		RoomRepo: rooms.NewRoomRepository(),
		UserRepo: users.NewUserRepository(),
	}
	ds.BookingRepo = booking.NewBookingRepository(ds.UserRepo, ds.RoomRepo)
	return ds
}

package services

import (
	"sync"

	"github.com/scottcagno/angular-refresher/cmd/roombooking/internal/booking"
	"github.com/scottcagno/angular-refresher/cmd/roombooking/internal/booking/rooms"
	"github.com/scottcagno/angular-refresher/cmd/roombooking/internal/booking/users"
)

var dataServiceOnce sync.Once

type DataService struct {
	RoomRepo    *rooms.RoomRepository
	UserRepo    *users.UserRepository
	BookingRepo *booking.BookingRepository
}

var DataServiceInstance *DataService

func NewDataService() *DataService {
	dataServiceOnce.Do(
		func() {
			DataServiceInstance = initDataServiceInstance()
		},
	)
	return DataServiceInstance
}

func initDataServiceInstance() *DataService {
	service := &DataService{
		RoomRepo: rooms.NewRoomRepository(),
		UserRepo: users.NewUserRepository(),
	}
	service.BookingRepo = booking.NewBookingRepository(service.UserRepo, service.RoomRepo)
	return service
}

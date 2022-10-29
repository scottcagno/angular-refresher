package services

import (
	"sync"

	"github.com/scottcagno/angular-refresher/pkg/roombooking/booking"
	"github.com/scottcagno/angular-refresher/pkg/roombooking/rooms"
	"github.com/scottcagno/angular-refresher/pkg/roombooking/users"
)

var once sync.Once

type DataService struct {
	RoomRepo    rooms.RoomRepository
	UserRepo    users.UserRepository
	BookingRepo booking.BookingRepository
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
		RoomRepo:    rooms.RoomRepository{},
		UserRepo:    users.UserRepository{},
		BookingRepo: booking.BookingRepository{},
	}
	ds.InitService()
	return ds
}

func (ds *DataService) InitService() {

	// initialize rooms repository
	ds.RoomRepo.Init()

	// add initial data to user repo
	ds.UserRepo.Init()

	// add initial data to booking repo
	ds.BookingRepo.Init()
}

func (ds *DataService) GetRepository(key string) any {
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

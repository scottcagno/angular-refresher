package services

import (
	"sync"

	"github.com/scottcagno/angular-refresher/cmd/roombooking/internal/booking"
	"github.com/scottcagno/angular-refresher/cmd/roombooking/internal/booking/rooms"
	"github.com/scottcagno/angular-refresher/cmd/roombooking/internal/booking/users"
	"github.com/scottcagno/angular-refresher/pkg/web/api"
)

var once sync.Once

type DataService struct {
	repos       map[string]api.Repository
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

func (ds *DataService) GetRepository(key string) api.Repository {
	repo, found := ds.repos[key]
	if !found {
		return nil
	}
	return repo
}

func (ds *DataService) AddRepository(key string, val api.Repository) {
	ds.repos[key] = val
}

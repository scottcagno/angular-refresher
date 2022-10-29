package main

import (
	"log"
	"net/http"

	"github.com/scottcagno/angular-refresher/pkg/roombooking/booking"
	"github.com/scottcagno/angular-refresher/pkg/roombooking/rooms"
	"github.com/scottcagno/angular-refresher/pkg/roombooking/services"
	"github.com/scottcagno/angular-refresher/pkg/roombooking/users"
	"github.com/scottcagno/angular-refresher/pkg/web/api"
)

func main() {

	// initialize global data service (contains ref to all repositories)
	ds := services.NewDataService()

	// initialize rooms controller (and inject the data service into it)
	roomCont := new(rooms.Controller)
	roomCont.Inject(ds)

	// initialize users controller (and inject the data service into it)
	userCont := new(users.Controller)
	userCont.Inject(ds)

	// initialize booking controller (and inject the data service into it)
	bookingCont := new(booking.Controller)
	bookingCont.Inject(ds)

	// initialize new rest api server
	restAPI := api.NewAPI("/api/", nil)

	// register controllers with api
	restAPI.Register("rooms", roomCont)
	restAPI.Register("users", userCont)
	restAPI.Register("bookings", bookingCont)

	// serve the api
	log.Fatal(http.ListenAndServe(":8080", restAPI))
}

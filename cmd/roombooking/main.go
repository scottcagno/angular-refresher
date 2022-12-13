package main

import (
	"log"
	"net/http"
	"time"

	"github.com/scottcagno/angular-refresher/cmd/roombooking/internal/booking"
	"github.com/scottcagno/angular-refresher/cmd/roombooking/internal/booking/rooms"
	"github.com/scottcagno/angular-refresher/cmd/roombooking/internal/booking/users"
	"github.com/scottcagno/angular-refresher/cmd/roombooking/services"
	"github.com/scottcagno/angular-refresher/pkg/web/api"
	"github.com/scottcagno/angular-refresher/pkg/web/api/middleware"
)

func main() {

	// initialize global data service (contains ref to all repositories)
	ds := services.NewDataService()

	// initialize rooms controller (and inject the data service into it)
	roomCont := &rooms.Controller{RoomRepository: ds.RoomRepo, BasicAuth: ds.BasicAuth}

	// initialize users controller
	userCont := &users.Controller{UserRepository: ds.UserRepo, BasicAuth: ds.BasicAuth}

	// initialize booking controller
	bookingCont := &booking.Controller{BookingRepository: ds.BookingRepo, BasicAuth: ds.BasicAuth}

	// initialize our cors handler
	// cors := middleware.CORSHandler(
	// 	&middleware.CORSConfig{
	// 		AllowOrigins:     "http://localhost:4200/api/**",
	// 		AllowMethods:     "GET,POST,PUT,DELETE",
	// 		AllowHeaders:     "",
	// 		AllowCredentials: false,
	// 		ExposeHeaders:    "",
	// 		MaxAge:           int(time.Duration(12 * time.Hour).Seconds()),
	// 	},
	// )

	apiConf := &api.APIConfig{
		CORS: &middleware.CORSConfig{
			AllowOrigins:     "http://localhost:4200/api/**",
			AllowMethods:     "GET,POST,PUT,DELETE",
			AllowHeaders:     "",
			AllowCredentials: false,
			ExposeHeaders:    "",
			MaxAge:           int(time.Duration(12 * time.Hour).Seconds()),
		},
		Auth: map[string]string{"admin": "secret"},
	}

	// initialize new rest api server
	restAPI := api.NewAPI("/api/", apiConf)

	// register controllers with api
	restAPI.Register("rooms", roomCont)
	restAPI.Register("users", userCont)
	restAPI.Register("bookings", bookingCont)
	restAPI.RegisterCustom("users/resetPassword", userCont)

	// serve the api
	log.Fatal(http.ListenAndServe(":8080", restAPI))
}

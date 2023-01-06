package main

import (
	"log"
	"net/http"
	"time"

	"github.com/scottcagno/angular-refresher/cmd/roombooking/internal/booking"
	"github.com/scottcagno/angular-refresher/cmd/roombooking/internal/booking/rooms"
	"github.com/scottcagno/angular-refresher/cmd/roombooking/internal/booking/users"
	"github.com/scottcagno/angular-refresher/cmd/roombooking/services"
	"github.com/scottcagno/angular-refresher/pkg/web"
	"github.com/scottcagno/angular-refresher/pkg/web/api"
	"github.com/scottcagno/angular-refresher/pkg/web/api/middleware"
)

func main() {

	// initialize global data service (contains ref to all repositories)
	ds := services.NewDataService()

	// initialize rooms controller (and inject the data service into it)
	roomCont := &rooms.Controller{RoomRepository: ds.RoomRepo}

	// initialize users controller
	userCont := &users.Controller{UserRepository: ds.UserRepo}

	// initialize booking controller
	bookingCont := &booking.Controller{BookingRepository: ds.BookingRepo}

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
			AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
			AllowHeaders:     "",
			AllowCredentials: false,
			ExposeHeaders:    "",
			MaxAge:           int(time.Duration(12 * time.Hour).Seconds()),
		},
	}

	// initialize new rest api server
	restAPI := api.NewAPI("/api/", apiConf)

	// initialize auth service
	authService := api.MakeAuthService(api.NewJWTAuthService(&web.SystemUser{
		Username: "admin",
		Password: "secret",
		Role:     "ROLE_ADMIN",
	},
		"cmd/roombooking/private_key.pem",
		"cmd/roombooking/public_key.pem",
	))

	// register controllers with api
	restAPI.RegisterAuthService("/api/auth", authService)
	restAPI.Register("rooms", roomCont, false)
	restAPI.Register("users", userCont, true)
	restAPI.Register("bookings", bookingCont, false)
	restAPI.RegisterCustom("users/resetPassword", userCont, true)

	// certFile := "cmd/roombooking/cert/CA/CA.pem"
	// keyFile := "cmd/roombooking/cert/CA/CA.key"

	// serve the api
	// log.Fatal(http.ListenAndServeTLS(":8080", certFile, keyFile, restAPI))

	log.Fatal(http.ListenAndServe(":8080", restAPI))
}

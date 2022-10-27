package main

import (
	"log"
	"net/http"

	"github.com/scottcagno/angular-refresher/pkg/roombooking/rooms"
	"github.com/scottcagno/angular-refresher/pkg/web/api"
)

func main() {
	mux := http.NewServeMux()
	api := api.NewAPI("/api/", mux)
	rc := new(rooms.RoomController)
	rc.InitFakeData()
	api.Register("rooms", rc)
	log.Fatal(http.ListenAndServe(":8080", api))
}

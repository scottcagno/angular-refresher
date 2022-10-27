package main

import (
	"log"
	"net/http"

	"github.com/scottcagno/angular-refresher/pkg/roombooking/rooms"
	"github.com/scottcagno/angular-refresher/pkg/web/api"
)

func main() {

	restAPI := api.NewAPI("/api/", nil)
	restAPI.Register("rooms", new(rooms.Rooms))

	log.Fatal(http.ListenAndServe(":8080", restAPI))
}

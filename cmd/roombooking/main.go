package main

import (
	"log"
	"net/http"

	"github.com/scottcagno/angular-refresher/pkg/roombooking/rooms"
)

func main() {

	mux := http.NewServeMux()

	rooms := rooms.NewRoomController(mux)
	rooms.InitFakeData()

	rooms.Get("/api/rooms", rooms.DefaultGetAllHandler())
	rooms.Get("/api/rooms/", rooms.DefaultGetOneHandler())
	// rooms.InitCRUD("/api/rooms")
	rooms.Get("/api/rooms/stats", rooms.StatsHandler())

	// api := roombooking.NewRestAPI("/api", ":8080")
	// log.Printf("Serving up API %s, on port %s\n", api.Version(), api.Addr())
	log.Fatal(http.ListenAndServe(":8080", mux))
}

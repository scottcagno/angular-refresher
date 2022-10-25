package main

import (
	"fmt"

	"github.com/scottcagno/angular-refresher/pkg/roombooking/api"
	"github.com/scottcagno/angular-refresher/pkg/roombooking/rooms"
)

func main() {

	svc := api.NewService[*rooms.Room](nil)
	fmt.Println(svc)

	// api := roombooking.NewRestAPI("/api", ":8080")
	// log.Printf("Serving up API %s, on port %s\n", api.Version(), api.Addr())
	// log.Fatal(http.ListenAndServe(api.Addr(), api))
}

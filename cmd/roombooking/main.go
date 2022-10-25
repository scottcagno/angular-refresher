package main

import (
	"log"
	"net/http"

	"github.com/scottcagno/angular-refresher/pkg/roombooking"
)

func main() {

	api := roombooking.NewRestAPI("/api", ":8080")

	log.Printf("Serving up API %s, on port %s\n", api.Version(), api.Addr())
	log.Fatal(http.ListenAndServe(api.Addr(), api))
}

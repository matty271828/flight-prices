package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/matty271828/flight-prices/amadeus"
	"github.com/matty271828/flight-prices/controller"
	"github.com/matty271828/flight-prices/server"
)

func main() {
	amadeusClient, err := amadeus.NewAmadeusClient()
	if err != nil {
		fmt.Printf("Error fetching getting amadeus client: %s\n", err)
		return
	}

	controller := controller.NewController(amadeusClient)
	server := server.NewServer(controller)

	log.Fatal(http.ListenAndServe(":8080", server))
}

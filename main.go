package main

import (
	"fmt"
	"net/http"

	"github.com/matty271828/flight-prices/amadeus"
	"github.com/matty271828/flight-prices/controller"
	"github.com/matty271828/flight-prices/server"
)

func main() {
	token, err := amadeus.GetAuthToken("UAGXkYYCKiVnZfUR0wflNqz9IK3upUea", "UV2vj7DHz3wJyhUG")
	if err != nil {
		fmt.Printf("Error fetching fetching amadeus token: %s\n", err)
		return
	}

	amadeusClient := amadeus.NewAmadeusClient(token)

	// Create a new controller with the Amadeus client
	controller := controller.NewController(amadeusClient)

	// Create a new server with the controller
	httpServer := server.NewServer(controller)

	err = http.ListenAndServe(":8080", httpServer)
	if err != nil {
		fmt.Println("failed to start server: %w", err)
	}
}

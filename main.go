package main

import (
	"fmt"
	"os"

	amadeus "github.com/matty271828/flight-prices/amadeus"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please provide a destination.")
		return
	}
	destination := os.Args[1]

	token, err := amadeus.GetAuthToken("UAGXkYYCKiVnZfUR0wflNqz9IK3upUea", "UV2vj7DHz3wJyhUG")
	if err != nil {
		fmt.Printf("Error fetching fetching amadeus token: %s\n", err)
		return
	}

	fmt.Println(token)

	flightInfo, err := amadeus.GetFlightInfo(destination, token)
	if err != nil {
		fmt.Printf("Error fetching flight info: %s\n", err)
		return
	}

	for _, flight := range flightInfo.Data {
		fmt.Println(flight.Destination)
		fmt.Println(flight.Price.Total)
	}
}

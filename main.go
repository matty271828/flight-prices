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

	flightInfo, err := amadeus.GetFlightInfo(destination)
	if err != nil {
		fmt.Printf("Error fetching flight info: %s\n", err)
		return
	}

	fmt.Println(flightInfo)
}

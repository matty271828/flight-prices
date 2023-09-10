package controller

import (
	"github.com/matty271828/flight-prices/amadeus"
	"github.com/matty271828/flight-prices/amadeus/flightinspiration"
)

type ControllerManager interface {
	FlightInspirationSearch(origin string) (*flightinspiration.FISResponse, error)

	FlightOffersSearch(oorigin, destination, departureDate, adults string) (*amadeus.ApiResponse, error)

	AirportSearch(keyword string) (*amadeus.ApiResponse, error)
}

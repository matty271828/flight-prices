package controller

import "github.com/matty271828/flight-prices/amadeus"

type ControllerManager interface {
	FlightInspirationSearch(origin string) (*amadeus.FISResponse, error)

	FlightOffersSearch(oorigin, destination, departureDate, adults string) (*amadeus.ApiResponse, error)

	AirportSearch(keyword string) (*amadeus.ApiResponse, error)
}

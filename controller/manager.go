package controller

import (
	"github.com/matty271828/flight-prices/amadeus/airportsearch"
	"github.com/matty271828/flight-prices/amadeus/flightinspiration"
	"github.com/matty271828/flight-prices/amadeus/flightoffers"
)

type ControllerManager interface {
	FlightInspirationSearch(origin string) (*flightinspiration.FISResponse, error)

	FlightOffersSearch(oorigin, destination, departureDate, adults string) (*flightoffers.FOSResponse, error)

	AirportSearch(keyword string) (*airportsearch.AirportSearchResponse, error)
}

package amadeus

import (
	"github.com/matty271828/flight-prices/amadeus/airportsearch"
	"github.com/matty271828/flight-prices/amadeus/flightinspiration"
	"github.com/matty271828/flight-prices/amadeus/flightoffers"
)

type AmadeusManager interface {
	FlightInspirationSearch(origin string) (*flightinspiration.FISResponse, error)

	FlightOffersSearch(origin, destination, departureDate, timeRange string) (*flightoffers.FOSResponse, error)

	AirportSearch(name string) (*airportsearch.AirportSearchResponse, error)
}

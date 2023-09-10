package amadeus

import "github.com/matty271828/flight-prices/amadeus/flightinspiration"

type AmadeusManager interface {
	FlightInspirationSearch(origin string) (*flightinspiration.FISResponse, error)

	FlightOffersSearch(origin, destination, departureDate, timeRange string) (*ApiResponse, error)

	AirportSearch(name string) (*ApiResponse, error)
}

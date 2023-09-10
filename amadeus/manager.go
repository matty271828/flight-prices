package amadeus

type AmadeusManager interface {
	FlightInspirationSearch(origin string) (*FISResponse, error)

	FlightOffersSearch(origin, destination, departureDate, timeRange string) (*ApiResponse, error)

	AirportSearch(name string) (*ApiResponse, error)
}

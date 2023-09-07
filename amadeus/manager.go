package amadeus

type AmadeusManager interface {
	FlightInspirationSearch(origin string) (*ApiResponse, error)

	FlightOffersSearch(origin, destination, departureDate, timeRange string) (*ApiResponse, error)
}

package amadeus

type AmadeusManager interface {
	FlightInspirationSearch(origin string) (*ApiResponse, error)
}

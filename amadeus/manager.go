package amadeus

type AmadeusManager interface {
	FlightOffersSearch(origin string) (*ApiResponse, error)
}

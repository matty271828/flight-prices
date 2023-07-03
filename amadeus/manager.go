package amadeus

type AmadeusManager interface {
	GetFlightInfo(origin string) (*ApiResponse, error)
}

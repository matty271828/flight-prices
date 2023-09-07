package amadeus

type ApiResponse struct {
	Data []FlightDestination `json:"data"`
}

type FlightDestination struct {
	Type          string `json:"type"`
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	DepartureDate string `json:"departureDate"`
	Price         Price  `json:"price"`
	Links         Links  `json:"links"`
}

type Price struct {
	Total string `json:"total"`
}

type Links struct {
	FlightDates  string `json:"flightDates"`
	FlightOffers string `json:"flightOffers"`
}

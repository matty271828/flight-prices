package flightinspiration

// FlightInspirationResponse represents the response structure for FlightInspirationSearch
type FISResponse struct {
	Data         []fisDestination `json:"data"`
	Dictionaries Dictionaries     `json:"dictionaries"`
	Meta         Meta             `json:"meta"`
}

// FlightDestination represents the flight destination details
type fisDestination struct {
	Type          string   `json:"type"`
	Origin        string   `json:"origin"`
	Destination   string   `json:"destination"`
	DepartureDate string   `json:"departureDate"`
	ReturnDate    string   `json:"returnDate"`
	FISPrice      fisPrice `json:"price"`
	FISLinks      fisLinks `json:"links"`
}

// Price represents the price details
type fisPrice struct {
	Total string `json:"total"`
}

// Links represents links related to the flight destination
type fisLinks struct {
	FlightDates  string `json:"flightDates"`
	FlightOffers string `json:"flightOffers"`
}

// Dictionaries contains information about currencies and locations
type Dictionaries struct {
	Currencies map[string]string         `json:"currencies"`
	Locations  map[string]LocationDetail `json:"locations"`
}

// LocationDetail provides details about a location
type LocationDetail struct {
	SubType      string `json:"subType"`
	DetailedName string `json:"detailedName"`
}

// Meta contains meta-information related to the search response
type Meta struct {
	Currency string                 `json:"currency"`
	Links    map[string]string      `json:"links"`
	Defaults map[string]interface{} `json:"defaults"`
}

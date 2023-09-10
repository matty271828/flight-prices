package amadeus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// FlightInspirationResponse represents the response structure for FlightInspirationSearch
type FISResponse struct {
	Data         []FlightDestination `json:"data"`
	Dictionaries Dictionaries        `json:"dictionaries"`
	Meta         Meta                `json:"meta"`
}

// FlightDestination represents the flight destination details
type FISDestination struct {
	Type          string `json:"type"`
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	DepartureDate string `json:"departureDate"`
	ReturnDate    string `json:"returnDate"`
	FISPrice      Price  `json:"price"`
	FISLinks      Links  `json:"links"`
}

// Price represents the price details
type FISPrice struct {
	Total string `json:"total"`
}

// Links represents links related to the flight destination
type FISLinks struct {
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

// FlightInspirationSearch finds the cheapest flight destinations from a specific city.
func (a *AmadeusClient) FlightInspirationSearch(origin string) (*FISResponse, error) {
	requestURL := fmt.Sprintf("%s/v1/shopping/flight-destinations?origin=%s&oneWay=true&nonStop=true", a.BaseURL, origin)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+a.Token)

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && checkStatusCode(resp) != nil {
		return nil, fmt.Errorf("unexpected status code %d: %v", resp.StatusCode, checkStatusCode(resp))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var fisResponse FISResponse
	err = json.Unmarshal(data, &fisResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling flight inspiration search response: %v", err)
	}

	return &fisResponse, nil
}

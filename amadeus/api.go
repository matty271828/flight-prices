package amadeus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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

func GetFlightInfo(departure, token string) (*ApiResponse, error) {
	requestURL := fmt.Sprintf("https://test.api.amadeus.com/v1/shopping/flight-destinations?origin=LON&oneWay=true&nonStop=true")

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(data, &apiResponse)
	if err != nil {
		return nil, err
	}

	return &apiResponse, nil
}

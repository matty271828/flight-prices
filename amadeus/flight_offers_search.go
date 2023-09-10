package amadeus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func (a *AmadeusClient) FlightOffersSearch(origin, destination, departureDate, adults string) (*ApiResponse, error) {
	requestURL := fmt.Sprintf(
		"%s/v1/shopping/flight-offers?originLocationCode=%s&destinationLocationCode=%s&departureDate=%s&adults=%s&max=20",
		a.BaseURL, origin, destination, departureDate, adults)

	// Log the request URL
	log.Printf("Sending request to URL: %s", requestURL)

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

	var apiResponse ApiResponse
	err = json.Unmarshal(data, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling API response: %v", err)
	}

	return &apiResponse, nil
}

package amadeus

import (
	"fmt"
	"net/http"
	"time"

	"github.com/matty271828/flight-prices/amadeus/airportsearch"
	"github.com/matty271828/flight-prices/amadeus/flightinspiration"
	"github.com/matty271828/flight-prices/amadeus/flightoffers"
)

type AmadeusClient struct {
	BaseURL string
	Token   string
}

// NewAmadeusClient initializes a new AmadeusClient with the provided config.
func NewAmadeusClient() (*AmadeusClient, error) {
	baseURL := GetBaseURL()

	cfg := GetAPIKeys()

	if cfg.ClientId == "" || cfg.ClientSecret == "" {
		return nil, fmt.Errorf("Missing credentials for Amadeus API")
	}

	token, err := GetAuthToken(cfg.ClientId, cfg.ClientSecret)
	if err != nil {
		return nil, fmt.Errorf("Error fetching Amadeus token: %w", err)
	}

	return &AmadeusClient{BaseURL: baseURL, Token: token}, nil
}

func (a *AmadeusClient) FlightInspirationSearch(origin string) (*flightinspiration.FISResponse, error) {
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

	return flightinspiration.ParseFISResponse(resp)
}

func (a *AmadeusClient) FlightOffersSearch(origin, destination, departureDate, adults string) (*flightoffers.FOSResponse, error) {
	requestURL := fmt.Sprintf(
		"%s/v2/shopping/flight-offers?originLocationCode=%s&destinationLocationCode=%s&departureDate=%s&adults=%s&max=20",
		a.BaseURL, origin, destination, departureDate, adults)

	// Create the HTTP request
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set the required headers
	req.Header.Set("Authorization", "Bearer "+a.Token)

	// Initialize the HTTP client with a timeout
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Execute the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %v", err)
	}
	defer resp.Body.Close()

	// Check the status code of the response
	if resp.StatusCode != http.StatusOK && checkStatusCode(resp) != nil {
		return nil, fmt.Errorf("unexpected status code %d: %v", resp.StatusCode, checkStatusCode(resp))
	}

	return flightoffers.ParseFOSResponse(resp)
}

// AirportSearch finds airports and cities that match a specific word or string of letters.
func (a *AmadeusClient) AirportSearch(keyword string) (*airportsearch.AirportSearchResponse, error) {
	requestURL := fmt.Sprintf("%s/v1/reference-data/locations?subType=CITY,AIRPORT&keyword=%s", a.BaseURL, keyword)

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
	return airportsearch.ParseAirportSearchResponse(resp)
}

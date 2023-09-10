package amadeus

import (
	"fmt"
	"net/http"
	"time"

	"github.com/matty271828/flight-prices/amadeus/flightinspiration"
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

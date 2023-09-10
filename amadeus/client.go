package amadeus

import (
	"fmt"
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

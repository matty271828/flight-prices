package amadeus

import (
	"fmt"
)

type AmadeusClient struct {
	Token string
}

// NewAmadeusClient initializes a new AmadeusClient with the provided config.
func NewAmadeusClient() (*AmadeusClient, error) {
	cfg := GetAPIKeys()

	if cfg.ClientId == "" || cfg.ClientSecret == "" {
		return nil, fmt.Errorf("Missing credentials for Amadeus API")
	}

	token, err := GetAuthToken(cfg.ClientId, cfg.ClientSecret)
	if err != nil {
		return nil, fmt.Errorf("Error fetching Amadeus token: %w", err)
	}

	return &AmadeusClient{Token: token}, nil
}

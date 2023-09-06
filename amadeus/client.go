package amadeus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type AmadeusClient struct {
	Token string
}

// Config struct to hold client configuration
type Config struct {
	ClientId     string
	ClientSecret string
}

func NewAmadeusClient(cfg Config) (*AmadeusClient, error) {
	if cfg.ClientId == "" || cfg.ClientSecret == "" {
		return nil, fmt.Errorf("Missing credentials for Amadeus API")
	}

	token, err := GetAuthToken(cfg.ClientId, cfg.ClientSecret)
	if err != nil {
		return nil, fmt.Errorf("Error fetching Amadeus token: %w", err)
	}

	return &AmadeusClient{Token: token}, nil
}

func checkStatusCode(resp *http.Response) error {
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var errorResponse AmadeusError
	err := json.Unmarshal(bodyBytes, &errorResponse)
	if err != nil {
		return fmt.Errorf("Error parsing error response: %v. Body: %s", err, string(bodyBytes))
	}

	// Construct a neat error message
	var errorMsgs []string
	for _, errDetail := range errorResponse.Errors {
		errorMsg := fmt.Sprintf(
			"Code: %d, Title: %s, Detail: %s, Status: %d",
			errDetail.Code,
			errDetail.Title,
			errDetail.Detail,
			errDetail.Status,
		)
		errorMsgs = append(errorMsgs, errorMsg)
	}
	return fmt.Errorf("Errors: %s", strings.Join(errorMsgs, " | "))
}

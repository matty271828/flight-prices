package flightinspiration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// FlightInspirationResponse represents the response structure for FlightInspirationSearch
type FISResponse struct {
	Data         []fisDestination `json:"data"`
	Dictionaries Dictionaries     `json:"dictionaries"`
	Meta         Meta             `json:"meta"`
}

// ParseFISResponse parses the response body into the FISResponse structure.
func ParseFISResponse(resp *http.Response) (*FISResponse, error) {
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var response FISResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling flight inspiration search response: %v", err)
	}

	return &response, nil
}

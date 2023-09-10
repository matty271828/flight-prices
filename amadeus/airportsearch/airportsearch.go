package airportsearch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type AirportSearchResponse struct {
	Meta Meta       `json:"meta"`
	Data []Location `json:"data"`
}

// ParseAirportSearchResponse parses the response body into the FOSResponse structure.
func ParseAirportSearchResponse(resp *http.Response) (*AirportSearchResponse, error) {
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var response AirportSearchResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling flight inspiration search response: %v", err)
	}

	return &response, nil
}

package flightoffers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type FOSResponse struct {
	Meta         Meta          `json:"meta"`
	Data         []FlightOffer `json:"data"`
	Dictionaries Dictionaries  `json:"dictionaries"`
}

// ParseFOSResponse parses the response body into the FOSResponse structure.
func ParseFOSResponse(resp *http.Response) (*FOSResponse, error) {
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var response FOSResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling flight inspiration search response: %v", err)
	}

	return &response, nil
}

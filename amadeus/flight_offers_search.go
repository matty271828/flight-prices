package amadeus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (a *AmadeusClient) FlightOffersSearch(origin, destination, departureDate, timeRange string) (*ApiResponse, error) {
	requestURL := fmt.Sprintf(
		"https://test.api.amadeus.com/v1/shopping/flight-offers?origin=%s&destination=%s&departureDate=%s&timeRange=%s",
		origin, destination, departureDate, timeRange)

	req, err := http.NewRequest("POST", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+a.Token)

	client := &http.Client{}
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

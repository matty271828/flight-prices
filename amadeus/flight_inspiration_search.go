package amadeus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (a *AmadeusClient) FlightInspirationSearch(origin string) (*ApiResponse, error) {
	requestURL := fmt.Sprintf("%s/v1/shopping/flight-destinations?origin=%s&oneWay=true&nonStop=true", a.BaseURL, origin)

	req, err := http.NewRequest("GET", requestURL, nil)
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

package amadeus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/matty271828/flight-prices/amadeus/flightinspiration"
)

// FlightInspirationSearch finds the cheapest flight destinations from a specific city.
func (a *AmadeusClient) FlightInspirationSearch2(origin string) (*flightinspiration.FISResponse, error) {
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

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var fisResponse flightinspiration.FISResponse
	err = json.Unmarshal(data, &fisResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling flight inspiration search response: %v", err)
	}

	return &fisResponse, nil
}

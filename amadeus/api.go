package amadeus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type FlightInfo struct {
	Departure string `json:"departure"`
	Arrival   string `json:"arrival"`
	Date      string `json:"date"`
	Time      string `json:"time"`
	Price     string `json:"price"`
}

func GetFlightInfo(departure string) (*FlightInfo, error) {
	requestURL := fmt.Sprintf("https://api.amadeus.com/v1/shopping/flight-offers?origin=%s", departure)

	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var flightInfo FlightInfo
	err = json.Unmarshal(data, &flightInfo)
	if err != nil {
		return nil, err
	}

	return &flightInfo, nil

}

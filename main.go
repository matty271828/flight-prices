package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type FlightInfo struct {
	Departure string `json:"departure"`
	Arrival   string `json:"arrival"`
	Date      string `json:"date"`
	Time      string `json:"time"`
	Price     string `json:"price"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please provide a destination.")
		return
	}
	destination := os.Args[1]

	requestURL := fmt.Sprintf("https://api.amadeus.com/v1/shopping/flight-offers?origin=MAD&destination=%s", destination)

	resp, err := http.Get(requestURL)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(resp.Body)
		var flightInfo FlightInfo
		err := json.Unmarshal([]byte(data), &flightInfo)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Departure: ", flightInfo.Departure)
			fmt.Println("Arrival: ", flightInfo.Arrival)
			fmt.Println("Date: ", flightInfo.Date)
			fmt.Println("Time: ", flightInfo.Time)
			fmt.Println("Price: ", flightInfo.Price)
		}
	}
}

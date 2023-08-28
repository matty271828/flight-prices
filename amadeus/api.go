package amadeus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type AmadeusClient struct {
	Token string
}

func NewAmadeusClient() (*AmadeusClient, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found")
	}

	clientId := os.Getenv("AMADEUS_API_KEY")
	clientSecret := os.Getenv("AMADEUS_API_SECRET")

	if clientId == "" || clientSecret == "" {
		return nil, fmt.Errorf("Environment variables for Amadeus API not set")
	}

	token, err := GetAuthToken(clientId, clientSecret)
	if err != nil {
		fmt.Printf("Error fetching fetching amadeus token: %s\n", err)
		return nil, err
	}

	return &AmadeusClient{Token: token}, nil
}

type ApiResponse struct {
	Data []FlightDestination `json:"data"`
}

type FlightDestination struct {
	Type          string `json:"type"`
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	DepartureDate string `json:"departureDate"`
	Price         Price  `json:"price"`
	Links         Links  `json:"links"`
}

type Price struct {
	Total string `json:"total"`
}

type Links struct {
	FlightDates  string `json:"flightDates"`
	FlightOffers string `json:"flightOffers"`
}

func (a *AmadeusClient) GetFlightInfo(origin string) (*ApiResponse, error) {
	requestURL := fmt.Sprintf("https://test.api.amadeus.com/v1/shopping/flight-destinations?origin=%s&oneWay=true&nonStop=true", origin)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+a.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(data, &apiResponse)
	if err != nil {
		return nil, err
	}

	return &apiResponse, nil
}

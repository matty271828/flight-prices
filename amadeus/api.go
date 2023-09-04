package amadeus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

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

	fmt.Printf("ClientId: %s, ClientSecret: %s\n", clientId, clientSecret)

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

func (a *AmadeusClient) FlightOffersSearch(origin string) (*ApiResponse, error) {
	requestURL := fmt.Sprintf("https://test.api.amadeus.com/v1/shopping/flight-destinations?origin=%s&oneWay=true&nonStop=true", origin)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+a.Token)
	log.Printf("Fetching flight info with token: %s", a.Token) // Debugging purposes.

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		var errorResponse AmadeusError
		err = json.Unmarshal(bodyBytes, &errorResponse)
		if err != nil {
			return nil, fmt.Errorf("unexpected status: %d. Error parsing error response: %v. Body: %s", resp.StatusCode, err, string(bodyBytes))
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

		return nil, fmt.Errorf("unexpected status: %d. Errors: %s", resp.StatusCode, strings.Join(errorMsgs, " | "))
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

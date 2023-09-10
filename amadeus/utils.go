package amadeus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// APIKeys struct to hold amadeus api keys
type APIKeys struct {
	ClientId     string
	ClientSecret string
}

func GetAPIKeys() APIKeys {
	// Initially assume using test keys
	clientId := os.Getenv("AMADEUS_API_TEST_KEY")
	clientSecret := os.Getenv("AMADEUS_API_TEST_SECRET")

	if os.Getenv("USE_PROD_API") == "true" {
		clientId = os.Getenv("AMADEUS_API_PROD_KEY")
		clientSecret = os.Getenv("AMADEUS_API_PROD_SECRET")
	}

	return APIKeys{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
}

// GetBaseURL is used to switch the client between test and prod calls
func GetBaseURL() string {
	if os.Getenv("USE_PROD_API") == "true" {
		return "https://api.amadeus.com"
	}
	return "https://test.api.amadeus.com"
}

func checkStatusCode(resp *http.Response) error {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to read response body: %v", err)
	}

	var errorResponse AmadeusError
	err = json.Unmarshal(bodyBytes, &errorResponse)
	if err != nil {
		return fmt.Errorf("Error parsing error response: %v. Body: %s", err, string(bodyBytes))
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
	return fmt.Errorf("Errors: %s", strings.Join(errorMsgs, " | "))
}

package amadeus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func loadConfig() Config {
	clientId := os.Getenv("AMADEUS_API_KEY")
	clientSecret := os.Getenv("AMADEUS_API_SECRET")

	log.Println("ClientId:", clientId)
	log.Println("ClientSecret:", clientSecret)

	return Config{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
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

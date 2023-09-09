package amadeus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func LoadConfig() Config {
	return Config{
		ClientId:     os.Getenv("AMADEUS_API_KEY"),
		ClientSecret: os.Getenv("AMADEUS_API_SECRET"),
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

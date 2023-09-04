package amadeus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type AuthResponse struct {
	Type      string `json:"type"`
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
}

func GetAuthToken(clientID string, clientSecret string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	req, err := http.NewRequest(http.MethodPost, "https://test.api.amadeus.com/v1/security/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("error creating token request: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error performing token request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status while getting token: %d. Response: %s", resp.StatusCode, string(bodyBytes))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading token response body: %v", err)
	}

	var authResponse AuthResponse
	err = json.Unmarshal(body, &authResponse)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling token response: %v", err)
	}

	if authResponse.Token == "" {
		return "", fmt.Errorf("received empty token. Full response: %s", string(body))
	}

	return authResponse.Token, nil
}

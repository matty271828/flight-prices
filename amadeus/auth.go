package amadeus

import (
	"encoding/json"
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
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var authResponse AuthResponse
	json.Unmarshal(body, &authResponse)

	return authResponse.Token, nil
}

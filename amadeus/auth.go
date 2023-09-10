package amadeus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var (
	mu           sync.Mutex
	currentToken string
	tokenExpiry  time.Time
)

type AuthResponse struct {
	Type      string    `json:"type"`
	Token     string    `json:"access_token"`
	ExpiresIn int       `json:"expires_in"`
	Expiry    time.Time `json:"-"` // This field is for internal use to compute token expiry
}

func GetAuthToken(clientID string, clientSecret string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	baseURL := GetBaseURL()
	tokenEndpoint := fmt.Sprintf("%s/v1/security/oauth2/token", baseURL)
	req, err := http.NewRequest(http.MethodPost, tokenEndpoint, strings.NewReader(data.Encode()))
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

	authResponse.Expiry = time.Now().Add(time.Second * time.Duration(authResponse.ExpiresIn))

	mu.Lock()
	currentToken = authResponse.Token
	tokenExpiry = authResponse.Expiry
	mu.Unlock()

	return authResponse.Token, nil
}

func EnsureValidTokenMiddleware(clientID string, clientSecret string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			tokenExpired := time.Now().After(tokenExpiry)
			mu.Unlock()

			if tokenExpired {
				_, err := GetAuthToken(clientID, clientSecret)
				if err != nil {
					http.Error(w, "Authentication failed", http.StatusUnauthorized)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

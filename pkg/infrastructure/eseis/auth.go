package eseis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// authRequest is the authentication request payload for Eseis API
type authRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	ClientID  string `json:"client_id"`
	GrantType string `json:"grant_type"`
	Scope     string `json:"scope"`
}

// authResponse is the authentication response payload from Eseis API
type authResponse struct {
	AccessToken  string `json:"access_token"`
	CreatedAt    int64  `json:"created_at"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

// authToken is the parsed authToken response from Eseis API
type authToken struct {
	accessToken  string
	expiresAt    time.Time
	refreshToken string
}

func (e *EseisClient) Authenticate() (*authToken, error) {
	requestBody := authRequest{
		Username:  e.config.Username,
		Password:  e.config.Password,
		ClientID:  e.config.ClientId,
		GrantType: "password",
		Scope:     "eseis",
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json payload: %w", err)
		// handle err
	}
	body := bytes.NewReader(requestBodyBytes)

	req, err := http.NewRequest("POST", e.buildURL("/v1/oauth/token"), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create authentication request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send authentication request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response code %d", resp.StatusCode)
	}

	authResponse := authResponse{}
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode authentication response: %w", err)
	}

	token := &authToken{
		accessToken:  authResponse.AccessToken,
		expiresAt:    time.Unix(authResponse.CreatedAt+authResponse.ExpiresIn, 0),
		refreshToken: authResponse.RefreshToken,
	}
	return token, nil
}

package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GooglePayload struct {
	SUB           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}

func ConvertGoogleToken(accessToken string) (*GooglePayload, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v3/userinfo?access_token=%s", accessToken))
	if err != nil {
		return nil, fmt.Errorf("error making request to Google API: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Directly unmarshal into GooglePayload struct
	var data GooglePayload
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, fmt.Errorf("error unmarshalling response into GooglePayload: %w", err)
	}

	// Check if the response contains an error field
	if data.Email == "" {
		return nil, errors.New("invalid token: missing user email")
	}

	return &data, nil
}

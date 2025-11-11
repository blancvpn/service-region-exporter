package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"region-exporter/config"
)

type RedditService struct {
	BaseService
	url string
}

func NewRedditService(client *http.Client) *RedditService {
	cfg := config.CLIConfig.Services.Reddit
	return &RedditService{
		BaseService: BaseService{
			name:    "Reddit",
			client:  client,
			enabled: cfg.Enabled,
		},
		url: cfg.URL,
	}
}

func (r *RedditService) CheckRegion() (string, error) {
	jsonBody := `{"scopes":["email"]}`
	tokenReq, err := http.NewRequest("POST", "https://www.reddit.com/auth/v2/oauth/access-token/loid", strings.NewReader(jsonBody))
	if err != nil {
		return "", err
	}

	tokenReq.Header.Set("Authorization", "Basic b2hYcG9xclpZdWIxa2c6")
	tokenReq.Header.Set("User-Agent", "Reddit/Version 2025.29.0/Build 2529021/Android 13")
	tokenReq.Header.Set("Content-Type", "application/json")

	tokenBody, err := r.makeRequest(tokenReq)
	if err != nil {
		return "", err
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(tokenBody, &tokenResp); err != nil {
		return "", err
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("failed to get access token")
	}

	gqlBody := `{"operationName":"UserLocation","variables":{},"extensions":{"persistedQuery":{"version":1,"sha256Hash":"f07de258c54537e24d7856080f662c1b1268210251e5789c8c08f20d76cc8ab2"}}}`
	locReq, err := http.NewRequest("POST", "https://gql-fed.reddit.com", strings.NewReader(gqlBody))
	if err != nil {
		return "", err
	}

	locReq.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	locReq.Header.Set("User-Agent", "Reddit/Version 2025.29.0/Build 2529021/Android 13")
	locReq.Header.Set("Content-Type", "application/json")

	locBody, err := r.makeRequest(locReq)
	if err != nil {
		return "", err
	}

	var locResp struct {
		Data struct {
			UserLocation struct {
				CountryCode string `json:"countryCode"`
			} `json:"userLocation"`
		} `json:"data"`
	}
	if err := json.Unmarshal(locBody, &locResp); err != nil {
		return "", err
	}

	if locResp.Data.UserLocation.CountryCode == "" {
		return "", fmt.Errorf("country code not found in response")
	}

	return locResp.Data.UserLocation.CountryCode, nil
}

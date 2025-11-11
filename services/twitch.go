package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"region-exporter/config"
)

type TwitchService struct {
	BaseService
	url      string
	clientID string
}

func NewTwitchService(client *http.Client) *TwitchService {
	cfg := config.CLIConfig.Services.Twitch
	return &TwitchService{
		BaseService: BaseService{
			name:    "Twitch",
			client:  client,
			enabled: cfg.Enabled,
		},
		url:      cfg.URL,
		clientID: cfg.Token,
	}
}

type twitchResponse []struct {
	Data struct {
		RequestInfo struct {
			CountryCode string `json:"countryCode"`
		} `json:"requestInfo"`
	} `json:"data"`
}

func (t *TwitchService) CheckRegion() (string, error) {
	jsonBody := `[{"operationName":"VerifyEmail_CurrentUser","variables":{},"extensions":{"persistedQuery":{"version":1,"sha256Hash":"f9e7dcdf7e99c314c82d8f7f725fab5f99d1df3d7359b53c9ae122deec590198"}}}]`

	req, err := http.NewRequest("POST", t.url, strings.NewReader(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Client-Id", t.clientID)
	req.Header.Set("Content-Type", "application/json")

	body, err := t.makeRequest(req)
	if err != nil {
		return "", err
	}

	var response twitchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(response) == 0 || response[0].Data.RequestInfo.CountryCode == "" {
		return "", fmt.Errorf("country code not found in response")
	}

	return response[0].Data.RequestInfo.CountryCode, nil
}

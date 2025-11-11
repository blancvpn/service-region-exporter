package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"region-exporter/config"
)

type TikTokService struct {
	BaseService
	url string
}

func NewTikTokService(client *http.Client) *TikTokService {
	cfg := config.CLIConfig.Services.TikTok
	return &TikTokService{
		BaseService: BaseService{
			name:    "TikTok",
			client:  client,
			enabled: cfg.Enabled,
		},
		url: cfg.URL,
	}
}

type tiktokResponse struct {
	Body struct {
		AppProps struct {
			Region string `json:"region"`
		} `json:"appProps"`
	} `json:"body"`
}

func (t *TikTokService) CheckRegion() (string, error) {
	req, err := http.NewRequest("GET", t.url, nil)
	if err != nil {
		return "", err
	}

	body, err := t.makeRequest(req)
	if err != nil {
		return "", err
	}

	var response tiktokResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	if response.Body.AppProps.Region == "" {
		return "", fmt.Errorf("region field is empty")
	}

	return response.Body.AppProps.Region, nil
}

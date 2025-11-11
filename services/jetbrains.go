package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"region-exporter/config"
)

type JetBrainsService struct {
	BaseService
	url string
}

func NewJetBrainsService(client *http.Client) *JetBrainsService {
	cfg := config.CLIConfig.Services.JetBrains
	return &JetBrainsService{
		BaseService: BaseService{
			name:    "JetBrains",
			client:  client,
			enabled: cfg.Enabled,
		},
		url: cfg.URL,
	}
}

type jetbrainsResponse struct {
	Code string `json:"code"`
}

func (j *JetBrainsService) CheckRegion() (string, error) {
	req, err := http.NewRequest("GET", j.url, nil)
	if err != nil {
		return "", err
	}

	body, err := j.makeRequest(req)
	if err != nil {
		return "", err
	}

	var response jetbrainsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	if response.Code == "" {
		return "", fmt.Errorf("code field is empty")
	}

	return response.Code, nil
}

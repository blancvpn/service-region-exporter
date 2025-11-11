package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"region-exporter/config"
)

type ChatGPTService struct {
	BaseService
	url    string
	apiKey string
}

func NewChatGPTService(client *http.Client) *ChatGPTService {
	cfg := config.CLIConfig.Services.ChatGPT
	return &ChatGPTService{
		BaseService: BaseService{
			name:    "ChatGPT",
			client:  client,
			enabled: cfg.Enabled,
		},
		url:    cfg.URL,
		apiKey: cfg.Token,
	}
}

type chatGPTResponse struct {
	DerivedFields struct {
		Country string `json:"country"`
	} `json:"derived_fields"`
}

func (c *ChatGPTService) CheckRegion() (string, error) {
	req, err := http.NewRequest("POST", c.url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Statsig-Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	body, err := c.makeRequest(req)
	if err != nil {
		return "", err
	}

	var response chatGPTResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	if response.DerivedFields.Country == "" {
		return "", fmt.Errorf("country field is empty")
	}

	return response.DerivedFields.Country, nil
}

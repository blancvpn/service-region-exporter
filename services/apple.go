package services

import (
	"net/http"

	"region-exporter/config"
)

type AppleService struct {
	BaseService
	url string
}

func NewAppleService(client *http.Client) *AppleService {
	cfg := config.CLIConfig.Services.Apple
	return &AppleService{
		BaseService: BaseService{
			name:    "Apple",
			client:  client,
			enabled: cfg.Enabled,
		},
		url: cfg.URL,
	}
}

func (a *AppleService) CheckRegion() (string, error) {
	req, err := http.NewRequest("GET", a.url, nil)
	if err != nil {
		return "", err
	}

	body, err := a.makeRequest(req)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

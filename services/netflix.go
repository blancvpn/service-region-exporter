package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"region-exporter/config"
)

type NetflixService struct {
	BaseService
	url   string
	token string
}

func NewNetflixService(client *http.Client) *NetflixService {
	cfg := config.CLIConfig.Services.Netflix
	return &NetflixService{
		BaseService: BaseService{
			name:    "Netflix",
			client:  client,
			enabled: cfg.Enabled,
		},
		url:   cfg.URL,
		token: cfg.Token,
	}
}

type netflixResponse struct {
	Client struct {
		Location struct {
			Country string `json:"country"`
		} `json:"location"`
	} `json:"client"`
}

func (n *NetflixService) CheckRegion() (string, error) {
	url := fmt.Sprintf("%s?https=true&token=%s&urlCount=1", n.url, n.token)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")

	body, err := n.makeRequest(req)
	if err != nil {
		return "", err
	}

	var response netflixResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	if response.Client.Location.Country == "" {
		return "", fmt.Errorf("country field is empty")
	}

	return response.Client.Location.Country, nil
}

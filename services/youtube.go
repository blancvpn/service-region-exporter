package services

import (
	"fmt"
	"net/http"
	"regexp"

	"region-exporter/config"
)

type YouTubeService struct {
	BaseService
	url string
}

func NewYouTubeService(client *http.Client) *YouTubeService {
	cfg := config.CLIConfig.Services.YouTube
	return &YouTubeService{
		BaseService: BaseService{
			name:    "YouTube",
			client:  client,
			enabled: cfg.Enabled,
		},
		url: cfg.URL,
	}
}

func (y *YouTubeService) CheckRegion() (string, error) {
	req, err := http.NewRequest("GET", y.url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	body, err := y.makeRequest(req)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`"countryCode":"(\w+)"`)
	matches := re.FindSubmatch(body)
	if len(matches) > 1 {
		return string(matches[1]), nil
	}

	return "", fmt.Errorf("countryCode not found in response")
}

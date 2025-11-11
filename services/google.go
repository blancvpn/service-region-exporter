package services

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"region-exporter/config"
)

type GoogleService struct {
	BaseService
	url string
}

func NewGoogleService(client *http.Client) *GoogleService {
	cfg := config.CLIConfig.Services.Google
	return &GoogleService{
		BaseService: BaseService{
			name:    "Google",
			client:  client,
			enabled: cfg.Enabled,
		},
		url: cfg.URL,
	}
}

func (g *GoogleService) CheckRegion() (string, error) {
	req, err := http.NewRequest("GET", g.url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US;q=0.9")

	body, err := g.makeRequest(req)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`<div class="yVZQTb">([^<(]+)`)
	matches := re.FindSubmatch(body)
	if len(matches) > 1 {
		country := strings.TrimSpace(string(matches[1]))

		if len(country) == 2 {
			return strings.ToUpper(country), nil
		}

		countryCode, err := getCountryCode(country)
		if err != nil {
			normalized := strings.ToLower(country)
			if code, ok := countryNameToCode[normalized]; ok {
				return code, nil
			}
			return "", fmt.Errorf("failed to convert country name to code: %s", country)
		}
		return countryCode, nil
	}

	return "", ErrRegionNotFound
}

var ErrRegionNotFound = fmt.Errorf("region not found in response")

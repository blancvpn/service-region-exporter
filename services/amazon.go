package services

import (
	"net/http"
	"regexp"

	"region-exporter/config"
)

type AmazonPrimeService struct {
	BaseService
	url string
}

func NewAmazonPrimeService(client *http.Client) *AmazonPrimeService {
	cfg := config.CLIConfig.Services.AmazonPrime
	return &AmazonPrimeService{
		BaseService: BaseService{
			name:    "Amazon Prime",
			client:  client,
			enabled: cfg.Enabled,
		},
		url: cfg.URL,
	}
}

func (a *AmazonPrimeService) CheckRegion() (string, error) {
	req, err := http.NewRequest("GET", a.url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")

	body, err := a.makeRequest(req)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`"currentTerritory":"(\w+)"`)
	if matches := re.FindSubmatch(body); len(matches) > 1 {
		region := string(matches[1])
		if len(region) > 2 {
			return region[:2], nil
		}
		return region, nil
	}

	return "", ErrRegionNotFound
}

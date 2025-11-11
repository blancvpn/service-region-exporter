package services

import (
	"net/http"
	"regexp"

	"region-exporter/config"
)

type DeezerService struct {
	BaseService
	url string
}

func NewDeezerService(client *http.Client) *DeezerService {
	cfg := config.CLIConfig.Services.Deezer
	return &DeezerService{
		BaseService: BaseService{
			name:    "Deezer",
			client:  client,
			enabled: cfg.Enabled,
		},
		url: cfg.URL,
	}
}

func (d *DeezerService) CheckRegion() (string, error) {
	req, err := http.NewRequest("GET", d.url, nil)
	if err != nil {
		return "", err
	}

	body, err := d.makeRequest(req)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`'country': '(.*)'`)
	if matches := re.FindSubmatch(body); len(matches) > 1 {
		return string(matches[1]), nil
	}

	return "", ErrRegionNotFound
}

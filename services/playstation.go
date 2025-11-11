package services

import (
	"net/http"
	"regexp"

	"region-exporter/config"
)

type PlayStationService struct {
	BaseService
	url string
}

func NewPlayStationService(client *http.Client) *PlayStationService {
	cfg := config.CLIConfig.Services.PlayStation
	return &PlayStationService{
		BaseService: BaseService{
			name:    "PlayStation",
			client:  client,
			enabled: cfg.Enabled,
		},
		url: cfg.URL,
	}
}

func (p *PlayStationService) CheckRegion() (string, error) {
	req, err := http.NewRequest("HEAD", p.url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	cookies := resp.Header.Get("Set-Cookie")
	re := regexp.MustCompile(`country=([A-Z]{2})`)
	if matches := re.FindStringSubmatch(cookies); len(matches) > 1 {
		return matches[1], nil
	}

	return "", ErrRegionNotFound
}

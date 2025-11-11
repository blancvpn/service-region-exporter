package services

import (
	"net/http"
	"regexp"

	"region-exporter/config"
)

type SteamService struct {
	BaseService
	url string
}

func NewSteamService(client *http.Client) *SteamService {
	cfg := config.CLIConfig.Services.Steam
	return &SteamService{
		BaseService: BaseService{
			name:    "Steam",
			client:  client,
			enabled: cfg.Enabled,
		},
		url: cfg.URL,
	}
}

func (s *SteamService) CheckRegion() (string, error) {
	req, err := http.NewRequest("HEAD", s.url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	cookies := resp.Header.Get("Set-Cookie")
	re := regexp.MustCompile(`steamCountry=([^;%]+)`)
	if matches := re.FindStringSubmatch(cookies); len(matches) > 1 {
		return matches[1], nil
	}

	return "", ErrRegionNotFound
}

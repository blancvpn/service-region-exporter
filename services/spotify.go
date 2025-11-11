package services

import (
	"net/http"
	"regexp"

	"region-exporter/config"
)

type SpotifyService struct {
	BaseService
	url string
}

func NewSpotifyService(client *http.Client) *SpotifyService {
	cfg := config.CLIConfig.Services.Spotify
	return &SpotifyService{
		BaseService: BaseService{
			name:    "Spotify",
			client:  client,
			enabled: cfg.Enabled,
		},
		url: cfg.URL,
	}
}

func (s *SpotifyService) CheckRegion() (string, error) {
	req, err := http.NewRequest("GET", s.url, nil)
	if err != nil {
		return "", err
	}

	body, err := s.makeRequest(req)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`"geoLocationCountryCode":"([^"]*)"`)
	matches := re.FindSubmatch(body)
	if len(matches) > 1 {
		return string(matches[1]), nil
	}

	return "", ErrRegionNotFound
}

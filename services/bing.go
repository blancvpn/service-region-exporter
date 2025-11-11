package services

import (
	"net/http"
	"regexp"

	"region-exporter/config"
)

type BingService struct {
	BaseService
	url string
}

func NewBingService(client *http.Client) *BingService {
	cfg := config.CLIConfig.Services.Bing
	return &BingService{
		BaseService: BaseService{
			name:    "Microsoft (Bing)",
			client:  client,
			enabled: cfg.Enabled,
		},
		url: cfg.URL,
	}
}

func (b *BingService) CheckRegion() (string, error) {
	req, err := http.NewRequest("GET", b.url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")

	body, err := b.makeRequest(req)
	if err != nil {
		return "", err
	}

	if regexp.MustCompile(`cn\.bing\.com`).Match(body) {
		return "CN", nil
	}

	re := regexp.MustCompile(`Region\s*:\s*"(\w{2})"`)
	if matches := re.FindSubmatch(body); len(matches) > 1 {
		region := string(matches[1])

		if region == "WW" {
			liveReq, err := http.NewRequest("GET", "https://login.live.com", nil)
			if err == nil {
				liveReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")
				liveBody, err := b.makeRequest(liveReq)
				if err == nil {
					liveRe := regexp.MustCompile(`"sRequestCountry":"(\w{2})"`)
					if liveMatches := liveRe.FindSubmatch(liveBody); len(liveMatches) > 1 {
						return string(liveMatches[1]), nil
					}
				}
			}
		}

		return region, nil
	}

	return "", ErrRegionNotFound
}

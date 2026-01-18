package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"region-exporter/config"
	"region-exporter/metrics"
)

const (
	abuseIPDBURL      = "https://api.abuseipdb.com/api/v2/check"
	abuseCheckInterval = time.Hour
)

type AbuseIPDBChecker struct {
	client    *http.Client
	token     string
	lastCheck time.Time
	mu        sync.Mutex
}

type abuseIPDBResponse struct {
	Data struct {
		AbuseConfidenceScore int `json:"abuseConfidenceScore"`
	} `json:"data"`
}

func NewAbuseIPDBChecker(client *http.Client) *AbuseIPDBChecker {
	return &AbuseIPDBChecker{
		client: client,
		token:  config.CLIConfig.Services.AbuseIPDB.Token,
	}
}

func (a *AbuseIPDBChecker) IsEnabled() bool {
	return a.token != ""
}

func (a *AbuseIPDBChecker) Check(ip string) {
	if !a.IsEnabled() {
		return
	}

	a.mu.Lock()
	if time.Since(a.lastCheck) < abuseCheckInterval {
		a.mu.Unlock()
		if config.CLIConfig.Verbose {
			log.Printf("AbuseIPDB check skipped (last check: %s ago)", time.Since(a.lastCheck).Round(time.Second))
		}
		return
	}
	a.lastCheck = time.Now()
	a.mu.Unlock()

	score, err := a.fetchAbuseScore(ip)
	if err != nil {
		log.Printf("✗ AbuseIPDB check failed: %v", err)
		return
	}

	log.Printf("✓ AbuseIPDB confidence score: %d", score)
	metrics.RecordAbuseScore(ip, float64(score))
}

func (a *AbuseIPDBChecker) fetchAbuseScore(ip string) (int, error) {
	req, err := http.NewRequest("GET", abuseIPDBURL, nil)
	if err != nil {
		return 0, err
	}

	q := req.URL.Query()
	q.Add("ipAddress", ip)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Key", a.token)
	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	var response abuseIPDBResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return 0, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return response.Data.AbuseConfidenceScore, nil
}

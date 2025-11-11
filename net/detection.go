package net

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"region-exporter/config"
	"region-exporter/models"
)

func DetectServerInfo(client *http.Client, cfg *config.CLI) (serverIP string, expectedRegion string, err error) {
	if cfg.Verbose {
		log.Println("Detecting server IP and region...")
	}

	url := cfg.Services.IPInfo.URL
	if cfg.Services.IPInfo.Token != "" {
		url = fmt.Sprintf("%s?token=%s", url, cfg.Services.IPInfo.Token)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to get IP info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response: %w", err)
	}

	var ipInfo models.IPInfoResponse
	if err := json.Unmarshal(body, &ipInfo); err != nil {
		return "", "", fmt.Errorf("failed to parse IP info: %w", err)
	}

	serverIP = ipInfo.IP

	if cfg.Check.ExpectedRegion == "" {
		expectedRegion = ipInfo.Country
		log.Printf("Auto-detected expected region: %s", expectedRegion)
	} else {
		expectedRegion = cfg.Check.ExpectedRegion
		log.Printf("Using configured expected region: %s", expectedRegion)
	}

	log.Printf("Server IP: %s", serverIP)

	return serverIP, expectedRegion, nil
}

package services

import (
	"log"
	"time"

	"region-exporter/config"
	"region-exporter/metrics"
	"region-exporter/models"
)

func CheckService(svc Service, cfg *config.CLI, serverIP, expectedRegion string) models.ServiceStatus {
	serviceName := svc.Name()

	if !svc.IsEnabled() {
		if cfg.Verbose {
			log.Printf("Service %s is disabled, skipping", serviceName)
		}
		return models.ServiceStatus{
			Name:           serviceName,
			ExpectedRegion: expectedRegion,
			Enabled:        false,
		}
	}

	startTime := time.Now()
	region, err := svc.CheckRegion()
	latencyMs := float64(time.Since(startTime).Milliseconds())

	status := models.ServiceStatus{
		Name:           serviceName,
		ExpectedRegion: expectedRegion,
		Enabled:        true,
		Latency:        latencyMs,
	}

	if err != nil {
		log.Printf("✗ %s check failed: %v", serviceName, err)
		status.DetectedRegion = "error"
		status.Match = false
		status.Error = err.Error()
		metrics.RecordMatch(serviceName, serverIP, "error", expectedRegion, false)
		metrics.RecordLatency(serviceName, serverIP, latencyMs)
	} else if region == "" {
		log.Printf("✗ %s returned empty region", serviceName)
		status.DetectedRegion = "unknown"
		status.Match = false
		metrics.RecordMatch(serviceName, serverIP, "unknown", expectedRegion, false)
		metrics.RecordLatency(serviceName, serverIP, latencyMs)
	} else {
		normalizedDetected := NormalizeCountry(region)
		normalizedExpected := NormalizeCountry(expectedRegion)
		matches := normalizedDetected == normalizedExpected

		status.DetectedRegion = region
		status.Match = matches

		if matches {
			log.Printf("✓ %s region matches: %s (latency: %dms)", serviceName, region, int64(latencyMs))
		} else {
			log.Printf("✗ %s region mismatch: got %s, expected %s (latency: %dms)", serviceName, region, expectedRegion, int64(latencyMs))
		}

		metrics.RecordMatch(serviceName, serverIP, region, expectedRegion, matches)
		metrics.RecordLatency(serviceName, serverIP, latencyMs)
	}

	return status
}

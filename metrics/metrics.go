package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ServiceRegionMatch = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ip_service_region_match",
			Help: "Service region detection match: 1 if detected region matches expected, 0 otherwise",
		},
		[]string{
			"service",
			"ip",
			"detected_region",
			"expected_region",
		},
	)

	ServiceLatency = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ip_service_latency_milliseconds",
			Help: "Service request latency in milliseconds",
		},
		[]string{
			"service",
			"ip",
		},
	)

	AbuseConfidenceScore = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ip_abuse_confidence_score",
			Help: "AbuseIPDB confidence score (0-100) indicating how abusive the IP is",
		},
		[]string{
			"ip",
		},
	)

	registry *prometheus.Registry
)

func Init() {
	registry = prometheus.NewRegistry()
	registry.MustRegister(ServiceRegionMatch)
	registry.MustRegister(ServiceLatency)
	registry.MustRegister(AbuseConfidenceScore)
}

func Handler() http.Handler {
	return promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	)
}

func Reset() {
	ServiceRegionMatch.Reset()
	ServiceLatency.Reset()
}

func RecordMatch(service, ip, detectedRegion, expectedRegion string, matches bool) {
	var value float64
	if matches {
		value = 1
	} else {
		value = 0
	}

	ServiceRegionMatch.WithLabelValues(
		service,
		ip,
		detectedRegion,
		expectedRegion,
	).Set(value)
}

func RecordLatency(service, ip string, latencyMs float64) {
	ServiceLatency.WithLabelValues(service, ip).Set(latencyMs)
}

func RecordAbuseScore(ip string, score float64) {
	AbuseConfidenceScore.WithLabelValues(ip).Set(score)
}

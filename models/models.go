package models

type IPInfoResponse struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
}

type ServiceStatus struct {
	Name           string  `json:"name"`
	DetectedRegion string  `json:"detected_region"`
	ExpectedRegion string  `json:"expected_region"`
	Match          bool    `json:"match"`
	Enabled        bool    `json:"enabled"`
	Error          string  `json:"error,omitempty"`
	Latency        float64 `json:"latency"`
}

type DashboardData struct {
	ServerIP       string          `json:"server_ip"`
	ExpectedRegion string          `json:"expected_region"`
	Services       []ServiceStatus `json:"services"`
	TotalServices  int             `json:"total_services"`
	MatchCount     int             `json:"match_count"`
	MismatchCount  int             `json:"mismatch_count"`
}

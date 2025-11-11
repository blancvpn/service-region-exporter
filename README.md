# Service Region Exporter

Service Region Exporter is a Prometheus exporter designed to monitor and verify the geographic region detection of various online services. The application periodically checks multiple services to determine which region they detect your server IP as belonging to, compares it with the expected region, and exposes metrics for monitoring.

## Features

- **Multi-Service Support**: Monitors region detection for 15+ popular services including Google, YouTube, Netflix, Spotify, ChatGPT, and more
- **Automatic Region Detection**: Automatically detects your server's IP address and expected region
- **Prometheus Metrics**: Exposes standard Prometheus metrics for integration with monitoring systems
- **REST API**: Provides JSON API endpoint for service status and region information
- **Configurable Checks**: Enable or disable individual service checks
- **Network Flexibility**: Support for SOCKS5 proxy, custom network interfaces, and IPv4/IPv6-only modes
- **Scheduled Checks**: Configurable interval for periodic region verification
- **Basic Authentication**: Optional HTTP basic authentication for metrics endpoint
- **Docker Support**: Official Docker image available

## Supported Services

The exporter supports region detection for the following services:

- Google Play
- YouTube
- ChatGPT
- Netflix
- Twitch
- Spotify
- Deezer
- Reddit
- Amazon Prime Video
- Apple
- Steam
- PlayStation
- TikTok
- JetBrains
- Microsoft (Bing)

## Installation

### Docker

The easiest way to run Service Region Exporter is using Docker:

```bash
docker run -d \
  --name region-exporter \
  -p 9999:9999 \
  -e CHECK_INTERVAL=60 \
  -e EXPECTED_REGION=US \
  ghcr.io/blancvpn/service-region-exporter:latest
```

### From Binary Release

1. Download the latest release for your architecture from [GitHub Releases](https://github.com/blancvpn/service-region-exporter/releases)
2. Extract the binary and make it executable:
   ```bash
   chmod +x region-exporter
   ```
3. Run the application:
   ```bash
   ./region-exporter
   ```

### From Source

1. Ensure you have Go 1.24+ installed
2. Clone the repository:
   ```bash
   git clone https://github.com/blancvpn/service-region-exporter.git
   cd service-region-exporter
   ```
3. Build the application:
   ```bash
   go build -o region-exporter .
   ```
4. Run the application:
   ```bash
   ./region-exporter
   ```

## Configuration

Service Region Exporter can be configured via command-line flags or environment variables. All configuration options support both methods.

### Basic Configuration

```bash
# Using command-line flags
./region-exporter \
  --metrics-host 0.0.0.0 \
  --metrics-port 9999 \
  --check-interval 60 \
  --check-expected-region US

# Using environment variables
export METRICS_HOST=0.0.0.0
export METRICS_PORT=9999
export CHECK_INTERVAL=60
export EXPECTED_REGION=US
./region-exporter
```

### Configuration Options

#### Metrics Server

- `--metrics-host` / `METRICS_HOST`: HTTP server host (default: `0.0.0.0`)
- `--metrics-port` / `METRICS_PORT`: HTTP server port (default: `9999`)
- `--metrics-username` / `METRICS_USERNAME`: Basic auth username (optional)
- `--metrics-password` / `METRICS_PASSWORD`: Basic auth password (optional)

#### Check Settings

- `--check-interval` / `CHECK_INTERVAL`: Interval between checks in seconds (default: `60`)
- `--check-timeout` / `CHECK_TIMEOUT`: HTTP request timeout in seconds (default: `10`)
- `--check-expected-region` / `EXPECTED_REGION`: Expected region code (e.g., `US`, `GB`). Auto-detected from server IP if not set

#### Network Settings

- `--network-proxy` / `NETWORK_PROXY`: SOCKS5 proxy address (format: `host:port`)
- `--network-interface` / `NETWORK_INTERFACE`: Network interface to use
- `--network-ipv4-only` / `NETWORK_IPV4_ONLY`: Use only IPv4 (default: `false`)
- `--network-ipv6-only` / `NETWORK_IPV6_ONLY`: Use only IPv6 (default: `false`)

#### Service Configuration

Each service can be individually enabled or disabled:

- `--service-google-enabled` / `GOOGLE_ENABLED`: Enable Google region check (default: `true`)
- `--service-youtube-enabled` / `YOUTUBE_ENABLED`: Enable YouTube region check (default: `true`)
- `--service-chatgpt-enabled` / `CHATGPT_ENABLED`: Enable ChatGPT region check (default: `true`)
- `--service-netflix-enabled` / `NETFLIX_ENABLED`: Enable Netflix region check (default: `true`)
- `--service-twitch-enabled` / `TWITCH_ENABLED`: Enable Twitch region check (default: `true`)
- `--service-spotify-enabled` / `SPOTIFY_ENABLED`: Enable Spotify region check (default: `true`)
- `--service-deezer-enabled` / `DEEZER_ENABLED`: Enable Deezer region check (default: `true`)
- `--service-reddit-enabled` / `REDDIT_ENABLED`: Enable Reddit region check (default: `true`)
- `--service-amazonprime-enabled` / `AMAZONPRIME_ENABLED`: Enable Amazon Prime region check (default: `true`)
- `--service-apple-enabled` / `APPLE_ENABLED`: Enable Apple region check (default: `true`)
- `--service-steam-enabled` / `STEAM_ENABLED`: Enable Steam region check (default: `true`)
- `--service-playstation-enabled` / `PLAYSTATION_ENABLED`: Enable PlayStation region check (default: `true`)
- `--service-tiktok-enabled` / `TIKTOK_ENABLED`: Enable TikTok region check (default: `true`)
- `--service-jetbrains-enabled` / `JETBRAINS_ENABLED`: Enable JetBrains region check (default: `true`)
- `--service-bing-enabled` / `BING_ENABLED`: Enable Microsoft (Bing) region check (default: `true`)

#### IP Detection

- `--service-ipinfo-url` / `IPINFO_URL`: IPInfo.io URL for IP detection (default: `https://ipinfo.io/json`)
- `--service-ipinfo-token` / `IPINFO_TOKEN`: IPInfo.io API token (optional, for higher rate limits)

### Advanced Configuration Example

```bash
./region-exporter \
  --metrics-host 0.0.0.0 \
  --metrics-port 9999 \
  --metrics-username admin \
  --metrics-password secure-password \
  --check-interval 60 \
  --check-timeout 15 \
  --check-expected-region US \
  --network-proxy socks5://127.0.0.1:1080 \
  --network-ipv4-only \
  --service-netflix-enabled false \
  --service-spotify-enabled false
```

## Usage

### Prometheus Metrics

The exporter exposes Prometheus metrics at the `/metrics` endpoint:

```bash
curl http://localhost:9999/metrics
```

#### Available Metrics

- `ip_service_region_match`: Service region detection match (1 if detected region matches expected, 0 otherwise)
  - Labels: `service`, `ip`, `detected_region`, `expected_region`
- `ip_service_latency_milliseconds`: Service request latency in milliseconds
  - Labels: `service`, `ip`

#### Prometheus Configuration

Add the following to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: "region-exporter"
    static_configs:
      - targets: ["localhost:9999"]
    basic_auth:
      username: "admin"
      password: "secure-password"
```

### REST API

The exporter provides a JSON API endpoint at `/api/status`:

```bash
curl http://localhost:9999/api/status
```

Response example:

```json
{
  "server_ip": "1.2.3.4",
  "expected_region": "US",
  "services": [
    {
      "name": "Google",
      "detected_region": "US",
      "expected_region": "US",
      "match": true,
      "enabled": true,
      "latency": 123.45
    },
    {
      "name": "Netflix",
      "detected_region": "GB",
      "expected_region": "US",
      "match": false,
      "enabled": true,
      "latency": 234.56
    }
  ],
  "total_services": 15,
  "match_count": 12,
  "mismatch_count": 3
}
```

### Health Check

A simple health check endpoint is available at `/health`:

```bash
curl http://localhost:9999/health
```

Returns `200 OK` if the service is running.

## Docker Compose Example

```yaml
version: "3.8"

services:
  region-exporter:
    image: ghcr.io/blancvpn/service-region-exporter:latest
    container_name: region-exporter
    ports:
      - "9999:9999"
    environment:
      - METRICS_HOST=0.0.0.0
      - METRICS_PORT=9999
      - CHECK_INTERVAL=60
      - EXPECTED_REGION=US
      - METRICS_USERNAME=admin
      - METRICS_PASSWORD=secure-password
    restart: unless-stopped
```

## Use Cases

- **VPN/Proxy Monitoring**: Verify that your VPN or proxy correctly routes traffic to the expected region
- **Geo-Restriction Testing**: Test whether services correctly detect your server's location
- **Service Availability**: Monitor which services are accessible from your server's region
- **Network Configuration**: Validate network routing and DNS configuration
- **Compliance Monitoring**: Ensure services detect the correct region for compliance purposes

## Tips

### Using with Proxy

If you need to check region detection through a proxy, configure the SOCKS5 proxy:

```bash
./region-exporter --network-proxy socks5://proxy.example.com:1080
```

### IPv4/IPv6 Configuration

Force IPv4 or IPv6 only:

```bash
# IPv4 only
./region-exporter --network-ipv4-only

# IPv6 only
./region-exporter --network-ipv6-only
```

### Custom Network Interface

Use a specific network interface:

```bash
./region-exporter --network-interface eth0
```

### Verbose Logging

Enable verbose logging for debugging:

```bash
./region-exporter --verbose
```

## Monitoring and Alerting

### Grafana Dashboard

You can create Grafana dashboards using the exposed metrics to visualize:

- Region match/mismatch rates per service
- Service latency trends
- Overall service availability

### Alerting Rules

Example Prometheus alerting rules:

```yaml
groups:
  - name: region_exporter
    rules:
      - alert: RegionMismatch
        expr: ip_service_region_match == 0
        for: 5m
        annotations:
          summary: "Service {{ $labels.service }} detected wrong region"
          description: "Service {{ $labels.service }} detected region {{ $labels.detected_region }} but expected {{ $labels.expected_region }}"

      - alert: HighServiceLatency
        expr: ip_service_latency_milliseconds > 5000
        for: 5m
        annotations:
          summary: "High latency for {{ $labels.service }}"
          description: "Service {{ $labels.service }} has latency of {{ $value }}ms"
```

## Contributing

We welcome contributions from the community! If you have ideas for improvements or have found a bug, please:

1. Create an issue on GitHub
2. Fork the repository
3. Create a feature branch
4. Make your changes
5. Submit a pull request

For major changes, please open an issue first to discuss what you would like to change.

## License

See [LICENSE](LICENSE) file for details.

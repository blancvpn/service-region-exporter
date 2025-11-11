# 🌍 Service Region Exporter

Prometheus exporter that monitors and verifies geographic region detection by various online services. It checks which region popular services detect for your server IP and exports metrics for monitoring.

## ✨ Features

- 🎯 **15+ services** — Google, YouTube, Netflix, Spotify, ChatGPT and more
- 🤖 **Auto region detection** — automatically detects your IP and expected region
- 📊 **Prometheus metrics** — standard metrics for monitoring system integration
- 🔌 **REST API** — JSON endpoint for service status
- ⚙️ **Flexible configuration** — enable/disable individual service checks
- 🌐 **Network support** — SOCKS5 proxy, custom interfaces, IPv4/IPv6
- ⏱️ **Scheduled checks** — configurable check intervals
- 🔐 **Basic Authentication** — optional metrics endpoint protection
- 🐳 **Docker ready** — official image available

## 🎯 Supported Services

Google • YouTube • ChatGPT • Netflix • Twitch • Spotify • Deezer • Reddit • Amazon Prime Video • Apple • Steam • PlayStation • TikTok • JetBrains • Microsoft Bing

## 🚀 Installation

### Docker (recommended)

```bash
docker run -d \
  --name region-exporter \
  -p 9999:9999 \
  -e CHECK_INTERVAL=60 \
  -e EXPECTED_REGION=US \
  blancvpn/service-region-exporter:latest
```

You can also use the image from GitHub Container Registry: `ghcr.io/blancvpn/service-region-exporter:latest`

### From release

```bash
# Download binary from GitHub Releases
chmod +x region-exporter
./region-exporter
```

### From source

```bash
git clone https://github.com/blancvpn/service-region-exporter.git
cd service-region-exporter
go build -o region-exporter .
./region-exporter
```

## ⚙️ Configuration

Configure via command-line flags or environment variables.

### Basic setup

```bash
# Flags
./region-exporter \
  --metrics-port=9999 \
  --check-interval=60 \
  --check-expected-region=US

# Environment variables
export METRICS_PORT=9999
export CHECK_INTERVAL=60
export EXPECTED_REGION=US
./region-exporter
```

### Configuration options

#### 📡 Metrics server

- `--metrics-host` / `METRICS_HOST` — HTTP server host (default: `0.0.0.0`)
- `--metrics-port` / `METRICS_PORT` — HTTP server port (default: `9999`)
- `--metrics-username` / `METRICS_USERNAME` — Basic Auth username
- `--metrics-password` / `METRICS_PASSWORD` — Basic Auth password

#### ⏱️ Check settings

- `--check-interval` / `CHECK_INTERVAL` — interval between checks in seconds (default: `60`)
- `--check-timeout` / `CHECK_TIMEOUT` — HTTP request timeout in seconds (default: `10`)
- `--check-expected-region` / `EXPECTED_REGION` — expected region code (`US`, `GB`, etc.)

#### 🌐 Network settings

- `--network-proxy` / `NETWORK_PROXY` — SOCKS5 proxy address (`host:port`)
- `--network-interface` / `NETWORK_INTERFACE` — network interface to use
- `--network-ipv4-only` / `NETWORK_IPV4_ONLY` — use IPv4 only
- `--network-ipv6-only` / `NETWORK_IPV6_ONLY` — use IPv6 only

#### 🎛️ Service toggles

Enable/disable individual services:

```bash
--service-google-enabled / GOOGLE_ENABLED
--service-youtube-enabled / YOUTUBE_ENABLED
--service-chatgpt-enabled / CHATGPT_ENABLED
--service-netflix-enabled / NETFLIX_ENABLED
# ... and others (all enabled by default)
```

### Advanced configuration example

```bash
./region-exporter \
  --metrics-port=9999 \
  --metrics-username=admin \
  --metrics-password=secure-pass \
  --check-interval=60 \
  --check-expected-region=US \
  --network-proxy=socks5://127.0.0.1:1080 \
  --network-ipv4-only=true   \
  --service-netflix-enabled=false
```

## 📊 Usage

### Prometheus metrics

Metrics available at `/metrics`:

```bash
curl http://localhost:9999/metrics
```

**Metrics:**

- `ip_service_region_match` — region match status (1 = match, 0 = mismatch)
  - Labels: `service`, `ip`, `detected_region`, `expected_region`
- `ip_service_latency_milliseconds` — request latency in milliseconds
  - Labels: `service`, `ip`

**Prometheus configuration:**

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

JSON endpoint at `/api/status`:

```bash
curl http://localhost:9999/api/status
```

**Example response:**

```json
{
  "server_ip": "1.2.3.4",
  "expected_region": "US",
  "services": [
    {
      "name": "Google",
      "detected_region": "US",
      "match": true,
      "enabled": true,
      "latency": 123.45
    }
  ],
  "total_services": 15,
  "match_count": 12,
  "mismatch_count": 3
}
```

### Health check

```bash
curl http://localhost:9999/health
```

## 🐳 Docker Compose

```yaml
version: "3.8"

services:
  region-exporter:
    image: blancvpn/service-region-exporter:latest
    container_name: region-exporter
    ports:
      - "9999:9999"
    environment:
      - CHECK_INTERVAL=60
      - EXPECTED_REGION=US
      - METRICS_USERNAME=admin
      - METRICS_PASSWORD=secure-password
    restart: unless-stopped
```

## 💡 Use Cases

- 🔒 **VPN/Proxy monitoring** — verify correct traffic routing
- 🌍 **Geo-restriction testing** — check service location detection
- ✅ **Service availability** — monitor accessibility from your region
- 🔧 **Network validation** — verify routing and DNS configuration
- 📋 **Compliance** — ensure correct region detection

## 🛠️ Tips

### Using with proxy

```bash
./region-exporter --network-proxy=socks5://proxy.example.com:1080
```

### IPv4/IPv6 only

```bash
./region-exporter --network-ipv4-only
./region-exporter --network-ipv6-only
```

### Custom network interface

```bash
./region-exporter --network-interface=eth0
```

## 🚨 Monitoring & Alerting

### Example alert rules

```yaml
groups:
  - name: region_exporter
    rules:
      - alert: RegionMismatch
        expr: ip_service_region_match == 0
        for: 5m
        annotations:
          summary: "{{ $labels.service }} detected wrong region"
          description: "Detected {{ $labels.detected_region }}, expected {{ $labels.expected_region }}"

      - alert: HighServiceLatency
        expr: ip_service_latency_milliseconds > 5000
        for: 5m
        annotations:
          summary: "High latency for {{ $labels.service }}"
          description: "Latency {{ $value }}ms"
```

## 🤝 Contributing

Contributions are welcome! To contribute:

1. Create an issue on GitHub
2. Fork the repository
3. Create a feature branch
4. Make your changes
5. Submit a pull request

For major changes, please open an issue first to discuss.

## 📄 License

See [LICENSE](LICENSE) file for details.

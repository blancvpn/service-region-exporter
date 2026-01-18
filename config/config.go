package config

import (
	"fmt"

	"github.com/alecthomas/kong"
)

var CLIConfig CLI

func Parse(version string) {
	ctx := kong.Parse(&CLIConfig,
		kong.Name("region-exporter"),
		kong.Description("Region Exporter: A Prometheus exporter for monitoring service region detection"),
		kong.Vars{
			"version": version,
		},
	)
	_ = ctx
}

type CLI struct {
	Metrics struct {
		Host     string `name:"host" help:"HTTP server host" default:"0.0.0.0" env:"METRICS_HOST"`
		Port     int    `name:"port" help:"HTTP server port" default:"9999" env:"METRICS_PORT"`
		Username string `name:"username" help:"Basic auth username (optional)" default:"" env:"METRICS_USERNAME"`
		Password string `name:"password" help:"Basic auth password (optional)" default:"" env:"METRICS_PASSWORD"`
	} `embed:"" prefix:"metrics-"`

	Check struct {
		Interval       int    `name:"interval" help:"Interval between checks in seconds" default:"60" env:"CHECK_INTERVAL"`
		Timeout        int    `name:"timeout" help:"HTTP request timeout in seconds" default:"10" env:"CHECK_TIMEOUT"`
		ExpectedRegion string `name:"expected-region" help:"Expected region code (e.g., US, GB). Auto-detected if not set" default:"" env:"EXPECTED_REGION"`
	} `embed:"" prefix:"check-"`

	Network struct {
		Proxy     string `name:"proxy" help:"SOCKS5 proxy address (host:port)" default:"" env:"NETWORK_PROXY"`
		Interface string `name:"interface" help:"Network interface to use" default:"" env:"NETWORK_INTERFACE"`
		IPv4Only  bool   `name:"ipv4-only" help:"Use only IPv4" default:"false" env:"NETWORK_IPV4_ONLY"`
		IPv6Only  bool   `name:"ipv6-only" help:"Use only IPv6" default:"false" env:"NETWORK_IPV6_ONLY"`
	} `embed:"" prefix:"network-"`

	Services struct {
		Google struct {
			Enabled bool   `name:"enabled" help:"Enable Google region check" default:"true" env:"GOOGLE_ENABLED"`
			URL     string `name:"url" help:"Google Play URL for region detection" default:"https://play.google.com/" env:"GOOGLE_URL"`
		} `embed:"" prefix:"google-"`

		YouTube struct {
			Enabled bool   `name:"enabled" help:"Enable YouTube region check" default:"true" env:"YOUTUBE_ENABLED"`
			URL     string `name:"url" help:"YouTube URL for region detection" default:"https://www.youtube.com" env:"YOUTUBE_URL"`
		} `embed:"" prefix:"youtube-"`

		ChatGPT struct {
			Enabled bool   `name:"enabled" help:"Enable ChatGPT region check" default:"true" env:"CHATGPT_ENABLED"`
			URL     string `name:"url" help:"ChatGPT API URL for region detection" default:"https://ab.chatgpt.com/v1/initialize" env:"CHATGPT_URL"`
			Token   string `name:"token" help:"ChatGPT API token" default:"client-zUdXdSTygXJdzoE0sWTkP8GKTVsUMF2IRM7ShVO2JAG" env:"CHATGPT_TOKEN"`
		} `embed:"" prefix:"chatgpt-"`

		Netflix struct {
			Enabled bool   `name:"enabled" help:"Enable Netflix region check" default:"true" env:"NETFLIX_ENABLED"`
			URL     string `name:"url" help:"Netflix API URL for region detection" default:"https://api.fast.com/netflix/speedtest/v2" env:"NETFLIX_URL"`
			Token   string `name:"token" help:"Netflix API token" default:"YXNkZmFzZGxmbnNkYWZoYXNkZmhrYWxm" env:"NETFLIX_TOKEN"`
		} `embed:"" prefix:"netflix-"`

		Twitch struct {
			Enabled bool   `name:"enabled" help:"Enable Twitch region check" default:"true" env:"TWITCH_ENABLED"`
			URL     string `name:"url" help:"Twitch GQL URL" default:"https://gql.twitch.tv/gql" env:"TWITCH_URL"`
			Token   string `name:"token" help:"Twitch Client ID token" default:"kimne78kx3ncx6brgo4mv6wki5h1ko" env:"TWITCH_TOKEN"`
		} `embed:"" prefix:"twitch-"`

		Spotify struct {
			Enabled bool   `name:"enabled" help:"Enable Spotify region check" default:"true" env:"SPOTIFY_ENABLED"`
			URL     string `name:"url" help:"Spotify URL" default:"https://accounts.spotify.com/status" env:"SPOTIFY_URL"`
		} `embed:"" prefix:"spotify-"`

		Deezer struct {
			Enabled bool   `name:"enabled" help:"Enable Deezer region check" default:"true" env:"DEEZER_ENABLED"`
			URL     string `name:"url" help:"Deezer URL" default:"https://www.deezer.com/en/offers" env:"DEEZER_URL"`
		} `embed:"" prefix:"deezer-"`

		Reddit struct {
			Enabled bool   `name:"enabled" help:"Enable Reddit region check" default:"true" env:"REDDIT_ENABLED"`
			URL     string `name:"url" help:"Reddit URL" default:"https://www.reddit.com" env:"REDDIT_URL"`
		} `embed:"" prefix:"reddit-"`

		AmazonPrime struct {
			Enabled bool   `name:"enabled" help:"Enable Amazon Prime region check" default:"true" env:"AMAZONPRIME_ENABLED"`
			URL     string `name:"url" help:"Amazon Prime URL" default:"https://www.primevideo.com" env:"AMAZONPRIME_URL"`
		} `embed:"" prefix:"amazonprime-"`

		Apple struct {
			Enabled bool   `name:"enabled" help:"Enable Apple region check" default:"true" env:"APPLE_ENABLED"`
			URL     string `name:"url" help:"Apple GCC URL" default:"https://gspe1-ssl.ls.apple.com/pep/gcc" env:"APPLE_URL"`
		} `embed:"" prefix:"apple-"`

		Steam struct {
			Enabled bool   `name:"enabled" help:"Enable Steam region check" default:"true" env:"STEAM_ENABLED"`
			URL     string `name:"url" help:"Steam URL" default:"https://store.steampowered.com" env:"STEAM_URL"`
		} `embed:"" prefix:"steam-"`

		PlayStation struct {
			Enabled bool   `name:"enabled" help:"Enable PlayStation region check" default:"true" env:"PLAYSTATION_ENABLED"`
			URL     string `name:"url" help:"PlayStation URL" default:"https://www.playstation.com" env:"PLAYSTATION_URL"`
		} `embed:"" prefix:"playstation-"`

		TikTok struct {
			Enabled bool   `name:"enabled" help:"Enable TikTok region check" default:"true" env:"TIKTOK_ENABLED"`
			URL     string `name:"url" help:"TikTok API URL" default:"https://www.tiktok.com/api/v1/web-cookie-privacy/config?appId=1988" env:"TIKTOK_URL"`
		} `embed:"" prefix:"tiktok-"`

		JetBrains struct {
			Enabled bool   `name:"enabled" help:"Enable JetBrains region check" default:"true" env:"JETBRAINS_ENABLED"`
			URL     string `name:"url" help:"JetBrains Geo URL" default:"https://data.services.jetbrains.com/geo" env:"JETBRAINS_URL"`
		} `embed:"" prefix:"jetbrains-"`

		Bing struct {
			Enabled bool   `name:"enabled" help:"Enable Microsoft (Bing) region check" default:"true" env:"BING_ENABLED"`
			URL     string `name:"url" help:"Bing Search URL" default:"https://www.bing.com/search?q=test" env:"BING_URL"`
		} `embed:"" prefix:"bing-"`

		IPInfo struct {
			URL   string `name:"url" help:"IPInfo.io URL for IP detection" default:"https://ipinfo.io/json" env:"IPINFO_URL"`
			Token string `name:"token" help:"IPInfo.io API token (optional)" default:"" env:"IPINFO_TOKEN"`
		} `embed:"" prefix:"ipinfo-"`

		AbuseIPDB struct {
			Token string `name:"token" help:"AbuseIPDB API token (enables abuse score check if set)" default:"" env:"ABUSEIPDB_TOKEN"`
		} `embed:"" prefix:"abuseipdb-"`
	} `embed:"" prefix:"service-"`

	Version VersionFlag `name:"version" help:"Print version information and quit"`
	Verbose bool        `name:"verbose" help:"Enable verbose logging" default:"false" env:"VERBOSE"`
}

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println("Region Exporter: A Prometheus exporter for monitoring services region detection")
	fmt.Printf("Version:\t %s\n", vars["version"])
	fmt.Printf("GitHub: https://github.com/blancvpn/service-region-exporter\n")
	fmt.Println("Made with ❤️ by BlancVPN")
	app.Exit(0)
	return nil
}

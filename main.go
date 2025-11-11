package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"

	"region-exporter/config"
	"region-exporter/metrics"
	"region-exporter/models"
	netutil "region-exporter/net"
	"region-exporter/services"
	"region-exporter/web"
)

var (
	version = "unknown"
)

type App struct {
	client         *http.Client
	services       []services.Service
	serverIP       string
	expectedRegion string
	servicesData   []models.ServiceStatus
}

func NewApp() *App {
	cfg := &config.CLIConfig

	transport, err := netutil.CreateTransport(cfg)
	if err != nil {
		log.Fatalf("Failed to create HTTP transport: %v", err)
	}

	client := &http.Client{
		Timeout:   time.Duration(cfg.Check.Timeout) * time.Second,
		Transport: transport,
	}

	app := &App{
		client: client,
	}

	app.services = services.GetAllServices(client)

	return app
}

func (a *App) runChecks() {
	cfg := &config.CLIConfig
	log.Println("Starting region checks...")

	a.servicesData = []models.ServiceStatus{}

	for _, svc := range a.services {
		status := services.CheckService(svc, cfg, a.serverIP, a.expectedRegion)
		if status.Enabled {
			a.servicesData = append(a.servicesData, status)
		}
	}

	web.UpdateData(a.serverIP, a.expectedRegion, a.servicesData)

	log.Println("Checks completed")
}

func (a *App) startScheduler(ctx context.Context) {
	cfg := &config.CLIConfig
	s := gocron.NewScheduler(time.UTC)

	interval := time.Duration(cfg.Check.Interval) * time.Second
	_, err := s.Every(interval).StartImmediately().Do(a.runChecks)
	if err != nil {
		log.Fatalf("Failed to schedule checks: %v", err)
	}

	log.Printf("Scheduler started with interval: %s", interval)

	s.StartAsync()

	<-ctx.Done()
	s.Stop()
	log.Println("Scheduler stopped")
}

func (a *App) Run() error {
	cfg := &config.CLIConfig

	serverIP, expectedRegion, err := netutil.DetectServerInfo(a.client, cfg)
	if err != nil {
		return err
	}
	a.serverIP = serverIP
	a.expectedRegion = expectedRegion

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	go a.startScheduler(ctx)

	return web.StartHTTPServer(ctx, cfg)
}

func main() {
	config.Parse(version)

	metrics.Init()

	app := NewApp()

	if err := app.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

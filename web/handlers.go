package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"region-exporter/config"
	"region-exporter/metrics"
	"region-exporter/models"
)

var (
	currentData models.DashboardData
	dataMutex   sync.RWMutex
)

func UpdateData(serverIP, expectedRegion string, servicesData []models.ServiceStatus) {
	dataMutex.Lock()
	defer dataMutex.Unlock()

	matchCount := 0
	mismatchCount := 0
	enabledCount := 0

	for _, svc := range servicesData {
		if svc.Enabled {
			enabledCount++
			if svc.Match {
				matchCount++
			} else {
				mismatchCount++
			}
		}
	}

	currentData = models.DashboardData{
		ServerIP:       serverIP,
		ExpectedRegion: expectedRegion,
		Services:       servicesData,
		TotalServices:  enabledCount,
		MatchCount:     matchCount,
		MismatchCount:  mismatchCount,
	}
}

func APIStatusHandler(w http.ResponseWriter, r *http.Request) {
	dataMutex.RLock()
	data := currentData
	dataMutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func BasicAuthMiddleware(next http.Handler, username, password string) http.Handler {
	if username == "" || password == "" {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != username || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func StartHTTPServer(ctx context.Context, cfg *config.CLI) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/status", APIStatusHandler)

	mux.Handle("/metrics", metrics.Handler())

	mux.HandleFunc("/health", HealthHandler)

	var handler http.Handler = mux
	if cfg.Metrics.Username != "" && cfg.Metrics.Password != "" {
		log.Printf("Basic authentication enabled for user: %s", cfg.Metrics.Username)
		authMux := http.NewServeMux()
		authMux.HandleFunc("/health", HealthHandler)
		authMux.Handle("/", BasicAuthMiddleware(mux, cfg.Metrics.Username, cfg.Metrics.Password))
		handler = authMux
	}

	addr := fmt.Sprintf("%s:%d", cfg.Metrics.Host, cfg.Metrics.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		<-ctx.Done()
		log.Println("Shutting down HTTP server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}()

	log.Printf("HTTP server starting on %s", addr)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server error: %w", err)
	}

	return nil
}

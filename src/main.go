package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// 1. Port configuration (12-Factor App)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 2. Initialize HTTP router
	mux := http.NewServeMux()

	// Expose Prometheus metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())

	// Application routes
	mux.HandleFunc("/healthz", HealthzHandler)
	mux.HandleFunc("/api/v1/hello", HelloHandler)
	mux.HandleFunc("/api/v1/toggle-error", ToggleErrorHandler)

	// Wrap router with RED metrics middleware
	handlerWithMetrics := MetricsMiddleware(mux)

	// 3. Configure HTTP server with production-grade timeouts
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handlerWithMetrics,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// 4. Start HTTP server in a background goroutine
	go func() {
		log.Printf("Server is running on port %s...", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// 5. Channel to listen for OS interrupt signals (SIGINT / SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-quit
	log.Println("Shutdown signal received, initiating graceful shutdown...")

	// 6. Graceful shutdown context with 10-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped.")
}

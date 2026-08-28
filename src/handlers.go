package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

// HealthzHandler: for kubernetes healthcheck
func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Header to application/json 
	w.WriteHeader(http.StatusOK)                       // 200
	_ = json.NewEncoder(w).Encode(map[string]string{    // { "status": "UP" } JSON
		"status": "UP",
	})
}

var simulateError int32 = 0

func HelloHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	if atomic.LoadInt32(&simulateError) == 1 {
		w.WriteHeader(http.StatusInternalServerError) //500

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "simulated_internal_server_error",
			"message": "Simulated test for canary rollback!",
		})
		return
	}

	w.WriteHeader(http.StatusOK) // 200
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Hello from Cloud-Native API!",
		"version":   "v2.0.0",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}


func ToggleErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	current := atomic.LoadInt32(&simulateError)
	var newState int32 = 0
	if current == 0 {
		newState = 1
	}

	atomic.StoreInt32(&simulateError, newState)
	
	statusText := "DISABLED (Returning 200 OK)"
	if newState == 1 {
		statusText = "ENABLED (Returning 500 Internal Server Error)"
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"simulate_error": fmt.Sprintf("%d", newState),
		"status": statusText,          
	})
}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *statusResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		sw := &statusResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // default 200
		}

		next.ServeHTTP(sw, r)

		duration := time.Since(start).Seconds()
		statusCodeStr := strconv.Itoa(sw.statusCode)

		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, statusCodeStr).Inc()
	})
}

	
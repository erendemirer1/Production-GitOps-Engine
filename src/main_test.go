package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

//Healthcheck test
func TestHealthzHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	HealthzHandler(rr,req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response json: %v", err)
	}

	if response["status"] != "UP" {
		t.Errorf("expected status 'UP', got '%s'", response["status"])
	}
}

//At normal mode must return 200, at error mode must 500 
func TestHelloHandler(t *testing.T) {
	//1-Normal mod test
	atomic.StoreInt32(&simulateError, 0)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/hello", nil)
	rr := httptest.NewRecorder()

	HelloHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 in normal mode, got %d", rr.Code)
	}
	//2-Error mod test
	atomic.StoreInt32(&simulateError, 1)
	rrError := httptest.NewRecorder()

	HelloHandler(rrError, req)

	if rrError.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500 in error mode, got %d", rrError.Code)
	}

	//back to default
	atomic.StoreInt32(&simulateError, 0)
}

//test for /toggle-error
func TestToggleErrorHandler(t *testing.T) {
	atomic.StoreInt32(&simulateError, 0)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/toggle-error", nil)
	rr := httptest.NewRecorder()
	ToggleErrorHandler(rr, req)
	if atomic.LoadInt32(&simulateError) != 1 {
		t.Errorf("expected simulateError to be 1 after toggle, got %d", simulateError)
	}

	ToggleErrorHandler(rr, req)
	if atomic.LoadInt32(&simulateError) != 0 {
		t.Errorf("expected simulateError to be 0 after second toggle, got %d", simulateError)
	}
}
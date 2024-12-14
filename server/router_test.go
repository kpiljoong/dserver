package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDynamicHandler(t *testing.T) {
	SetRouter(nil) // Ensure no router is configured initially
	handler := &DynamicHandler{}

	// Test with no router configured
	req := httptest.NewRequest("GET", "/test", nil)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d but got %d", http.StatusInternalServerError, resp.Code)
	}
	if resp.Body.String() != "No router configured\n" {
		t.Errorf("Expected body 'No router configured' but got %s", resp.Body.String())
	}

	// Update the router and re-test
	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	SetRouter(mux)

	req = httptest.NewRequest("GET", "/test", nil)
	resp = httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("Expected status %d but got %d", http.StatusOK, resp.Code)
	}
	if resp.Body.String() != "OK" {
		t.Errorf("Expected body %s but got %s", "OK", resp.Body.String())
	}
}

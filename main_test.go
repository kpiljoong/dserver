package main

import (
  "os"
  "net/http"
  "net/http/httptest"
  "sync/atomic"
  "testing"
  "time"
)

var TestConfig = []RouteConfig{
  {
    Path: "/api/v1/heatlh",
    Method: "GET",
    StatusCode: 200,
    ResponseBody: `{"status": "healthy"}`,
    ContentType: "application/json",
  },
  {
    Path: "/api/v1/resource",
    Method: "GET",
    StatusCode: 200,
    ResponseBody: `{"message": "Resource retrieved successfully"}`,
    ContentType: "application/json",
  },
  {
    Path: "/api/v1/resource",
    Method: "POST",
    StatusCode: 201,
    ResponseBody: `{"message": "Resource created successfully"}`,
    ContentType: "application/json",
  },
}

func TestWatchConfigFile(t *testing.T) {
  // Use a temp file for testing
  tmpFile, err := os.CreateTemp("", "config-*.toml")
  if err != nil {
    t.Fatalf("Failed to create temp file: %v", err)
  }
  defer os.Remove(tmpFile.Name())

  var reloadCalled int32

  mockReloadConfig := func() {
    atomic.AddInt32(&reloadCalled, 1)
  }

  go func() {
    watchConfigFile(tmpFile.Name(), mockReloadConfig)
  }()

  // Simulate file changes
  time.Sleep(100 * time.Millisecond)
  os.WriteFile(tmpFile.Name(), []byte("test"), 0644)
  time.Sleep(200 * time.Millisecond)

  os.Rename(tmpFile.Name(), tmpFile.Name()+"-renamed")
  time.Sleep(200 * time.Millisecond)

  os.WriteFile(tmpFile.Name()+"-renamed", []byte("test content"), 0644)
  time.Sleep(200 * time.Millisecond)

  reloadCount := atomic.LoadInt32(&reloadCalled)
  if reloadCount != 3 {
    t.Errorf("Expected reload count of 3, got %d", reloadCount)
  }
}

func TestHealthEndpoint(t *testing.T) {
  handler := createHandler([]RouteConfig{TestConfig[0]})

  req := httptest.NewRequest("GET", "/api/v1/health", nil)
  w := httptest.NewRecorder()

  handler(w, req)

  resp := w.Result()
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
  }

  expectedContentType := "application/json"
  if resp.Header.Get("Content-Type") != expectedContentType {
    t.Errorf("Expected content type %s, got %s", expectedContentType, resp.Header.Get("Content-Type"))
  }
}

// TestResourceEndpointGet verifies the /api/v1/resource endpoint for GET requests.
func TestResourceEndpointGet(t *testing.T) {
	handler := createHandler([]RouteConfig{TestConfig[1]})

	req := httptest.NewRequest("GET", "/api/v1/resource", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	expectedContentType := "application/json"
	if resp.Header.Get("Content-Type") != expectedContentType {
		t.Errorf("Expected Content-Type %s, got %s", expectedContentType, resp.Header.Get("Content-Type"))
	}
}

// TestResourceEndpointPost verifies the /api/v1/resource endpoint for POST requests.
func TestResourceEndpointPost(t *testing.T) {
	handler := createHandler([]RouteConfig{TestConfig[2]})

	req := httptest.NewRequest("POST", "/api/v1/resource", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	expectedContentType := "application/json"
	if resp.Header.Get("Content-Type") != expectedContentType {
		t.Errorf("Expected Content-Type %s, got %s", expectedContentType, resp.Header.Get("Content-Type"))
	}
}

// TestInvalidMethod tests an unsupported method on /api/v1/resource.
func TestInvalidMethod(t *testing.T) {
	handler := createHandler([]RouteConfig{TestConfig[1], TestConfig[2]})

	req := httptest.NewRequest("PUT", "/api/v1/resource", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", resp.StatusCode)
	}
}

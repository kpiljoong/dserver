package config

import (
	"os"
	"testing"
)

func TestReloadConfig(t *testing.T) {
	tempFile, err := os.CreateTemp("", "config-*.toml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	configContent := `
[[routes]]
path = "/api/v1/resource"
method = "GET"
status_code = 200
response_body = "Hello, World!"
`
	if _, err := tempFile.Write([]byte(configContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Reload the config
	err = ReloadConfig(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to reload config: %v", err)
	}

	// Verify the loaded config
	ConfigMutex.RLock()
	defer ConfigMutex.RUnlock()

	if len(CurrentConfig.Routes) != 1 {
		t.Fatalf("Expected 1 route, got %d", len(CurrentConfig.Routes))
	}

	route := CurrentConfig.Routes[0]
	if route.Path != "/api/v1/resource" {
		t.Errorf("Expected path /api/v1/resource, got %s", route.Path)
	}
	if route.Method != "GET" {
		t.Errorf("Expected method GET, got %s", route.Method)
	}
	if route.ResponseBody != "Hello, World!" {
		t.Errorf("Expected response body 'Hello, World!', got %s", route.ResponseBody)
	}
}

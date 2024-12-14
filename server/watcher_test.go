package server

import (
	"os"
	"testing"
	"time"
)

func TestWatchConfigFile(t *testing.T) {
	// Create a temporary file to simulate the config file
	tempFile, err := os.CreateTemp("", "config-*.toml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	content := `
[[routes]]
path = "/api/v1/resource"
method = "GET"
status_code = 200
response_body = "Hello, Watcher!"
`

	if _, err := tempFile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	reloaded := make(chan bool)

	go WatchConfigFile(tempFile.Name(), func() error {
		reloaded <- true
		return nil
	}, true)

	// Modify the file
	time.Sleep(100 * time.Millisecond)
	if err := os.WriteFile(tempFile.Name(), []byte(content), 0644); err != nil {
		t.Fatalf("Failed to modify temp file: %v", err)
	}

	select {
	case <-reloaded:
		// Success
	case <-time.After(1 * time.Second):
		t.Error("Expected reload event, but it didn't occur")
	}
}

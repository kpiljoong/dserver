package utils

import (
	"net/url"
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	key := "TEST_ENV_VAR"
	fallback := "default_value"

	value := GetEnv(key, fallback)
	if value != fallback {
		t.Errorf("Expected %s, got %s", fallback, value)
	}

	expectedValue := "env_value"
	os.Setenv(key, expectedValue)
	defer os.Unsetenv(key)

	value = GetEnv(key, fallback)
	if value != expectedValue {
		t.Errorf("Expected %s, got %s", expectedValue, value)
	}
}

func TestMatchQueryParams(t *testing.T) {
	query := url.Values{
		"type": []string{"active"},
	}

	expectedParams := map[string]string{"type": "active"}
	if !MatchQueryParams(query, expectedParams) {
		t.Error("Expected query params to match")
	}

	expectedParams["type"] = "inactive"
	if MatchQueryParams(query, expectedParams) {
		t.Error("Expected query params not to match")
	}

	// Test empty expectedParams
	if !MatchQueryParams(query, map[string]string{}) {
		t.Error("Expected match with empty expectedParams")
	}

	// Test empty query
	if MatchQueryParams(url.Values{}, expectedParams) {
		t.Error("Expected no match with empty query")
	}

	if !MatchQueryParams(url.Values{}, map[string]string{}) {
		t.Error("Expected match with both empty")
	}
}

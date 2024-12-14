package utils

import (
	"log"
	"net/http"
	"os"
	"time"
)

func ApplyDelay(delayMs int) {
	if delayMs > 0 {
		time.Sleep(time.Duration(delayMs) * time.Millisecond)
	}
}

func WriteResponse(w http.ResponseWriter, statusCode int, contentType string, body string) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	w.Write([]byte(body))
}

// MatchQueryParams checks if the request's query parameters match the expected query parameters.
func MatchQueryParams(requestParams map[string][]string, expectedParams map[string]string) bool {
	// If expectedParams is empty, always return true
	if len(expectedParams) == 0 {
		return true
	}

	for key, value := range expectedParams {
		if reqValues, ok := requestParams[key]; !ok || len(reqValues) == 0 || reqValues[0] != value {
			log.Printf("No match for key: %s, value: %s\n", key, value)
			return false
		}
	}
	return true
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/pelletier/go-toml"
)

type QueryMatch map[string]string

type ResponseConfig struct {
	Query        QueryMatch `toml:"query"`
	StatusCode   *int       `toml:"status_code"`
	ContentType  *string    `toml:"content_type"`
	ResponseBody string     `toml:"response_body"`
	DelayMs      *int       `toml:"delay_ms"`
}

type RouteConfig struct {
	Path         string           `toml:"path"`
	Method       string           `toml:"method"`
	StatusCode   int              `toml:"status_code"`
	ContentType  string           `toml:"content_type"`
	DelayMs      int              `toml:"delay_ms"`      // Optional delay in milliseconds
	ResponseBody string           `toml:"response_body"` // Fallback response body
	Responses    []ResponseConfig `toml:"responses"`
}

type Config struct {
	Routes []RouteConfig `toml:"routes"`
}

var (
	ConfigMutex   sync.RWMutex
	CurrentConfig *Config
)

func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func ReloadConfig(filePath string) error {
	newConfig, err := LoadConfig(filePath)
	if err != nil {
		return err
	}

	ConfigMutex.Lock()
	defer ConfigMutex.Unlock()
	CurrentConfig = newConfig
	return nil
}

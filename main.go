package main

import (
  "flag"
  "fmt"
  "log"
  "net/http"
  "os"
  "time"

  "github.com/pelletier/go-toml"
)

type RouteConfig struct {
  Path string `toml:"path"`
  Method string `toml:"method"`
  StatusCode int `toml:"status_code"`
  ResponseBody string `toml:"response_body"`
  ContentType string `toml:"content_type"`
  DelayMs int `toml:"delay_ms"` // Optional delay in milliseconds
}

type Config struct {
  Routes []RouteConfig `toml:"routes"`
}

func loadConfig(filePath string) (*Config, error) {
  file, err := os.ReadFile(filePath)
  if err != nil {
    return nil, err
  }

  var config Config
  err = toml.Unmarshal(file, &config)
  if err != nil {
    return nil, err
  }
  return &config, nil
}

func createHandler(routeConfigs []RouteConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
    for _, route := range routeConfigs {
      if r.Method == route.Method {
        if route.DelayMs > 0 {
          time.Sleep(time.Duration(route.DelayMs) * time.Millisecond)
        }

        w.Header().Set("Content-Type", route.ContentType)
        w.WriteHeader(route.StatusCode)
        _, err := w.Write([]byte(route.ResponseBody))
        if err != nil {
          log.Printf("Error writing response: %v\n", err)
        }
        return
      }
    }
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getEnv(key, fallback string) string {
  if value, ok := os.LookupEnv(key); ok {
    return value
  }
  return fallback
}

func main() {
  // Use environment variables with a fallback to flags
  defaultConfigPath := getEnv("DSERVER_CONFIG", "config.toml")
  defaultPort := getEnv("DSERVER_PORT", "8080")

  configPath := flag.String("config", defaultConfigPath, "Path to the config file")
  port := flag.String("port", defaultPort, "Port to run the server on")
  verbose := flag.Bool("verbose", false, "Enable verbose logging")

  flag.Parse()

  config, err := loadConfig(*configPath)
  if err != nil {
    log.Fatalf("Failed to load config: %v\n", err)
  }

  if *verbose {
    log.Printf("Using configuration file: %s", *configPath)
    log.Printf("Server will run on port: %s", *port)
  }

  // Group routes by path
  routeMap := make(map[string][]RouteConfig)
  for _, route := range config.Routes {
    routeMap[route.Path] = append(routeMap[route.Path], route)
  }

  // Register handlers
  for path, routes := range routeMap {
    if *verbose {
      log.Printf("Registering route: %s\n", path)
    }
    http.HandleFunc(path, createHandler(routes))
  }

  fmt.Printf("DServer running on port %s\n", *port)
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}


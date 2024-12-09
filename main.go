package main

import (
  "flag"
  "fmt"
  "log"
  "net/http"
  "os"
  "sync"
  "time"

  "github.com/fsnotify/fsnotify"
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

var (
  mux = http.NewServeMux()
  configPath string
  port string
  configMutex sync.RWMutex
  currentConfig *Config
  verbose bool
)

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

func reloadConfig() {
  configMutex.Lock()
  defer configMutex.Unlock()

  config, err := loadConfig(configPath)
  if err != nil {
    log.Printf("Failed to reload config: %v\n", err)
    return
  }

  if verbose {
    log.Println("Reloading config...")
  }
  currentConfig = config
  registerRoutes()
}

func registerRoutes() {
  mux = http.NewServeMux()

  routeMap := make(map[string][]RouteConfig)
  for _, route := range currentConfig.Routes {
    routeMap[route.Path] = append(routeMap[route.Path], route)
  }

  // Register handlers
  for path, routes := range routeMap {
    if verbose {
      log.Printf("Registering route: %s\n", path)
    }
    mux.HandleFunc(path, createHandler(routes))
  }
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

func watchConfigFile(filePath string, reloadConfigFunc func()) {
  watcher, err := fsnotify.NewWatcher()
  if err != nil {
    log.Fatalf("Failed to create watcher: %v\n", err)
  }
  defer watcher.Close()

  addFileToWatcher(watcher, filePath)

  for {
    select {
    case event := <-watcher.Events:
      if event.Op&fsnotify.Write == fsnotify.Write {
        if verbose {
          log.Println("Config file modified. Reloading...")
        }
        reloadConfigFunc()
      } else if event.Op&fsnotify.Rename == fsnotify.Rename {
        if verbose {
          log.Println("Config file renamed. Re-adding watcher...")
        }
        time.Sleep(100 * time.Millisecond)
        addFileToWatcher(watcher, configPath)
        reloadConfigFunc()
      } else if event.Op&fsnotify.Create == fsnotify.Create {
        if verbose {
          log.Println("Config file recreated. Reloading...")
        }
        addFileToWatcher(watcher, configPath)
        reloadConfigFunc()
      }
    case err := <-watcher.Errors:
      log.Printf("Watcher error: %v\n", err)
    }
  }
}

func addFileToWatcher(watcher *fsnotify.Watcher, path string) {
  err := watcher.Add(path)
  if err != nil {
    log.Printf("Failed to add file to watcher: %v\n", err)
  } else {
    if verbose {
      log.Printf("Watching file: %s", path)
    }
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

  flag.StringVar(&configPath, "config", defaultConfigPath, "Path to the config file")
  flag.StringVar(&port, "port", defaultPort, "Port to run the server on")
  flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")

  flag.Parse()

  var err error
  currentConfig, err = loadConfig(configPath)
  if err != nil {
    log.Fatalf("Failed to load config: %v\n", err)
  }
  if verbose {
    log.Printf("Using configuration file: %s", configPath)
    log.Printf("Server will run on port: %s", port)
  }

  registerRoutes()

  go watchConfigFile(configPath, reloadConfig)

  fmt.Printf("DServer running on port %s\n", port)
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    configMutex.RLock()
    defer configMutex.RUnlock()
    mux.ServeHTTP(w, r)
  })))
}


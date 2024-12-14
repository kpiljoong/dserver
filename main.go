package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"dserver/config"
	"dserver/server"
	"dserver/utils"
)

var verbose bool

func main() {
	// Use environment variables with a fallback to flags
	defaultConfigPath := utils.GetEnv("DSERVER_CONFIG", "config.toml")
	defaultPort := utils.GetEnv("DSERVER_PORT", "8080")

	configPath := flag.String("config", defaultConfigPath, "Path to the config file")
	port := flag.String("port", defaultPort, "Port to run the server on")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	flag.Parse()

	// Load initial config
	if err := config.ReloadConfig(*configPath); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if verbose {
		log.Printf("Configuration loaded: %s", *configPath)
	}

	server.RegisterRoutes(config.CurrentConfig, verbose)

	// Watch the configuration file for changes
	go server.WatchConfigFile(*configPath, func() error {
		err := config.ReloadConfig(*configPath)
		if err == nil {
			server.RegisterRoutes(config.CurrentConfig, verbose)
		}
		return err
	}, verbose)

	address := fmt.Sprintf(":%s", *port)
	log.Printf("DServer running on %s", address)
	log.Fatal(http.ListenAndServe(address, &server.DynamicHandler{}))
}

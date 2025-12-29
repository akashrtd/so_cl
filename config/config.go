package config

// Package config handles configuration loading from environment variables
// and default values for the so_cl application.

import (
	"os"
	"strconv"
)

// Config represents the configuration for so_cl application.
type Config struct {
	// DataDir is the directory where so_cl stores data
	DataDir string
	// ListenPort is the port to listen on for SSB connections
	ListenPort int
	// NetworkKey is the SSB network key (default: empty for mainnet)
	NetworkKey string
	// EnableLANDiscovery enables UDP broadcast for local peer discovery
	EnableLANDiscovery bool
	// LogLevel is the logging level (debug, info, warn, error)
	LogLevel string
	// Debug enables verbose logging
	Debug bool
}

// DefaultConfig returns a configuration with reasonable defaults.
func DefaultConfig() Config {
	return Config{
		DataDir:            homeDir() + "/.so_cl",
		ListenPort:         8008,
		NetworkKey:         "", // Use default SSB network
		EnableLANDiscovery: true,
		LogLevel:           "info",
		Debug:              false,
	}
}

// Load loads configuration from environment variables.
// Environment variables:
//   - SO_CL_DATA_DIR: Data directory (default: ~/.so_cl)
//   - SO_CL_PORT: Listen port (default: 8008)
//   - SO_CL_NETWORK_KEY: SSB network key (default: empty)
//   - SO_CL_ENABLE_LAN_DISCOVERY: Enable LAN discovery (default: true)
//   - SO_CL_LOG_LEVEL: Log level (default: info)
//   - SO_CL_DEBUG: Enable debug logging (default: false)
func Load() Config {
	cfg := DefaultConfig()

	if dataDir := os.Getenv("SO_CL_DATA_DIR"); dataDir != "" {
		cfg.DataDir = dataDir
	}

	if port := os.Getenv("SO_CL_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.ListenPort = p
		}
	}

	if networkKey := os.Getenv("SO_CL_NETWORK_KEY"); networkKey != "" {
		cfg.NetworkKey = networkKey
	}

	if enableLAN := os.Getenv("SO_CL_ENABLE_LAN_DISCOVERY"); enableLAN != "" {
		cfg.EnableLANDiscovery = enableLAN == "true" || enableLAN == "1"
	}

	if logLevel := os.Getenv("SO_CL_LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}

	if debug := os.Getenv("SO_CL_DEBUG"); debug != "" {
		cfg.Debug = debug == "true" || debug == "1"
	}

	return cfg
}

// homeDir returns the user's home directory.
func homeDir() string {
	if dir := os.Getenv("HOME"); dir != "" {
		return dir
	}
	return os.Getenv("USERPROFILE") // Windows
}

package scuttlego

// Package scuttlego provides a wrapper around the scuttlego library
// for easier integration with the so_cl application.
//
// This package handles:
// - Service lifecycle initialization
// - Configuration management
// - Adapters to convert between scuttlego and so_cl types

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// Service wraps the scuttlego service for use in so_cl.
// It provides a simplified interface for common operations like
// publishing posts, following peers, and querying messages.
type Service struct {
	// config holds the service configuration
	config Config
	// logger is the structured logger for this service
	logger *zap.Logger
}

// Config represents the configuration for the scuttlego service.
type Config struct {
	// DataDir is the directory where scuttlego stores data
	DataDir string
	// ListenPort is the port to listen on for SSB connections
	ListenPort int
	// NetworkKey is the SSB network key (default: SSB network)
	NetworkKey string
	// EnableLANDiscovery enables UDP broadcast for local peer discovery
	EnableLANDiscovery bool
}

// DefaultConfig returns a reasonable default configuration.
func DefaultConfig() Config {
	return Config{
		DataDir:            "~/.so_cl/data",
		ListenPort:         8008,
		NetworkKey:         "", // Use default SSB network key
		EnableLANDiscovery: true,
	}
}

// NewService creates a new scuttlego service with the given configuration.
//
// It initializes the scuttlego service with all required components:
// - BadgerDB storage
// - Network listener
// - EBT replication
// - LAN discovery (if enabled)
//
// Returns an error if service initialization fails.
func NewService(cfg Config, logger *zap.Logger) (*Service, error) {
	logger.Info("Initializing scuttlego service",
		zap.String("data_dir", cfg.DataDir),
		zap.Int("port", cfg.ListenPort),
	)

	// TODO: Implement actual scuttlego service initialization
	// This requires:
	// 1. Setting up BadgerDB
	// 2. Configuring network settings
	// 3. Creating service.Application
	// 4. Setting up listener
	// 5. Setting up EBT replication
	// 6. Setting up LAN discovery (if enabled)

	return &Service{
		config: cfg,
		logger: logger,
	}, nil
}

// Run starts the scuttlego service and blocks until context is cancelled.
//
// This method:
// 1. Starts the network listener
// 2. Starts EBT replication
// 3. Starts LAN discovery (if enabled)
// 4. Blocks until context is cancelled
// 5. Gracefully shuts down all components
//
// Returns an error if any component fails to start.
func (s *Service) Run(ctx context.Context) error {
	s.logger.Info("Starting scuttlego service")

	// TODO: Implement actual service run logic
	// This requires:
	// 1. Starting the listener
	// 2. Starting the EBT replication loop
	// 3. Starting the LAN discovery (if enabled)
	// 4. Handling shutdown on context cancellation

	<-ctx.Done()
	s.logger.Info("Stopping scuttlego service")
	return nil
}

// Close gracefully shuts down the scuttlego service.
func (s *Service) Close() error {
	s.logger.Info("Closing scuttlego service")

	// TODO: Implement graceful shutdown
	// This requires:
	// 1. Closing the listener
	// 2. Waiting for in-flight operations
	// 3. Closing BadgerDB
	// 4. Flushing any buffered data

	return nil
}

// Publish publishes a new post to the SSB feed.
// The text must be 1-280 characters and contain only ASCII characters.
//
// Returns the message reference if successful, or an error if publishing fails.
func (s *Service) Publish(text string) (string, error) {
	if len(text) == 0 {
		return "", fmt.Errorf("empty post text")
	}
	if len(text) > 280 {
		return "", fmt.Errorf("post text exceeds 280 character limit")
	}

	s.logger.Info("Publishing post",
		zap.String("text", text),
		zap.Int("length", len(text)),
	)

	// TODO: Implement publish using scuttlego PublishRaw command
	// This requires:
	// 1. Building SSB post content
	// 2. Calling PublishRaw command
	// 3. Returning message reference

	return "", nil
}

// Follow starts following the given peer feed.
// The feedRef is a SSB feed reference (e.g., "@alice...").
//
// Returns an error if following fails.
func (s *Service) Follow(feedRef string) error {
	s.logger.Info("Following peer",
		zap.String("feed_ref", feedRef),
	)

	// TODO: Implement follow using scuttlego Follow command
	// This requires:
	// 1. Parsing feed reference
	// 2. Calling Follow command
	// 3. Returning error if it fails

	return nil
}

// Connect connects to a peer at the given address.
// The address is in multiserver format (e.g., "net:127.0.0.1:8008~shs:@alice...").
//
// Returns an error if connection fails.
func (s *Service) Connect(address string) error {
	s.logger.Info("Connecting to peer",
		zap.String("address", address),
	)

	// TODO: Implement connect using scuttlego Connect command
	// This requires:
	// 1. Parsing address
	// 2. Calling Connect command
	// 3. Returning error if it fails

	return nil
}

// GetRecentMessages retrieves recent messages from the feed.
// The limit specifies the maximum number of messages to retrieve.
//
// Returns a slice of message references, or an error if retrieval fails.
func (s *Service) GetRecentMessages(limit int) ([]string, error) {
	s.logger.Info("Retrieving recent messages",
		zap.Int("limit", limit),
	)

	// TODO: Implement query using scuttlego ReceiveLog query
	// This requires:
	// 1. Calling ReceiveLog query
	// 2. Iterating over messages
	// 3. Returning message references

	return []string{}, nil
}

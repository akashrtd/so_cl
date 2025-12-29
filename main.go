package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourusername/so_cl/config"
	"github.com/yourusername/so_cl/scuttlego"
	"github.com/yourusername/so_cl/ui"
)

var (
	// version is the current so_cl version
	version = "0.0.1-dev"

	// buildTime is when this binary was built
	buildTime = "unknown"
)

func main() {
	// Parse configuration
	cfg := config.Load()

	// Initialize logger
	logger := initLogger(cfg.Debug)

	// Log startup
	logger.Info("Starting so_cl",
		zap.String("version", version),
		zap.String("build_time", buildTime),
	)

	// Create context with graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	go handleShutdownSignals(ctx, cancel, logger)

	// Initialize scuttlego service
	scuttlegoService, err := scuttlego.NewService(
		scuttlego.DefaultConfig(),
		logger,
	)
	if err != nil {
		logger.Fatal("Failed to initialize scuttlego service",
			zap.Error(err),
		)
	}

	// Start scuttlego service
	go func() {
		if err := scuttlegoService.Run(ctx); err != nil {
			logger.Error("Scuttlego service error",
				zap.Error(err),
			)
		}
	}()

	// Create Bubble Tea model
	model := ui.NewSoClModel()

	// Create Bubble Tea program
	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Run Bubble Tea
	if _, err := program.Run(); err != nil {
		logger.Error("Bubble Tea error",
			zap.Error(err),
		)
		os.Exit(1)
	}

	// Clean shutdown
	logger.Info("Shutting down so_cl")

	// Close scuttlego service
	if err := scuttlegoService.Close(); err != nil {
		logger.Error("Failed to close scuttlego service",
			zap.Error(err),
		)
		os.Exit(1)
	}

	logger.Info("so_cl shut down successfully")
}

// initLogger initializes zap logger based on debug flag.
func initLogger(debug bool) *zap.Logger {
	var zapConfig zap.Config
	if debug {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}

	logger, err := zapConfig.Build()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	return logger
}

// handleShutdownSignals handles SIGINT and SIGTERM for graceful shutdown.
func handleShutdownSignals(ctx context.Context, cancel context.CancelFunc, logger *zap.Logger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	logger.Info("Received shutdown signal",
		zap.String("signal", sig.String()),
	)
	cancel()
}

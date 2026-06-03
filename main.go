package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourusername/so_cl/config"
	"github.com/yourusername/so_cl/scuttlego"
	"github.com/yourusername/so_cl/ui"
)

var (
	// version is the current so_cl version
	version = "0.1.5"

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
	go handleShutdownSignals(cancel, logger)

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
			if err != context.Canceled {
				logger.Error("Scuttlego service error",
					zap.Error(err),
				)
			}
		}
	}()

	// Create Bubble Tea model
	model := ui.NewSoClModel(scuttlegoService)

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
// Logs to a file to avoid interfering with the TUI.
func initLogger(debug bool) *zap.Logger {
	home, _ := os.UserHomeDir()
	logDir := filepath.Join(home, ".so_cl")
	os.MkdirAll(logDir, 0o755)
	logPath := filepath.Join(logDir, "so_cl.log")

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		panic(fmt.Sprintf("Failed to open log file: %v", err))
	}

	var encoder zapcore.Encoder
	var level zapcore.Level
	if debug {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		level = zapcore.DebugLevel
	} else {
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		level = zapcore.InfoLevel
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(logFile), level)
	return zap.New(core)
}

// handleShutdownSignals handles SIGINT and SIGTERM for graceful shutdown.
func handleShutdownSignals(cancel context.CancelFunc, logger *zap.Logger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	logger.Info("Received shutdown signal",
		zap.String("signal", sig.String()),
	)
	cancel()
}

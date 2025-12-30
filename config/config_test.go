package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	t.Run("returns config with default values", func(t *testing.T) {
		cfg := DefaultConfig()

		require.NotEmpty(t, cfg.DataDir)
		require.Contains(t, cfg.DataDir, "/.so_cl")
		require.Equal(t, 8008, cfg.ListenPort)
		require.Empty(t, cfg.NetworkKey)
		require.True(t, cfg.EnableLANDiscovery)
		require.Equal(t, "info", cfg.LogLevel)
		require.False(t, cfg.Debug)
	})
}

func TestLoad(t *testing.T) {
	// Clean up environment variables before each test
	cleanup := func() {
		os.Unsetenv("SO_CL_DATA_DIR")
		os.Unsetenv("SO_CL_PORT")
		os.Unsetenv("SO_CL_NETWORK_KEY")
		os.Unsetenv("SO_CL_ENABLE_LAN_DISCOVERY")
		os.Unsetenv("SO_CL_LOG_LEVEL")
		os.Unsetenv("SO_CL_DEBUG")
	}

	t.Run("loads default config when no env vars set", func(t *testing.T) {
		cleanup()
		defer cleanup()

		cfg := Load()

		require.NotEmpty(t, cfg.DataDir)
		require.Contains(t, cfg.DataDir, "/.so_cl")
		require.Equal(t, 8008, cfg.ListenPort)
		require.Empty(t, cfg.NetworkKey)
		require.True(t, cfg.EnableLANDiscovery)
		require.Equal(t, "info", cfg.LogLevel)
		require.False(t, cfg.Debug)
	})

	t.Run("loads custom data directory", func(t *testing.T) {
		cleanup()
		defer cleanup()

		os.Setenv("SO_CL_DATA_DIR", "/custom/path")

		cfg := Load()

		require.Equal(t, "/custom/path", cfg.DataDir)
	})

	t.Run("loads custom port", func(t *testing.T) {
		cleanup()
		defer cleanup()

		os.Setenv("SO_CL_PORT", "9000")

		cfg := Load()

		require.Equal(t, 9000, cfg.ListenPort)
	})

	t.Run("ignores invalid port", func(t *testing.T) {
		cleanup()
		defer cleanup()

		os.Setenv("SO_CL_PORT", "invalid")

		cfg := Load()

		// Should use default port when invalid
		require.Equal(t, 8008, cfg.ListenPort)
	})

	t.Run("loads custom network key", func(t *testing.T) {
		cleanup()
		defer cleanup()

		os.Setenv("SO_CL_NETWORK_KEY", "custom_network_key")

		cfg := Load()

		require.Equal(t, "custom_network_key", cfg.NetworkKey)
	})

	t.Run("loads LAN discovery enabled with true", func(t *testing.T) {
		cleanup()
		defer cleanup()

		os.Setenv("SO_CL_ENABLE_LAN_DISCOVERY", "true")

		cfg := Load()

		require.True(t, cfg.EnableLANDiscovery)
	})

	t.Run("loads LAN discovery enabled with 1", func(t *testing.T) {
		cleanup()
		defer cleanup()

		os.Setenv("SO_CL_ENABLE_LAN_DISCOVERY", "1")

		cfg := Load()

		require.True(t, cfg.EnableLANDiscovery)
	})

	t.Run("loads LAN discovery disabled with false", func(t *testing.T) {
		cleanup()
		defer cleanup()

		os.Setenv("SO_CL_ENABLE_LAN_DISCOVERY", "false")

		cfg := Load()

		require.False(t, cfg.EnableLANDiscovery)
	})

	t.Run("loads LAN discovery disabled with 0", func(t *testing.T) {
		cleanup()
		defer cleanup()

		os.Setenv("SO_CL_ENABLE_LAN_DISCOVERY", "0")

		cfg := Load()

		require.False(t, cfg.EnableLANDiscovery)
	})

	t.Run("loads custom log level", func(t *testing.T) {
		cleanup()
		defer cleanup()

		os.Setenv("SO_CL_LOG_LEVEL", "debug")

		cfg := Load()

		require.Equal(t, "debug", cfg.LogLevel)
	})

	t.Run("loads debug enabled with true", func(t *testing.T) {
		cleanup()
		defer cleanup()

		os.Setenv("SO_CL_DEBUG", "true")

		cfg := Load()

		require.True(t, cfg.Debug)
	})

	t.Run("loads debug enabled with 1", func(t *testing.T) {
		cleanup()
		defer cleanup()

		os.Setenv("SO_CL_DEBUG", "1")

		cfg := Load()

		require.True(t, cfg.Debug)
	})

	t.Run("loads debug disabled with false", func(t *testing.T) {
		cleanup()
		defer cleanup()

		os.Setenv("SO_CL_DEBUG", "false")

		cfg := Load()

		require.False(t, cfg.Debug)
	})

	t.Run("loads debug disabled with 0", func(t *testing.T) {
		cleanup()
		defer cleanup()

		os.Setenv("SO_CL_DEBUG", "0")

		cfg := Load()

		require.False(t, cfg.Debug)
	})

	t.Run("loads multiple custom values", func(t *testing.T) {
		cleanup()
		defer cleanup()

		os.Setenv("SO_CL_DATA_DIR", "/test/path")
		os.Setenv("SO_CL_PORT", "9001")
		os.Setenv("SO_CL_NETWORK_KEY", "test_key")
		os.Setenv("SO_CL_ENABLE_LAN_DISCOVERY", "false")
		os.Setenv("SO_CL_LOG_LEVEL", "warn")
		os.Setenv("SO_CL_DEBUG", "true")

		cfg := Load()

		require.Equal(t, "/test/path", cfg.DataDir)
		require.Equal(t, 9001, cfg.ListenPort)
		require.Equal(t, "test_key", cfg.NetworkKey)
		require.False(t, cfg.EnableLANDiscovery)
		require.Equal(t, "warn", cfg.LogLevel)
		require.True(t, cfg.Debug)
	})
}

func TestHomeDir(t *testing.T) {
	t.Run("returns HOME env var when set", func(t *testing.T) {
		originalHome := os.Getenv("HOME")
		defer func() {
			if originalHome != "" {
				os.Setenv("HOME", originalHome)
			} else {
				os.Unsetenv("HOME")
			}
		}()

		os.Setenv("HOME", "/test/home")

		home := homeDir()

		require.Equal(t, "/test/home", home)
	})

	t.Run("returns USERPROFILE env var when HOME not set", func(t *testing.T) {
		originalHome := os.Getenv("HOME")
		originalUserProfile := os.Getenv("USERPROFILE")
		defer func() {
			if originalHome != "" {
				os.Setenv("HOME", originalHome)
			} else {
				os.Unsetenv("HOME")
			}
			if originalUserProfile != "" {
				os.Setenv("USERPROFILE", originalUserProfile)
			} else {
				os.Unsetenv("USERPROFILE")
			}
		}()

		os.Unsetenv("HOME")
		os.Setenv("USERPROFILE", "/test/profile")

		home := homeDir()

		require.Equal(t, "/test/profile", home)
	})
}

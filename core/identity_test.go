package core

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewIdentity(t *testing.T) {
	t.Run("creates valid identity", func(t *testing.T) {
		identity, err := NewIdentity()

		require.NoError(t, err, "NewIdentity should not return an error")
		require.NotNil(t, identity, "Identity should not be nil")
		require.NotEmpty(t, identity.FeedRef, "FeedRef should not be empty")
		require.NotEmpty(t, identity.PrivateKey, "PrivateKey should not be empty")
		require.NotEmpty(t, identity.PublicKey, "PublicKey should not be empty")
		require.Contains(t, identity.FeedRef, ".ed25519", "FeedRef should contain .ed25519")
	})

	t.Run("creates different identities", func(t *testing.T) {
		identity1, err1 := NewIdentity()
		identity2, err2 := NewIdentity()

		require.NoError(t, err1)
		require.NoError(t, err2)
		require.NotEqual(t, identity1.FeedRef, identity2.FeedRef, "FeedRefs should be different")
		require.NotEqual(t, identity1.PrivateKey, identity2.PrivateKey, "PrivateKeys should be different")
		require.NotEqual(t, identity1.PublicKey, identity2.PublicKey, "PublicKeys should be different")
	})
}

func TestIdentity_String(t *testing.T) {
	t.Run("returns feed reference", func(t *testing.T) {
		identity, _ := NewIdentity()

		str := identity.String()

		require.Equal(t, identity.FeedRef, str, "String() should return FeedRef")
	})
}

func TestIdentity_ExportSeed(t *testing.T) {
	t.Run("exports seed bytes", func(t *testing.T) {
		identity, _ := NewIdentity()

		seed, err := identity.ExportSeed()

		require.NoError(t, err, "ExportSeed should not return an error")
		require.Len(t, seed, 64, "Seed should be 64 bytes (Ed25519 private key)")
	})

	t.Run("exports valid hex seed", func(t *testing.T) {
		identity, _ := NewIdentity()

		seed, err := identity.ExportSeed()

		require.NoError(t, err)
		// Verify seed is valid hex by decoding it
		decoded, err := hex.DecodeString(identity.PrivateKey)
		require.NoError(t, err, "PrivateKey should be valid hex")
		require.Equal(t, seed, decoded, "Exported seed should match decoded PrivateKey")
	})
}

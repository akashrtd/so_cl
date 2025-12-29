package core

// Identity represents a Secure Scuttlebutt user identity.
// This is a stub implementation that will be updated when
// scuttlego is fully integrated.
//
// TODO: Update to use scuttlego's identity.Private and identity.Public types.
// Currently using placeholder types.

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
)

// Identity represents a Secure Scuttlebutt user identity.
type Identity struct {
	// FeedRef is the SSB feed reference (e.g., "@alice...")
	FeedRef string
	// PrivateKey is the Ed25519 private key (hex encoded)
	PrivateKey string
	// PublicKey is the Ed25519 public key (hex encoded)
	PublicKey string
}

// NewIdentity creates a new SSB identity (Ed25519 keypair).
// This generates a new random keypair and wraps it in an Identity struct.
//
// Returns an error if identity generation fails.
func NewIdentity() (*Identity, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Ed25519 keypair: %w", err)
	}

	return &Identity{
		FeedRef:    "@" + hex.EncodeToString(pub[:])[:8] + ".ed25519",
		PrivateKey: hex.EncodeToString(priv),
		PublicKey:  hex.EncodeToString(pub),
	}, nil
}

// String returns the string representation of the feed reference.
func (i *Identity) String() string {
	return i.FeedRef
}

// ExportSeed exports the private key seed for backup.
// This returns the raw private key bytes that can be encrypted and stored.
//
// Returns an error if export fails.
func (i *Identity) ExportSeed() ([]byte, error) {
	seed, err := hex.DecodeString(i.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}

	// For Ed25519, the private key is 64 bytes
	// This is the "seed" that can be backed up
	return seed[:], nil
}

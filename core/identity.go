package core

// Identity represents a Secure Scuttlebutt user identity.
// This is a placeholder wrapper that will be updated when
// scuttlego is properly integrated.
type Identity struct {
	// FeedRef is the SSB feed reference (e.g., "@alice...")
	FeedRef string
	// PrivateKey is the Ed25519 private key (for signing)
	// This is handled by scuttlego internally
	PrivateKey string
}

// NewIdentity creates a new SSB identity.
// This will be implemented using scuttlego's identity generation.
func NewIdentity() (*Identity, error) {
	// TODO: Implement using scuttlego identity.NewPrivate()
	return &Identity{
		FeedRef:    "",
		PrivateKey: "",
	}, nil
}

// String returns the string representation of the feed reference.
func (i *Identity) String() string {
	return i.FeedRef
}

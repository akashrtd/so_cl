package core

// Common types used across so_cl application.

// SoClPost represents a social media post in the so_cl feed.
type SoClPost struct {
	// Ref is the SSB message reference (e.g., "%abc123.sha256")
	Ref string
	// Author is the SSB feed reference of the author (e.g., "@alice...")
	Author string
	// Text is the post content (1-280 ASCII characters)
	Text string
	// Timestamp is when the post was created (Unix timestamp in seconds)
	Timestamp int64
	// Sequence is the message sequence number in the author's feed
	Sequence int64
	// Root is the message reference this post replies to (if any)
	Root string
	// Branch is the message reference this post replies to (if any)
	Branch string
	// Tags are hashtags extracted from the post
	Tags []string
	// Mentions are @mentions extracted from the post
	Mentions []string
}

// SoClPeer represents a connected SSB peer.
type SoClPeer struct {
	// FeedRef is the SSB feed reference of the peer
	FeedRef string
	// Address is the network address of the peer
	Address string
	// Connected indicates if the peer is currently connected
	Connected bool
	// LastSeen is the timestamp of last activity from this peer
	LastSeen int64
	// Following indicates if the current user follows this peer
	Following bool
	// Follower indicates if this peer follows the current user
	Follower bool
}

// SoClProfile represents a user's profile information.
type SoClProfile struct {
	// FeedRef is the SSB feed reference of the user
	FeedRef string
	// Username is the display name (optional)
	Username string
	// PFP is the ASCII profile picture (6x6 colored ANSI)
	PFP string
	// Bio is the user's bio/description (optional)
	Bio string
	// FollowingCount is the number of peers the user follows
	FollowingCount int
	// FollowersCount is the number of peers following the user
	FollowersCount int
	// PostCount is the total number of posts by the user
	PostCount int
}

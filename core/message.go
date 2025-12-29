package core

// Post represents a Secure Scuttlebutt social media post.
type Post struct {
	// Ref is the SSB message reference (e.g., "%abc123.sha256")
	Ref string
	// Author is the SSB feed reference of the post author (e.g., "@alice...")
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
	// Tags is the list of hashtags extracted from the post
	Tags []string
	// Mentions is the list of @mentions extracted from the post
	Mentions []string
}

// NewPost creates a new Post struct with the given parameters.
func NewPost(ref, author, text string, timestamp, sequence int64) *Post {
	return &Post{
		Ref:       ref,
		Author:    author,
		Text:      text,
		Timestamp: timestamp,
		Sequence:  sequence,
		Tags:      []string{},
		Mentions:  []string{},
	}
}

// Reply creates a new Post that replies to an existing post.
// Root is the original message reference, Branch is the specific message being replied to.
func Reply(post *Post, text string, timestamp int64) *Post {
	return &Post{
		Ref:       "", // Generated when published
		Author:    "", // Current user's identity
		Text:      text,
		Timestamp: timestamp,
		Root:      post.Ref,
		Branch:    post.Ref,
		Tags:      []string{},
		Mentions:  []string{},
	}
}

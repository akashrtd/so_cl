package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPost(t *testing.T) {
	t.Run("creates post with all fields", func(t *testing.T) {
		post := NewPost("@msg.ref", "@author", "test text", 123456, 1)

		require.Equal(t, "@msg.ref", post.Ref)
		require.Equal(t, "@author", post.Author)
		require.Equal(t, "test text", post.Text)
		require.Equal(t, int64(123456), post.Timestamp)
		require.Equal(t, int64(1), post.Sequence)
		require.Empty(t, post.Tags, "Tags should be initialized empty")
		require.Empty(t, post.Mentions, "Mentions should be initialized empty")
	})

	t.Run("allows empty ref", func(t *testing.T) {
		post := NewPost("", "@author", "test", 123456, 1)
		require.Empty(t, post.Ref, "Empty ref should be allowed")
	})

	t.Run("allows empty author", func(t *testing.T) {
		post := NewPost("@msg.ref", "", "test", 123456, 1)
		require.Empty(t, post.Author, "Empty author should be allowed")
	})
}

func TestReply(t *testing.T) {
	t.Run("creates reply with root and branch", func(t *testing.T) {
		original := NewPost("@msg.ref", "@author", "original", 123456, 1)
		reply := Reply(original, "reply text", 123457)

		require.Equal(t, "@msg.ref", reply.Root, "Root should match original post ref")
		require.Equal(t, "@msg.ref", reply.Branch, "Branch should match original post ref")
		require.Equal(t, "reply text", reply.Text)
		require.Equal(t, int64(123457), reply.Timestamp)
		require.Empty(t, reply.Tags, "Tags should be initialized empty")
		require.Empty(t, reply.Mentions, "Mentions should be initialized empty")
	})

	t.Run("creates reply with empty root and branch", func(t *testing.T) {
		post := NewPost("", "", "", 0, 0)
		reply := Reply(post, "text", 123457)

		require.Empty(t, reply.Root, "Root should be empty")
		require.Empty(t, reply.Branch, "Branch should be empty")
	})

	t.Run("author is empty for reply", func(t *testing.T) {
		original := NewPost("@msg.ref", "@author", "original", 123456, 1)
		reply := Reply(original, "reply text", 123457)

		require.Empty(t, reply.Author, "Author should be empty for reply")
	})
}

func TestPost_TagsInitialization(t *testing.T) {
	t.Run("initializes empty tags and mentions", func(t *testing.T) {
		post := NewPost("@msg.ref", "@author", "test", 123456, 1)

		require.NotNil(t, post.Tags, "Tags should not be nil")
		require.Empty(t, post.Tags, "Tags should be empty")
		require.NotNil(t, post.Mentions, "Mentions should not be nil")
		require.Empty(t, post.Mentions, "Mentions should be empty")
	})
}

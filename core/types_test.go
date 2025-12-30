package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSoClPost_Fields(t *testing.T) {
	t.Run("all fields are accessible", func(t *testing.T) {
		post := SoClPost{
			Ref:       "@msg.ref",
			Author:    "@author",
			Text:      "test",
			Timestamp: 123456,
			Sequence:  1,
			Root:      "@root.ref",
			Branch:    "@branch.ref",
			Tags:      []string{"golang"},
			Mentions:  []string{"@user"},
			LikeCount: 5,
		}

		require.Equal(t, "@msg.ref", post.Ref)
		require.Equal(t, "@author", post.Author)
		require.Equal(t, "test", post.Text)
		require.Equal(t, int64(123456), post.Timestamp)
		require.Equal(t, int64(1), post.Sequence)
		require.Equal(t, "@root.ref", post.Root)
		require.Equal(t, "@branch.ref", post.Branch)
		require.Len(t, post.Tags, 1)
		require.Equal(t, "golang", post.Tags[0])
		require.Len(t, post.Mentions, 1)
		require.Equal(t, "@user", post.Mentions[0])
		require.Equal(t, 5, post.LikeCount)
	})
}

func TestSoClPeer_Fields(t *testing.T) {
	t.Run("all fields are accessible", func(t *testing.T) {
		peer := SoClPeer{
			FeedRef:   "@peer.ed25519",
			Address:   "net:127.0.0.1:8008~shs:@peer.ed25519",
			Connected: true,
			LastSeen:  1234567890,
			Following: true,
			Follower:  false,
		}

		require.Equal(t, "@peer.ed25519", peer.FeedRef)
		require.Equal(t, "net:127.0.0.1:8008~shs:@peer.ed25519", peer.Address)
		require.True(t, peer.Connected)
		require.Equal(t, int64(1234567890), peer.LastSeen)
		require.True(t, peer.Following)
		require.False(t, peer.Follower)
	})
}

func TestSoClProfile_Fields(t *testing.T) {
	t.Run("all fields are accessible", func(t *testing.T) {
		profile := SoClProfile{
			FeedRef:        "@user.ed25519",
			Username:       "alice",
			PFP:            "colored ascii art",
			Bio:            "Hello, I'm Alice!",
			FollowingCount: 10,
			FollowersCount: 25,
			PostCount:      100,
		}

		require.Equal(t, "@user.ed25519", profile.FeedRef)
		require.Equal(t, "alice", profile.Username)
		require.Equal(t, "colored ascii art", profile.PFP)
		require.Equal(t, "Hello, I'm Alice!", profile.Bio)
		require.Equal(t, 10, profile.FollowingCount)
		require.Equal(t, 25, profile.FollowersCount)
		require.Equal(t, 100, profile.PostCount)
	})
}

func TestVote_Fields(t *testing.T) {
	t.Run("all fields are accessible", func(t *testing.T) {
		vote := Vote{
			Ref:        "@vote.msg",
			Author:     "@voter.ed25519",
			PostRef:    "@post.msg",
			Expression: "like",
			Timestamp:  1234567890,
		}

		require.Equal(t, "@vote.msg", vote.Ref)
		require.Equal(t, "@voter.ed25519", vote.Author)
		require.Equal(t, "@post.msg", vote.PostRef)
		require.Equal(t, "like", vote.Expression)
		require.Equal(t, int64(1234567890), vote.Timestamp)
	})
}

func TestNotification_Fields(t *testing.T) {
	t.Run("all fields are accessible", func(t *testing.T) {
		notification := Notification{
			Type:      "mention",
			From:      "@user.ed25519",
			PostRef:   "@post.msg",
			Text:      "@user mentioned you",
			Timestamp: 1234567890,
			Read:      false,
		}

		require.Equal(t, "mention", notification.Type)
		require.Equal(t, "@user.ed25519", notification.From)
		require.Equal(t, "@post.msg", notification.PostRef)
		require.Equal(t, "@user mentioned you", notification.Text)
		require.Equal(t, int64(1234567890), notification.Timestamp)
		require.False(t, notification.Read)
	})
}

func TestTrendingHashtag_Fields(t *testing.T) {
	t.Run("all fields are accessible", func(t *testing.T) {
		hashtag := TrendingHashtag{
			Name:  "golang",
			Count: 42,
		}

		require.Equal(t, "golang", hashtag.Name)
		require.Equal(t, 42, hashtag.Count)
	})
}

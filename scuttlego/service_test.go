package scuttlego

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yourusername/so_cl/core"
	"go.uber.org/zap"
)

type mockScuttlegoService struct {
	publishResult string
	publishErr    error
	messages      []Message
	peers         []Peer
	ebtActive     bool
	ebtSent       int
	ebtReceived   int
	trending      []core.TrendingHashtag
	mentions      []string
	identity      string
	following     []string
	followers     []string
}

func (m *mockScuttlegoService) Publish(text string) (string, error) {
	return m.publishResult, m.publishErr
}

func (m *mockScuttlegoService) Reply(text, root, branch string) (string, error) {
	return "@reply.msg", nil
}

func (m *mockScuttlegoService) React(postRef, expression string) (string, error) {
	return "@vote.msg", nil
}

func (m *mockScuttlegoService) GetRecentMessages(limit int) ([]Message, error) {
	if len(m.messages) > limit {
		return m.messages[:limit], nil
	}
	return m.messages, nil
}

func (m *mockScuttlegoService) Follow(feedRef string) error {
	// Add to following list for testing
	m.following = append(m.following, feedRef)
	return nil
}

func (m *mockScuttlegoService) Unfollow(feedRef string) error {
	// Remove from following list for testing
	newFollowing := make([]string, 0, len(m.following))
	for _, f := range m.following {
		if f != feedRef {
			newFollowing = append(newFollowing, f)
		}
	}
	m.following = newFollowing
	return nil
}

func (m *mockScuttlegoService) Connect(address string) error {
	return nil
}

func (m *mockScuttlegoService) RedeemInvite(inviteCode string) error {
	return nil
}

func (m *mockScuttlegoService) GetPeers() ([]Peer, error) {
	return m.peers, nil
}

func (m *mockScuttlegoService) GetEBTStatus() (bool, int, int) {
	return m.ebtActive, m.ebtSent, m.ebtReceived
}

func (m *mockScuttlegoService) GetTopHashtags(n int) ([]core.TrendingHashtag, error) {
	return m.trending, nil
}

func (m *mockScuttlegoService) GetMentions(feedRef string) ([]string, error) {
	return m.mentions, nil
}

func (m *mockScuttlegoService) GetIdentity() string {
	return m.identity
}

func (m *mockScuttlegoService) GetFollowing(feedRef string) ([]string, error) {
	return m.following, nil
}

func (m *mockScuttlegoService) GetFollowers(feedRef string) ([]string, error) {
	return m.followers, nil
}

func (m *mockScuttlegoService) IsFollowing(follower, following string) (bool, error) {
	for _, f := range m.following {
		if f == following {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockScuttlegoService) GetFollowingCount(feedRef string) (int, error) {
	return len(m.following), nil
}

func (m *mockScuttlegoService) GetFollowersCount(feedRef string) (int, error) {
	return len(m.followers), nil
}

func (m *mockScuttlegoService) SearchPosts(query string) ([]string, error) {
	// Simple mock - return empty results
	return []string{}, nil
}

func (m *mockScuttlegoService) FilterByHashtag(hashtag string) ([]string, error) {
	// Simple mock - return empty results
	return []string{}, nil
}

func (m *mockScuttlegoService) FilterByAuthor(author string) ([]string, error) {
	// Simple mock - return empty results
	return []string{}, nil
}

func setupTestService(t *testing.T) *Service {
	t.Helper()

	tmpDir := t.TempDir()
	cfg := Config{
		DataDir:            tmpDir,
		ListenPort:         0,
		EnableLANDiscovery: false,
	}

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	svc, err := NewService(cfg, logger)
	require.NoError(t, err)
	require.NotNil(t, svc)

	t.Cleanup(func() {
		if err := svc.Close(); err != nil {
			t.Logf("Failed to close service: %v", err)
		}
	})

	return svc
}

func startTestService(t *testing.T, svc *Service) context.CancelFunc {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := svc.Run(ctx); err != nil && err != context.Canceled {
			t.Logf("Service error: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	return cancel
}

func TestNewService(t *testing.T) {
	t.Run("creates service with valid config", func(t *testing.T) {
		svc := setupTestService(t)
		require.NotNil(t, svc)
	})

	t.Run("creates data directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		cfg := Config{
			DataDir:    tmpDir,
			ListenPort: 0,
		}

		logger, _ := zap.NewDevelopment()
		svc, err := NewService(cfg, logger)
		require.NoError(t, err)
		require.NotNil(t, svc)
		defer svc.Close()

		_, err = os.Stat(filepath.Join(tmpDir, "badger"))
		require.NoError(t, err)
	})

	t.Run("generates unique identity", func(t *testing.T) {
		svc1 := setupTestService(t)
		svc2 := setupTestService(t)

		cancel1 := startTestService(t, svc1)
		defer cancel1()

		cancel2 := startTestService(t, svc2)
		defer cancel2()

		_, err1 := svc1.Publish("test 1")
		require.NoError(t, err1)

		_, err2 := svc2.Publish("test 2")
		require.NoError(t, err2)

		messages1, _ := svc1.GetRecentMessages(10)
		messages2, _ := svc2.GetRecentMessages(10)
		require.Len(t, messages1, 1)
		require.Len(t, messages2, 1)
		require.Equal(t, "test 1", messages1[0].Text)
		require.Equal(t, "test 2", messages2[0].Text)
	})
}

func TestPublish(t *testing.T) {
	t.Run("publishes valid post", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		ref, err := svc.Publish("hello world")
		require.NoError(t, err)
		require.NotEmpty(t, ref)
	})

	t.Run("rejects empty post", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		_, err := svc.Publish("")
		require.Error(t, err)
		require.Contains(t, err.Error(), "empty post text")
	})

	t.Run("rejects post exceeding 280 characters", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		longText := ""
		for i := 0; i < 281; i++ {
			longText += "a"
		}

		_, err := svc.Publish(longText)
		require.Error(t, err)
		require.Contains(t, err.Error(), "exceeds 280 character limit")
	})

	t.Run("accepts post at exactly 280 characters", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		text := ""
		for i := 0; i < 280; i++ {
			text += "a"
		}

		ref, err := svc.Publish(text)
		require.NoError(t, err)
		require.NotEmpty(t, ref)
	})

	t.Run("persists published post", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		text := "test post"
		_, err := svc.Publish(text)
		require.NoError(t, err)

		messages, err := svc.GetRecentMessages(10)
		require.NoError(t, err)
		require.Len(t, messages, 1)
		require.Equal(t, text, messages[0].Text)
	})

	t.Run("publishes multiple posts", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		for i := 0; i < 5; i++ {
			text := "post"
			_, err := svc.Publish(text)
			require.NoError(t, err)
		}

		messages, err := svc.GetRecentMessages(10)
		require.NoError(t, err)
		require.Len(t, messages, 5)
	})
}

func TestFollow(t *testing.T) {
	t.Run("rejects invalid feed reference", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		err := svc.Follow("invalid")
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid feed reference")
	})

	t.Run("rejects feed reference with invalid key", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		err := svc.Follow("@test.ed25519")
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid feed reference")
	})
}

func TestConnect(t *testing.T) {
	t.Run("validates address format", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		err := svc.Connect("invalid-address")
		require.Error(t, err)
		require.Contains(t, err.Error(), "expected ~shs: separator")
	})

	t.Run("validates address with ~shs: separator", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		addr := "net:127.0.0.1:8008~shs:@test.key.ed25519"

		err := svc.Connect(addr)
		require.Error(t, err)
		require.Contains(t, err.Error(), "could not parse identity")
	})
}

func TestGetRecentMessages(t *testing.T) {
	t.Run("returns empty messages initially", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		messages, err := svc.GetRecentMessages(10)
		require.NoError(t, err)
		require.Empty(t, messages)
	})

	t.Run("retrieves published messages", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		text1 := "first post"
		text2 := "second post"

		_, err1 := svc.Publish(text1)
		require.NoError(t, err1)

		_, err2 := svc.Publish(text2)
		require.NoError(t, err2)

		messages, err := svc.GetRecentMessages(10)
		require.NoError(t, err)
		require.Len(t, messages, 2)
		require.Equal(t, text1, messages[0].Text)
		require.Equal(t, text2, messages[1].Text)
	})

	t.Run("respects limit parameter", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		for i := 0; i < 10; i++ {
			_, err := svc.Publish("post")
			require.NoError(t, err)
		}

		messages, err := svc.GetRecentMessages(5)
		require.NoError(t, err)
		require.Len(t, messages, 5)
	})

	t.Run("handles limit larger than available messages", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		_, err := svc.Publish("post")
		require.NoError(t, err)

		messages, err := svc.GetRecentMessages(100)
		require.NoError(t, err)
		require.Len(t, messages, 1)
	})
}

func TestClose(t *testing.T) {
	t.Run("closes service successfully", func(t *testing.T) {
		svc := setupTestService(t)

		err := svc.Close()
		require.NoError(t, err)
	})

	t.Run("can be called multiple times", func(t *testing.T) {
		svc := setupTestService(t)

		err1 := svc.Close()
		err2 := svc.Close()

		require.NoError(t, err1)
		require.NoError(t, err2)
	})
}

func TestIntegration_PublishAndRetrieve(t *testing.T) {
	t.Run("full publish and retrieve cycle", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		texts := []string{"hello world", "test post", "another message"}

		for _, text := range texts {
			ref, err := svc.Publish(text)
			require.NoError(t, err)
			require.NotEmpty(t, ref)
			time.Sleep(10 * time.Millisecond)
		}

		messages, err := svc.GetRecentMessages(10)
		require.NoError(t, err)
		t.Logf("Retrieved %d messages", len(messages))
		for i, msg := range messages {
			t.Logf("Message %d: Author=%s, Text=%q, Time=%d", i, msg.Author, msg.Text, msg.Time)
		}
		require.Len(t, messages, len(texts))

		for i, expectedText := range texts {
			require.Equal(t, expectedText, messages[i].Text)
			require.NotEmpty(t, messages[i].Author)
			require.GreaterOrEqual(t, messages[i].Time, 0)
		}
	})
}

func TestReply(t *testing.T) {
	t.Run("publishes reply with root and branch", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		root := "%abc123.sha256"
		branch := "%def456.sha256"
		text := "reply text"

		ref, err := svc.Reply(text, root, branch)
		require.NoError(t, err)
		require.NotEmpty(t, ref)
	})

	t.Run("rejects empty reply", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		_, err := svc.Reply("", "%abc123.sha256", "%def456.sha256")
		require.Error(t, err)
		require.Contains(t, err.Error(), "empty reply text")
	})

	t.Run("rejects reply exceeding 280 characters", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		longText := ""
		for i := 0; i < 281; i++ {
			longText += "a"
		}

		_, err := svc.Reply(longText, "%abc123.sha256", "%def456.sha256")
		require.Error(t, err)
		require.Contains(t, err.Error(), "exceeds 280 character limit")
	})
}

func TestReact(t *testing.T) {
	t.Run("publishes reaction", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		postRef := "%abc123.sha256"
		expression := "like"

		ref, err := svc.React(postRef, expression)
		require.NoError(t, err)
		require.NotEmpty(t, ref)
	})

	t.Run("handles empty post ref", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		// React doesn't validate postRef - it just publishes the vote message
		ref, err := svc.React("", "like")
		require.NoError(t, err)
		require.NotEmpty(t, ref)
	})
}

func TestUnfollow(t *testing.T) {
	t.Run("publishes unfollow message", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		feedRef := "@alice.ed25519"

		err := svc.Unfollow(feedRef)
		require.NoError(t, err)
	})
}

func TestGetIdentity(t *testing.T) {
	t.Run("returns feed reference", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		identity := svc.GetIdentity()
		require.NotEmpty(t, identity)
		// GetIdentity returns the public key string (base64-encoded), not with .ed25519 suffix
		require.NotContains(t, identity, ".ed25519")
	})
}

func TestRedeemInvite(t *testing.T) {
	t.Run("handles invalid invite", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		inviteCode := "invalid-invite"

		err := svc.RedeemInvite(inviteCode)
		require.Error(t, err)
	})

	t.Run("handles malformed invite", func(t *testing.T) {
		svc := setupTestService(t)
		cancel := startTestService(t, svc)
		defer cancel()

		// This is not a valid invite code format
		inviteCode := "ssb:feed/invite/code"

		err := svc.RedeemInvite(inviteCode)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to parse invite code")
	})
}

func TestGetTopHashtags_NoIndexer(t *testing.T) {
	t.Run("returns empty list when indexer is nil", func(t *testing.T) {
		tmpDir := t.TempDir()
		cfg := Config{
			DataDir:    tmpDir,
			ListenPort: 0,
		}

		logger, _ := zap.NewDevelopment()
		svc, err := NewService(cfg, logger)
		require.NoError(t, err)
		defer svc.Close()

		hashtags, err := svc.GetTopHashtags(10)
		require.NoError(t, err)
		require.Empty(t, hashtags)
	})
}

func TestGetMentions_NoIndexer(t *testing.T) {
	t.Run("returns empty list when indexer is nil", func(t *testing.T) {
		tmpDir := t.TempDir()
		cfg := Config{
			DataDir:    tmpDir,
			ListenPort: 0,
		}

		logger, _ := zap.NewDevelopment()
		svc, err := NewService(cfg, logger)
		require.NoError(t, err)
		defer svc.Close()

		mentions, err := svc.GetMentions("@alice")
		require.NoError(t, err)
		require.Empty(t, mentions)
	})
}

func TestGetFollowing_NoIndexer(t *testing.T) {
	t.Run("returns empty list when indexer is nil", func(t *testing.T) {
		tmpDir := t.TempDir()
		cfg := Config{
			DataDir:    tmpDir,
			ListenPort: 0,
		}

		logger, _ := zap.NewDevelopment()
		svc, err := NewService(cfg, logger)
		require.NoError(t, err)
		defer svc.Close()

		following, err := svc.GetFollowing("@alice")
		require.NoError(t, err)
		require.Empty(t, following)
	})
}

func TestGetFollowers_NoIndexer(t *testing.T) {
	t.Run("returns empty list when indexer is nil", func(t *testing.T) {
		tmpDir := t.TempDir()
		cfg := Config{
			DataDir:    tmpDir,
			ListenPort: 0,
		}

		logger, _ := zap.NewDevelopment()
		svc, err := NewService(cfg, logger)
		require.NoError(t, err)
		defer svc.Close()

		followers, err := svc.GetFollowers("@alice")
		require.NoError(t, err)
		require.Empty(t, followers)
	})
}

func TestSearchPosts_NoIndexer(t *testing.T) {
	t.Run("returns empty list when indexer is nil", func(t *testing.T) {
		tmpDir := t.TempDir()
		cfg := Config{
			DataDir:    tmpDir,
			ListenPort: 0,
		}

		logger, _ := zap.NewDevelopment()
		svc, err := NewService(cfg, logger)
		require.NoError(t, err)
		defer svc.Close()

		results, err := svc.SearchPosts("test")
		require.NoError(t, err)
		require.Empty(t, results)
	})
}

func TestFilterByHashtag_NoIndexer(t *testing.T) {
	t.Run("returns empty list when indexer is nil", func(t *testing.T) {
		tmpDir := t.TempDir()
		cfg := Config{
			DataDir:    tmpDir,
			ListenPort: 0,
		}

		logger, _ := zap.NewDevelopment()
		svc, err := NewService(cfg, logger)
		require.NoError(t, err)
		defer svc.Close()

		results, err := svc.FilterByHashtag("golang")
		require.NoError(t, err)
		require.Empty(t, results)
	})
}

func TestFilterByAuthor_NoIndexer(t *testing.T) {
	t.Run("returns empty list when indexer is nil", func(t *testing.T) {
		tmpDir := t.TempDir()
		cfg := Config{
			DataDir:    tmpDir,
			ListenPort: 0,
		}

		logger, _ := zap.NewDevelopment()
		svc, err := NewService(cfg, logger)
		require.NoError(t, err)
		defer svc.Close()

		results, err := svc.FilterByAuthor("alice")
		require.NoError(t, err)
		require.Empty(t, results)
	})
}

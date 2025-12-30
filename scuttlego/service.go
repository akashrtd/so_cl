package scuttlego

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/planetary-social/scuttlego/service"
	"github.com/planetary-social/scuttlego/service/app/commands"
	"github.com/planetary-social/scuttlego/service/app/common"
	"github.com/planetary-social/scuttlego/service/app/queries"
	scuttlegodi "github.com/planetary-social/scuttlego/service/di"
	"github.com/planetary-social/scuttlego/service/domain/feeds/message"
	"github.com/planetary-social/scuttlego/service/domain/identity"
	"github.com/planetary-social/scuttlego/service/domain/invites"
	"github.com/planetary-social/scuttlego/service/domain/network"
	"github.com/planetary-social/scuttlego/service/domain/refs"
	"github.com/yourusername/so_cl/core"
	"github.com/yourusername/so_cl/indexes"
	"go.uber.org/zap"
)

type Service struct {
	config  Config
	logger  *zap.Logger
	indexer *indexes.Indexer

	svc     *service.Service
	cleanup func()

	// EBT replication status
	ebtReplicating bool
}

type Config struct {
	DataDir            string
	ListenPort         int
	EnableLANDiscovery bool
}

func DefaultConfig() Config {
	home, _ := os.UserHomeDir()
	return Config{
		DataDir:            filepath.Join(home, ".so_cl", "data"),
		ListenPort:         8008,
		EnableLANDiscovery: true,
	}
}

func NewService(cfg Config, logger *zap.Logger) (*Service, error) {
	logger.Info("Initializing scuttlego service",
		zap.String("data_dir", cfg.DataDir),
		zap.Int("port", cfg.ListenPort),
	)

	listenAddr := fmt.Sprintf(":%d", cfg.ListenPort)

	scConfig := service.Config{
		DataDirectory: cfg.DataDir,
		ListenAddress: listenAddr,
	}
	scConfig.SetDefaults()

	privateIdentity, err := identity.NewPrivate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate identity: %w", err)
	}

	svc, cleanup, err := scuttlegodi.BuildService(privateIdentity, scConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build scuttlego service: %w", err)
	}

	// Open BadgerDB for indexes
	indexesDir := filepath.Join(cfg.DataDir, "indexes")
	indexDB, err := badger.Open(badger.DefaultOptions(indexesDir))
	if err != nil {
		logger.Warn("Failed to open indexes database, indexing disabled",
			zap.Error(err),
		)
		return &Service{
			config:  cfg,
			logger:  logger,
			svc:     &svc,
			cleanup: cleanup,
		}, nil
	}

	return &Service{
		config:  cfg,
		logger:  logger,
		indexer: indexes.NewIndexer(indexDB),
		svc:     &svc,
		cleanup: cleanup,
	}, nil
}

func (s *Service) Run(ctx context.Context) error {
	s.logger.Info("Starting scuttlego service")

	errCh := make(chan error, 1)

	go func() {
		errCh <- s.svc.Run(ctx)
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("Context cancelled, stopping service")
		return ctx.Err()
	case err := <-errCh:
		return fmt.Errorf("scuttlego service error: %w", err)
	}
}

func (s *Service) Close() error {
	s.logger.Info("Closing scuttlego service")

	if s.cleanup != nil {
		s.cleanup()
	}

	if s.indexer != nil {
		if err := s.indexer.Close(); err != nil {
			s.logger.Error("Failed to close indexer",
				zap.Error(err),
			)
		}
	}

	return nil
}

func (s *Service) Publish(text string) (string, error) {
	if len(text) == 0 {
		return "", fmt.Errorf("empty post text")
	}
	if len(text) > 280 {
		return "", fmt.Errorf("post text exceeds 280 character limit")
	}

	s.logger.Info("Publishing post",
		zap.String("text", text),
		zap.Int("length", len(text)),
	)

	content := map[string]interface{}{
		"type": "post",
		"text": text,
	}

	contentJSON, err := json.Marshal(content)
	if err != nil {
		return "", fmt.Errorf("failed to marshal content: %w", err)
	}

	rawContent, err := message.NewRawContent(contentJSON)
	if err != nil {
		return "", fmt.Errorf("failed to create raw content: %w", err)
	}

	cmd, err := commands.NewPublishRaw(rawContent.Bytes())
	if err != nil {
		return "", fmt.Errorf("failed to create publish command: %w", err)
	}

	msgRef, err := s.svc.App.Commands.PublishRaw.Handle(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to publish: %w", err)
	}

	// Index the post for hashtags and mentions
	if s.indexer != nil {
		ref := msgRef.String()
		if err := s.indexer.IndexPost(ref, text); err != nil {
			s.logger.Warn("Failed to index post",
				zap.String("ref", ref),
				zap.Error(err),
			)
		}
	}

	return msgRef.String(), nil
}

// Reply publishes a reply to an existing post.
// The root parameter is the original message reference, branch is the specific message being replied to.
func (s *Service) Reply(text, root, branch string) (string, error) {
	if len(text) == 0 {
		return "", fmt.Errorf("empty reply text")
	}
	if len(text) > 280 {
		return "", fmt.Errorf("reply text exceeds 280 character limit")
	}

	s.logger.Info("Publishing reply",
		zap.String("text", text),
		zap.String("root", root),
		zap.String("branch", branch),
	)

	content := map[string]interface{}{
		"type":   "post",
		"text":   text,
		"root":   root,
		"branch": branch,
	}

	contentJSON, err := json.Marshal(content)
	if err != nil {
		return "", fmt.Errorf("failed to marshal content: %w", err)
	}

	rawContent, err := message.NewRawContent(contentJSON)
	if err != nil {
		return "", fmt.Errorf("failed to create raw content: %w", err)
	}

	cmd, err := commands.NewPublishRaw(rawContent.Bytes())
	if err != nil {
		return "", fmt.Errorf("failed to create publish command: %w", err)
	}

	msgRef, err := s.svc.App.Commands.PublishRaw.Handle(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to publish reply: %w", err)
	}

	// Index the reply for hashtags and mentions
	if s.indexer != nil {
		ref := msgRef.String()
		if err := s.indexer.IndexPost(ref, text); err != nil {
			s.logger.Warn("Failed to index reply",
				zap.String("ref", ref),
				zap.Error(err),
			)
		}
	}

	return msgRef.String(), nil
}

// React publishes a reaction (vote) to a post.
// The expression parameter is the vote expression (e.g., "like", "❤️").
func (s *Service) React(postRef, expression string) (string, error) {
	s.logger.Info("Publishing reaction",
		zap.String("post_ref", postRef),
		zap.String("expression", expression),
	)

	content := map[string]interface{}{
		"type": "vote",
		"vote": map[string]interface{}{
			"link":       postRef,
			"value":      1,
			"expression": expression,
		},
	}

	contentJSON, err := json.Marshal(content)
	if err != nil {
		return "", fmt.Errorf("failed to marshal content: %w", err)
	}

	rawContent, err := message.NewRawContent(contentJSON)
	if err != nil {
		return "", fmt.Errorf("failed to create raw content: %w", err)
	}

	cmd, err := commands.NewPublishRaw(rawContent.Bytes())
	if err != nil {
		return "", fmt.Errorf("failed to create publish command: %w", err)
	}

	msgRef, err := s.svc.App.Commands.PublishRaw.Handle(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to publish reaction: %w", err)
	}

	return msgRef.String(), nil
}

func (s *Service) Follow(feedRef string) error {
	s.logger.Info("Following peer",
		zap.String("feed_ref", feedRef),
	)

	peerIdentity, err := refs.NewIdentity(feedRef)
	if err != nil {
		return fmt.Errorf("invalid feed reference: %w", err)
	}

	cmd := commands.Follow{Target: peerIdentity}
	err = s.svc.App.Commands.Follow.Handle(cmd)
	if err != nil {
		return fmt.Errorf("failed to follow: %w", err)
	}

	return nil
}

func (s *Service) Connect(address string) error {
	s.logger.Info("Connecting to peer",
		zap.String("address", address),
	)

	sep := "~shs:"
	idx := strings.Index(address, sep)
	if idx < 0 {
		return fmt.Errorf("invalid address format, expected ~shs: separator")
	}

	identityString := address[idx+len(sep):]
	peerIdentity, err := refs.NewIdentity(identityString)
	if err != nil {
		return fmt.Errorf("could not parse identity from address: %w", err)
	}

	addr := network.NewAddress(address)

	cmd := commands.Connect{
		Remote:  refs.MustNewIdentityFromPublic(peerIdentity.Identity()).Identity(),
		Address: addr,
	}

	ctx := context.Background()
	connectErr := s.svc.App.Commands.Connect.Handle(ctx, cmd)
	if connectErr != nil {
		return fmt.Errorf("failed to connect: %w", connectErr)
	}

	return nil
}

func (s *Service) GetRecentMessages(limit int) ([]Message, error) {
	s.logger.Info("Retrieving recent messages",
		zap.Int("limit", limit),
	)

	startSeq, err := common.NewReceiveLogSequence(0)
	if err != nil {
		return nil, err
	}

	query, err := queries.NewReceiveLog(startSeq, limit)
	if err != nil {
		return nil, err
	}

	messages, err := s.svc.App.Queries.ReceiveLog.Handle(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	result := make([]Message, 0, len(messages))
	for _, logMsg := range messages {
		author := logMsg.Message.Author().String()

		if logMsg.Message.Content().IsZero() {
			result = append(result, Message{
				Ref:       "",
				Author:    author,
				Text:      "(raw content)",
				Time:      logMsg.Sequence.Int(),
				Root:      "",
				Branch:    "",
				LikeCount: 0,
			})
			continue
		}

		postContent := make(map[string]interface{})
		if err := json.Unmarshal(logMsg.Message.Content().Raw().Bytes(), &postContent); err == nil {
			// Parse post content
			text := ""
			root := ""
			branch := ""

			if t, ok := postContent["text"].(string); ok {
				text = t
			}
			if r, ok := postContent["root"].(string); ok {
				root = r
			}
			if b, ok := postContent["branch"].(string); ok {
				branch = b
			}

			// Check if this is a post with text
			if text != "" {
				result = append(result, Message{
					Ref:       "",
					Author:    author,
					Text:      text,
					Time:      logMsg.Sequence.Int(),
					Root:      root,
					Branch:    branch,
					LikeCount: 0,
				})
				continue
			}

			// Check content type for other message types
			contentType := ""
			if ct, ok := postContent["type"].(string); ok {
				contentType = ct
			}
			result = append(result, Message{
				Ref:       "",
				Author:    author,
				Text:      fmt.Sprintf("(%s)", contentType),
				Time:      logMsg.Sequence.Int(),
				Root:      "",
				Branch:    "",
				LikeCount: 0,
			})
			continue
		}

		known, hasKnown := logMsg.Message.Content().KnownContent()
		if hasKnown {
			result = append(result, Message{
				Ref:       "",
				Author:    author,
				Text:      fmt.Sprintf("%s", known.Type()),
				Time:      logMsg.Sequence.Int(),
				Root:      "",
				Branch:    "",
				LikeCount: 0,
			})
			continue
		}

		result = append(result, Message{
			Ref:       "",
			Author:    author,
			Text:      "(unknown content)",
			Time:      logMsg.Sequence.Int(),
			Root:      "",
			Branch:    "",
			LikeCount: 0,
		})
	}

	return result, nil
}

type Message struct {
	Ref       string
	Author    string
	Text      string
	Time      int
	Root      string
	Branch    string
	LikeCount int
}

// RedeemInvite redeems an invite code to follow a peer.
// Invite code format: ssb:feed/invite/...
func (s *Service) RedeemInvite(inviteCode string) error {
	s.logger.Info("Redeeming invite code",
		zap.String("invite", inviteCode),
	)

	// Parse invite code using invites package
	invite, err := invites.NewInviteFromString(inviteCode)
	if err != nil {
		return fmt.Errorf("failed to parse invite code: %w", err)
	}

	// Redeem the invite (requires context)
	cmd := commands.RedeemInvite{Invite: invite}
	ctx := context.Background()
	err = s.svc.App.Commands.RedeemInvite.Handle(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to redeem invite: %w", err)
	}

	return nil
}

// GetPeers returns the list of connected peers.
func (s *Service) GetPeers() ([]Peer, error) {
	s.logger.Info("Getting connected peers")

	// Get network status (Status query takes no arguments)
	status, err := s.svc.App.Queries.Status.Handle()
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	peers := make([]Peer, 0, len(status.Peers))
	for _, p := range status.Peers {
		peers = append(peers, Peer{
			Address: p.Identity.String(),
			State:   "connected",
		})
	}

	return peers, nil
}

// Peer represents a connected peer.
type Peer struct {
	Address string
	State   string
}

// GetEBTStatus returns EBT replication status.
// EBT is automatic in scuttlego, this provides monitoring.
func (s *Service) GetEBTStatus() (bool, int, int) {
	// EBT is always running in scuttlego
	// Return status for UI display
	return true, 0, 0
}

// GetTopHashtags retrieves the top N hashtags by count.
func (s *Service) GetTopHashtags(n int) ([]core.TrendingHashtag, error) {
	if s.indexer == nil {
		return []core.TrendingHashtag{}, nil
	}

	hashtags, err := s.indexer.GetTopHashtags(n)
	if err != nil {
		return nil, fmt.Errorf("failed to get top hashtags: %w", err)
	}

	// Convert indexes.HashtagCount to core.TrendingHashtag
	result := make([]core.TrendingHashtag, len(hashtags))
	for i, h := range hashtags {
		result[i] = core.TrendingHashtag{
			Name:  h.Name,
			Count: h.Count,
		}
	}

	return result, nil
}

// GetMentions retrieves all post references mentioning a given user.
func (s *Service) GetMentions(feedRef string) ([]string, error) {
	if s.indexer == nil {
		return []string{}, nil
	}

	mentions, err := s.indexer.GetMentions(feedRef)
	if err != nil {
		return nil, fmt.Errorf("failed to get mentions: %w", err)
	}

	return mentions, nil
}

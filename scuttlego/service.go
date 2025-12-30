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
	"github.com/planetary-social/scuttlego/service/domain/network"
	"github.com/planetary-social/scuttlego/service/domain/refs"
	"github.com/yourusername/so_cl/indexes"
	"go.uber.org/zap"
)

type Service struct {
	config  Config
	logger  *zap.Logger
	indexer *indexes.Indexer

	svc     *service.Service
	cleanup func()
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
				Author: author,
				Text:   "(raw content)",
				Time:   logMsg.Sequence.Int(),
			})
			continue
		}

		postContent := make(map[string]interface{})
		if err := json.Unmarshal(logMsg.Message.Content().Raw().Bytes(), &postContent); err == nil {
			if text, ok := postContent["text"].(string); ok {
				result = append(result, Message{
					Author: author,
					Text:   text,
					Time:   logMsg.Sequence.Int(),
				})
				continue
			}

			contentType := ""
			if ct, ok := postContent["type"].(string); ok {
				contentType = ct
			}
			result = append(result, Message{
				Author: author,
				Text:   fmt.Sprintf("(%s)", contentType),
				Time:   logMsg.Sequence.Int(),
			})
			continue
		}

		known, hasKnown := logMsg.Message.Content().KnownContent()
		if hasKnown {
			result = append(result, Message{
				Author: author,
				Text:   fmt.Sprintf("%s", known.Type()),
				Time:   logMsg.Sequence.Int(),
			})
			continue
		}

		result = append(result, Message{
			Author: author,
			Text:   "(unknown content)",
			Time:   logMsg.Sequence.Int(),
		})
	}

	return result, nil
}

type Message struct {
	Author string
	Text   string
	Time   int
}

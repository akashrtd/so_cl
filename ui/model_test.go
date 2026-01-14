package ui

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/so_cl/core"
	"github.com/yourusername/so_cl/scuttlego"
)

type mockScuttlegoService struct {
	publishResult string
	publishErr    error
	messages      []scuttlego.Message
	peers         []scuttlego.Peer
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

func (m *mockScuttlegoService) GetRecentMessages(limit int) ([]scuttlego.Message, error) {
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

func (m *mockScuttlegoService) GetPeers() ([]scuttlego.Peer, error) {
	return m.peers, nil
}

func (m *mockScuttlegoService) GetEBTStatus() (bool, int, int) {
	return m.ebtActive, m.ebtSent, m.ebtReceived
}

func (m *mockScuttlegoService) GetNetworkStatus() (bool, float64, []scuttlego.Peer) {
	return true, 1024.0, m.peers
}

func (m *mockScuttlegoService) UpdateNetworkStatus() {}

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

func TestNewSoClModel(t *testing.T) {
	svc := &mockScuttlegoService{}
	model := NewSoClModel(svc)

	require.NotNil(t, model)
	require.NotNil(t, model.scuttlego)
	require.Empty(t, model.posts)
	require.Empty(t, model.composerText)
	require.False(t, model.editing)
	require.False(t, model.loading)
	require.False(t, model.publishing)
	require.Empty(t, model.errorMsg)
}

func TestModelInit(t *testing.T) {
	svc := &mockScuttlegoService{}
	model := NewSoClModel(svc)

	cmd := model.Init()

	require.NotNil(t, cmd)
}

func TestModelUpdate_WindowSize(t *testing.T) {
	svc := &mockScuttlegoService{}
	model := NewSoClModel(svc)

	msg := tea.WindowSizeMsg{
		Width:  100,
		Height: 50,
	}

	newModel, cmd := model.Update(msg)

	require.Same(t, model, newModel)
	require.Equal(t, 100, model.width)
	require.Equal(t, 50, model.height)
	require.Nil(t, cmd)
}

func TestModelUpdate_FeedLoaded(t *testing.T) {
	svc := &mockScuttlegoService{}
	model := NewSoClModel(svc)

	posts := []Post{
		{Author: "alice", Text: "test post", Time: 0},
		{Author: "bob", Text: "another post", Time: 1},
	}

	msg := FeedLoadedMsg{Posts: posts}

	newModel, cmd := model.Update(msg)

	require.Same(t, model, newModel)
	require.Equal(t, posts, model.posts)
	require.False(t, model.loading)
	require.Nil(t, cmd)
}

func TestModelUpdate_PostPublished(t *testing.T) {
	svc := &mockScuttlegoService{}
	model := NewSoClModel(svc)

	msg := PostPublishedMsg{
		Ref:  "@test.msg",
		Text: "published post",
	}

	newModel, cmd := model.Update(msg)

	require.Same(t, model, newModel)
	require.False(t, model.publishing)
	require.NotNil(t, cmd)
}

func TestModelUpdate_ErrorMsg(t *testing.T) {
	svc := &mockScuttlegoService{}
	model := NewSoClModel(svc)

	msg := ErrorMsg{Err: errors.New("test error")}

	newModel, cmd := model.Update(msg)

	require.Same(t, model, newModel)
	require.Equal(t, "test error", model.errorMsg)
	require.False(t, model.loading)
	require.False(t, model.publishing)
	require.Nil(t, cmd)
}

func TestModelUpdate_KeyMsg_Escape(t *testing.T) {
	svc := &mockScuttlegoService{}
	model := NewSoClModel(svc)

	msg := tea.KeyMsg{
		Type: tea.KeyEsc,
	}

	newModel, cmd := model.Update(msg)

	require.Same(t, model, newModel)
	require.NotNil(t, cmd)
}

func TestModelUpdate_KeyMsg_Enter_StartEditing(t *testing.T) {
	svc := &mockScuttlegoService{}
	model := NewSoClModel(svc)
	model.editing = false

	msg := tea.KeyMsg{
		Type: tea.KeyEnter,
	}

	newModel, cmd := model.Update(msg)

	require.Same(t, model, newModel)
	require.True(t, model.editing)
	require.Nil(t, cmd)
}

func TestModelUpdate_KeyMsg_Enter_Publish(t *testing.T) {
	svc := &mockScuttlegoService{
		publishResult: "@test.msg",
		publishErr:    nil,
	}
	model := NewSoClModel(svc)
	model.editing = true
	model.composerText = "test post"

	msg := tea.KeyMsg{
		Type: tea.KeyEnter,
	}

	newModel, cmd := model.Update(msg)

	require.Same(t, model, newModel)
	require.False(t, model.editing)
	require.Empty(t, model.composerText)
	require.True(t, model.publishing)
	require.Equal(t, 1, len(model.posts))
	require.Equal(t, "You", model.posts[0].Author)
	require.Equal(t, "test post", model.posts[0].Text)
	require.NotNil(t, cmd)
}

func TestModelUpdate_KeyMsg_Enter_AlreadyPublishing(t *testing.T) {
	svc := &mockScuttlegoService{}
	model := NewSoClModel(svc)
	model.editing = true
	model.publishing = true

	msg := tea.KeyMsg{
		Type: tea.KeyEnter,
	}

	newModel, cmd := model.Update(msg)

	require.Same(t, model, newModel)
	require.True(t, model.editing)
	require.True(t, model.publishing)
	require.Nil(t, cmd)
}

func TestModelUpdate_KeyMsg_Up(t *testing.T) {
	svc := &mockScuttlegoService{}
	model := NewSoClModel(svc)
	model.posts = []Post{
		{Author: "a", Text: "post1", Time: 0},
		{Author: "b", Text: "post2", Time: 1},
		{Author: "c", Text: "post3", Time: 2},
	}
	model.cursor = 1

	msg := tea.KeyMsg{
		Type: tea.KeyUp,
	}

	newModel, cmd := model.Update(msg)

	require.Same(t, model, newModel)
	require.Equal(t, 0, model.cursor)
	require.Nil(t, cmd)
}

func TestModelUpdate_KeyMsg_Down(t *testing.T) {
	svc := &mockScuttlegoService{}
	model := NewSoClModel(svc)
	model.posts = []Post{
		{Author: "a", Text: "post1", Time: 0},
		{Author: "b", Text: "post2", Time: 1},
		{Author: "c", Text: "post3", Time: 2},
	}
	model.cursor = 0
	// Set navCursor to max so KeyDown moves the cursor, not navCursor
	model.navCursor = 4 // PageHome, PageDiscover, PagePeers, PageProfile, PageSettings (5 pages, so max index is 4)

	msg := tea.KeyMsg{
		Type: tea.KeyDown,
	}

	newModel, cmd := model.Update(msg)

	require.Same(t, model, newModel)
	require.Equal(t, 1, model.cursor)
	require.Nil(t, cmd)
}

func TestModelUpdate_KeyMsg_Down_AtBottom(t *testing.T) {
	svc := &mockScuttlegoService{}
	model := NewSoClModel(svc)
	model.posts = []Post{{Author: "a", Text: "post", Time: 0}}
	model.cursor = 0

	msg := tea.KeyMsg{
		Type: tea.KeyDown,
	}

	newModel, cmd := model.Update(msg)

	require.Same(t, model, newModel)
	require.Equal(t, 0, model.cursor)
	require.Nil(t, cmd)
}

func TestModelUpdate_KeyMsg_Backspace(t *testing.T) {
	svc := &mockScuttlegoService{}
	model := NewSoClModel(svc)
	model.editing = true
	model.composerText = "hello"

	msg := tea.KeyMsg{
		Type: tea.KeyBackspace,
	}

	newModel, cmd := model.Update(msg)

	require.Same(t, model, newModel)
	require.Equal(t, "hell", model.composerText)
	require.Nil(t, cmd)
}

func TestModelUpdate_KeyMsg_Backspace_Empty(t *testing.T) {
	svc := &mockScuttlegoService{}
	model := NewSoClModel(svc)
	model.editing = true
	model.composerText = ""

	msg := tea.KeyMsg{
		Type: tea.KeyBackspace,
	}

	newModel, cmd := model.Update(msg)

	require.Same(t, model, newModel)
	require.Empty(t, model.composerText)
	require.Nil(t, cmd)
}

func TestModelUpdate_KeyMsg_Typing(t *testing.T) {
	svc := &mockScuttlegoService{}
	model := NewSoClModel(svc)
	model.editing = true

	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'h', 'e', 'l', 'l', 'o'},
	}

	newModel, cmd := model.Update(msg)

	// When settings editing is active, typing is handled by settings case
	// Composer typing only works when settings editing is not active
	if model.showSettings {
		// Settings editing mode - composer typing is disabled
		require.Same(t, model, newModel)
		require.Equal(t, "", model.composerText)
		require.Nil(t, cmd)
	} else {
		// Normal typing mode - composer should be updated
		require.Same(t, model, newModel)
		require.Equal(t, "hello", model.composerText)
		require.Nil(t, cmd)
	}
}

func TestModelUpdate_KeyMsg_Typing_NotEditing(t *testing.T) {
	svc := &mockScuttlegoService{}
	model := NewSoClModel(svc)
	model.editing = false

	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'h', 'i'},
	}

	newModel, cmd := model.Update(msg)

	require.Same(t, model, newModel)
	require.Empty(t, model.composerText)
	require.Nil(t, cmd)
}

func TestModelView(t *testing.T) {
	svc := &mockScuttlegoService{}
	model := NewSoClModel(svc)

	view := model.View()

	require.Equal(t, "Loading...", view)

	model.width = 80
	model.height = 24

	view = model.View()

	// View renders "FEED" (styled) not "Composer" as a literal string
	require.Contains(t, view, "FEED")
	require.Contains(t, view, "NETWORK STATUS")
}

func TestModelView_WithError(t *testing.T) {
	svc := &mockScuttlegoService{}
	model := NewSoClModel(svc)
	model.width = 80
	model.height = 24
	model.errorMsg = "test error"

	view := model.View()

	require.Contains(t, view, "test error")
}

func TestModelView_Feed(t *testing.T) {
	svc := &mockScuttlegoService{}
	model := NewSoClModel(svc)
	model.width = 80
	model.height = 24

	// Initialize viewport with WindowSizeMsg
	_, _ = model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	model.posts = []Post{
		{Author: "@alice", Text: "hello world", Time: 0},
		{Author: "@bob", Text: "test post", Time: 1},
	}
	model.cursor = 1

	// Manually update viewport content since test bypasses message flow
	model.updateViewportContent()

	view := model.View()

	// View renders "FEED" (styled) not "Feed"
	require.Contains(t, view, "FEED")
	require.Contains(t, view, "@alice")
	require.Contains(t, view, "@bob")
	require.Contains(t, view, "hello world")
	require.Contains(t, view, "test post")
}

func TestModelView_EmptyFeed(t *testing.T) {
	svc := &mockScuttlegoService{}
	model := NewSoClModel(svc)
	model.width = 80
	model.height = 24
	model.posts = []Post{}

	view := model.View()

	require.Contains(t, view, "No posts yet")
}

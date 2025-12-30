package ui

// Package ui provides ├── Bubble Tea UI for so_cl.
// It handles:
// - Model (state machine)
// - View (rendering)
// - Update (event handling)

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourusername/so_cl/core"
	"github.com/yourusername/so_cl/scuttlego"
)

const maxPostsInMemory = 100

type ScuttlegoService interface {
	Publish(text string) (string, error)
	GetRecentMessages(limit int) ([]scuttlego.Message, error)
	Follow(feedRef string) error
	Connect(address string) error
	RedeemInvite(inviteCode string) error
	GetPeers() ([]scuttlego.Peer, error)
	GetEBTStatus() (bool, int, int)
}

type Post struct {
	Author string
	Text   string
	Time   int
	pfp    string
}

type FeedLoadedMsg struct {
	Posts []Post
}

type PostPublishedMsg struct {
	Ref  string
	Text string
}

type NewMessageMsg struct {
	Post scuttlego.Message
}

type ErrorMsg struct {
	Err error
}

type PeersLoadedMsg struct {
	Peers []scuttlego.Peer
}

type InviteRedeemedMsg struct {
	InviteCode string
}

type FollowedMsg struct {
	FeedRef string
}

// SoClModel is Bubble Tea model for so_cl.
// It holds all application state:
// - Posts (feed)
// - Composer (post input)
// - Peers (sidebar)
// - Configuration
type SoClModel struct {
	scuttlego ScuttlegoService
	// posts is the feed of social media posts
	posts []Post
	// composerText holds the text input for creating new posts
	composerText string
	// peers is the list of connected SSB peers
	peers []string
	// width is the terminal width
	width int
	// height is the terminal height
	height int
	// cursor shows which post is currently selected
	cursor int
	// editing shows if user is currently typing a post
	editing bool
	// loading shows if feed is currently loading
	loading bool
	// publishing shows if a post is being published
	publishing bool
	// error shows any error message
	errorMsg string
	// peers is the list of connected peers
	connectedPeers []scuttlego.Peer
	// showPeers shows if peer sidebar is visible
	showPeers bool
	// inviteInput shows if user is entering an invite code
	inviteInput string
	// showInviteInput shows if invite input is visible
	showInviteInput bool
	// followInput shows if user is entering a feed ref to follow
	followInput string
	// showFollowInput shows if follow input is visible
	showFollowInput bool
}

// NewSoClModel creates a new SoClModel with default state.
func NewSoClModel(svc ScuttlegoService) *SoClModel {
	return &SoClModel{
		scuttlego:       svc,
		posts:           []Post{},
		composerText:    "",
		peers:           []string{},
		width:           0,
		height:          0,
		cursor:          0,
		editing:         false,
		loading:         false,
		publishing:      false,
		errorMsg:        "",
		connectedPeers:  []scuttlego.Peer{},
		showPeers:       false,
		inviteInput:     "",
		showInviteInput: false,
		followInput:     "",
		showFollowInput: false,
	}
}

// Init is called at the start of the Bubble Tea program.
// It performs initial setup like initial size calculation.
func (m *SoClModel) Init() tea.Cmd {
	return tea.Batch(
		m.loadFeed(),
		m.loadPeers(),
	)
}

// Update handles incoming events from Bubble Tea.
// It processes:
// - Key messages (keyboard input)
// - Window size changes
// - Custom messages (from async operations)
func (m *SoClModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, m.handleKeyMsg(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case FeedLoadedMsg:
		m.posts = msg.Posts
		m.loading = false
		m.publishing = false
		return m, nil
	case PostPublishedMsg:
		m.publishing = false
		return m, m.loadFeed()
	case PeersLoadedMsg:
		m.connectedPeers = msg.Peers
		return m, nil
	case InviteRedeemedMsg:
		m.showInviteInput = false
		m.inviteInput = ""
		m.errorMsg = "Invite redeemed successfully"
		return m, nil
	case FollowedMsg:
		m.showFollowInput = false
		m.followInput = ""
		m.errorMsg = "Followed successfully"
		return m, nil
	case NewMessageMsg:
		// New message received from EBT replication
		// Add to feed with optimistic UI update
		newPost := Post{
			Author: msg.Post.Author,
			Text:   msg.Post.Text,
			Time:   len(m.posts),
			pfp:    core.GeneratePFP(msg.Post.Author),
		}
		m.posts = append([]Post{newPost}, m.posts...)

		// Trim to max posts in memory
		if len(m.posts) > maxPostsInMemory {
			m.posts = m.posts[:maxPostsInMemory]
		}

		// Show notification
		m.errorMsg = fmt.Sprintf("New message from %s", msg.Post.Author)
		return m, nil
	case ErrorMsg:
		m.errorMsg = msg.Err.Error()
		m.loading = false
		m.publishing = false
		return m, nil
	default:
		return m, nil
	}
}

// handleKeyMsg handles keyboard input.
func (m *SoClModel) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return tea.Quit

	case tea.KeyF1:
		// Toggle peer sidebar
		m.showPeers = !m.showPeers
		return nil

	case tea.KeyF2:
		// Toggle invite input
		m.showInviteInput = !m.showInviteInput
		m.inviteInput = ""
		return nil

	case tea.KeyF3:
		// Toggle follow input
		m.showFollowInput = !m.showFollowInput
		m.followInput = ""
		return nil

	case tea.KeyEnter:
		if m.showInviteInput {
			// Redeem invite
			inviteCode := m.inviteInput
			m.inviteInput = ""
			m.showInviteInput = false
			return m.redeemInvite(inviteCode)
		} else if m.showFollowInput {
			// Follow peer
			feedRef := m.followInput
			m.followInput = ""
			m.showFollowInput = false
			return m.followPeer(feedRef)
		} else if m.editing && !m.publishing {
			// Publish post
			text := m.composerText
			m.composerText = ""
			m.editing = false
			m.publishing = true

			// Optimistic UI update: add post immediately
			newPost := Post{
				Author: "You",
				Text:   text,
				Time:   len(m.posts),
				pfp:    core.GeneratePFP("optimistic"),
			}
			m.posts = append([]Post{newPost}, m.posts...)

			// Trim to max posts in memory
			if len(m.posts) > maxPostsInMemory {
				m.posts = m.posts[:maxPostsInMemory]
			}

			m.cursor = 0

			return m.publishPost(text)
		} else if m.editing {
			// Already publishing, ignore
			return nil
		} else {
			// Start editing
			m.editing = true
			return nil
		}

	case tea.KeyUp:
		// Navigate up in feed
		if m.cursor > 0 {
			m.cursor--
		}
		return nil

	case tea.KeyDown:
		// Navigate down in feed
		if m.cursor < len(m.posts)-1 {
			m.cursor++
		}
		return nil

	case tea.KeyBackspace:
		// Delete last character
		if m.editing && len(m.composerText) > 0 {
			m.composerText = m.composerText[:len(m.composerText)-1]
		}
		if m.showInviteInput && len(m.inviteInput) > 0 {
			m.inviteInput = m.inviteInput[:len(m.inviteInput)-1]
		}
		if m.showFollowInput && len(m.followInput) > 0 {
			m.followInput = m.followInput[:len(m.followInput)-1]
		}
		return nil

	default:
		// Typing
		if m.editing && len(msg.Runes) > 0 {
			m.composerText += string(msg.Runes)
		}
		if m.showInviteInput && len(msg.Runes) > 0 {
			m.inviteInput += string(msg.Runes)
		}
		if m.showFollowInput && len(msg.Runes) > 0 {
			m.followInput += string(msg.Runes)
		}
		return nil
	}
}

func (m *SoClModel) loadFeed() tea.Cmd {
	return func() tea.Msg {
		messages, err := m.scuttlego.GetRecentMessages(100)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		posts := make([]Post, len(messages))
		for i, msg := range messages {
			posts[i] = Post{
				Author: msg.Author,
				Text:   msg.Text,
				Time:   msg.Time,
				pfp:    core.GeneratePFP(msg.Author),
			}
		}

		return FeedLoadedMsg{Posts: posts}
	}
}

func (m *SoClModel) publishPost(text string) tea.Cmd {
	return func() tea.Msg {
		ref, err := m.scuttlego.Publish(text)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return PostPublishedMsg{Ref: ref, Text: text}
	}
}

func (m *SoClModel) loadPeers() tea.Cmd {
	return func() tea.Msg {
		peers, err := m.scuttlego.GetPeers()
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return PeersLoadedMsg{Peers: peers}
	}
}

func (m *SoClModel) redeemInvite(inviteCode string) tea.Cmd {
	return func() tea.Msg {
		err := m.scuttlego.RedeemInvite(inviteCode)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return InviteRedeemedMsg{InviteCode: inviteCode}
	}
}

func (m *SoClModel) followPeer(feedRef string) tea.Cmd {
	return func() tea.Msg {
		err := m.scuttlego.Follow(feedRef)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return FollowedMsg{FeedRef: feedRef}
	}
}

// View renders the TUI to the terminal.
// It uses basic formatting (Lip Gloss to be added).
func (m *SoClModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	if m.errorMsg != "" {
		return fmt.Sprintf("Error: %s\n\nPress Ctrl+C to exit", m.errorMsg)
	}

	// TODO: Implement full UI rendering with Lip Gloss
	// This should render:
	// 1. Header
	// 2. Feed area
	// 3. Composer area
	// 4. Sidebar (peers, trending)

	result := m.renderFeed() + "\n\n" + m.renderComposer()

	if m.showInviteInput {
		result += "\n\n" + m.renderInviteInput()
	}

	if m.showFollowInput {
		result += "\n\n" + m.renderFollowInput()
	}

	if m.showPeers {
		result += "\n\n" + m.renderSidebar()
	}

	result += "\n\nF1: Toggle Peers | F2: Invite | F3: Follow"

	return result
}

// renderFeed renders the social media feed.
func (m *SoClModel) renderFeed() string {
	if m.loading {
		return "Loading feed..."
	}

	if len(m.posts) == 0 {
		return "No posts yet. Press Enter to type a post."
	}

	var feed string
	feed += "=== Feed ===\n\n"

	for i, post := range m.posts {
		prefix := "  "
		if i == m.cursor {
			prefix = "> "
		}
		feed += fmt.Sprintf("%s%d. %s\n", prefix, i+1, post.Author)
		feed += post.pfp
		feed += fmt.Sprintf("%s\n\n", post.Text)
	}

	return feed
}

// renderComposer renders the post input area.
func (m *SoClModel) renderComposer() string {
	prefix := "  "
	if m.editing {
		prefix = "> "
	}

	status := ""
	if m.editing {
		status = "(Editing)"
	} else {
		status = "Press Enter to start typing"
	}

	return fmt.Sprintf("=== Composer ===\n%s%s (%d/280)\n%s",
		prefix,
		m.composerText,
		len(m.composerText),
		status,
	)
}

// renderSidebar renders the peer list sidebar.
func (m *SoClModel) renderSidebar() string {
	if len(m.connectedPeers) == 0 {
		return "=== Peers ===\n\nNo peers connected"
	}

	var sidebar string
	sidebar += "=== Peers ===\n\n"

	for _, peer := range m.connectedPeers {
		sidebar += fmt.Sprintf("  - %s [%s]\n", peer.Address, peer.State)
	}

	return sidebar
}

// renderInviteInput renders the invite code input area.
func (m *SoClModel) renderInviteInput() string {
	return fmt.Sprintf("=== Invite Code ===\n> %s\n\nPress Enter to redeem, F2 to cancel",
		m.inviteInput,
	)
}

// renderFollowInput renders the follow peer input area.
func (m *SoClModel) renderFollowInput() string {
	return fmt.Sprintf("=== Follow Peer ===\n> %s\n\nPress Enter to follow, F3 to cancel",
		m.followInput,
	)
}

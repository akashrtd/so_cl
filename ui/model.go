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

type ErrorMsg struct {
	Err error
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
}

// NewSoClModel creates a new SoClModel with default state.
func NewSoClModel(svc ScuttlegoService) *SoClModel {
	return &SoClModel{
		scuttlego:    svc,
		posts:        []Post{},
		composerText: "",
		peers:        []string{},
		width:        0,
		height:       0,
		cursor:       0,
		editing:      false,
		loading:      false,
		publishing:   false,
		errorMsg:     "",
	}
}

// Init is called at the start of the Bubble Tea program.
// It performs initial setup like initial size calculation.
func (m *SoClModel) Init() tea.Cmd {
	return m.loadFeed()
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

	case tea.KeyEnter:
		if m.editing && !m.publishing {
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
		return nil

	default:
		// Typing
		if m.editing && len(msg.Runes) > 0 {
			m.composerText += string(msg.Runes)
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

	return m.renderFeed() + "\n\n" + m.renderComposer() + "\n\n" + m.renderSidebar()
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
	if len(m.peers) == 0 {
		return "=== Peers ===\n\nNo peers connected"
	}

	var sidebar string
	sidebar += "=== Peers ===\n\n"

	for _, peer := range m.peers {
		sidebar += fmt.Sprintf("  - %s\n", peer)
	}

	return sidebar
}

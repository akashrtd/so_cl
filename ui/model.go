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
	Reply(text, root, branch string) (string, error)
	React(postRef, expression string) (string, error)
	GetRecentMessages(limit int) ([]scuttlego.Message, error)
	Follow(feedRef string) error
	Connect(address string) error
	RedeemInvite(inviteCode string) error
	GetPeers() ([]scuttlego.Peer, error)
	GetEBTStatus() (bool, int, int)
	GetTopHashtags(n int) ([]core.TrendingHashtag, error)
	GetMentions(feedRef string) ([]string, error)
}

type Post struct {
	Ref       string
	Author    string
	Text      string
	Time      int
	pfp       string
	Root      string
	Branch    string
	LikeCount int
}

type FeedLoadedMsg struct {
	Posts []Post
}

type PostPublishedMsg struct {
	Ref  string
	Text string
}

type ReplyPublishedMsg struct {
	Ref    string
	Text   string
	Root   string
	Branch string
}

type ReactionPublishedMsg struct {
	Ref        string
	PostRef    string
	Expression string
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

type TrendingLoadedMsg struct {
	Hashtags []core.TrendingHashtag
}

type MentionsLoadedMsg struct {
	Mentions []string
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
	// showTrending shows if trending sidebar is visible
	showTrending bool
	// showMentions shows if mentions list is visible
	showMentions bool
	// trendingHashtags is the list of trending hashtags
	trendingHashtags []core.TrendingHashtag
	// mentions is the list of mentions for the current user
	mentions []string
	// unreadMentions is the count of unread mentions
	unreadMentions int
	// replyInput shows if user is entering a reply
	replyInput string
	// showReplyInput shows if reply input is visible
	showReplyInput bool
	// replyingTo is the post ref being replied to
	replyingTo string
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
		scuttlego:        svc,
		posts:            []Post{},
		composerText:     "",
		peers:            []string{},
		width:            0,
		height:           0,
		cursor:           0,
		editing:          false,
		loading:          false,
		publishing:       false,
		errorMsg:         "",
		connectedPeers:   []scuttlego.Peer{},
		showPeers:        false,
		showTrending:     false,
		showMentions:     false,
		trendingHashtags: []core.TrendingHashtag{},
		mentions:         []string{},
		unreadMentions:   0,
		replyInput:       "",
		showReplyInput:   false,
		replyingTo:       "",
		inviteInput:      "",
		showInviteInput:  false,
		followInput:      "",
		showFollowInput:  false,
	}
}

// Init is called at the start of the Bubble Tea program.
// It performs initial setup like initial size calculation.
func (m *SoClModel) Init() tea.Cmd {
	return tea.Batch(
		m.loadFeed(),
		m.loadPeers(),
		m.loadTrending(),
		m.loadMentions(),
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
	case TrendingLoadedMsg:
		m.trendingHashtags = msg.Hashtags
		return m, nil
	case MentionsLoadedMsg:
		m.mentions = msg.Mentions
		m.unreadMentions = len(msg.Mentions)
		return m, nil
	case ReplyPublishedMsg:
		m.showReplyInput = false
		m.replyInput = ""
		m.replyingTo = ""
		m.errorMsg = "Reply published successfully"
		return m, m.loadFeed()
	case ReactionPublishedMsg:
		m.errorMsg = "Reaction published successfully"
		return m, m.loadFeed()
	case NewMessageMsg:
		// New message received from EBT replication
		// Add to feed with optimistic UI update
		newPost := Post{
			Ref:       msg.Post.Ref,
			Author:    msg.Post.Author,
			Text:      msg.Post.Text,
			Time:      len(m.posts),
			pfp:       core.GeneratePFP(msg.Post.Author),
			Root:      msg.Post.Root,
			Branch:    msg.Post.Branch,
			LikeCount: msg.Post.LikeCount,
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

	case tea.KeyF4:
		// Toggle reply input (reply to selected post)
		if len(m.posts) > 0 && m.cursor < len(m.posts) {
			m.showReplyInput = !m.showReplyInput
			if m.showReplyInput {
				m.replyingTo = m.posts[m.cursor].Ref
				m.replyInput = ""
			} else {
				m.replyingTo = ""
				m.replyInput = ""
			}
		}
		return nil

	case tea.KeyF5:
		// Like the selected post
		if len(m.posts) > 0 && m.cursor < len(m.posts) {
			postRef := m.posts[m.cursor].Ref
			if postRef != "" {
				return m.reactToPost(postRef, "like")
			}
		}
		return nil

	case tea.KeyF6:
		// Toggle trending sidebar
		m.showTrending = !m.showTrending
		if m.showTrending {
			return m.loadTrending()
		}
		return nil

	case tea.KeyF7:
		// Toggle mentions list
		m.showMentions = !m.showMentions
		if m.showMentions {
			m.unreadMentions = 0
			return m.loadMentions()
		}
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
		} else if m.showReplyInput {
			// Publish reply
			text := m.replyInput
			root := m.replyingTo
			branch := m.replyingTo
			m.replyInput = ""
			m.showReplyInput = false
			m.replyingTo = ""
			return m.replyPost(text, root, branch)
		} else if m.editing && !m.publishing {
			// Publish post
			text := m.composerText
			m.composerText = ""
			m.editing = false
			m.publishing = true

			// Optimistic UI update: add post immediately
			newPost := Post{
				Ref:       "",
				Author:    "You",
				Text:      text,
				Time:      len(m.posts),
				pfp:       core.GeneratePFP("optimistic"),
				Root:      "",
				Branch:    "",
				LikeCount: 0,
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
		if m.showReplyInput && len(m.replyInput) > 0 {
			m.replyInput = m.replyInput[:len(m.replyInput)-1]
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
		if m.showReplyInput && len(msg.Runes) > 0 {
			m.replyInput += string(msg.Runes)
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
				Ref:       msg.Ref,
				Author:    msg.Author,
				Text:      msg.Text,
				Time:      msg.Time,
				pfp:       core.GeneratePFP(msg.Author),
				Root:      msg.Root,
				Branch:    msg.Branch,
				LikeCount: msg.LikeCount,
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

func (m *SoClModel) loadTrending() tea.Cmd {
	return func() tea.Msg {
		hashtags, err := m.scuttlego.GetTopHashtags(10)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return TrendingLoadedMsg{Hashtags: hashtags}
	}
}

func (m *SoClModel) loadMentions() tea.Cmd {
	return func() tea.Msg {
		// TODO: Get current user's feed ref from scuttlego service
		// For now, use empty string to test
		mentions, err := m.scuttlego.GetMentions("")
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return MentionsLoadedMsg{Mentions: mentions}
	}
}

func (m *SoClModel) replyPost(text, root, branch string) tea.Cmd {
	return func() tea.Msg {
		ref, err := m.scuttlego.Reply(text, root, branch)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return ReplyPublishedMsg{Ref: ref, Text: text, Root: root, Branch: branch}
	}
}

func (m *SoClModel) reactToPost(postRef, expression string) tea.Cmd {
	return func() tea.Msg {
		ref, err := m.scuttlego.React(postRef, expression)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return ReactionPublishedMsg{Ref: ref, PostRef: postRef, Expression: expression}
	}
}

// View renders the TUI to the terminal.
// It uses basic formatting (Lip Gloss to be added).
func (m *SoClModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// TODO: Implement full UI rendering with Lip Gloss
	// This should render:
	// 1. Header
	// 2. Feed area
	// 3. Composer area
	// 4. Sidebar (peers, trending)

	result := m.renderHeader()
	result += "\n\n"
	result += m.renderFeed()
	result += "\n\n"
	result += m.renderComposer()

	if m.showReplyInput {
		result += "\n\n" + m.renderReplyInput()
	}

	if m.showInviteInput {
		result += "\n\n" + m.renderInviteInput()
	}

	if m.showFollowInput {
		result += "\n\n" + m.renderFollowInput()
	}

	if m.showPeers {
		result += "\n\n" + m.renderSidebar()
	}

	if m.showTrending {
		result += "\n\n" + m.renderTrending()
	}

	if m.showMentions {
		result += "\n\n" + m.renderMentions()
	}

	result += "\n\n" + m.renderFooter()

	return result
}

// renderHeader renders the header with notification indicator.
func (m *SoClModel) renderHeader() string {
	header := "=== so_cl ===="
	if m.unreadMentions > 0 {
		header += fmt.Sprintf(" [%d new mentions]", m.unreadMentions)
	}
	return header
}

// renderFooter renders the footer with key bindings.
func (m *SoClModel) renderFooter() string {
	footer := "F1: Peers | F2: Invite | F3: Follow | F4: Reply | F5: Like | F6: Trending | F7: Mentions"
	if m.errorMsg != "" {
		footer += fmt.Sprintf("\n%s", m.errorMsg)
	}
	return footer
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

		// Indent replies (thread hierarchy)
		indent := ""
		if post.Root != "" {
			indent = "  "
		}

		feed += fmt.Sprintf("%s%s%d. %s", indent, prefix, i+1, post.Author)
		if post.LikeCount > 0 {
			feed += fmt.Sprintf(" [%d likes]", post.LikeCount)
		}
		feed += "\n"
		feed += indent + post.pfp
		feed += fmt.Sprintf("%s\n\n", indent+post.Text)
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

// renderReplyInput renders the reply input area.
func (m *SoClModel) renderReplyInput() string {
	return fmt.Sprintf("=== Reply ===\n> %s (%d/280)\n\nPress Enter to reply, F4 to cancel",
		m.replyInput,
		len(m.replyInput),
	)
}

// renderTrending renders the trending hashtags sidebar.
func (m *SoClModel) renderTrending() string {
	if len(m.trendingHashtags) == 0 {
		return "=== Trending ===\n\nNo trending hashtags"
	}

	var sidebar string
	sidebar += "=== Trending ===\n\n"

	for i, tag := range m.trendingHashtags {
		sidebar += fmt.Sprintf("  %d. #%s (%d)\n", i+1, tag.Name, tag.Count)
	}

	return sidebar
}

// renderMentions renders the mentions list.
func (m *SoClModel) renderMentions() string {
	if len(m.mentions) == 0 {
		return "=== Mentions ===\n\nNo mentions"
	}

	var list string
	list += "=== Mentions ===\n\n"

	for i, mention := range m.mentions {
		list += fmt.Sprintf("  %d. %s\n", i+1, mention)
	}

	return list
}

package ui

// Package ui provides |-- Bubble Tea UI for so_cl.
// It handles:
// - Model (state machine)
// - View (rendering)
// - Update (event handling)

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/so_cl/core"
	"github.com/yourusername/so_cl/scuttlego"
)

// Page represents the different pages in the application
type Page string

const (
	PageHome     Page = "home"
	PageDiscover Page = "discover"
	PagePeers    Page = "peers"
	PageProfile  Page = "profile"
	PageSettings Page = "settings"
)

const maxPostsInMemory = 100

// LipGloss styles for the TUI
var (
	// Colors
	primaryColor   = lipgloss.Color("#00ff00") // Green
	secondaryColor = lipgloss.Color("#0088ff") // Blue
	accentColor    = lipgloss.Color("#ffaa00") // Orange
	dimColor       = lipgloss.Color("#666666") // Dim gray
	borderColor    = lipgloss.Color("#444444") // Border gray

	// Styles
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#222222")).
			Padding(0, 1)

	navStyle = lipgloss.NewStyle().
			Width(20).
			Border(lipgloss.NormalBorder()).
			BorderForeground(borderColor).
			Padding(1)

	navItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Padding(0, 1)

	navActiveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(primaryColor).
			Padding(0, 1).
			Bold(true)

	mainStyle = lipgloss.NewStyle().
			Width(50).
			Border(lipgloss.NormalBorder()).
			BorderForeground(borderColor).
			Padding(1)

	sidebarStyle = lipgloss.NewStyle().
			Width(30).
			Border(lipgloss.NormalBorder()).
			BorderForeground(borderColor).
			Padding(1)

	sectionTitleStyle = lipgloss.NewStyle().
				Foreground(secondaryColor).
				Bold(true).
				Underline(true)

	postStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(dimColor).
			Padding(0, 1).
			MarginBottom(1)

	selectedPostStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor).
				Padding(0, 1).
				MarginBottom(1)

	buttonStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Padding(0, 1)

	footerStyle = lipgloss.NewStyle().
			Foreground(dimColor).
			Padding(0, 1)
)

type ScuttlegoService interface {
	Publish(text string) (string, error)
	Reply(text, root, branch string) (string, error)
	React(postRef, expression string) (string, error)
	GetRecentMessages(limit int) ([]scuttlego.Message, error)
	Follow(feedRef string) error
	Unfollow(feedRef string) error
	GetIdentity() string
	Connect(address string) error
	RedeemInvite(inviteCode string) error
	GetPeers() ([]scuttlego.Peer, error)
	GetEBTStatus() (bool, int, int)
	GetTopHashtags(n int) ([]core.TrendingHashtag, error)
	GetMentions(feedRef string) ([]string, error)
	GetFollowing(feedRef string) ([]string, error)
	GetFollowers(feedRef string) ([]string, error)
	IsFollowing(follower, following string) (bool, error)
	GetFollowingCount(feedRef string) (int, error)
	GetFollowersCount(feedRef string) (int, error)
	SearchPosts(query string) ([]string, error)
	FilterByHashtag(hashtag string) ([]string, error)
	FilterByAuthor(author string) ([]string, error)
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

type UnfollowedMsg struct {
	FeedRef string
}

type FollowGraphLoadedMsg struct {
	Following []string
	Followers []string
}

type SearchResultsLoadedMsg struct {
	Results []string
	Query   string
}

type ProfileLoadedMsg struct {
	Profile core.SoClProfile
}

type SettingsLoadedMsg struct {
	Username     string
	LANDiscovery bool
}

type UsernameUpdatedMsg struct {
	Username string
}

type PFPRegeneratedMsg struct {
	PFP string
}

type LANDiscoveryToggledMsg struct {
	Enabled bool
}

// SoClModel is Bubble Tea model for so_cl.
// It holds all application state:
// - Posts (feed)
// - Composer (post input)
// - Peers (sidebar)
// - Configuration
// - Page navigation
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
	// showFollowGraph shows if follow graph sidebar is visible
	showFollowGraph bool
	// following is the list of peers the current user follows
	following []string
	// followers is the list of peers following the current user
	followers []string
	// showSearchInput shows if search input is visible
	showSearchInput bool
	// searchQuery is the current search query
	searchQuery string
	// searchResults is the list of post references matching search
	searchResults []string
	// searchFilterType is the type of filter (text, hashtag, author)
	searchFilterType string
	// showSearchResults shows if search results are visible
	showSearchResults bool
	// showProfile shows if profile view is visible
	showProfile bool
	// profile is the currently viewed profile
	profile core.SoClProfile
	// showSettings shows if settings view is visible
	showSettings bool
	// settingsUsername is the current username
	settingsUsername string
	// settingsLANDiscovery is the current LAN discovery setting
	settingsLANDiscovery bool
	// settingsInput shows if user is entering a setting value
	settingsInput string
	// settingsEditing shows which setting is being edited
	settingsEditing string
	// currentPage is the currently active page
	currentPage Page
	// navCursor is the cursor position in the navigation menu
	navCursor int
	// appVersion is the application version
	appVersion string
}

// NewSoClModel creates a new SoClModel with default state.
func NewSoClModel(svc ScuttlegoService) *SoClModel {
	return &SoClModel{
		scuttlego:            svc,
		posts:                []Post{},
		composerText:         "",
		peers:                []string{},
		width:                0,
		height:               0,
		cursor:               0,
		editing:              false,
		loading:              false,
		publishing:           false,
		errorMsg:             "",
		connectedPeers:       []scuttlego.Peer{},
		showPeers:            false,
		showTrending:         false,
		showMentions:         false,
		trendingHashtags:     []core.TrendingHashtag{},
		mentions:             []string{},
		unreadMentions:       0,
		replyInput:           "",
		showReplyInput:       false,
		replyingTo:           "",
		inviteInput:          "",
		showInviteInput:      false,
		followInput:          "",
		showFollowInput:      false,
		showFollowGraph:      false,
		following:            []string{},
		followers:            []string{},
		showSearchInput:      false,
		searchQuery:          "",
		searchResults:        []string{},
		searchFilterType:     "text",
		showSearchResults:    false,
		showProfile:          false,
		profile:              core.SoClProfile{},
		showSettings:         false,
		settingsUsername:     "",
		settingsLANDiscovery: true,
		settingsInput:        "",
		settingsEditing:      "",
		currentPage:          PageHome,
		navCursor:            0,
		appVersion:           "v0.1.5",
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
		m.loadFollowGraph(),
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
		return m, m.loadFollowGraph()
	case UnfollowedMsg:
		m.errorMsg = "Unfollowed successfully"
		return m, m.loadFollowGraph()
	case FollowGraphLoadedMsg:
		m.following = msg.Following
		m.followers = msg.Followers
		return m, nil
	case SearchResultsLoadedMsg:
		m.searchResults = msg.Results
		m.searchQuery = msg.Query
		m.showSearchResults = true
		return m, nil
	case ProfileLoadedMsg:
		m.profile = msg.Profile
		m.showProfile = true
		return m, nil
	case SettingsLoadedMsg:
		m.settingsUsername = msg.Username
		m.settingsLANDiscovery = msg.LANDiscovery
		return m, nil
	case UsernameUpdatedMsg:
		m.settingsUsername = msg.Username
		m.errorMsg = "Username updated successfully"
		return m, nil
	case PFPRegeneratedMsg:
		m.errorMsg = "PFP regenerated successfully"
		return m, nil
	case LANDiscoveryToggledMsg:
		m.settingsLANDiscovery = msg.Enabled
		m.errorMsg = fmt.Sprintf("LAN discovery %s", map[bool]string{true: "enabled", false: "disabled"}[msg.Enabled])
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
	// If editing any input field, handle input first
	if m.editing || m.showReplyInput || m.showInviteInput || m.showFollowInput || m.showSearchInput || (m.showSettings && m.settingsEditing != "") {
		return m.handleInputKeyMsg(msg)
	}

	// Otherwise, handle navigation
	return m.handleNavigationKeyMsg(msg)
}

// handleInputKeyMsg handles keyboard input when editing text fields.
func (m *SoClModel) handleInputKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		// Cancel editing
		m.editing = false
		m.showReplyInput = false
		m.showInviteInput = false
		m.showFollowInput = false
		m.showSearchInput = false
		if m.showSettings {
			m.settingsEditing = ""
			m.settingsInput = ""
		}
		return nil

	case tea.KeyEnter:
		if m.showSettings && m.settingsEditing != "" {
			// Save setting value
			value := m.settingsInput
			m.settingsInput = ""
			m.settingsEditing = ""
			switch m.settingsEditing {
			case "username":
				return m.updateUsername(value)
			case "pfp":
				return m.regeneratePFP()
			case "lan":
				enabled := value == "y" || value == "yes" || value == "1"
				return m.toggleLANDiscovery(enabled)
			}
			return nil
		} else if m.showSearchInput {
			// Perform search
			query := m.searchQuery
			m.searchQuery = ""
			m.showSearchInput = false
			return m.performSearch(query)
		} else if m.showInviteInput {
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
		if m.showSearchInput && len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
		}
		if m.showSettings && len(m.settingsInput) > 0 {
			m.settingsInput = m.settingsInput[:len(m.settingsInput)-1]
		}
		return nil

	case tea.KeyTab:
		// Switch search filter type
		if m.showSearchInput {
			types := []string{"text", "hashtag", "author"}
			for i, t := range types {
				if t == m.searchFilterType {
					m.searchFilterType = types[(i+1)%len(types)]
					break
				}
			}
		}
		return nil

	case tea.KeyRunes:
		// Handle single-character settings commands
		if len(msg.Runes) == 1 && m.showSettings && m.settingsEditing == "" {
			switch msg.Runes[0] {
			case '1':
				// Edit username setting
				m.settingsEditing = "username"
				m.settingsInput = m.settingsUsername
				return nil
			case '2':
				// Toggle LAN discovery setting
				m.settingsEditing = "lan"
				m.settingsInput = map[bool]string{true: "y", false: "n"}[m.settingsLANDiscovery]
				return nil
			case '3':
				// Regenerate PFP setting
				return m.regeneratePFP()
			}
		}
		// Fall through to typing
		fallthrough

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
		if m.showSearchInput && len(msg.Runes) > 0 {
			m.searchQuery += string(msg.Runes)
		}
		if m.showSettings && len(msg.Runes) > 0 && m.settingsEditing != "" {
			m.settingsInput += string(msg.Runes)
		}
		return nil
	}
}

// handleNavigationKeyMsg handles keyboard input for navigation.
func (m *SoClModel) handleNavigationKeyMsg(msg tea.KeyMsg) tea.Cmd {
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
		// Like selected post
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

	case tea.KeyF8:
		// Toggle search input
		m.showSearchInput = !m.showSearchInput
		m.searchQuery = ""
		return nil

	case tea.KeyF9:
		// Toggle profile view
		m.showProfile = !m.showProfile
		if m.showProfile {
			myFeedRef := m.scuttlego.GetIdentity()
			return m.loadProfile(myFeedRef)
		}
		return nil

	case tea.KeyF10:
		// Toggle settings view
		m.showSettings = !m.showSettings
		if m.showSettings {
			return m.loadSettings()
		}
		return nil

	case tea.KeyF11:
		// Toggle follow graph sidebar
		m.showFollowGraph = !m.showFollowGraph
		if m.showFollowGraph {
			return m.loadFollowGraph()
		}
		return nil

	case tea.KeyEnter:
		// Start editing post
		if !m.editing {
			m.editing = true
			return nil
		}
		return nil

	case tea.KeyUp:
		// Navigate up in navigation menu or feed
		if m.navCursor > 0 {
			m.navCursor--
			// Update current page
			pages := []Page{PageHome, PageDiscover, PagePeers, PageProfile, PageSettings}
			m.currentPage = pages[m.navCursor]
		} else if m.cursor > 0 {
			m.cursor--
		}
		return nil

	case tea.KeyDown:
		// Navigate down in navigation menu or feed
		pages := []Page{PageHome, PageDiscover, PagePeers, PageProfile, PageSettings}
		if m.navCursor < len(pages)-1 {
			m.navCursor++
			// Update current page
			m.currentPage = pages[m.navCursor]
		} else if m.cursor < len(m.posts)-1 {
			m.cursor++
		}
		return nil

	default:
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

func (m *SoClModel) unfollowPeer(feedRef string) tea.Cmd {
	return func() tea.Msg {
		err := m.scuttlego.Unfollow(feedRef)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return UnfollowedMsg{FeedRef: feedRef}
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

func (m *SoClModel) loadFollowGraph() tea.Cmd {
	return func() tea.Msg {
		myFeedRef := m.scuttlego.GetIdentity()
		following, err := m.scuttlego.GetFollowing(myFeedRef)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		followers, err := m.scuttlego.GetFollowers(myFeedRef)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return FollowGraphLoadedMsg{
			Following: following,
			Followers: followers,
		}
	}
}

func (m *SoClModel) performSearch(query string) tea.Cmd {
	return func() tea.Msg {
		var results []string
		var err error

		switch m.searchFilterType {
		case "text":
			results, err = m.scuttlego.SearchPosts(query)
		case "hashtag":
			results, err = m.scuttlego.FilterByHashtag(query)
		case "author":
			results, err = m.scuttlego.FilterByAuthor(query)
		default:
			results, err = m.scuttlego.SearchPosts(query)
		}

		if err != nil {
			return ErrorMsg{Err: err}
		}

		return SearchResultsLoadedMsg{
			Results: results,
			Query:   query,
		}
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

// View renders the TUI to the terminal using Lip Gloss styling.
// It renders a 3-column layout: NAV (20%), MAIN CONTENT (50%), SIDEBAR (30%).
func (m *SoClModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Render header
	header := m.renderHeader()

	// Render 3-column layout
	navPanel := m.renderNavPanel()
	mainPanel := m.renderMainPanel()
	sidebarPanel := m.renderSidebarPanel()

	// Join panels horizontally
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		navPanel,
		mainPanel,
		sidebarPanel,
	)

	// Render footer
	footer := m.renderFooter()

	// Combine all parts
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		footer,
	)
}

// renderHeader renders the header with app name, version, and notification indicator.
func (m *SoClModel) renderHeader() string {
	header := fmt.Sprintf("{＊} so_cl │ Social P2P Network %s", m.appVersion)
	if m.unreadMentions > 0 {
		header += fmt.Sprintf(" [%d new mentions]", m.unreadMentions)
	}
	return headerStyle.Render(header)
}

// renderNavPanel renders the left navigation panel.
func (m *SoClModel) renderNavPanel() string {
	pages := []Page{PageHome, PageDiscover, PagePeers, PageProfile, PageSettings}
	var navItems []string

	for i, page := range pages {
		var item string
		switch page {
		case PageHome:
			item = "[Home]"
		case PageDiscover:
			item = "[Discover]"
		case PagePeers:
			item = "[Peers]"
		case PageProfile:
			item = "[Profile]"
		case PageSettings:
			item = "[Settings]"
		}

		if i == m.navCursor {
			navItems = append(navItems, navActiveStyle.Render(item))
		} else {
			navItems = append(navItems, navItemStyle.Render(item))
		}
	}

	return navStyle.Render(lipgloss.JoinVertical(lipgloss.Left, navItems...))
}

// renderMainPanel renders the main content area based on current page.
func (m *SoClModel) renderMainPanel() string {
	var content string

	switch m.currentPage {
	case PageHome:
		content = m.renderHomePage()
	case PageDiscover:
		content = m.renderDiscoverPage()
	case PagePeers:
		content = m.renderPeersPage()
	case PageProfile:
		content = m.renderProfilePage()
	case PageSettings:
		content = m.renderSettingsPage()
	}

	return mainStyle.Render(content)
}

// renderSidebarPanel renders the right sidebar with network status and ASCII art gallery.
func (m *SoClModel) renderSidebarPanel() string {
	var sections []string

	// Network status section
	networkSection := m.renderNetworkStatus()
	sections = append(sections, networkSection)

	// ASCII art gallery section
	artGallery := m.renderArtGallery()
	sections = append(sections, artGallery)

	return sidebarStyle.Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

// renderNetworkStatus renders the network status section.
func (m *SoClModel) renderNetworkStatus() string {
	var status strings.Builder
	status.WriteString(sectionTitleStyle.Render("NETWORK STATUS"))
	status.WriteString("\n\n")

	// Connection status
	connected := len(m.connectedPeers) > 0
	connStatus := "Disconnected"
	if connected {
		connStatus = "Connected"
	}
	status.WriteString(fmt.Sprintf("Status: %s\n", connStatus))

	// Peers count
	status.WriteString(fmt.Sprintf("Peers: %d\n", len(m.connectedPeers)))

	// Speed (placeholder for now)
	speed := "0.0 NB/s"
	if connected {
		speed = "1.2 NB/s"
	}
	status.WriteString(fmt.Sprintf("Speed: %s\n", speed))

	return status.String()
}

// renderArtGallery renders the ASCII art gallery section.
func (m *SoClModel) renderArtGallery() string {
	var gallery strings.Builder
	gallery.WriteString(sectionTitleStyle.Render("ASCII ART GALLERY"))
	gallery.WriteString("\n\n")

	if len(m.posts) > 0 {
		// Show up to 3 recent ASCII art previews
		previewCount := 3
		if len(m.posts) < previewCount {
			previewCount = len(m.posts)
		}

		for i := 0; i < previewCount; i++ {
			preview := m.posts[i].pfp
			if preview != "" {
				gallery.WriteString(preview)
				gallery.WriteString("\n")
			}
		}
	} else {
		gallery.WriteString("No art yet\n")
	}

	return gallery.String()
}

// renderHomePage renders the home page with feed.
func (m *SoClModel) renderHomePage() string {
	var page strings.Builder
	page.WriteString(sectionTitleStyle.Render("FEED"))
	page.WriteString("\n\n")

	if m.loading {
		page.WriteString("Loading feed...")
		return page.String()
	}

	if len(m.posts) == 0 {
		page.WriteString("No posts yet. Press Enter to type a post.")
		return page.String()
	}

	// Render posts
	for i, post := range m.posts {
		postStyle := postStyle
		if i == m.cursor {
			postStyle = selectedPostStyle
		}

		// Format post
		postText := fmt.Sprintf("@%s: %s", post.Author, post.Text)
		if post.LikeCount > 0 {
			postText += fmt.Sprintf("\n<reply> <share> <like> %d", post.LikeCount)
		} else {
			postText += "\n<reply> <share> <like>"
		}

		postText += fmt.Sprintf("\n%d min ago", post.Time)

		page.WriteString(postStyle.Render(postText))
		page.WriteString("\n")
	}

	// Add composer at bottom
	page.WriteString("\n")
	page.WriteString(m.renderComposer())

	return page.String()
}

// renderDiscoverPage renders the discover page with trending topics and popular posts.
func (m *SoClModel) renderDiscoverPage() string {
	var page strings.Builder
	page.WriteString(sectionTitleStyle.Render("DISCOVER"))
	page.WriteString("\n\n")

	// Trending topics section
	page.WriteString(sectionTitleStyle.Render("TRENDING TOPICS"))
	page.WriteString("\n\n")

	if len(m.trendingHashtags) > 0 {
		for i, tag := range m.trendingHashtags {
			indicator := "  "
			if i == 0 {
				indicator = "* " // Hot
			} else if i == 1 {
				indicator = "^ " // Rising
			} else if i == 2 {
				indicator = "+ " // Growing
			}
			page.WriteString(fmt.Sprintf("%s#%s %d\n", indicator, tag.Name, tag.Count))
		}
	} else {
		page.WriteString("No trending topics yet\n")
	}

	// Popular posts section
	page.WriteString("\n")
	page.WriteString(sectionTitleStyle.Render("POPULAR POSTS"))
	page.WriteString("\n\n")

	if len(m.posts) > 0 {
		// Show top post
		topPost := m.posts[0]
		page.WriteString(fmt.Sprintf("@%s: %s\n", topPost.Author, topPost.Text))
		if topPost.LikeCount > 0 {
			page.WriteString(fmt.Sprintf("<like> %d <share>\n", topPost.LikeCount))
		}
	} else {
		page.WriteString("No popular posts yet\n")
	}

	// New peers to follow section
	page.WriteString("\n")
	page.WriteString(sectionTitleStyle.Render("NEW PEERS TO FOLLOW"))
	page.WriteString("\n\n")

	if len(m.followers) > 0 {
		for i, follower := range m.followers {
			if i < 3 {
				page.WriteString(fmt.Sprintf("[@%s] • %d posts\n", truncateFeedRef(follower), 0))
			}
		}
	} else {
		page.WriteString("No new peers to follow\n")
	}

	return page.String()
}

// renderPeersPage renders the peers page with connected peers and statistics.
func (m *SoClModel) renderPeersPage() string {
	var page strings.Builder
	page.WriteString(sectionTitleStyle.Render("PEERS"))
	page.WriteString("\n\n")

	// Connected peers section
	page.WriteString(sectionTitleStyle.Render("CONNECTED PEERS"))
	page.WriteString("\n\n")

	if len(m.connectedPeers) > 0 {
		for _, peer := range m.connectedPeers {
			indicator := "● " // Active
			if peer.State != "connected" {
				indicator = "○ " // Offline
			}
			page.WriteString(fmt.Sprintf("%s%s\n", indicator, peer.Address))
			page.WriteString(fmt.Sprintf("  %s\n", peer.State))
			page.WriteString("  [Message] [Follow] [Info]\n")
		}
	} else {
		page.WriteString("No connected peers\n")
	}

	// Peer statistics section
	page.WriteString("\n")
	page.WriteString(sectionTitleStyle.Render("PEER STATISTICS"))
	page.WriteString("\n\n")

	page.WriteString(fmt.Sprintf("Total Connections: %d\n", len(m.connectedPeers)))
	page.WriteString(fmt.Sprintf("Active Sessions: %d\n", len(m.connectedPeers)))
	page.WriteString("Data Transferred: 0 GB\n")
	page.WriteString("Network Health: 100%\n")

	return page.String()
}

// renderProfilePage renders the profile page with user info and statistics.
func (m *SoClModel) renderProfilePage() string {
	var page strings.Builder
	page.WriteString(sectionTitleStyle.Render("PROFILE"))
	page.WriteString("\n\n")

	// Profile card
	page.WriteString("╔═════════════════════╗\n")
	page.WriteString(fmt.Sprintf("║  %s  ║\n", m.profile.PFP))
	page.WriteString("╚═════════════════════╝\n")
	page.WriteString(fmt.Sprintf("@%s\n", m.profile.Username))
	if m.profile.Bio != "" {
		page.WriteString(fmt.Sprintf("%s\n", m.profile.Bio))
	}

	// Statistics section
	page.WriteString("\n")
	page.WriteString(sectionTitleStyle.Render("PROFILE STATISTICS"))
	page.WriteString("\n\n")

	page.WriteString(fmt.Sprintf("Posts: %d    Following: %d\n", m.profile.PostCount, m.profile.FollowingCount))
	page.WriteString(fmt.Sprintf("Followers: %d  Likes: 0\n", m.profile.FollowersCount))
	page.WriteString("ASCII Arts: 0  Reposts: 0\n")

	return page.String()
}

// renderSettingsPage renders the settings page.
func (m *SoClModel) renderSettingsPage() string {
	var page strings.Builder
	page.WriteString(sectionTitleStyle.Render("SETTINGS"))
	page.WriteString("\n\n")

	page.WriteString(fmt.Sprintf("1. Username: %s\n", m.settingsUsername))
	page.WriteString(fmt.Sprintf("2. LAN Discovery: %s\n", map[bool]string{true: "Enabled", false: "Disabled"}[m.settingsLANDiscovery]))
	page.WriteString("3. Regenerate PFP\n")

	if m.settingsEditing != "" {
		page.WriteString(fmt.Sprintf("\nEditing: %s\n> %s\n", m.settingsEditing, m.settingsInput))
		page.WriteString("Press Enter to save, Esc to cancel")
	} else {
		page.WriteString("\nPress 1-3 to edit")
	}

	return page.String()
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

	return fmt.Sprintf("%s%s (%d/280)\n%s",
		prefix,
		m.composerText,
		len(m.composerText),
		status,
	)
}

// renderFooter renders the footer with key bindings and error messages.
func (m *SoClModel) renderFooter() string {
	var footer strings.Builder
	footer.WriteString("↑/↓: Navigate | Enter: Edit | F1: Peers | F2: Invite | F3: Follow | F4: Reply | F5: Like | F6: Trending | F7: Mentions | F8: Search | F9: Profile | F10: Settings | F11: Follow Graph | Ctrl+C: Quit")

	if m.errorMsg != "" {
		footer.WriteString("\n")
		footer.WriteString(buttonStyle.Render(m.errorMsg))
	}

	return footerStyle.Render(footer.String())
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

// renderFollowGraph renders the follow/followers sidebar.
func (m *SoClModel) renderFollowGraph() string {
	var sidebar string
	sidebar += "=== Follow Graph ===\n\n"

	sidebar += fmt.Sprintf("Following: %d\n", len(m.following))
	for i, f := range m.following {
		sidebar += fmt.Sprintf("  %d. %s\n", i+1, truncateFeedRef(f))
	}

	sidebar += fmt.Sprintf("\nFollowers: %d\n", len(m.followers))
	for i, f := range m.followers {
		sidebar += fmt.Sprintf("  %d. %s\n", i+1, truncateFeedRef(f))
	}

	return sidebar
}

// truncateFeedRef truncates a feed reference for display.
func truncateFeedRef(feedRef string) string {
	if len(feedRef) <= 20 {
		return feedRef
	}
	return feedRef[:20] + "..."
}

// renderSearchInput renders the search input area.
func (m *SoClModel) renderSearchInput() string {
	return fmt.Sprintf("=== Search ===\n> %s\n\nFilter: %s (Tab to switch)\nPress Enter to search, F8 to cancel",
		m.searchQuery,
		m.searchFilterType,
	)
}

// renderSearchResults renders the search results.
func (m *SoClModel) renderSearchResults() string {
	if len(m.searchResults) == 0 {
		return fmt.Sprintf("=== Search Results ===\n\nNo results for: %s\n\nPress F8 to search again", m.searchQuery)
	}

	var results string
	results += fmt.Sprintf("=== Search Results (%s) ===\n\n", m.searchFilterType)

	for i, ref := range m.searchResults {
		results += fmt.Sprintf("  %d. %s\n", i+1, ref)
	}

	results += "\nPress F8 to search again"
	return results
}

// loadProfile loads a user's profile.
func (m *SoClModel) loadProfile(feedRef string) tea.Cmd {
	return func() tea.Msg {
		// Get following/followers counts
		followingCount, _ := m.scuttlego.GetFollowingCount(feedRef)
		followersCount, _ := m.scuttlego.GetFollowersCount(feedRef)

		// Get user's posts (simplified - just count for now)
		// TODO: Get actual posts by author
		postCount := 0

		profile := core.SoClProfile{
			FeedRef:        feedRef,
			Username:       feedRef,
			PFP:            core.GeneratePFP(feedRef),
			Bio:            "",
			FollowingCount: followingCount,
			FollowersCount: followersCount,
			PostCount:      postCount,
		}

		return ProfileLoadedMsg{Profile: profile}
	}
}

// renderProfile renders the profile view.
func (m *SoClModel) renderProfile() string {
	var profile string
	profile += "=== Profile ===\n\n"
	profile += fmt.Sprintf("Feed Ref: %s\n", m.profile.FeedRef)
	profile += fmt.Sprintf("PFP:\n%s\n", m.profile.PFP)
	profile += fmt.Sprintf("Following: %d\n", m.profile.FollowingCount)
	profile += fmt.Sprintf("Followers: %d\n", m.profile.FollowersCount)
	profile += fmt.Sprintf("Posts: %d\n", m.profile.PostCount)

	if m.profile.Bio != "" {
		profile += fmt.Sprintf("Bio: %s\n", m.profile.Bio)
	}

	profile += "\nPress F9 to close profile"
	return profile
}

// loadSettings loads the current settings.
func (m *SoClModel) loadSettings() tea.Cmd {
	return func() tea.Msg {
		// Get current username from BadgerDB (simplified)
		// TODO: Implement proper settings storage
		username := m.settingsUsername
		if username == "" {
			username = m.scuttlego.GetIdentity()
		}

		return SettingsLoadedMsg{
			Username:     username,
			LANDiscovery: m.settingsLANDiscovery,
		}
	}
}

// renderSettings renders the settings view.
func (m *SoClModel) renderSettings() string {
	var settings string
	settings += "=== Settings ===\n\n"

	settings += fmt.Sprintf("1. Username: %s\n", m.settingsUsername)
	settings += fmt.Sprintf("2. LAN Discovery: %s\n", map[bool]string{true: "Enabled", false: "Disabled"}[m.settingsLANDiscovery])
	settings += fmt.Sprintf("3. Regenerate PFP\n")

	if m.settingsEditing != "" {
		settings += fmt.Sprintf("\nEditing: %s\n> %s\n", m.settingsEditing, m.settingsInput)
		settings += "Press Enter to save, Esc to cancel"
	} else {
		settings += "\nPress 1-3 to edit, F10 to close"
	}

	return settings
}

// updateUsername updates the username.
func (m *SoClModel) updateUsername(username string) tea.Cmd {
	return func() tea.Msg {
		// TODO: Store username in BadgerDB
		return UsernameUpdatedMsg{Username: username}
	}
}

// regeneratePFP regenerates the ASCII profile picture.
func (m *SoClModel) regeneratePFP() tea.Cmd {
	return func() tea.Msg {
		myFeedRef := m.scuttlego.GetIdentity()
		pfp := core.GeneratePFP(myFeedRef)
		return PFPRegeneratedMsg{PFP: pfp}
	}
}

// toggleLANDiscovery toggles LAN discovery.
func (m *SoClModel) toggleLANDiscovery(enabled bool) tea.Cmd {
	return func() tea.Msg {
		// TODO: Implement actual LAN discovery toggle
		return LANDiscoveryToggledMsg{Enabled: enabled}
	}
}

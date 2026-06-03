package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/viewport"
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

type FeedLoadedMsg struct {
	Posts []core.Post
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

type NetworkStatusUpdatedMsg struct {
	Connected bool
	Speed     float64
}

// SoClModel is Bubble Tea model for so_cl.
// It holds all application state:
// - Posts (feed)
// - Composer (post input)
// - Peers (sidebar)
// - Configuration
// - Page navigation
type SoClModel struct {
	scuttlego scuttlego.ScuttlegoService
	// posts is the feed of social media posts
	posts []core.Post
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
	settingsEditing   string
	currentPage       Page
	navCursor         int
	appVersion        string
	internetConnected bool
	networkSpeed      float64
	viewport          viewport.Model
}

// NewSoClModel creates a new SoClModel with default state.
func NewSoClModel(svc scuttlego.ScuttlegoService) *SoClModel {
	return &SoClModel{
		scuttlego:            svc,
		posts:                []core.Post{},
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
		internetConnected:    false,
		networkSpeed:         0,
		viewport:             viewport.New(0, 0),
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
		tea.Tick(time.Second*10, func(t time.Time) tea.Msg {
			return t
		}),
	)
}

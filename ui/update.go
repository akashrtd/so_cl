package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourusername/so_cl/core"
)

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
		headerHeight := 3
		footerHeight := 3
		mainHeight := m.height - headerHeight - footerHeight

		navWidth := 20
		sidebarWidth := 30
		mainWidth := m.width - navWidth - sidebarWidth - 4

		m.viewport = viewport.New(mainWidth, mainHeight)
		return m, nil
	case FeedLoadedMsg:
		m.posts = msg.Posts
		m.loading = false
		m.publishing = false
		m.updateViewportContent()
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
		newPost := core.Post{
			Ref:       msg.Post.Ref,
			Author:    msg.Post.Author,
			Text:      msg.Post.Text,
			Timestamp: int64(msg.Post.Time),
			PFP:       core.GeneratePFP(msg.Post.Author),
			Root:      msg.Post.Root,
			Branch:    msg.Post.Branch,
			LikeCount: msg.Post.LikeCount,
		}
		m.posts = append([]core.Post{newPost}, m.posts...)

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
	case NetworkStatusUpdatedMsg:
		m.internetConnected = msg.Connected
		m.networkSpeed = msg.Speed
		return m, nil
	case time.Time:
		m.scuttlego.UpdateNetworkStatus()
		connected, speed, _ := m.scuttlego.GetNetworkStatus()
		m.internetConnected = connected
		m.networkSpeed = speed
		return m, tea.Tick(time.Second*10, func(t time.Time) tea.Msg {
			return t
		})
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
			newPost := core.Post{
				Ref:       "",
				Author:    "You",
				Text:      text,
				Timestamp: int64(len(m.posts)),
				PFP:       core.GeneratePFP("optimistic"),
				Root:      "",
				Branch:    "",
				LikeCount: 0,
			}
			m.posts = append([]core.Post{newPost}, m.posts...)

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
			pages := []Page{PageHome, PageDiscover, PagePeers, PageProfile, PageSettings}
			m.currentPage = pages[m.navCursor]
		} else if m.currentPage == PageHome {
			if m.cursor > 0 {
				m.cursor--
				m.viewport.LineUp(1)
			}
		}
		return nil

	case tea.KeyDown:
		// Navigate down in navigation menu or feed
		pages := []Page{PageHome, PageDiscover, PagePeers, PageProfile, PageSettings}
		if m.navCursor < len(pages)-1 {
			m.navCursor++
			m.currentPage = pages[m.navCursor]
		} else if m.currentPage == PageHome {
			if m.cursor < len(m.posts)-1 {
				m.cursor++
				m.viewport.LineDown(1)
			}
		}
		return nil

	default:
		return nil
	}
}

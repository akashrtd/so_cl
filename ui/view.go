package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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

// renderModalOverlay renders the modal overlay for input dialogs.
func (m *SoClModel) renderModalOverlay() string {
	if m.showInviteInput {
		return m.renderInviteInput()
	}
	if m.showFollowInput {
		return m.renderFollowInput()
	}
	if m.showReplyInput {
		return m.renderReplyInput()
	}
	if m.showSearchInput {
		return m.renderSearchInput()
	}
	if m.showSearchResults {
		return m.renderSearchResults()
	}
	if m.showPeers {
		return m.renderSidebar()
	}
	if m.showTrending {
		return m.renderTrending()
	}
	if m.showMentions {
		return m.renderMentions()
	}
	if m.showFollowGraph {
		return m.renderFollowGraph()
	}
	if m.showProfile {
		return m.renderProfile()
	}
	if m.showSettings && m.currentPage != PageSettings {
		return m.renderSettings()
	}
	return ""
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

	connStatus := "Disconnected"
	if m.internetConnected {
		connStatus = "Connected"
	}
	status.WriteString(fmt.Sprintf("Status: %s\n", connStatus))

	status.WriteString(fmt.Sprintf("Peers: %d\n", len(m.connectedPeers)))

	speedStr := "0.0 KB/s"
	if m.networkSpeed > 0 {
		speedStr = fmt.Sprintf("%.1f KB/s", m.networkSpeed)
	}
	status.WriteString(fmt.Sprintf("Speed: %s\n", speedStr))

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
			preview := m.posts[i].PFP
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

func (m *SoClModel) updateViewportContent() {
	if len(m.posts) == 0 || m.viewport.Height == 0 {
		return
	}

	var feedContent strings.Builder
	for i, post := range m.posts {
		style := postStyle
		if i == m.cursor {
			style = selectedPostStyle
		}

		postText := fmt.Sprintf("@%s: %s", post.Author, post.Text)
		if post.LikeCount > 0 {
			postText += fmt.Sprintf("\n<reply> <share> <like> %d", post.LikeCount)
		} else {
			postText += "\n<reply> <share> <like>"
		}

		postText += fmt.Sprintf("\n%d min ago", post.Timestamp)

		feedContent.WriteString(style.Render(postText))
		feedContent.WriteString("\n")
	}

	m.viewport.SetContent(feedContent.String())
}

// renderPageHeader renders the page header.
func (m *SoClModel) renderPageHeader(title string) string {
	return sectionTitleStyle.Render(title) + "\n\n"
}

// truncateFeedRef truncates a feed reference for display.
func truncateFeedRef(feedRef string) string {
	if len(feedRef) <= 20 {
		return feedRef
	}
	return feedRef[:20] + "..."
}

// truncateString truncates a string to the given length.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
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
		feed += indent + post.PFP
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

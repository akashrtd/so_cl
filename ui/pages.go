package ui

import (
	"fmt"
	"strings"
)

// renderHomePage renders the home page with feed.
func (m *SoClModel) renderHomePage() string {
	if m.loading {
		return sectionTitleStyle.Render("FEED") + "\n\nLoading feed..."
	}

	if len(m.posts) == 0 {
		return sectionTitleStyle.Render("FEED") + "\n\nNo posts yet. Press Enter to type a post."
	}

	// Build page content
	var page strings.Builder
	page.WriteString(sectionTitleStyle.Render("FEED"))
	page.WriteString("\n\n")
	page.WriteString(m.viewport.View())
	page.WriteString("\n\n")
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

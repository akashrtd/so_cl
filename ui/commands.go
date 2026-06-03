package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourusername/so_cl/core"
)

func (m *SoClModel) loadFeed() tea.Cmd {
	return func() tea.Msg {
		messages, err := m.scuttlego.GetRecentMessages(100)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		posts := make([]core.Post, len(messages))
		for i, msg := range messages {
			posts[i] = core.Post{
				Ref:       msg.Ref,
				Author:    msg.Author,
				Text:      msg.Text,
				Timestamp: int64(msg.Time),
				PFP:       core.GeneratePFP(msg.Author),
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

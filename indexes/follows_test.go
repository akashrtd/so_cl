package indexes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndexFollow(t *testing.T) {
	t.Run("indexes follow relationship", func(t *testing.T) {
		indexer := setupTestIndexer(t)

		follower := "@alice.ed25519"
		following := "@bob.ed25519"

		err := indexer.IndexFollow(follower, following)
		require.NoError(t, err)

		// Verify following list
		followingList, err := indexer.GetFollowing(follower)
		require.NoError(t, err)
		require.Len(t, followingList, 1)
		require.Equal(t, following, followingList[0])

		// Verify followers list
		followersList, err := indexer.GetFollowers(following)
		require.NoError(t, err)
		require.Len(t, followersList, 1)
		require.Equal(t, follower, followersList[0])
	})

	t.Run("ignores duplicate follows", func(t *testing.T) {
		indexer := setupTestIndexer(t)

		follower := "@alice.ed25519"
		following := "@bob.ed25519"

		err1 := indexer.IndexFollow(follower, following)
		err2 := indexer.IndexFollow(follower, following)
		require.NoError(t, err1)
		require.NoError(t, err2)

		followingList, err := indexer.GetFollowing(follower)
		require.NoError(t, err)
		require.Len(t, followingList, 1)
	})
}

func TestUnfollow(t *testing.T) {
	t.Run("removes follow relationship", func(t *testing.T) {
		indexer := setupTestIndexer(t)

		follower := "@alice.ed25519"
		following := "@bob.ed25519"

		indexer.IndexFollow(follower, following)
		err := indexer.Unfollow(follower, following)
		require.NoError(t, err)

		followingList, err := indexer.GetFollowing(follower)
		require.NoError(t, err)
		require.Empty(t, followingList)

		followersList, err := indexer.GetFollowers(following)
		require.NoError(t, err)
		require.Empty(t, followersList)
	})

	t.Run("handles non-existent follow", func(t *testing.T) {
		indexer := setupTestIndexer(t)

		err := indexer.Unfollow("@alice.ed25519", "@bob.ed25519")
		require.NoError(t, err)
	})
}

func TestIsFollowing(t *testing.T) {
	t.Run("returns true when following", func(t *testing.T) {
		indexer := setupTestIndexer(t)

		follower := "@alice.ed25519"
		following := "@bob.ed25519"

		indexer.IndexFollow(follower, following)

		isFollowing, err := indexer.IsFollowing(follower, following)
		require.NoError(t, err)
		require.True(t, isFollowing)
	})

	t.Run("returns false when not following", func(t *testing.T) {
		indexer := setupTestIndexer(t)

		isFollowing, err := indexer.IsFollowing("@alice.ed25519", "@bob.ed25519")
		require.NoError(t, err)
		require.False(t, isFollowing)
	})
}

func TestGetFollowing(t *testing.T) {
	t.Run("returns empty list when not following anyone", func(t *testing.T) {
		indexer := setupTestIndexer(t)

		following, err := indexer.GetFollowing("@alice.ed25519")
		require.NoError(t, err)
		require.Empty(t, following)
	})

	t.Run("returns following list", func(t *testing.T) {
		indexer := setupTestIndexer(t)

		follower := "@alice.ed25519"
		following1 := "@bob.ed25519"
		following2 := "@charlie.ed25519"

		indexer.IndexFollow(follower, following1)
		indexer.IndexFollow(follower, following2)

		following, err := indexer.GetFollowing(follower)
		require.NoError(t, err)
		require.Len(t, following, 2)
		require.Contains(t, following, following1)
		require.Contains(t, following, following2)
	})
}

func TestGetFollowers(t *testing.T) {
	t.Run("returns empty list when no followers", func(t *testing.T) {
		indexer := setupTestIndexer(t)

		followers, err := indexer.GetFollowers("@bob.ed25519")
		require.NoError(t, err)
		require.Empty(t, followers)
	})

	t.Run("returns followers list", func(t *testing.T) {
		indexer := setupTestIndexer(t)

		following := "@bob.ed25519"
		follower1 := "@alice.ed25519"
		follower2 := "@charlie.ed25519"

		indexer.IndexFollow(follower1, following)
		indexer.IndexFollow(follower2, following)

		followers, err := indexer.GetFollowers(following)
		require.NoError(t, err)
		require.Len(t, followers, 2)
		require.Contains(t, followers, follower1)
		require.Contains(t, followers, follower2)
	})
}

func TestGetFollowingCount(t *testing.T) {
	t.Run("returns correct count", func(t *testing.T) {
		indexer := setupTestIndexer(t)

		follower := "@alice.ed25519"

		indexer.IndexFollow(follower, "@bob.ed25519")
		indexer.IndexFollow(follower, "@charlie.ed25519")
		indexer.IndexFollow(follower, "@dave.ed25519")

		count, err := indexer.GetFollowingCount(follower)
		require.NoError(t, err)
		require.Equal(t, 3, count)
	})

	t.Run("returns zero when not following anyone", func(t *testing.T) {
		indexer := setupTestIndexer(t)

		count, err := indexer.GetFollowingCount("@alice.ed25519")
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})
}

func TestGetFollowersCount(t *testing.T) {
	t.Run("returns correct count", func(t *testing.T) {
		indexer := setupTestIndexer(t)

		following := "@bob.ed25519"

		indexer.IndexFollow("@alice.ed25519", following)
		indexer.IndexFollow("@charlie.ed25519", following)
		indexer.IndexFollow("@dave.ed25519", following)

		count, err := indexer.GetFollowersCount(following)
		require.NoError(t, err)
		require.Equal(t, 3, count)
	})

	t.Run("returns zero when no followers", func(t *testing.T) {
		indexer := setupTestIndexer(t)

		count, err := indexer.GetFollowersCount("@alice.ed25519")
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})
}

func TestGetFollowRelationship(t *testing.T) {
	t.Run("returns timestamp when following", func(t *testing.T) {
		indexer := setupTestIndexer(t)

		follower := "@alice.ed25519"
		following := "@bob.ed25519"

		indexer.IndexFollow(follower, following)

		timestamp, err := indexer.GetFollowRelationship(follower, following)
		require.NoError(t, err)
		require.Greater(t, timestamp, int64(0))
	})

	t.Run("returns zero when not following", func(t *testing.T) {
		indexer := setupTestIndexer(t)

		timestamp, err := indexer.GetFollowRelationship("@alice.ed25519", "@bob.ed25519")
		require.NoError(t, err)
		require.Equal(t, int64(0), timestamp)
	})
}

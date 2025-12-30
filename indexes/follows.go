package indexes

// Package indexes provides follow graph tracking for so_cl.
// It maintains custom indexes on top of BadgerDB for:
// - Following relationships
// - Follower relationships
// - Follow counts

import (
	"fmt"
	"strings"
	"time"

	badger "github.com/dgraph-io/badger/v3"
)

// FollowRelationship represents a follow relationship between two peers.
type FollowRelationship struct {
	// Follower is the feed reference of the follower
	Follower string
	// Following is the feed reference of the peer being followed
	Following string
	// Timestamp is when the follow occurred
	Timestamp int64
}

// IndexFollow indexes a follow relationship.
// This updates BadgerDB with:
// - Following list for the follower
// - Follower list for the followed
// - Follow counts
//
// Returns an error if indexing fails.
func (idx *Indexer) IndexFollow(follower, following string) error {
	timestamp := time.Now().Unix()

	return idx.db.Update(func(txn *badger.Txn) error {
		// Add to follower's following list
		followingKey := []byte("following:" + follower)
		var followingList []string
		item, err := txn.Get(followingKey)
		if err == nil {
			item.Value(func(val []byte) error {
				followingList = strings.Split(string(val), ",")
				return nil
			})
		}

		// Check if already following
		for _, f := range followingList {
			if f == following {
				return nil // Already following
			}
		}

		followingList = append(followingList, following)
		if err := txn.Set(followingKey, []byte(strings.Join(followingList, ","))); err != nil {
			return err
		}

		// Add to followed's follower list
		followerKey := []byte("followers:" + following)
		var followerList []string
		item, err = txn.Get(followerKey)
		if err == nil {
			item.Value(func(val []byte) error {
				followerList = strings.Split(string(val), ",")
				return nil
			})
		}

		// Check if already in followers list
		for _, f := range followerList {
			if f == follower {
				return nil // Already in followers list
			}
		}

		followerList = append(followerList, follower)
		if err := txn.Set(followerKey, []byte(strings.Join(followerList, ","))); err != nil {
			return err
		}

		// Store follow relationship with timestamp
		relKey := []byte(fmt.Sprintf("follow:%s:%s", follower, following))
		if err := txn.Set(relKey, []byte(fmt.Sprintf("%d", timestamp))); err != nil {
			return err
		}

		return nil
	})
}

// Unfollow removes a follow relationship.
// This updates BadgerDB by removing:
// - The following from the follower's following list
// - The follower from the followed's follower list
//
// Returns an error if removal fails.
func (idx *Indexer) Unfollow(follower, following string) error {
	return idx.db.Update(func(txn *badger.Txn) error {
		// Remove from follower's following list
		followingKey := []byte("following:" + follower)
		var followingList []string
		item, err := txn.Get(followingKey)
		if err == nil {
			item.Value(func(val []byte) error {
				followingList = strings.Split(string(val), ",")
				return nil
			})
		}

		newFollowingList := make([]string, 0, len(followingList))
		for _, f := range followingList {
			if f != following {
				newFollowingList = append(newFollowingList, f)
			}
		}

		if len(newFollowingList) > 0 {
			if err := txn.Set(followingKey, []byte(strings.Join(newFollowingList, ","))); err != nil {
				return err
			}
		} else {
			if err := txn.Delete(followingKey); err != nil {
				return err
			}
		}

		// Remove from followed's follower list
		followerKey := []byte("followers:" + following)
		var followerList []string
		item, err = txn.Get(followerKey)
		if err == nil {
			item.Value(func(val []byte) error {
				followerList = strings.Split(string(val), ",")
				return nil
			})
		}

		newFollowerList := make([]string, 0, len(followerList))
		for _, f := range followerList {
			if f != follower {
				newFollowerList = append(newFollowerList, f)
			}
		}

		if len(newFollowerList) > 0 {
			if err := txn.Set(followerKey, []byte(strings.Join(newFollowerList, ","))); err != nil {
				return err
			}
		} else {
			if err := txn.Delete(followerKey); err != nil {
				return err
			}
		}

		// Remove follow relationship
		relKey := []byte(fmt.Sprintf("follow:%s:%s", follower, following))
		if err := txn.Delete(relKey); err != nil {
			return err
		}

		return nil
	})
}

// GetFollowing retrieves the list of peers that a user is following.
//
// Returns a slice of feed references, or an error if retrieval fails.
func (idx *Indexer) GetFollowing(feedRef string) ([]string, error) {
	var following []string

	key := []byte("following:" + feedRef)

	err := idx.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil // Not following anyone
			}
			return err
		}

		item.Value(func(val []byte) error {
			following = strings.Split(string(val), ",")
			return nil
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return following, nil
}

// GetFollowers retrieves the list of peers following a user.
//
// Returns a slice of feed references, or an error if retrieval fails.
func (idx *Indexer) GetFollowers(feedRef string) ([]string, error) {
	var followers []string

	key := []byte("followers:" + feedRef)

	err := idx.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil // No followers
			}
			return err
		}

		item.Value(func(val []byte) error {
			followers = strings.Split(string(val), ",")
			return nil
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return followers, nil
}

// IsFollowing checks if a user is following another user.
//
// Returns true if following, false otherwise, or an error if check fails.
func (idx *Indexer) IsFollowing(follower, following string) (bool, error) {
	followingList, err := idx.GetFollowing(follower)
	if err != nil {
		return false, err
	}

	for _, f := range followingList {
		if f == following {
			return true, nil
		}
	}

	return false, nil
}

// GetFollowingCount returns the number of peers a user is following.
//
// Returns the count, or an error if retrieval fails.
func (idx *Indexer) GetFollowingCount(feedRef string) (int, error) {
	following, err := idx.GetFollowing(feedRef)
	if err != nil {
		return 0, err
	}
	return len(following), nil
}

// GetFollowersCount returns the number of followers a user has.
//
// Returns the count, or an error if retrieval fails.
func (idx *Indexer) GetFollowersCount(feedRef string) (int, error) {
	followers, err := idx.GetFollowers(feedRef)
	if err != nil {
		return 0, err
	}
	return len(followers), nil
}

// GetFollowRelationship retrieves the follow relationship between two peers.
//
// Returns the timestamp of the follow, or an error if retrieval fails.
// Returns 0 if not following.
func (idx *Indexer) GetFollowRelationship(follower, following string) (int64, error) {
	var timestamp int64

	key := []byte(fmt.Sprintf("follow:%s:%s", follower, following))

	err := idx.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil // Not following
			}
			return err
		}

		item.Value(func(val []byte) error {
			_, err := fmt.Sscanf(string(val), "%d", &timestamp)
			return err
		})

		return nil
	})

	if err != nil {
		return 0, err
	}

	return timestamp, nil
}

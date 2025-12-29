package indexes

// Package indexes provides indexing functionality for so_cl.
// It maintains custom indexes on top of BadgerDB for:
// - Hashtag counting
// - Mention tracking
// - Trending metrics

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	badger "github.com/dgraph-io/badger/v3"
)

// Indexer maintains custom indexes for so_cl.
type Indexer struct {
	db *badger.DB
}

// NewIndexer creates a new indexer with a BadgerDB connection.
//
// Returns an error if BadgerDB cannot be opened.
func NewIndexer(db *badger.DB) *Indexer {
	return &Indexer{db: db}
}

// Close closes the indexer (no-op for BadgerDB as it's managed externally).
func (idx *Indexer) Close() error {
	// BadgerDB is managed externally, no-op here
	return nil
}

// ExtractHashtags extracts hashtag references from text.
// Hashtags are words starting with # (e.g., "#golang", "#scuttlebutt").
//
// Returns a slice of hashtags (without # prefix).
func ExtractHashtags(text string) []string {
	re := regexp.MustCompile(`#(\w+)`)
	matches := re.FindAllStringSubmatch(text, -1)

	tags := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			tags = append(tags, match[1])
		}
	}

	return tags
}

// ExtractMentions extracts @mentions from text.
// Mentions are @-prefixed feed references (e.g., "@alice", "@bob...").
//
// Returns a slice of feed references (without @ prefix).
func ExtractMentions(text string) []string {
	re := regexp.MustCompile(`@([\w\./]+)`)
	matches := re.FindAllStringSubmatch(text, -1)

	mentions := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			mentions = append(mentions, match[1])
		}
	}

	return mentions
}

// IndexPost indexes a post for hashtags and mentions.
// This updates BadgerDB with:
// - Hashtag counts
// - Mention queues
//
// Returns an error if indexing fails.
func (idx *Indexer) IndexPost(postRef string, text string) error {
	tags := ExtractHashtags(text)
	mentions := ExtractMentions(text)

	err := idx.db.Update(func(txn *badger.Txn) error {
		// Index hashtags
		for _, tag := range tags {
			key := []byte("hashtag:" + tag)

			var count int
			item, err := txn.Get(key)
			if err == nil {
				item.Value(func(val []byte) error {
					_, err := fmt.Sscanf(string(val), "%d", &count)
					return err
				})
			}

			count++
			data := fmt.Sprintf("%d", count)

			if err := txn.Set(key, []byte(data)); err != nil {
				return err
			}
		}

		// Index mentions
		for _, mention := range mentions {
			key := []byte("mention:" + mention)

			var mentions []string
			item, err := txn.Get(key)
			if err == nil {
				item.Value(func(val []byte) error {
					mentions = strings.Split(string(val), ",")
					return nil
				})
			}

			mentions = append(mentions, postRef)
			data := strings.Join(mentions, ",")

			if err := txn.Set(key, []byte(data)); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// GetTopHashtags retrieves the top N hashtags by count.
//
// Returns a slice of (hashtag, count) pairs, or an error if retrieval fails.
func (idx *Indexer) GetTopHashtags(n int) ([]HashtagCount, error) {
	var results []HashtagCount

	err := idx.db.View(func(txn *badger.Txn) error {
		prefix := []byte("hashtag:")
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()

			var count int
			item.Value(func(val []byte) error {
				_, err := fmt.Sscanf(string(val), "%d", &count)
				return err
			})

			tag := strings.TrimPrefix(string(item.Key()), "hashtag:")
			results = append(results, HashtagCount{
				Name:  tag,
				Count: count,
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort by count (descending)
	sortHashtagCounts(results)

	// Return top N
	if len(results) > n {
		return results[:n], nil
	}

	return results, nil
}

// GetMentions retrieves all post references mentioning a given user.
// The feedRef is the SSB feed reference of the user (e.g., "@alice...").
//
// Returns a slice of post references, or an error if retrieval fails.
func (idx *Indexer) GetMentions(feedRef string) ([]string, error) {
	var mentions []string

	key := []byte("mention:" + feedRef)

	err := idx.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil // No mentions
			}
			return err
		}

		item.Value(func(val []byte) error {
			mentions = strings.Split(string(val), ",")
			return nil
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return mentions, nil
}

// HashtagCount represents a hashtag with its usage count.
type HashtagCount struct {
	// Name is the hashtag (without # prefix)
	Name string
	// Count is the number of times this hashtag was used
	Count int
}

// sortHashtagCounts sorts hashtag counts by count in descending order.
func sortHashtagCounts(counts []HashtagCount) {
	for i := 0; i < len(counts); i++ {
		for j := i + 1; j < len(counts); j++ {
			if counts[i].Count < counts[j].Count {
				counts[i], counts[j] = counts[j], counts[i]
			}
		}
	}
}

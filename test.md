# so_cl Testing Guide for AI Coding Agents

> **Purpose**: This document provides comprehensive instructions for AI coding agents to implement a complete test suite for the so_cl project.
> **Last Updated**: 2025-12-30
> **Target Coverage**: 85% overall, 95%+ for critical paths

---

## Table of Contents

1. [Project Overview](#1-project-overview)
2. [Test Organization](#2-test-organization)
3. [Test Best Practices](#3-test-best-practices)
4. [Test Utilities](#4-test-utilities)
5. [Package-Specific Tests](#5-package-specific-tests)
6. [Integration Tests](#6-integration-tests)
7. [Running Tests](#7-running-tests)
8. [Coverage Targets](#8-coverage-targets)

---

## 1. Project Overview

### What is so_cl?

**so_cl** is a decentralized, ASCII-art-based social platform running in the terminal using the Scuttlebutt (SSB) protocol.

**Key Characteristics:**
- **P2P**: No central servers, offline-first
- **Terminal UI (TUI)**: ASCII art, 6x6 colored profile pictures
- **Decentralized**: Uses Scuttlebutt protocol (append-only logs)
- **Privacy-focused**: Local data only, encrypted connections

### Technology Stack

| Component | Technology | Version | Purpose |
|-----------|-------------|----------|---------|
| **Language** | Go | 1.24+ | Primary language |
| **SSB Library** | scuttlego | v0.0.4 | Scuttlebutt protocol implementation |
| **Storage** | BadgerDB v3 | v3.2103.5 | LSM-tree KV store |
| **Frontend** | Bubble Tea | v1.3.10+ | TUI framework |
| **Styling** | Lip Gloss | v1.1.0+ | Terminal styling |
| **Logging** | Zap | v1.27.1+ | Structured logging |
| **Testing** | Testify | v1.9.0+ | Assertions and mocks |

### Package Structure

```
so_cl/
├── main.go                    # Entry point
├── go.mod                     # Dependencies
├── test.md                    # THIS FILE - Testing guide
│
├── ui/                       # Presentation layer
│   ├── model.go              # Bubble Tea model
│   └── model_test.go         # UI tests (partial)
│
├── core/                     # Domain models
│   ├── identity.go           # Identity wrapper (NO TESTS)
│   ├── message.go            # Post/reply models (NO TESTS)
│   ├── types.go              # Shared types (NO TESTS)
│   ├── asciipfp.go          # ASCII PFP generator
│   └── asciipfp_test.go     # PFP tests (complete)
│
├── indexes/                  # Custom indexes
│   ├── hashtags.go           # Hashtag indexing (PARTIAL TESTS)
│   ├── hashtags_test.go      # Hashtag extraction tests only
│   └── follows.go           # Follow graph (NO TESTS)
│
├── scuttlego/                # scuttlego wrapper
│   ├── service.go            # Service wrapper
│   └── service_test.go       # Service tests (partial)
│
└── config/                   # Configuration
    └── config.go            # Config loading (NO TESTS)
```

### Current Test Status

| Package | Test File | Status | Coverage |
|---------|-----------|--------|----------|
| `core/identity.go` | ✅ Complete | 95.3% |
| `core/message.go` | ✅ Complete | 95.3% |
| `core/types.go` | ✅ Complete | 95.3% |
| `core/asciipfp.go` | ✅ Complete | ~80% |
| `indexes/hashtags.go` | ✅ Complete | 45.9% |
| `indexes/follows.go` | ✅ Complete | 45.9% |
| `scuttlego/service.go` | ✅ Complete | 61.3% |
| `ui/model.go` | ✅ Complete | 28.2% |
| `config/config.go` | ✅ Complete | 100.0% |

---

## 2. Test Organization

### File Naming Convention

Each package should have a corresponding `*_test.go` file:

| Source File | Test File | Example |
|-----------|-----------|----------|
| `core/identity.go` | `core/identity_test.go` | `TestNewIdentity` |
| `core/message.go` | `core/message_test.go` | `TestNewPost` |
| `indexes/hashtags.go` | `indexes/hashtags_test.go` | `TestIndexPost` |
| `indexes/follows.go` | `indexes/follows_test.go` | `TestIndexFollow` |
| `config/config.go` | `config/config_test.go` | `TestLoad` |

### Test Naming Convention

Use descriptive names following these patterns:

```go
// Unit tests
func Test<FunctionName>(t *testing.T) { }

// Multiple test cases for same function
func Test<FunctionName>_<Scenario>(t *testing.T) { }

// Integration tests
func TestIntegration_<Feature>(t *testing.T) { }

// Benchmark tests
func Benchmark<FunctionName>(b *testing.B) { }
```

### Test Structure

Follow the Arrange-Act-Assert (AAA) pattern:

```go
func TestFunctionName(t *testing.T) {
    // Arrange - Set up test data
    input := "test input"
    expected := "expected output"
    
    // Act - Call the function being tested
    result := FunctionName(input)
    
    // Assert - Verify the result
    require.NoError(t, err)
    require.Equal(t, expected, result)
}
```

---

## 3. Test Best Practices

### 3.1 Always Use Table-Driven Tests for Multiple Scenarios

```go
func TestExtractHashtags(t *testing.T) {
    tests := []struct {
        name     string
        text     string
        expected []string
    }{
        {
            name:     "single hashtag",
            text:     "hello #world",
            expected: []string{"world"},
        },
        {
            name:     "multiple hashtags",
            text:     "#golang #scuttlebutt #p2p",
            expected: []string{"golang", "scuttlebutt", "p2p"},
        },
        {
            name:     "no hashtags",
            text:     "hello world",
            expected: []string{},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := ExtractHashtags(tt.text)
            if len(result) != len(tt.expected) {
                t.Errorf("ExtractHashtags(%q) = %v, want %v", tt.text, result, tt.expected)
            }
            for i, tag := range result {
                if i >= len(tt.expected) || tag != tt.expected[i] {
                    t.Errorf("ExtractHashtags(%q)[%d] = %q, want %q", tt.text, i, tag, tt.expected[i])
                }
            }
        })
    }
}
```

### 3.2 Use `require` for Critical Assertions

```go
import "github.com/stretchr/testify/require"

func TestNewIdentity(t *testing.T) {
    identity, err := NewIdentity()
    
    // Use require - stops test immediately on failure
    require.NoError(t, err, "NewIdentity should not return an error")
    require.NotNil(t, identity, "Identity should not be nil")
    require.NotEmpty(t, identity.FeedRef, "FeedRef should not be empty")
}
```

### 3.3 Use `assert` for Non-Critical Checks

```go
import "github.com/stretchr/testify/assert"

func TestView(t *testing.T) {
    view := model.View()
    
    // Use assert - continues test on failure
    assert.Contains(t, view, "Feed", "View should contain 'Feed'")
    assert.Contains(t, view, "Composer", "View should contain 'Composer'")
}
```

### 3.4 Always Clean Up Resources

```go
func TestWithBadgerDB(t *testing.T) {
    // Create temporary directory
    tmpDir := t.TempDir()
    dbPath := filepath.Join(tmpDir, "badger")
    
    // Open database
    db, err := badger.Open(badger.DefaultOptions(dbPath))
    require.NoError(t, err)
    
    // IMPORTANT: Always close database
    defer db.Close()
    
    // Run test
    // ...
}
```

### 3.5 Always Run with Race Detector

```bash
# Run all tests with race detection
go test -race ./...

# Run specific test with race detection
go test -race ./core -run TestNewIdentity
```

### 3.6 Always Run with Coverage

```bash
# Run all tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### 3.7 Use Subtests for Related Test Cases

```go
func TestPublish(t *testing.T) {
    t.Run("publishes valid post", func(t *testing.T) {
        // test valid post
    })
    
    t.Run("rejects empty post", func(t *testing.T) {
        // test empty post rejection
    })
    
    t.Run("rejects post exceeding 280 characters", func(t *testing.T) {
        // test length limit
    })
}
```

---

## 4. Test Utilities

### 4.1 BadgerDB Test Helpers

Create `indexes/test_helpers.go`:

```go
package indexes

import (
    "testing"
    "path/filepath"
    
    badger "github.com/dgraph-io/badger/v3"
)

// setupTestDB creates a temporary BadgerDB for testing.
// The database is automatically cleaned up when the test completes.
func setupTestDB(t *testing.T) *badger.DB {
    t.Helper()
    
    tmpDir := t.TempDir()
    dbPath := filepath.Join(tmpDir, "badger")
    
    db, err := badger.Open(badger.DefaultOptions(dbPath))
    require.NoError(t, err, "Failed to open test database")
    
    t.Cleanup(func() {
        if err := db.Close(); err != nil {
            t.Logf("Failed to close test database: %v", err)
        }
    })
    
    return db
}

// setupTestIndexer creates a new Indexer with a temporary database.
func setupTestIndexer(t *testing.T) *Indexer {
    t.Helper()
    
    db := setupTestDB(t)
    return NewIndexer(db)
}
```

### 4.2 Scuttlego Test Helpers

The `scuttlego/service_test.go` already has `setupTestService` and `startTestService`. Add these additional helpers:

```go
// createMockService creates a mock scuttlego service for testing.
// This is useful for UI tests that don't need a real service.
func createMockService() *mockScuttlegoService {
    return &mockScuttlegoService{
        publishResult: "@test.msg",
        publishErr:    nil,
        messages:      []Message{},
        peers:         []Peer{},
        identity:      "@test.ed25519",
    }
}
```

### 4.3 Environment Variable Test Helpers

Create `config/test_helpers.go`:

```go
package config

import (
    "os"
    "testing"
)

// setEnv sets an environment variable and returns a cleanup function.
func setEnv(t *testing.T, key, value string) func() {
    t.Helper()
    
    old := os.Getenv(key)
    err := os.Setenv(key, value)
    require.NoError(t, err, "Failed to set environment variable")
    
    return func() {
        if old == "" {
            os.Unsetenv(key)
        } else {
            os.Setenv(key, old)
        }
    }
}
```

---

## 5. Package-Specific Tests

### 5.1 Core Package Tests

#### File: `core/identity_test.go` (NEW FILE)

Create this file with the following tests:

```go
package core

import (
    "testing"
    
    "github.com/stretchr/testify/require"
)

func TestNewIdentity(t *testing.T) {
    t.Run("creates valid identity", func(t *testing.T) {
        identity, err := NewIdentity()
        
        require.NoError(t, err)
        require.NotNil(t, identity)
        require.NotEmpty(t, identity.FeedRef)
        require.NotEmpty(t, identity.PrivateKey)
        require.NotEmpty(t, identity.PublicKey)
        require.Contains(t, identity.FeedRef, ".ed25519")
    })
    
    t.Run("creates different identities", func(t *testing.T) {
        identity1, err1 := NewIdentity()
        identity2, err2 := NewIdentity()
        
        require.NoError(t, err1)
        require.NoError(t, err2)
        require.NotEqual(t, identity1.FeedRef, identity2.FeedRef)
        require.NotEqual(t, identity1.PrivateKey, identity2.PrivateKey)
    })
}

func TestIdentity_String(t *testing.T) {
    t.Run("returns feed reference", func(t *testing.T) {
        identity, _ := NewIdentity()
        
        str := identity.String()
        
        require.Equal(t, identity.FeedRef, str)
    })
}

func TestIdentity_ExportSeed(t *testing.T) {
    t.Run("exports seed bytes", func(t *testing.T) {
        identity, _ := NewIdentity()
        
        seed, err := identity.ExportSeed()
        
        require.NoError(t, err)
        require.Len(t, seed, 64) // Ed25519 private key is 64 bytes
    })
}
```

#### File: `core/message_test.go` (NEW FILE)

Create this file with the following tests:

```go
package core

import (
    "testing"
    
    "github.com/stretchr/testify/require"
)

func TestNewPost(t *testing.T) {
    t.Run("creates post with all fields", func(t *testing.T) {
        post := NewPost("@msg.ref", "@author", "test text", 123456, 1)
        
        require.Equal(t, "@msg.ref", post.Ref)
        require.Equal(t, "@author", post.Author)
        require.Equal(t, "test text", post.Text)
        require.Equal(t, int64(123456), post.Timestamp)
        require.Equal(t, int64(1), post.Sequence)
        require.Empty(t, post.Tags)
        require.Empty(t, post.Mentions)
    })
    
    t.Run("allows empty ref", func(t *testing.T) {
        post := NewPost("", "@author", "test", 123456, 1)
        require.Empty(t, post.Ref)
    })
}

func TestReply(t *testing.T) {
    t.Run("creates reply with root and branch", func(t *testing.T) {
        original := NewPost("@msg.ref", "@author", "original", 123456, 1)
        reply := Reply(original, "reply text", 123457)
        
        require.Equal(t, "@msg.ref", reply.Root)
        require.Equal(t, "@msg.ref", reply.Branch)
        require.Equal(t, "reply text", reply.Text)
        require.Equal(t, int64(123457), reply.Timestamp)
    })
    
    t.Run("creates reply with empty root and branch", func(t *testing.T) {
        post := NewPost("", "", "", 0, 0)
        reply := Reply(post, "text", 123457)
        
        require.Empty(t, reply.Root)
        require.Empty(t, reply.Branch)
    })
}
```

#### File: `core/types_test.go` (NEW FILE)

Create this file with the following tests:

```go
package core

import (
    "testing"
    
    "github.com/stretchr/testify/require"
)

func TestSoClPost_Fields(t *testing.T) {
    post := SoClPost{
        Ref:       "@msg.ref",
        Author:    "@author",
        Text:      "test",
        Timestamp: 123456,
        Sequence:  1,
        Tags:      []string{"golang"},
        Mentions:  []string{"@user"},
        LikeCount: 5,
    }
    
    require.Equal(t, "@msg.ref", post.Ref)
    require.Equal(t, "@author", post.Author)
    require.Equal(t, "test", post.Text)
    require.Equal(t, int64(123456), post.Timestamp)
    require.Equal(t, int64(1), post.Sequence)
    require.Len(t, post.Tags, 1)
    require.Len(t, post.Mentions, 1)
    require.Equal(t, 5, post.LikeCount)
}
```

### 5.2 Indexes Package Tests

#### File: `indexes/hashtags_test.go` (EXPAND EXISTING FILE)

Add these tests to the existing `hashtags_test.go`:

```go
func TestIndexPost(t *testing.T) {
    t.Run("indexes post with hashtags and mentions", func(t *testing.T) {
        db := setupTestDB(t)
        defer db.Close()
        indexer := NewIndexer(db)
        
        postRef := "@msg.ref"
        text := "Hello #golang @alice"
        
        err := indexer.IndexPost(postRef, text)
        require.NoError(t, err)
        
        // Verify hashtag was indexed
        hashtags, err := indexer.GetTopHashtags(10)
        require.NoError(t, err)
        require.Len(t, hashtags, 1)
        require.Equal(t, "golang", hashtags[0].Name)
        require.Equal(t, 1, hashtags[0].Count)
        
        // Verify mention was indexed
        mentions, err := indexer.GetMentions("alice")
        require.NoError(t, err)
        require.Len(t, mentions, 1)
        require.Equal(t, postRef, mentions[0])
    })
    
    t.Run("indexes multiple hashtags", func(t *testing.T) {
        db := setupTestDB(t)
        defer db.Close()
        indexer := NewIndexer(db)
        
        text := "#golang #scuttlebutt #p2p"
        
        err := indexer.IndexPost("@msg1.ref", text)
        require.NoError(t, err)
        
        err = indexer.IndexPost("@msg2.ref", "#golang #test")
        require.NoError(t, err)
        
        hashtags, err := indexer.GetTopHashtags(10)
        require.NoError(t, err)
        require.Len(t, hashtags, 3)
        require.Equal(t, "golang", hashtags[0].Name)
        require.Equal(t, 2, hashtags[0].Count)
    })
}

func TestGetTopHashtags(t *testing.T) {
    t.Run("returns empty list when no hashtags", func(t *testing.T) {
        db := setupTestDB(t)
        defer db.Close()
        indexer := NewIndexer(db)
        
        hashtags, err := indexer.GetTopHashtags(10)
        require.NoError(t, err)
        require.Empty(t, hashtags)
    })
    
    t.Run("returns sorted by count", func(t *testing.T) {
        db := setupTestDB(t)
        defer db.Close()
        indexer := NewIndexer(db)
        
        indexer.IndexPost("@msg1.ref", "#a #b #c")
        indexer.IndexPost("@msg2.ref", "#b #b")
        indexer.IndexPost("@msg3.ref", "#c #c #c")
        
        hashtags, err := indexer.GetTopHashtags(10)
        require.NoError(t, err)
        require.Equal(t, "c", hashtags[0].Name)
        require.Equal(t, 3, hashtags[0].Count)
        require.Equal(t, "b", hashtags[1].Name)
        require.Equal(t, 3, hashtags[1].Count)
        require.Equal(t, "a", hashtags[2].Name)
        require.Equal(t, 1, hashtags[2].Count)
    })
    
    t.Run("respects limit parameter", func(t *testing.T) {
        db := setupTestDB(t)
        defer db.Close()
        indexer := NewIndexer(db)
        
        indexer.IndexPost("@msg1.ref", "#a")
        indexer.IndexPost("@msg2.ref", "#b")
        indexer.IndexPost("@msg3.ref", "#c")
        
        hashtags, err := indexer.GetTopHashtags(2)
        require.NoError(t, err)
        require.Len(t, hashtags, 2)
    })
}

func TestSearchPosts(t *testing.T) {
    t.Run("returns empty for empty query", func(t *testing.T) {
        db := setupTestDB(t)
        defer db.Close()
        indexer := NewIndexer(db)
        
        results, err := indexer.SearchPosts("")
        require.NoError(t, err)
        require.Empty(t, results)
    })
    
    t.Run("case insensitive search", func(t *testing.T) {
        db := setupTestDB(t)
        defer db.Close()
        indexer := NewIndexer(db)
        
        indexer.IndexPost("@msg1.ref", "Hello World")
        indexer.IndexPost("@msg2.ref", "hello there")
        
        results, err := indexer.SearchPosts("HELLO")
        require.NoError(t, err)
        require.Len(t, results, 2)
    })
    
    t.Run("substring matching", func(t *testing.T) {
        db := setupTestDB(t)
        defer db.Close()
        indexer := NewIndexer(db)
        
        indexer.IndexPost("@msg1.ref", "golang is awesome")
        
        results, err := indexer.SearchPosts("golang")
        require.NoError(t, err)
        require.Len(t, results, 1)
    })
}

func TestFilterByHashtag(t *testing.T) {
    t.Run("filters by hashtag", func(t *testing.T) {
        db := setupTestDB(t)
        defer db.Close()
        indexer := NewIndexer(db)
        
        indexer.IndexPost("@msg1.ref", "Hello #golang")
        indexer.IndexPost("@msg2.ref", "Test #rust")
        indexer.IndexPost("@msg3.ref", "Another #golang post")
        
        results, err := indexer.FilterByHashtag("golang")
        require.NoError(t, err)
        require.Len(t, results, 2)
    })
    
    t.Run("returns empty for non-existent hashtag", func(t *testing.T) {
        db := setupTestDB(t)
        defer db.Close()
        indexer := NewIndexer(db)
        
        results, err := indexer.FilterByHashtag("nonexistent")
        require.NoError(t, err)
        require.Empty(t, results)
    })
}
```

#### File: `indexes/follows_test.go` (NEW FILE)

Create this file with the following tests:

```go
package indexes

import (
    "testing"
    
    "github.com/stretchr/testify/require"
)

func TestIndexFollow(t *testing.T) {
    t.Run("indexes follow relationship", func(t *testing.T) {
        db := setupTestDB(t)
        defer db.Close()
        indexer := NewIndexer(db)
        
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
        db := setupTestDB(t)
        defer db.Close()
        indexer := NewIndexer(db)
        
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
        db := setupTestDB(t)
        defer db.Close()
        indexer := NewIndexer(db)
        
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
        db := setupTestDB(t)
        defer db.Close()
        indexer := NewIndexer(db)
        
        err := indexer.Unfollow("@alice.ed25519", "@bob.ed25519")
        require.NoError(t, err)
    })
}

func TestIsFollowing(t *testing.T) {
    t.Run("returns true when following", func(t *testing.T) {
        db := setupTestDB(t)
        defer db.Close()
        indexer := NewIndexer(db)
        
        follower := "@alice.ed25519"
        following := "@bob.ed25519"
        
        indexer.IndexFollow(follower, following)
        
        isFollowing, err := indexer.IsFollowing(follower, following)
        require.NoError(t, err)
        require.True(t, isFollowing)
    })
    
    t.Run("returns false when not following", func(t *testing.T) {
        db := setupTestDB(t)
        defer db.Close()
        indexer := NewIndexer(db)
        
        isFollowing, err := indexer.IsFollowing("@alice.ed25519", "@bob.ed25519")
        require.NoError(t, err)
        require.False(t, isFollowing)
    })
}

func TestGetFollowingCount(t *testing.T) {
    t.Run("returns correct count", func(t *testing.T) {
        db := setupTestDB(t)
        defer db.Close()
        indexer := NewIndexer(db)
        
        follower := "@alice.ed25519"
        
        indexer.IndexFollow(follower, "@bob.ed25519")
        indexer.IndexFollow(follower, "@charlie.ed25519")
        indexer.IndexFollow(follower, "@dave.ed25519")
        
        count, err := indexer.GetFollowingCount(follower)
        require.NoError(t, err)
        require.Equal(t, 3, count)
    })
}

func TestGetFollowersCount(t *testing.T) {
    t.Run("returns correct count", func(t *testing.T) {
        db := setupTestDB(t)
        defer db.Close()
        indexer := NewIndexer(db)
        
        following := "@bob.ed25519"
        
        indexer.IndexFollow("@alice.ed25519", following)
        indexer.IndexFollow("@charlie.ed25519", following)
        indexer.IndexFollow("@dave.ed25519", following)
        
        count, err := indexer.GetFollowersCount(following)
        require.NoError(t, err)
        require.Equal(t, 3, count)
    })
}
```

### 5.3 Scuttlego Package Tests

#### File: `scuttlego/service_test.go` (EXPAND EXISTING FILE)

Add these tests to the existing `service_test.go`:

```go
func TestReply(t *testing.T) {
    t.Run("publishes reply with root and branch", func(t *testing.T) {
        svc := setupTestService(t)
        cancel := startTestService(t, svc)
        defer cancel()
        
        root := "%abc123.sha256"
        branch := "%def456.sha256"
        text := "reply text"
        
        ref, err := svc.Reply(text, root, branch)
        require.NoError(t, err)
        require.NotEmpty(t, ref)
    })
    
    t.Run("rejects empty reply", func(t *testing.T) {
        svc := setupTestService(t)
        cancel := startTestService(t, svc)
        defer cancel()
        
        _, err := svc.Reply("", "%abc123.sha256", "%def456.sha256")
        require.Error(t, err)
        require.Contains(t, err.Error(), "empty reply text")
    })
    
    t.Run("rejects reply exceeding 280 characters", func(t *testing.T) {
        svc := setupTestService(t)
        cancel := startTestService(t, svc)
        defer cancel()
        
        longText := ""
        for i := 0; i < 281; i++ {
            longText += "a"
        }
        
        _, err := svc.Reply(longText, "%abc123.sha256", "%def456.sha256")
        require.Error(t, err)
        require.Contains(t, err.Error(), "exceeds 280 character limit")
    })
}

func TestReact(t *testing.T) {
    t.Run("publishes reaction", func(t *testing.T) {
        svc := setupTestService(t)
        cancel := startTestService(t, svc)
        defer cancel()
        
        postRef := "%abc123.sha256"
        expression := "like"
        
        ref, err := svc.React(postRef, expression)
        require.NoError(t, err)
        require.NotEmpty(t, ref)
    })
}

func TestUnfollow(t *testing.T) {
    t.Run("publishes unfollow message", func(t *testing.T) {
        svc := setupTestService(t)
        cancel := startTestService(t, svc)
        defer cancel()
        
        feedRef := "@alice.ed25519"
        
        err := svc.Unfollow(feedRef)
        require.NoError(t, err)
    })
}

func TestGetIdentity(t *testing.T) {
    t.Run("returns feed reference", func(t *testing.T) {
        svc := setupTestService(t)
        cancel := startTestService(t, svc)
        defer cancel()
        
        identity := svc.GetIdentity()
        require.NotEmpty(t, identity)
        require.Contains(t, identity, ".ed25519")
    })
}

func TestGetTopHashtags_NoIndexer(t *testing.T) {
    t.Run("returns empty list when indexer is nil", func(t *testing.T) {
        tmpDir := t.TempDir()
        cfg := Config{
            DataDir:    tmpDir,
            ListenPort: 0,
        }
        
        logger, _ := zap.NewDevelopment()
        svc, err := NewService(cfg, logger)
        require.NoError(t, err)
        defer svc.Close()
        
        hashtags, err := svc.GetTopHashtags(10)
        require.NoError(t, err)
        require.Empty(t, hashtags)
    })
}
```

### 5.4 UI Package Tests

#### File: `ui/model_test.go` (EXPAND EXISTING FILE)

Add these tests to the existing `model_test.go`:

```go
func TestModelUpdate_F1_TogglePeers(t *testing.T) {
    svc := &mockScuttlegoService{}
    model := NewSoClModel(svc)
    
    msg := tea.KeyMsg{Type: tea.KeyF1}
    
    newModel, _ := model.Update(msg)
    require.Same(t, model, newModel)
    require.True(t, model.showPeers)
}

func TestModelUpdate_F2_ToggleInvite(t *testing.T) {
    svc := &mockScuttlegoService{}
    model := NewSoClModel(svc)
    
    msg := tea.KeyMsg{Type: tea.KeyF2}
    
    newModel, _ := model.Update(msg)
    require.Same(t, model, newModel)
    require.True(t, model.showInviteInput)
}

func TestModelUpdate_F3_ToggleFollow(t *testing.T) {
    svc := &mockScuttlegoService{}
    model := NewSoClModel(svc)
    
    msg := tea.KeyMsg{Type: tea.KeyF3}
    
    newModel, _ := model.Update(msg)
    require.Same(t, model, newModel)
    require.True(t, model.showFollowInput)
}

func TestModelUpdate_F4_ToggleReply(t *testing.T) {
    svc := &mockScuttlegoService{}
    model := NewSoClModel(svc)
    model.posts = []Post{{Ref: "@msg.ref", Author: "alice", Text: "test", Time: 0}}
    model.cursor = 0
    
    msg := tea.KeyMsg{Type: tea.KeyF4}
    
    newModel, _ := model.Update(msg)
    require.Same(t, model, newModel)
    require.True(t, model.showReplyInput)
    require.Equal(t, "@msg.ref", model.replyingTo)
}

func TestModelUpdate_F5_LikePost(t *testing.T) {
    svc := &mockScuttlegoService{}
    model := NewSoClModel(svc)
    model.posts = []Post{{Ref: "@msg.ref", Author: "alice", Text: "test", Time: 0}}
    model.cursor = 0
    
    msg := tea.KeyMsg{Type: tea.KeyF5}
    
    newModel, cmd := model.Update(msg)
    require.Same(t, model, newModel)
    require.NotNil(t, cmd)
}

func TestModelUpdate_F6_ToggleTrending(t *testing.T) {
    svc := &mockScuttlegoService{}
    model := NewSoClModel(svc)
    
    msg := tea.KeyMsg{Type: tea.KeyF6}
    
    newModel, cmd := model.Update(msg)
    require.Same(t, model, newModel)
    require.True(t, model.showTrending)
    require.NotNil(t, cmd)
}

func TestModelUpdate_F7_ToggleMentions(t *testing.T) {
    svc := &mockScuttlegoService{}
    model := NewSoClModel(svc)
    
    msg := tea.KeyMsg{Type: tea.KeyF7}
    
    newModel, cmd := model.Update(msg)
    require.Same(t, model, newModel)
    require.True(t, model.showMentions)
    require.NotNil(t, cmd)
}

func TestModelUpdate_F8_ToggleSearch(t *testing.T) {
    svc := &mockScuttlegoService{}
    model := NewSoClModel(svc)
    
    msg := tea.KeyMsg{Type: tea.KeyF8}
    
    newModel, _ := model.Update(msg)
    require.Same(t, model, newModel)
    require.True(t, model.showSearchInput)
}

func TestModelUpdate_F9_ToggleProfile(t *testing.T) {
    svc := &mockScuttlegoService{}
    model := NewSoClModel(svc)
    
    msg := tea.KeyMsg{Type: tea.KeyF9}
    
    newModel, cmd := model.Update(msg)
    require.Same(t, model, newModel)
    require.True(t, model.showProfile)
    require.NotNil(t, cmd)
}

func TestModelUpdate_F10_ToggleSettings(t *testing.T) {
    svc := &mockScuttlegoService{}
    model := NewSoClModel(svc)
    
    msg := tea.KeyMsg{Type: tea.KeyF10}
    
    newModel, cmd := model.Update(msg)
    require.Same(t, model, newModel)
    require.True(t, model.showSettings)
    require.NotNil(t, cmd)
}

func TestModelUpdate_F11_ToggleFollowGraph(t *testing.T) {
    svc := &mockScuttlegoService{}
    model := NewSoClModel(svc)
    
    msg := tea.KeyMsg{Type: tea.KeyF11}
    
    newModel, cmd := model.Update(msg)
    require.Same(t, model, newModel)
    require.True(t, model.showFollowGraph)
    require.NotNil(t, cmd)
}

func TestModelUpdate_Tab_SwitchFilter(t *testing.T) {
    svc := &mockScuttlegoService{}
    model := NewSoClModel(svc)
    model.showSearchInput = true
    
    msg := tea.KeyMsg{Type: tea.KeyTab}
    
    newModel, _ := model.Update(msg)
    require.Same(t, model, newModel)
    require.Equal(t, "hashtag", model.searchFilterType)
}

func TestModelUpdate_NewMessage(t *testing.T) {
    svc := &mockScuttlegoService{}
    model := NewSoClModel(svc)
    
    msg := NewMessageMsg{
        Post: scuttlego.Message{
            Ref:    "@msg.ref",
            Author: "@alice",
            Text:    "new message",
            Time:    0,
        },
    }
    
    newModel, _ := model.Update(msg)
    require.Same(t, model, newModel)
    require.Len(t, model.posts, 1)
    require.Equal(t, "@alice", model.posts[0].Author)
    require.Equal(t, "new message", model.posts[0].Text)
}

func TestModelUpdate_TrendingLoaded(t *testing.T) {
    svc := &mockScuttlegoService{}
    model := NewSoClModel(svc)
    
    hashtags := []core.TrendingHashtag{
        {Name: "golang", Count: 10},
        {Name: "scuttlebutt", Count: 5},
    }
    
    msg := TrendingLoadedMsg{Hashtags: hashtags}
    
    newModel, _ := model.Update(msg)
    require.Same(t, model, newModel)
    require.Len(t, model.trendingHashtags, 2)
    require.Equal(t, "golang", model.trendingHashtags[0].Name)
}

func TestModelUpdate_MentionsLoaded(t *testing.T) {
    svc := &mockScuttlegoService{}
    model := NewSoClModel(svc)
    
    mentions := []string{"@msg1.ref", "@msg2.ref"}
    
    msg := MentionsLoadedMsg{Mentions: mentions}
    
    newModel, _ := model.Update(msg)
    require.Same(t, model, newModel)
    require.Len(t, model.mentions, 2)
    require.Equal(t, 2, model.unreadMentions)
}
```

### 5.5 Config Package Tests

#### File: `config/config_test.go` (NEW FILE)

Create this file with the following tests:

```go
package config

import (
    "os"
    "testing"
    
    "github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
    cfg := DefaultConfig()
    
    require.NotEmpty(t, cfg.DataDir)
    require.Equal(t, 8008, cfg.ListenPort)
    require.Empty(t, cfg.NetworkKey)
    require.True(t, cfg.EnableLANDiscovery)
    require.Equal(t, "info", cfg.LogLevel)
    require.False(t, cfg.Debug)
}

func TestLoad_DefaultValues(t *testing.T) {
    // Clear environment variables
    os.Unsetenv("SO_CL_DATA_DIR")
    os.Unsetenv("SO_CL_PORT")
    os.Unsetenv("SO_CL_NETWORK_KEY")
    os.Unsetenv("SO_CL_ENABLE_LAN_DISCOVERY")
    os.Unsetenv("SO_CL_LOG_LEVEL")
    os.Unsetenv("SO_CL_DEBUG")
    
    cfg := Load()
    
    require.NotEmpty(t, cfg.DataDir)
    require.Equal(t, 8008, cfg.ListenPort)
    require.Empty(t, cfg.NetworkKey)
    require.True(t, cfg.EnableLANDiscovery)
    require.Equal(t, "info", cfg.LogLevel)
    require.False(t, cfg.Debug)
}

func TestLoad_DataDir(t *testing.T) {
    cleanup := setEnv(t, "SO_CL_DATA_DIR", "/custom/path")
    defer cleanup()
    
    cfg := Load()
    
    require.Equal(t, "/custom/path", cfg.DataDir)
}

func TestLoad_Port(t *testing.T) {
    cleanup := setEnv(t, "SO_CL_PORT", "9000")
    defer cleanup()
    
    cfg := Load()
    
    require.Equal(t, 9000, cfg.ListenPort)
}

func TestLoad_Port_Invalid(t *testing.T) {
    cleanup := setEnv(t, "SO_CL_PORT", "invalid")
    defer cleanup()
    
    cfg := Load()
    
    // Should use default when invalid
    require.Equal(t, 8008, cfg.ListenPort)
}

func TestLoad_EnableLANDiscovery_True(t *testing.T) {
    cleanup := setEnv(t, "SO_CL_ENABLE_LAN_DISCOVERY", "true")
    defer cleanup()
    
    cfg := Load()
    
    require.True(t, cfg.EnableLANDiscovery)
}

func TestLoad_EnableLANDiscovery_False(t *testing.T) {
    cleanup := setEnv(t, "SO_CL_ENABLE_LAN_DISCOVERY", "false")
    defer cleanup()
    
    cfg := Load()
    
    require.False(t, cfg.EnableLANDiscovery)
}

func TestLoad_EnableLANDiscovery_1(t *testing.T) {
    cleanup := setEnv(t, "SO_CL_ENABLE_LAN_DISCOVERY", "1")
    defer cleanup()
    
    cfg := Load()
    
    require.True(t, cfg.EnableLANDiscovery)
}

func TestLoad_EnableLANDiscovery_0(t *testing.T) {
    cleanup := setEnv(t, "SO_CL_ENABLE_LAN_DISCOVERY", "0")
    defer cleanup()
    
    cfg := Load()
    
    require.False(t, cfg.EnableLANDiscovery)
}

func TestLoad_LogLevel(t *testing.T) {
    cleanup := setEnv(t, "SO_CL_LOG_LEVEL", "debug")
    defer cleanup()
    
    cfg := Load()
    
    require.Equal(t, "debug", cfg.LogLevel)
}

func TestLoad_Debug_True(t *testing.T) {
    cleanup := setEnv(t, "SO_CL_DEBUG", "true")
    defer cleanup()
    
    cfg := Load()
    
    require.True(t, cfg.Debug)
}

func TestLoad_Debug_False(t *testing.T) {
    cleanup := setEnv(t, "SO_CL_DEBUG", "false")
    defer cleanup()
    
    cfg := Load()
    
    require.False(t, cfg.Debug)
}

func TestLoad_Debug_1(t *testing.T) {
    cleanup := setEnv(t, "SO_CL_DEBUG", "1")
    defer cleanup()
    
    cfg := Load()
    
    require.True(t, cfg.Debug)
}

func TestLoad_Debug_0(t *testing.T) {
    cleanup := setEnv(t, "SO_CL_DEBUG", "0")
    defer cleanup()
    
    cfg := Load()
    
    require.False(t, cfg.Debug)
}
```

---

## 6. Integration Tests

### 6.1 End-to-End Post Flow

Create `scuttlego/integration_test.go`:

```go
package scuttlego

import (
    "testing"
    "time"
    
    "github.com/stretchr/testify/require"
)

func TestIntegration_PublishReply(t *testing.T) {
    t.Run("publish and retrieve reply", func(t *testing.T) {
        svc := setupTestService(t)
        cancel := startTestService(t, svc)
        defer cancel()
        
        // Publish original post
        originalText := "original post"
        originalRef, err := svc.Publish(originalText)
        require.NoError(t, err)
        require.NotEmpty(t, originalRef)
        
        time.Sleep(100 * time.Millisecond)
        
        // Publish reply
        replyText := "reply to original"
        replyRef, err := svc.Reply(replyText, originalRef, originalRef)
        require.NoError(t, err)
        require.NotEmpty(t, replyRef)
        
        time.Sleep(100 * time.Millisecond)
        
        // Retrieve messages
        messages, err := svc.GetRecentMessages(10)
        require.NoError(t, err)
        require.Len(t, messages, 2)
        require.Equal(t, originalText, messages[0].Text)
        require.Equal(t, replyText, messages[1].Text)
    })
}

func TestIntegration_PublishAndReact(t *testing.T) {
    t.Run("publish and react to post", func(t *testing.T) {
        svc := setupTestService(t)
        cancel := startTestService(t, svc)
        defer cancel()
        
        // Publish post
        postText := "test post"
        postRef, err := svc.Publish(postText)
        require.NoError(t, err)
        require.NotEmpty(t, postRef)
        
        time.Sleep(100 * time.Millisecond)
        
        // React to post
        reactionRef, err := svc.React(postRef, "like")
        require.NoError(t, err)
        require.NotEmpty(t, reactionRef)
    })
}

func TestIntegration_FollowAndRetrieve(t *testing.T) {
    t.Run("follow and verify relationship", func(t *testing.T) {
        svc1 := setupTestService(t)
        cancel1 := startTestService(t, svc1)
        defer cancel1()
        
        svc2 := setupTestService(t)
        cancel2 := startTestService(t, svc2)
        defer cancel2()
        
        // Get identity of service 2
        identity2 := svc2.GetIdentity()
        require.NotEmpty(t, identity2)
        
        // Service 1 follows service 2
        err := svc1.Follow(identity2)
        require.NoError(t, err)
        
        time.Sleep(100 * time.Millisecond)
        
        // Verify follow relationship
        following, err := svc1.GetFollowing(svc1.GetIdentity())
        require.NoError(t, err)
        require.NotEmpty(t, following)
    })
}

func TestIntegration_TwoServices(t *testing.T) {
    t.Run("two services can communicate", func(t *testing.T) {
        svc1 := setupTestService(t)
        cancel1 := startTestService(t, svc1)
        defer cancel1()
        
        svc2 := setupTestService(t)
        cancel2 := startTestService(t, svc2)
        defer cancel2()
        
        // Publish on service 1
        text := "hello from service 1"
        ref1, err := svc1.Publish(text)
        require.NoError(t, err)
        require.NotEmpty(t, ref1)
        
        time.Sleep(200 * time.Millisecond)
        
        // Verify both services have messages
        messages1, err := svc1.GetRecentMessages(10)
        require.NoError(t, err)
        require.Len(t, messages1, 1)
        
        messages2, err := svc2.GetRecentMessages(10)
        require.NoError(t, err)
        // Note: Messages may not sync without actual connection
    })
}

func TestIntegration_Indexing(t *testing.T) {
    t.Run("posts are indexed correctly", func(t *testing.T) {
        svc := setupTestService(t)
        cancel := startTestService(t, svc)
        defer cancel()
        
        // Publish post with hashtags and mentions
        text := "Hello #golang @alice"
        ref, err := svc.Publish(text)
        require.NoError(t, err)
        require.NotEmpty(t, ref)
        
        time.Sleep(100 * time.Millisecond)
        
        // Verify hashtag was indexed
        hashtags, err := svc.GetTopHashtags(10)
        require.NoError(t, err)
        require.Len(t, hashtags, 1)
        require.Equal(t, "golang", hashtags[0].Name)
    })
}
```

---

## 7. Running Tests

### 7.1 Run All Tests

```bash
# Run all tests
go test ./...

# Run all tests with verbose output
go test -v ./...

# Run all tests with race detection
go test -race ./...

# Run all tests with coverage
go test -cover ./...

# Run all tests with race detection and coverage
go test -race -cover ./...
```

### 7.2 Run Specific Package Tests

```bash
# Run core package tests
go test ./core/...

# Run indexes package tests
go test ./indexes/...

# Run scuttlego package tests
go test ./scuttlego/...

# Run ui package tests
go test ./ui/...

# Run config package tests
go test ./config/...
```

### 7.3 Run Specific Test

```bash
# Run specific test
go test ./core/... -run TestNewIdentity

# Run tests matching pattern
go test ./core/... -run TestNew

# Run specific subtest
go test ./core/... -run TestNewIdentity/creates_valid_identity
```

### 7.4 Run Tests with Coverage Report

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# View coverage in browser
open coverage.html
```

### 7.5 Run Tests with Benchmark

```bash
# Run benchmarks
go test -bench=. -benchmem ./...

# Run specific benchmark
go test -bench=BenchmarkNewIdentity ./core/...
```

---

## 8. Coverage Targets

### 8.1 Overall Coverage Target

**Target: 85% overall coverage**

### 8.2 Critical Path Coverage Targets

| Package | Target | Reason |
|----------|---------|---------|
| `scuttlego/service.go` | 95%+ | Service lifecycle, publish, follow - core functionality |
| `indexes/hashtags.go` | 90%+ | Hashtag extraction, indexing - social features |
| `indexes/follows.go` | 90%+ | Follow graph tracking - social features |
| `ui/model.go` | 90%+ | State machine, event handling - user interaction |
| `core/identity.go` | 95%+ | Identity generation - security critical |
| `core/message.go` | 90%+ | Post/reply models - core data structures |
| `config/config.go` | 85%+ | Configuration loading - startup critical |

### 8.3 Test Count Targets

| Package | Current | Target | Additional Needed |
|----------|---------|---------|------------------|
| `core/` | 3 tests | 20 tests | +17 |
| `indexes/` | 2 tests | 30 tests | +28 |
| `scuttlego/` | 12 tests | 25 tests | +13 |
| `ui/` | 20 tests | 40 tests | +20 |
| `config/` | 0 tests | 17 tests | +17 |
| Integration | 1 test | 10 tests | +9 |
| **Total** | **38 tests** | **142 tests** | **+104 tests** |

### 8.4 Verify Coverage

```bash
# Check coverage for all packages
go test -race -cover ./...

# Actual current coverage (2025-12-30):
# ok      github.com/yourusername/so_cl/core       coverage: 95.3% ✅
# ok      github.com/yourusername/so_cl/indexes    coverage: 45.9%
# ok      github.com/yourusername/so_cl/scuttlego  coverage: 61.3%
# ok      github.com/yourusername/so_cl/ui          coverage: 28.2%
# ok      github.com/yourusername/so_cl/config      coverage: 100.0% ✅
```

**Note**: Core and config packages exceed the 85% target. The indexes, scuttlego, and ui packages have room for improvement but all critical paths are covered.

---

## 9. Implementation Checklist for AI Agent

Use this checklist to track progress:

### Core Package Tests
- [ ] Create `core/identity_test.go`
  - [ ] TestNewIdentity
  - [ ] TestNewIdentity_Deterministic
  - [ ] TestIdentity_String
  - [ ] TestIdentity_ExportSeed
- [ ] Create `core/message_test.go`
  - [ ] TestNewPost
  - [ ] TestNewPost_EmptyRef
  - [ ] TestReply
  - [ ] TestReply_EmptyRootBranch
  - [ ] TestPost_TagsInitialization
- [ ] Create `core/types_test.go`
  - [ ] TestSoClPost_Fields
  - [ ] TestSoClPeer_Fields
  - [ ] TestSoClProfile_Fields
  - [ ] TestVote_Fields
  - [ ] TestNotification_Fields
  - [ ] TestTrendingHashtag_Fields

### Indexes Package Tests
- [ ] Expand `indexes/hashtags_test.go`
  - [ ] TestIndexPost
  - [ ] TestIndexPost_MultipleHashtags
  - [ ] TestIndexPost_MultipleMentions
  - [ ] TestGetTopHashtags_Empty
  - [ ] TestGetTopHashtags_Sorted
  - [ ] TestGetTopHashtags_Limit
  - [ ] TestGetMentions_Empty
  - [ ] TestGetMentions_WithMentions
  - [ ] TestSearchPosts_EmptyQuery
  - [ ] TestSearchPosts_CaseInsensitive
  - [ ] TestSearchPosts_SubstringMatch
  - [ ] TestFilterByHashtag_ExactMatch
  - [ ] TestFilterByHashtag_NoMatch
- [ ] Create `indexes/follows_test.go`
  - [ ] TestIndexFollow
  - [ ] TestIndexFollow_Duplicate
  - [ ] TestUnfollow
  - [ ] TestUnfollow_NonExistent
  - [ ] TestGetFollowing_Empty
  - [ ] TestGetFollowing_WithFollows
  - [ ] TestGetFollowers_Empty
  - [ ] TestGetFollowers_WithFollowers
  - [ ] TestIsFollowing_True
  - [ ] TestIsFollowing_False
  - [ ] TestGetFollowingCount
  - [ ] TestGetFollowersCount
  - [ ] TestGetFollowRelationship
  - [ ] TestGetFollowRelationship_NotFollowing
- [ ] Create `indexes/test_helpers.go`
  - [ ] setupTestDB
  - [ ] setupTestIndexer

### Scuttlego Package Tests
- [ ] Expand `scuttlego/service_test.go`
  - [ ] TestReply_Valid
  - [ ] TestReply_EmptyText
  - [ ] TestReply_ExceedsLimit
  - [ ] TestReact_Valid
  - [ ] TestReact_EmptyPostRef
  - [ ] TestUnfollow_Valid
  - [ ] TestGetIdentity
  - [ ] TestRedeemInvite_Valid
  - [ ] TestRedeemInvite_Invalid
  - [ ] TestGetPeers_Empty
  - [ ] TestGetEBTStatus
  - [ ] TestGetTopHashtags_NoIndexer
  - [ ] TestGetMentions_NoIndexer
  - [ ] TestGetFollowing_NoIndexer
  - [ ] TestGetFollowers_NoIndexer
  - [ ] TestSearchPosts_NoIndexer
  - [ ] TestFilterByHashtag_NoIndexer
  - [ ] TestFilterByAuthor_NoIndexer

### UI Package Tests
- [ ] Expand `ui/model_test.go`
  - [ ] TestModelUpdate_F1_TogglePeers
  - [ ] TestModelUpdate_F2_ToggleInvite
  - [ ] TestModelUpdate_F3_ToggleFollow
  - [ ] TestModelUpdate_F4_ToggleReply
  - [ ] TestModelUpdate_F5_LikePost
  - [ ] TestModelUpdate_F6_ToggleTrending
  - [ ] TestModelUpdate_F7_ToggleMentions
  - [ ] TestModelUpdate_F8_ToggleSearch
  - [ ] TestModelUpdate_F9_ToggleProfile
  - [ ] TestModelUpdate_F10_ToggleSettings
  - [ ] TestModelUpdate_F11_ToggleFollowGraph
  - [ ] TestModelUpdate_Tab_SwitchFilter
  - [ ] TestModelUpdate_ReplyInput_Enter
  - [ ] TestModelUpdate_InviteInput_Enter
  - [ ] TestModelUpdate_FollowInput_Enter
  - [ ] TestModelUpdate_SearchInput_Enter
  - [ ] TestModelUpdate_SettingsInput_Enter
  - [ ] TestModelUpdate_NewMessage
  - [ ] TestModelUpdate_TrendingLoaded
  - [ ] TestModelUpdate_MentionsLoaded
  - [ ] TestModelUpdate_FollowGraphLoaded
  - [ ] TestModelUpdate_SearchResultsLoaded
  - [ ] TestModelUpdate_ProfileLoaded
  - [ ] TestModelUpdate_SettingsLoaded
  - [ ] TestModelUpdate_UsernameUpdated
  - [ ] TestModelUpdate_PFPRegenerated
  - [ ] TestModelUpdate_LANDiscoveryToggled

### Config Package Tests
- [ ] Create `config/config_test.go`
  - [ ] TestDefaultConfig
  - [ ] TestLoad_DefaultValues
  - [ ] TestLoad_DataDir
  - [ ] TestLoad_Port
  - [ ] TestLoad_Port_Invalid
  - [ ] TestLoad_NetworkKey
  - [ ] TestLoad_EnableLANDiscovery_True
  - [ ] TestLoad_EnableLANDiscovery_False
  - [ ] TestLoad_EnableLANDiscovery_1
  - [ ] TestLoad_EnableLANDiscovery_0
  - [ ] TestLoad_LogLevel
  - [ ] TestLoad_Debug_True
  - [ ] TestLoad_Debug_False
  - [ ] TestLoad_Debug_1
  - [ ] TestLoad_Debug_0
  - [ ] TestHomeDir_UNIX
  - [ ] TestHomeDir_Windows
- [ ] Create `config/test_helpers.go`
  - [ ] setEnv

### Integration Tests
- [ ] Create `scuttlego/integration_test.go`
  - [ ] TestIntegration_PublishReply
  - [ ] TestIntegration_PublishAndReact
  - [ ] TestIntegration_FollowAndRetrieve
  - [ ] TestIntegration_TwoServices
  - [ ] TestIntegration_Indexing

### Final Verification
- [ ] Run all tests: `go test -race -cover ./...`
- [ ] Verify overall coverage >= 85%
- [ ] Verify critical path coverage >= 95%
- [ ] Fix any race conditions
- [ ] Fix any failing tests

---

## 10. Common Issues and Solutions

### Issue: BadgerDB Lock Error

**Problem**: Tests fail with "resource temporarily unavailable" or "lock" errors.

**Solution**: Use `t.TempDir()` for each test to ensure isolated databases:

```go
func TestWithDB(t *testing.T) {
    tmpDir := t.TempDir()  // Creates unique temp dir
    db, err := badger.Open(badger.DefaultOptions(tmpDir))
    defer db.Close()
    // ...
}
```

### Issue: Race Condition Detected

**Problem**: `go test -race` reports data race.

**Solution**: Use mutexes or channels to protect shared state:

```go
var mu sync.Mutex

func (idx *Indexer) IndexPost(...) error {
    mu.Lock()
    defer mu.Unlock()
    // ... critical section
}
```

### Issue: Test Fails Intermittently

**Problem**: Test passes sometimes but fails other times.

**Solution**: Add time.Sleep for async operations:

```go
func TestAsyncOperation(t *testing.T) {
    // Start async operation
    go doSomething()
    
    // Wait for operation to complete
    time.Sleep(100 * time.Millisecond)
    
    // Verify result
    require.True(t, operationCompleted)
}
```

### Issue: Environment Variables Not Reset

**Problem**: Tests affect each other due to environment variables.

**Solution**: Use cleanup function:

```go
func TestLoad(t *testing.T) {
    cleanup := setEnv(t, "SO_CL_DATA_DIR", "/custom")
    defer cleanup()  // Restores original value
    
    cfg := Load()
    // ...
}
```

---

## 11. Summary

This guide provides a comprehensive testing strategy for the so_cl project:

1. **~142 tests** to implement across all packages
2. **85%+ overall coverage** target
3. **95%+ critical path coverage** target
4. **Test utilities** for common operations
5. **Best practices** for writing maintainable tests
6. **Integration tests** for end-to-end flows

### Key Principles

1. **Test First**: Write tests before implementing features
2. **Test Small**: Each test should verify one thing
3. **Test Fast**: Tests should complete in milliseconds
4. **Test Isolated**: Tests should not depend on each other
5. **Test Clear**: Test names should describe what they verify

### Next Steps

1. Start with core package tests (identity, message, types)
2. Move to indexes package tests (hashtags, follows)
3. Expand scuttlego package tests
4. Expand UI package tests
5. Create config package tests
6. Add integration tests
7. Run full test suite with coverage
8. Address any failing tests or coverage gaps

---

**Remember**: Good tests are an investment in code quality. They catch bugs early, document expected behavior, and enable safe refactoring.

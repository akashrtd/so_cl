# 🤖 so_cl AI Agent Reference Guide

> **Purpose**: This file serves as a memory bank, reference point, and operational guide for AI agents working on the so_cl project. It ensures continuity across sessions and maintains high-quality coding standards.

**Last Updated**: 2025-12-30
**Project Phase**: Phase 1 (P2P) ✅ Phase 0 Complete
**Status**: Phase 1 P2P features partially implemented

---

## 📋 **TABLE OF CONTENTS**

1. [Project Context](#1-project-context)
2. [Technology Stack](#2-technology-stack)
3. [Architecture Overview](#3-architecture-overview)
4. [AI Agent Workflow](#4-ai-agent-workflow)
5. [Precautions & Constraints](#5-precautions--constraints)
6. [Best Practices](#6-best-practices)
7. [Editable TODO List](#7-editable-todo-list)
8. [Testing Strategy](#8-testing-strategy)
9. [Common Tasks](#9-common-tasks)
10. [API Quick Reference](#10-api-quick-reference)
11. [Decision Matrix](#11-decision-matrix)
12. [Debugging Guidelines](#12-debugging-guidelines)
13. [Pause & Continue](#13-pause--continue)
14. [Success Metrics](#14-success-metrics)

---

## 1. **PROJECT CONTEXT**

### What is so_cl?

**so_cl** is a decentralized, ASCII-art-based social platform running in the terminal using the Scuttlebutt (SSB) protocol.

**Key Characteristics:**
- **P2P**: No central servers, offline-first
- **Terminal UI (TUI)**: ASCII art, 6x6 colored profile pictures
- **Decentralized**: Uses Scuttlebutt protocol (append-only logs)
- **Privacy-focused**: Local data only, encrypted connections
- **Open source**: MIT licensed

### Project Goals

| Phase | Duration | Features | Status |
|-------|----------|----------|--------|
| **Phase 0** | Week 1-2 | Identity, posts, feed view, ASCII PFP | ✅ Complete |
| **Phase 1** | Week 3-4 | P2P replication, follows, LAN discovery | 🔲 TODO |
| **Phase 2** | Week 5-6 | Replies, likes, hashtags, mentions | 🔲 TODO |
| **Phase 3** | Week 7+ | Search, profiles, settings, web UI | 🔲 TODO |

### Target Users
- Developers who love terminal tools
- Privacy advocates
- Offline-first enthusiasts
- SSB/Scuttlebutt community

---

## 2. **TECHNOLOGY STACK**

### Core Dependencies

| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| **Language** | Go | 1.23+ | Primary language |
| **SSB Library** | scuttlego | v0.0.4 | Scuttlebutt protocol implementation |
| **Storage** | BadgerDB v3 | v3.2103.5 | LSM-tree KV store (used by scuttlego) |
| **Frontend** | Bubble Tea | v1.3.10+ | TUI framework |
| **Styling** | Lip Gloss | v1.1.0+ | Terminal styling |
| **Logging** | Zap | v1.27.1+ | Structured logging |

### Transitive Dependencies (~40+)

Critical ones to be aware of:
- `github.com/dgraph-io/badger/v3` - Storage engine
- `github.com/planetary-social/scuttlego` - SSB protocol
- `github.com/ssbc/go-secretstream` - Crypto (NaCl)
- `github.com/ssbc/margaret` - Feed logs
- `go.uber.org/zap` - Logging
- `github.com/charmbracelet/bubbletea` - TUI
- `github.com/charmbracelet/lipgloss` - Styling

### Build Targets

```
Binary Size: 30-50MB (static, CGO_ENABLED=0)
Memory Usage: 100-200MB (typical load)
Platforms: Linux (amd64/arm64), macOS (amd64/arm64), Windows (amd64)
Build Status: ✅ Successful (22MB with CGO_ENABLED=0)
```

---

## 3. **ARCHITECTURE OVERVIEW**

### System Layers (Bottom-Up)

```
┌─────────────────────────────────────────┐
│  Bubble Tea TUI Event Loop               │  User interface
│  (60fps rendering, keyboard input)       │
├─────────────────────────────────────────┤
│  so_cl Model (State Machine)             │  Application logic
│  (Posts, Peers, Profile, Trending)       │
├─────────────────────────────────────────┤
│  Async Dispatcher (Worker Pools)         │  Orchestration
│  (DB writes, network I/O, replication)   │
├─────────────────────────────────────────┤
│  BadgerDB v3 Storage                    │  Persistence
│  - scuttlego internal (SSB data)        │  Managed by scuttlego
│  - so_cl indexes (hashtags, mentions)    │  Managed by so_cl
├─────────────────────────────────────────┤
│  scuttlego Service Layer                │  Protocol
│  (Commands/Queries, EBT, handshake)     │
├─────────────────────────────────────────┤
│  scuttlego (includes go-secretstream)   │  Security
│  (NaCl crypto, Ed25519 signatures)       │
└─────────────────────────────────────────┘
```

### Key Design Patterns

**1. Two-Layer Storage**
```
scuttlego Internal (BadgerDB):
  - SSB messages (canonical format)
  - Feed logs, peer identities
  - EBT replication state

so_cl Application Indexes (BadgerDB):
  - Hashtag counts (for trending)
  - Mention notifications
  - Trending metrics
```

**2. scuttlego Commands/Queries Pattern**
```go
// Commands (write operations)
scuttlego.App.Commands.PublishRaw.Handle(...)
scuttlego.App.Commands.Follow.Handle(...)
scuttlego.App.Commands.Connect.Handle(...)

// Queries (read operations)
scuttlego.App.Queries.GetMessage.Handle(...)
scuttlego.App.Queries.ReceiveLog.Handle(...)
scuttlego.App.Queries.Status.Handle(...)
```

**3. Async Operations via Bubble Tea**
```go
// Update() returns tea.Cmd
func (m *SoClModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    case tea.KeyMsg:
        case "enter":
            return m, m.publishPost(text)
}

// tea.Cmd spawns goroutine and returns message
func (m *SoClModel) publishPost(text string) tea.Cmd {
    return func() tea.Msg {
        // Async work here
        return PostPublishedMsg{Ref: msgRef}
    }
}
```

### File Structure

```
so_cl/
├── main.go                    # Entry point
├── go.mod                     # Dependencies
├── agent.md                   # THIS FILE - AI agent reference
├── production.md             # Full product specification
│
├── ui/                       # Presentation layer
│   ├── model.go              # Bubble Tea model
│   ├── view.go               # Rendering
│   └── update.go             # Event handlers
│
├── core/                     # Domain models
│   ├── identity.go           # Identity wrapper
│   ├── message.go            # Post/reply models
│   └── types.go              # Shared types
│
├── async/                    # Async orchestration
│   ├── dispatcher.go         # Worker pool
│   ├── db_worker.go          # BadgerDB workers
│   └── net_worker.go         # P2P sync
│
├── scuttlego/                # scuttlego wrapper
│   ├── service.go            # Service lifecycle
│   ├── config.go             # Configuration
│   └── adapter.go            # so_cl ↔ scuttlego bridge
│
├── indexes/                  # Custom indexes
│   ├── hashtags.go           # Hashtag counting
│   ├── mentions.go           # Mentions queue
│   └── trending.go           # Trending metrics
│
├── protocol/                 # Protocol adapters
│   ├── scuttlego.go          # scuttlego client wrapper
│   └── adapter.go            # Bridge to state machine
│
├── config/                   # Configuration
│   ├── config.go             # Load from env/file
│   └── defaults.go           # Default values
│
└── tests/                    # Tests
    ├── model_test.go         # UI tests
    ├── dispatcher_test.go    # Async tests
    ├── indexes_test.go        # Index tests
    └── scuttlego_test.go    # Integration tests
```

---

## 4. **AI AGENT WORKFLOW**

### Daily Workflow

**Morning Standup**
```
1. Read agent.md (this file) → Refresh context
2. Check TODO list → Pick next task
3. Review recent changes → Understand state
4. Choose highest-impact task → Start work
```

**During Development**
```
✅ Read relevant code first
✅ Write test FIRST (TDD)
✅ Implement feature
✅ Run: go test -race ./...
✅ Run: go vet ./...
✅ Run: golangci-lint run
✅ Run: CGO_ENABLED=0 go build
✅ Update TODO list
   - Run `go test -race` before committing
   - Maintain 85%+ coverage (95%+ for critical paths)

2. **Never modify scuttlego internals**
   - Use adapters/wrappers instead
   - Don't patch scuttlego code
   - Use Commands/Queries pattern

3. **Never hardcode secrets**
   - Use environment variables
   - Add sensitive values to `.gitignore`
   - Use proper secret management

4. **Never implement custom crypto**
   - Use scuttlego's go-secretstream
   - Don't implement Ed25519 signing manually
   - Don't create your own encryption

5. **Never break SSB protocol**
   - Don't modify message format
   - Don't change signature verification
   - Don't alter EBT logic

### ⚠️ **ALWAYS ASK BEFORE** (Requires Approval)

| Action | When to Ask | Reason |
|--------|------------|--------|
| **Protocol Changes** | Modify EBT, handshake, message format | Affects all peers, breaking change |
| **Schema Changes** | Add/modify scuttlego internal schema | Breaks backwards compatibility |
| **Security Decisions** | Rate limiting, access control, crypto params | Human intent required |
| **New Dependencies** | Add external Go library | Supply chain + reproducibility |
| **Performance Tradeoffs** | Cache entire feed, disable security | Long-term impact |
| **scuttlego Config** | Change replication scheduler, EBT params | Complex internals |

### ✅ **ALWAYS DO** (Autonomous Actions)

1. **UI rendering** (Lip Gloss, colors, layouts)
2. **Validation logic** (280-char limit, ASCII-only)
3. **Error handling** (try-catch, logging, formatted errors)
4. **Performance optimization** (caching, goroutines for I/O)
5. **Testing** (unit, integration, E2E tests)
6. **Documentation** (godoc, comments, README)
7. **Build/deploy** (Docker, goreleaser, CI/CD)

---

## 6. **BEST PRACTICES**

### Code Style

**Go Conventions:**
- Use `gofmt` (enforced by linter)
- Exported types: `PascalCase`
- Private types: `camelCase`
- Constants: `UPPER_SNAKE_CASE`
- Interfaces: `er` suffix (`Reader`, `Writer`)

**Error Handling:**
```go
// ✅ Good: Wrap errors with context
return fmt.Errorf("failed to publish post: %w", err)

// ❌ Bad: Generic errors
return errors.New("error")

// ✅ Good: Use structured logging
zap.L().Error("publish failed",
    zap.String("text", text),
    zap.Error(err),
)
```

**Async Operations:**
```go
// ✅ Good: Use channels for coordination
resultCh := make(chan PostResult)
go func() {
    result := publishPost(text)
    resultCh <- result
}()

// ❌ Bad: Blocking in UI thread
result := publishPost(text)  // Blocks TUI!
```

### Testing

**Unit Tests (Fast, Isolated):**
```go
func TestPublishPost_ValidText(t *testing.T) {
    // Arrange
    m := NewTestModel()
    text := "hello world"

    // Act
    err := m.PublishPost(text)

    // Assert
    require.Nil(t, err)
    require.Equal(t, 1, len(m.posts))
}
```

**Integration Tests (Real DB):**
```go
func TestPublishPost_WithScuttlego(t *testing.T) {
    // Start test service
    svc := startTestScuttlego(t)
    defer svc.Close()

    // Test real API
    msgRef, err := svc.App.Commands.PublishRaw.Handle(...)

    // Verify
    require.Nil(t, err)
    require.NotEmpty(t, msgRef)
}
```

**Race Detection:**
```bash
# ALWAYS run with -race
go test -race ./...
```

### Memory Management

**BadgerDB Usage:**
```go
// ✅ Good: Close transactions
err := db.Update(func(txn *badger.Txn) error {
    return txn.Set(key, val)
})

// ✅ Good: Use View for reads
err := db.View(func(txn *badger.Txn) error {
    item, _ := txn.Get(key)
    return item.Value(func(v []byte) error {
        return nil
    })
})

// ❌ Bad: Don't forget to close
txn := db.NewTransaction(true)
txn.Set(key, val)  // Leak if not committed/discarded!
```

**Goroutine Management:**
```go
// ✅ Good: Use context for cancellation
ctx, cancel := context.WithCancel(context.Background())
go func() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            doWork()
        case <-ctx.Done():
            return  // Clean exit
        }
    }
}()
cancel()  // Trigger cleanup
```

### Logging

**Structured Logging:**
```go
import "go.uber.org/zap"

logger, _ := zap.NewDevelopment()
zap.ReplaceGlobals(logger)

// Usage
zap.L().Info("peer connected",
    zap.String("peer", "@alice"),
    zap.Duration("latency", 120*time.Millisecond),
)
```

**Log Levels:**
- `Debug`: Detailed diagnostics (disabled in production)
- `Info`: Normal operation (connection, publish, sync)
- `Warn`: Unexpected but recoverable (retry, fallback)
- `Error`: Operation failed (write error, network timeout)

---

## 7. **EDITABLE TODO LIST**

### 📋 **Phase 0: MVP (Week 1-2)**

- [x] **Initialize project structure**
  - [x] Create `go.mod` with dependencies
  - [x] Set up directory structure (ui/, core/, async/, etc.)
  - [x] Create `scuttlego/service.go` wrapper
  - [x] Create `indexes/` package skeleton
  - [x] Add `.gitignore`, `LICENSE`, `README.md`

 - [x] **Identity generation**
   - [x] Generate SSB keypair (Ed25519)
   - [x] Store in BadgerDB (via scuttlego)
   - [x] Create `core/identity.go` wrapper
   - [x] Add ASCII PFP generator (6x6, ANSI colors)
   - [x] Test key generation and signing

 - [x] **Post creation**
    - [x] Build TUI composer (Bubble Tea input)
    - [x] Validate 280-char limit
    - [x] Build SSB post content (`message.NewRawContent`)
    - [x] Publish via `scuttlego.App.Commands.PublishRaw`
    - [x] Index post for hashtags/mentions
    - [x] Optimistic UI update
    - [x] Test publish flow

  - [x] **Feed view**
    - [x] Query posts via `scuttlego.App.Queries.ReceiveLog`
    - [x] Format posts for TUI display
    - [x] Show ASCII PFP (6x6 colors)
    - [x] Show timestamp, author, text
    - [x] Implement pagination (100 posts max in memory)
    - [x] Test feed rendering

  - [x] **Static binary build**
    - [x] Create `Makefile`
    - [x] Build with `CGO_ENABLED=0`
    - [x] Verify <50MB size (22MB achieved)
    - [x] Test locally (macOS)
    - [x] Add `goreleaser` config

### 📋 **Phase 1: P2P (Week 3-4)**

- [ ] **Follow via invite codes**
  - [ ] Parse invite codes (`domain.NewInvite`)
  - [ ] Implement `scuttlego.App.Commands.RedeemInvite`
  - [ ] Generate invite codes for self
  - [ ] Add follow/unfollow UI
  - [ ] Test follow flow

- [ ] **Peer connections**
  - [ ] Implement `scuttlego.App.Commands.Connect`
  - [ ] Handle incoming connections
  - [ ] Show peer list sidebar
  - [ ] Display connection status
  - [ ] Test peer-to-peer connection

- [ ] **EBT replication**
  - [ ] Verify EBT is automatic (via scuttlego)
  - [ ] Monitor replication status
  - [ ] Show sync progress in UI
  - [ ] Test replication with 2+ peers

- [ ] **LAN discovery (UDP broadcast)**
  - [ ] Start UDP advertiser (scuttlego builtin)
  - [ ] Discover local peers
  - [ ] Auto-connect to discovered peers
  - [ ] Test on local network

- [ ] **Real-time feed sync**
  - [ ] Subscribe to new message events
  - [ ] Update feed when new messages arrive
  - [ ] Show "new messages available" indicator
  - [ ] Test real-time updates

### 📋 **Phase 2: Social (Week 5-6)**

- [ ] **Replies (threaded messages)**
  - [ ] Build reply composer
  - [ ] Link to parent message (root/branch)
  - [ ] Show thread hierarchy in feed
  - [ ] Test reply flow

- [ ] **Like/repost reactions**
  - [ ] Add reaction button in UI
  - [ ] Publish reaction messages (type: "vote")
  - [ ] Show reaction counts on posts
  - [ ] Test reactions

- [ ] **Hashtag indexing**
  - [ ] Extract hashtags from posts
  - [ ] Store in BadgerDB indexes
  - [ ] Count hashtag usage
  - [ ] Test hashtag extraction

- [ ] **Trending sidebar**
  - [ ] Query top 10 hashtags from indexes
  - [ ] Display in sidebar
  - [ ] Update in real-time
  - [ ] Test trending display

- [ ] **@mentions + notifications**
  - [ ] Extract @mentions from posts
  - [ ] Store in mention queue (BadgerDB)
  - [ ] Show notification indicator
  - [ ] Display mention list
  - [ ] Test mentions

- [ ] **Follow/unfollow graph**
  - [ ] Track follow status per peer
  - [ ] Show following/followers count
  - [ ] Display follow relationships
  - [ ] Test follow graph

### 📋 **Phase 3: Polish (Week 7+)**

- [ ] **Search + filters**
  - [ ] Search posts by text
  - [ ] Filter by hashtag
  - [ ] Filter by author
  - [ ] Test search

- [ ] **Profile pages**
  - [ ] Show user profile (PFP, bio, stats)
  - [ ] Show user's posts
  - [ ] Show followers/following
  - [ ] Test profile view

- [ ] **Settings (username, PFP, discovery)**
  - [ ] Change username (store in BadgerDB)
  - [ ] Regenerate ASCII PFP
  - [ ] Toggle LAN discovery
  - [ ] Test settings

- [ ] **Export/backup identity**
  - [ ] Export keypair to file
  - [ ] Encrypt export (password)
  - [ ] Import from backup
  - [ ] Test backup/restore

- [ ] **Web UI**
  - [ ] Implement same state machine in `net/http`
  - [ ] Serve from same process
  - [ ] Share scuttlego service
  - [ ] Test web UI

---

## 8. **TESTING STRATEGY**

### Test Pyramid

```
        /\
       /  \        E2E Tests (TUI mock)
      /────\       3-5 tests, 30s to run
     /      \
    /────────\     Integration Tests
   /          \    (scuttlego + BadgerDB)
  /────────────\   15-20 tests, 5s to run
 /              \
/______________\  Unit Tests
     Base        (Functions)
                 50-70 tests, 1s to run
```

### Coverage Targets

| Module | Target | Critical Path |
|--------|--------|---------------|
| `scuttlego/service.go` | 95%+ | ✅ YES (service lifecycle) |
| `protocol/adapter.go` | 90%+ | ✅ YES (bridge) |
| `indexes/hashtags.go` | 90% | - |
| `async/dispatcher.go` | 95%+ | ✅ YES (goroutines) |
| `ui/model.go` | 90% | - |
| **Overall** | **85%** | **Critical 95%+** |

### Test Categories

**Unit Tests (70%):**
```bash
# Run all unit tests
go test ./... -run "^Test[^Integration][^E2E]" -v

# Run with race detector
go test -race ./...
```

**Integration Tests (20%):**
```go
func TestIntegration_PublishAndSync(t *testing.T) {
    // Start 2 scuttlego services
    svc1 := startTestScuttlego(t)
    svc2 := startTestScuttlego(t)
    defer svc1.Close()
    defer svc2.Close()

    // Connect peers
    connectPeers(svc1, svc2)

    // Publish on svc1
    msgRef := publishPost(svc1, "hello")

    // Wait for sync
    time.Sleep(1 * time.Second)

    // Verify on svc2
    msg, _ := svc2.App.Queries.GetMessage.Handle(
        queries.GetMessage{Ref: msgRef},
    )
    require.Equal(t, "hello", msg.Content.Text)
}
```

**E2E Tests (10%):**
```go
func TestE2E_UserFlow(t *testing.T) {
    // Start full TUI model
    m := NewTestModel()

    // Simulate user input
    m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("hello")})
    m.Update(tea.KeyMsg{Type: tea.KeyEnter})

    // Wait for async operations
    time.Sleep(100 * time.Millisecond)

    // Verify state
    require.Greater(t, len(m.cachedFeed), 0)
    require.Equal(t, "hello", m.cachedFeed[0].Text)
}
```

### Running Tests

```bash
# All tests with coverage
go test -race -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package
go test -race ./scuttlego/ -v

# Run specific test
go test -v ./ui/ -run TestModelUpdate
```

### Race Conditions

**ALWAYS run with `-race`:**
```bash
go test -race ./...
```

**Common race patterns:**
```go
// ❌ RACE: Shared state without mutex
var counter int
go func() {
    counter++  // Data race!
}()

// ✅ FIXED: Use atomic or mutex
var counter int64
go func() {
    atomic.AddInt64(&counter, 1)
}()
```

---

## 9. **COMMON TASKS**

### Task 1: Create a New Package

```bash
# 1. Create directory
mkdir pkgname

# 2. Create files
touch pkgname/pkgname.go
touch pkgname/pkgname_test.go

# 3. Add imports and basic structure
```

```go
// pkgname/pkgname.go
package pkgname

type PkgName struct {
    field string
}

func NewPkgName(field string) *PkgName {
    return &PkgName{field: field}
}

func (p *PkgName) DoSomething() error {
    return nil
}
```

### Task 2: Add a New Bubble Tea View

```go
// ui/view.go
func (m *SoClModel) View() string {
    return lipgloss.JoinVertical(
        lipgloss.Left,
        m.renderHeader(),
        m.renderFeed(),
        m.renderComposer(),
        m.renderSidebar(),
    )
}

func (m *SoClModel) renderFeed() string {
    var feed strings.Builder
    for i, post := range m.feed {
        feed.WriteString(fmt.Sprintf(
            "[%d] %s: %s\n",
            i+1, post.Author, post.Text,
        ))
    }
    return feed.String()
}
```

### Task 3: Add Async Operation

```go
// async/dispatcher.go
type Dispatcher struct {
    scuttlego *scuttlego.Service
    writeQ    chan WriteTask
}

func (d *Dispatcher) Publish(text string) tea.Cmd {
    return func() tea.Msg {
        // Build content
        content, _ := message.NewRawContent([]byte(
            fmt.Sprintf(`{"type":"post","text":"%s"}`, text),
        ))

        // Publish via scuttlego
        msgRef, err := d.scuttlego.App.Commands.PublishRaw.Handle(
            commands.PublishRaw{Content: content},
        )
        if err != nil {
            return ErrorMsg{err}
        }

        return PostPublishedMsg{Ref: msgRef}
    }
}
```

### Task 4: Add BadgerDB Index

```go
// indexes/hashtags.go
func (idx *Indexer) IndexHashtags(msg *Message, tags []string) error {
    return idx.db.Update(func(txn *badger.Txn) error {
        for _, tag := range tags {
            key := []byte("hashtag:" + tag)

            // Get current count
            item, _ := txn.Get(key)
            var count int
            if item != nil {
                item.Value(func(val []byte) error {
                    json.Unmarshal(val, &count)
                    return nil
                })
            }

            // Increment
            count++
            data, _ := json.Marshal(count)

            return txn.Set(key, data)
        }
        return nil
    })
}
```

### Task 5: Add scuttlego Query

```go
// ui/model.go
func (m *SoClModel) loadFeed() tea.Cmd {
    return func() tea.Msg {
        // Query scuttlego
        receiveLog, err := m.scuttlego.App.Queries.ReceiveLog.Handle(
            queries.ReceiveLog{Limit: 100},
        )
        if err != nil {
            return ErrorMsg{err}
        }

        // Convert to so_cl posts
        posts := make([]Post, 0)
        for msg := range receiveLog.Messages {
            posts = append(posts, Post{
                Author: msg.Author.String(),
                Text:   msg.Content.Text,
                Time:   msg.Timestamp,
            })
        }

        return FeedLoadedMsg{Posts: posts}
    }
}
```

---

## 10. **API QUICK REFERENCE**

### scuttlego Commands (Write Operations)

```go
import (
    "github.com/planetary-social/scuttlego/service/app/commands"
    "github.com/planetary-social/scuttlego/service/domain/transport/message"
)

// Publish a post
content, _ := message.NewRawContent([]byte(`{
    "type": "post",
    "text": "hello world"
}`))
msgRef, err := svc.App.Commands.PublishRaw.Handle(
    commands.PublishRaw{Content: content},
)

// Follow a peer
err = svc.App.Commands.Follow.Handle(
    commands.Follow{Feed: peerIdentity},
)

// Connect to a peer
err = svc.App.Commands.Connect.Handle(
    commands.Connect{Address: "net:127.0.0.1:8008~shs:@alice/..."},
)

// Redeem an invite code
invite, _ := domain.NewInvite("ssb:feed/invite/...")
err = svc.App.Commands.RedeemInvite.Handle(
    commands.RedeemInvite{Invite: invite},
)
```

### scuttlego Queries (Read Operations)

```go
import "github.com/planetary-social/scuttlego/service/app/queries"

// Get a specific message
msg, err := svc.App.Queries.GetMessage.Handle(
    queries.GetMessage{Ref: msgRef},
)

// Get recent messages (feed)
receiveLog, err := svc.App.Queries.ReceiveLog.Handle(
    queries.ReceiveLog{Limit: 100},
)
for msg := range receiveLog.Messages {
    // Process messages
}

// Get service status
status, err := svc.App.Queries.Status.Handle(queries.Status{})
identity := status.Identity
following := status.Following
```

### Bubble Tea Events

```go
import tea "github.com/charmbracelet/bubbletea"

// Key messages
case tea.KeyMsg:
    switch msg.Type {
    case tea.KeyEnter:
        // Enter pressed
    case tea.KeyCtrlC:
        // Exit
    case tea.KeyRunes:
        // Typing
    }

// Custom messages
type PostPublishedMsg struct {
    Ref refs.Message
}

type FeedLoadedMsg struct {
    Posts []Post
}

type ErrorMsg struct {
    Err error
}
```

### BadgerDB Operations

```go
import badger "github.com/dgraph-io/badger/v3"

// Open database
db, err := badger.Open(badger.DefaultOptions("./data"))
defer db.Close()

// Write transaction
err = db.Update(func(txn *badger.Txn) error {
    return txn.Set(key, value)
})

// Read transaction
err = db.View(func(txn *badger.Txn) error {
    item, err := txn.Get(key)
    return item.Value(func(v []byte) error {
        return nil
    })
})

// Batch write
err = db.Update(func(txn *badger.Txn) error {
    for k, v := range items {
        if err := txn.Set([]byte(k), []byte(v)); err != nil {
            return err
        }
    }
    return nil
})
```

---

## 11. **DECISION MATRIX**

### When to Code Autonomously ✅

| Scenario | Action | Example |
|----------|--------|---------|
| **UI rendering** | Code it | Add colors, format posts |
| **Validation** | Code it | Check 280-char limit |
| **Error handling** | Code it | Wrap errors, log with context |
| **Performance** | Code it | Cache, goroutines |
| **Testing** | Code it | Write unit/integration tests |
| **Documentation** | Code it | Godoc, comments |

### When to Ask for Approval ⏸️

| Scenario | Why Ask |
|----------|----------|
| **Protocol changes** | Affects all peers |
| **Schema changes** | Breaks backwards compatibility |
| **Security decisions** | Human intent required |
| **New dependencies** | Supply chain risk |
| **Performance tradeoffs** | Long-term impact |
| **scuttlego config** | Complex internals |

---

## 12. **DEBUGGING GUIDELINES**

### Common Issues

**Issue: Posts not syncing between peers**

**Diagnosis:**
```bash
# Check EBT status
grep "replication" ~/.so_cl/log.txt

# Verify connection
grep "peer connected" ~/.so_cl/log.txt

# Check scuttlego logs
DEBUG=1 ./so_cl 2>&1 | grep -i "ebt\|replicate"
```

**Fix:**
- Verify peers are connected
- Check EBT is running (automatic in scuttlego)
- Ensure both peers follow each other

---

**Issue: High memory usage (>200MB)**

**Diagnosis:**
```go
// Check goroutine count
fmt.Println("Goroutines:", runtime.NumGoroutine())

// Check BadgerDB cache
// (scuttlego manages this internally)

// Check feed size
fmt.Println("Feed size:", len(m.feed))
```

**Fix:**
- Implement pagination (max 100 posts in memory)
- Reduce BadgerDB cache size (scuttlego config)
- Clear old feed entries

---

**Issue: Slow UI (<30fps)**

**Diagnosis:**
```bash
# Profile with pprof
go tool pprof http://localhost:6060/debug/pprof/profile

# Check render time
# Add timing in View()
start := time.Now()
view := m.View()
fmt.Println("Render time:", time.Since(start))
```

**Fix:**
- Cache rendered strings
- Reduce post count in feed
- Optimize View() logic

---

**Issue: scuttlego service won't start**

**Diagnosis:**
```bash
# Check BadgerDB permissions
ls -la ~/.so_cl/data/

# Check port availability
lsof -i :8008

# Check scuttlego logs
./so_cl 2>&1 | grep -i "error\|panic"
```

**Fix:**
- Verify data directory is writable
- Ensure port 8008 is available
- Check scuttlego configuration

---

### Debugging Tools

**Logging:**
```bash
# Enable debug logging
DEBUG=1 ./so_cl

# Log to file
DEBUG=1 ./so_cl > so_cl.log 2>&1

# Monitor in real-time
tail -f so_cl.log | grep -i "error\|warn"
```

**Profiling:**
```go
import _ "net/http/pprof"

go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()

# Access: http://localhost:6060/debug/pprof/
```

**Memory:**
```bash
# Check memory usage
ps aux | grep so_cl

# Memory profile
go tool pprof http://localhost:6060/debug/pprof/heap
```

---

## 13. **PAUSE & CONTINUE**

### How to Pause

1. **Save current state:**
   ```bash
   # Commit current work
   git add .
   git commit -m "WIP: [task name]"

   # Create checkpoint branch
   git branch checkpoint-$(date +%Y%m%d-%H%M%S)
   ```

2. **Update agent.md:**
   - Update "Session State" (above)
   - Update TODO list status
   - Note blockers

3. **Clean up:**
   ```bash
   # Stop running processes
   pkill so_cl

   # Clear temporary files
   rm -rf /tmp/so_cl-*
   ```

### How to Continue

1. **Read agent.md** (this file) - Refresh context
2. **Check TODO list** - Pick next task
3. **Review recent commits:**
   ```bash
   git log --oneline -10
   ```
4. **Review branch status:**
   ```bash
   git status
   git diff
   ```
5. **Start working** - Follow daily workflow (Section 4)

### Session State Template

```markdown
## Session State

**Date**: 2025-12-29
**Phase**: Phase 0 (MVP) → Phase 1 ready
**Current Task**: Phase 0 completed - scuttlego integration done
**Last Completed**: Full scuttlego service integration with Publish, Follow, Connect, GetRecentMessages
**Blockers**: None
**Next Priority**: Start Phase 1 - P2P features (invite codes, peer connections, EBT replication)

**Files Modified**:
- [x] `scuttlego/service.go` - Full scuttlego integration with BadgerDB, network, EBT
- [x] `core/asciipfp.go` - 6x6 ANSI art PFP generator with deterministic seeding
- [x] `go.mod` - Added scuttlego v0.0.4 and all dependencies

**Tests Written**:
- [x] `core/asciipfp_test.go` - Tests for PFP generation (determinism, rendering)

**Next Session Goals**:
  | 2025-12-29 | Full scuttlego service integration (Publish, Follow, Connect, GetRecentMessages) | ✅ Complete | All Phase 0 MVP features implemented
  1. [ ] Add tests for scuttlego service wrapper
  2. [ ] Integrate ASCII PFP into feed rendering
  3. [ ] Implement optimistic UI updates for publishing
```

---

## 14. **SCUTTLEGO INTEGRATION NOTES**

### Service Initialization
The scuttlego service is initialized via `scuttlegodi.BuildService(privateIdentity, config)` which:
- Generates new Ed25519 keypair automatically
- Initializes BadgerDB at `config.DataDirectory`
- Sets up network listener on `config.ListenAddress`
- Configures EBT replication
- Returns cleanup function for graceful shutdown

### Key API Patterns

**Publishing Posts:**
```go
content := map[string]interface{}{
    "type": "post",
    "text": text,
}
contentJSON, _ := json.Marshal(content)
rawContent, _ := message.NewRawContent(contentJSON)
cmd, _ := commands.NewPublishRaw(rawContent.Bytes())
msgRef, _ := svc.App.Commands.PublishRaw.Handle(cmd)
```

**Following Peers:**
```go
peerIdentity, _ := refs.NewIdentity("@alice.key.ed25519")
cmd := commands.Follow{Target: peerIdentity}
_ = svc.App.Commands.Follow.Handle(cmd)
```

**Connecting:**
```go
// Multiserver address format: "net:host:port~shs:@alice.key.ed25519"
addr := network.NewAddress(address)
cmd := commands.Connect{
    Remote:  refs.MustNewIdentityFromPublic(peerIdentity).Identity(),
    Address: addr,
}
_ = svc.App.Commands.Connect.Handle(ctx, cmd)
```

**Querying Feed:**
```go
startSeq, _ := common.NewReceiveLogSequence(0)
query, _ := queries.NewReceiveLog(startSeq, limit)
messages, _ := svc.App.Queries.ReceiveLog.Handle(query)
for _, logMsg := range messages {
    author := logMsg.Message.Author().String()
    // Parse content manually as Known() is unexported
}
```

### Type Mappings
- `identity.Private` → `refs.Identity` (via `identity.Public()`)
- `message.RawContent` → `[]byte` (via `Bytes()`)
- Content parsing requires manual JSON unmarshal of `Content().Raw().Bytes()`

---

## 15. **SUCCESS METRICS**

### Code Quality

| Metric | Target | How to Measure |
|--------|--------|----------------|
| **Test Coverage** | ≥85% (critical 95%+) | `go test -cover ./...` |
| **Race Conditions** | 0 | `go test -race ./...` |
| **Lint Errors** | 0 | `golangci-lint run` |
| **Build Success** | 100% | `CGO_ENABLED=0 go build` |
| **Binary Size** | ≤50MB | `ls -lh so_cl` |
| **Memory Usage** | ≤200MB | `ps aux \| grep so_cl` |

### Feature Delivery

| Metric | Target | Deadline |
|--------|--------|----------|
| **Phase 0** | MVP features | Week 2 |
| **Phase 1** | P2P features | Week 4 |
| **Phase 2** | Social features | Week 6 |
| **Phase 3** | Polish features | Week 7+ |

### User Experience

| Metric | Target | How to Measure |
|--------|--------|----------------|
| **UI Responsiveness** | ≥60fps | Manual testing |
| **Startup Time** | ≤5s | `time ./so_cl` |
| **Post Publish Time** | ≤1s | Manual testing |
| **Sync Time** | ≤30s (first) | Manual testing |

---

## **QUICK REFERENCE CARD**

### Essential Commands

```bash
# Development
go mod tidy                    # Update dependencies
go test -race -cover ./...     # Test with coverage
go vet ./...                   # Static analysis
golangci-lint run              # Linting
CGO_ENABLED=0 go build         # Static build

# Scuttlego
./so_cl                        # Run application
DEBUG=1 ./so_cl               # Debug mode

# Git
git status                     # Check state
git log --oneline -5          # Recent commits
git diff                       # Show changes
```

### Key Files

| File | Purpose |
|------|---------|
| `agent.md` | THIS FILE - AI agent reference |
| `production.md` | Full product specification |
| `go.mod` | Go dependencies |
| `scuttlego/service.go` | scuttlego service wrapper |
| `ui/model.go` | Bubble Tea model |

### Key Imports

```go
// SSB protocol
import "github.com/planetary-social/scuttlego/service"
import "github.com/planetary-social/scuttlego/service/app/commands"
import "github.com/planetary-social/scuttlego/service/app/queries"

// Storage
import badger "github.com/dgraph-io/badger/v3"

// UI
import tea "github.com/charmbracelet/bubbletea"
import "github.com/charmbracelet/lipgloss"

// Logging
import "go.uber.org/zap"
```

---

## **EMERGENCY CONTACT**

If you encounter critical issues:

1. **Check logs** (`~/.so_cl/log.txt`)
2. **Review this file** (agent.md)
3. **Read production.md** (full spec)
4. **Check scuttlego docs** (github.com/planetary-social/scuttlego)
5. **Pause and document** (Section 13)

---

**Remember**: This is a **living document**. Update it as you learn and progress. Good luck! 🚀

---

**END OF agent.md**

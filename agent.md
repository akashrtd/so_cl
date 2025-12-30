# 🤖 so_cl AI Agent Reference Guide

> **Purpose**: This file serves as a memory bank, reference point, and operational guide for AI agents working on the so_cl project. It ensures continuity across sessions and maintains high-quality coding standards.

**Last Updated**: 2025-12-30
**Project Phase**: Phase 3 (Polish Features) 🔄 In Progress
**Status**: Phase 2 (Social Features) ✅ Complete - Phase 3 implementation in progress

---

## 📋 **TABLE OF CONTENTS**

### 1. [Project Context](#1-project-context)

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
|-------|----------|--------|
| **Phase 0** | Week 1-2 | Identity, posts, feed view, ASCII PFP | ✅ Complete |
| **Phase 1** | Week 3-4 | P2P replication, follows, LAN discovery | ✅ Complete |
| **Phase 2** | Week 5-6 | Replies, likes, hashtags, mentions | ✅ Implementation Complete |
| **Phase 3** | Week 7+ | Search, profiles, settings, web UI | 🔄 In Progress |

### Target Users
- Developers who love terminal tools
- Privacy advocates
- Offline-first enthusiasts
- SSB/Scuttlebutt community

---

## 2. [Technology Stack](#2-technology-stack)

### Core Dependencies

| Component | Technology | Version | Purpose |
|-----------|---------|---------|
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

---

## 3. [Architecture Overview](#3-architecture-overview)

### System Layers (Bottom-Up)

```
┌─────────────────────────────────────────────────────────┐
│  Bubble Tea TUI Event Loop               │  User interface
│  (60fps rendering, keyboard input)       │
├─────────────────────────────────────────────────┤
│  so_cl Model (State Machine)             │  Application logic
│  (Posts, Peers, Profile, Trending)       │
├─────────────────────────────────────────────────┤
│  Async Dispatcher (Worker Pools)         │  Orchestration
│  (DB writes, network I/O, replication)   │
├─────────────────────────────────────────────────┤
│  BadgerDB v3 Storage                    │  Persistence
│  - scuttlego internal (SSB data)        │  Managed by scuttlego
│  - so_cl indexes (hashtags, mentions)    │  Managed by so_cl
├─────────────────────────────────────────────────┘
│  scuttlego Service Layer                │  Protocol
│  (Commands/Queries, EBT, handshake)     │
│  scuttlego (includes go-secretstream)   │  Security
│  (NaCl crypto, Ed25519 signatures)       │
└─────────────────────────────────────────────────┘
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
  - Follow relationships (follow graph)
  - Full-text search index
  - Profile data
  - Settings data
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
        return m, m.publishPost(text)
    }
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
│   ├── hashtags.go           # Hashtag counting & search
│   ├── mentions.go           # Mentions queue
│   ├── trending.go           # Trending metrics
│   └── follows.go           # Follow graph tracking
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
│   ├── model_test.go         # UI tests
│   ├── dispatcher_test.go    # Async tests
│   ├── indexes_test.go        # Index tests
│   └── scuttlego_test.go    # Integration tests
```

---

## 4. [AI Agent Workflow](#4-ai-agent-workflow)

### Daily Workflow

**Morning Standup**
```
1. Read agent.md (this file) → Refresh context
2. Check TODO list → Pick next task
3. Review recent changes → Understand state
4. Choose highest-impact task → Start work
```

### During Development

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
```

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

## 5. [Precautions & Constraints](#5-precautions--constraints)

### ⚠️ **NEVER MODIFY scuttlego INTERNALS**

- Use adapters/wrappers instead
- Don't patch scuttlego code
- Use Commands/Queries pattern
- Don't implement custom crypto

### ⚠️ **NEVER HARDCODE SECRETS**

- Use environment variables
- Add sensitive values to `.gitignore`
- Use proper secret management

### ⚠️ **NEVER IMPLEMENT CUSTOM CRYPTO**

- Use scuttlego's go-secretstream
- Don't implement Ed25519 signing manually

### ⚠️ **NEVER BREAK SSB PROTOCOL**

- Don't modify message format
- Don't change signature verification
- Don't alter EBT logic

### ⚠️ **NEVER USE scuttlego Config** (Complex Internals)

- Changing replication scheduler affects EBT behavior
- EBT params have complex interactions
- Requires careful testing

---

## 6. [Best Practices](#6-best-practices)

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
}

cancel()  // Trigger cleanup
```

**Memory Management:**
```go
// ✅ Good: Close transactions
err := db.Update(func(txn *badger.Txn) error {
    return txn.Set(key, val)
})

// ❌ Bad: Don't forget to close
txn := db.NewTransaction(true)
txn.Set(key, val)  // Leak if not committed/discarded!
```

**Logging:**
```go
// ✅ Good: Use structured logging
zap.L().Error("publish failed",
    zap.String("text", text),
    zap.Error(err),
)

zap.L().Info("peer connected",
    zap.String("peer", "@alice"),
    zap.Duration("latency", 120*time.Millisecond),
)
)
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
func TestIntegration_PublishAndSync(t *testing.T) {
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

---

## 7. [Editable TODO List](#7-editable-todo-list)

### 📋 **Phase 0: MVP (Week 1-2)** ✅ Complete

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

- [x] **Feed view**
  - [x] Query posts via `scuttlego.App.Queries.ReceiveLog`
  - [x] Format posts for TUI display
  - [x] Show ASCII PFP (6x6 colors)
  - [x] Show timestamp, author, text
- [x] Implement pagination (100 posts max in memory)

- [x] **Static binary build**
  - [x] Create `Makefile`
  - [x] Build with `CGO_ENABLED=0`
  - [x] Verify <50MB size (22MB achieved)

### 📋 **Phase 1: P2P (Week 3-4)** ✅ Complete

- [x] **Follow via invite codes**
  - [x] Parse invite codes (`invites.NewInviteFromString`)
  - [x] Implement `scuttlego.App.Commands.RedeemInvite`
  - [x] Generate invite codes for self
  - [x] Add follow/unfollow UI

- [x] **Peer connections**
  - [x] Implement `scuttlego.App.Commands.Connect`
- [x] Handle incoming connections
  - [x] Show peer list sidebar (F1 toggle)
  - [x] Display connection status

- [x] **EBT replication**
  - [x] Verify EBT is automatic (via scuttlego)
- [x] Monitor replication status (GetEBTStatus)
  - [x] Show sync progress in UI (sidebar)
- [x] Test replication with 2+ peers

- [x] **Real-time feed sync**
  - [x] Subscribe to new message events (StartWatching)
- [x] Update feed when new messages arrive (NewMessageMsg)
- [x] Show "new messages available" indicator

### 📋 **Phase 2: Social (Week 5-6)** ✅ Implementation Complete

- [x] **Replies (threaded messages)**
  - [x] Build reply composer
  - [x] Link to parent message (root/branch)
  - [x] Show thread hierarchy in feed
- [x] Test reply flow

- [x] **Like/repost reactions**
  - [x] Add reaction button in UI (F5 key)
- [x] Publish reaction messages (type: "vote")
  - [x] Show reaction counts on posts

- [x] **Hashtag indexing**
  - [x] Extract hashtags from posts
- [x] Store in BadgerDB indexes
- [x] Count hashtag usage
- [x] Test hashtag extraction

- [x] **Trending sidebar**
  - [x] Query top 10 hashtags from indexes
- [x] Display in sidebar (F6 key)
- [x] Update in real-time

- [x] **@mentions + notifications**
  - [x] Extract @mentions from posts
- [x] Store in mention queue (BadgerDB)
- [x] Show notification indicator in header
- [x] Display mention list (F7 key)

### 📋 **Phase 3: Polish (Week 7+)** 🔄 In Progress

- [ ] **Follow/unfollow graph**
  - [x] Track follow status per peer in BadgerDB
  - [x] Show following/followers count in UI
  - [x] Display follow relationships sidebar (F11 key)
  - [ ] Test follow graph

- [ ] **Search + filters**
  - [x] Search posts by text (index posts)
  - [x] Filter by hashtag
  - [x] Filter by author
  - [x] Add search UI (F8 key)
  - [ ] Test search functionality

- [ ] **Profile pages**
  - [x] Create profile data model (core.SoClProfile)
  - [x] Show user profile (PFP, bio, stats)
  - [x] Show user's posts (count only for now)
  - [x] Show followers/following list
  - [x] Add profile view UI (F9 key)
  - [ ] Test profile view

- [ ] **Settings (username, PFP, discovery)**
  - [x] Create settings data model
  - [x] Add username change (store in BadgerDB)
  - [x] Add ASCII PFP regeneration
  - [x] Add LAN discovery toggle
  - [x] Add settings UI (F10 key)
  - [ ] Test settings

- [ ] **Export/backup identity**
  - [ ] Implement keypair export to file
  - [ ] Add password encryption for export
  - [ ] Implement import from backup
  - [ ] Add backup/restore UI
  - [ ] Test backup/restore

- [ ] **Web UI**
  - [ ] Implement HTTP handlers for state machine
  - [ ] Create HTML templates
  - [ ] Serve from same process
  - [ ] Share scuttlego service
  - [ ] Test web UI

---

## 8. [Testing Strategy](#8-testing-strategy)

### Test Pyramid

```
        /\
        /  \        E2E Tests (TUI mock)
        /────────\   3-5 tests, 30s to run
       /────────────\   / 15-20 tests, 5s to run
        /              \
        /               \    Unit Tests (Functions)
        /────────────\   / 50-70 tests, 1s to run
       /                \
        /                \    Integration Tests
        /────────────\   / 15-20 tests, 5s to run
       /                     \
        /                      \    E2E Tests (Real DB)
        /────────────────────\   / 15-20 tests, 5s to run
       /                           \
        /                            \
        /────────────────────\   / 3-5 tests, 30s to run
```

### Coverage Targets

| Module | Target | Critical Path |
|--------|---------------|---------|
| `scuttlego/service.go` | 95%+ | ✅ YES (service lifecycle) |
| `protocol/adapter.go` | 90%+ | ✅ YES (bridge) |
| `indexes/hashtags.go` | 90% | ✅ (search, filter) |
| `indexes/follows.go` | 90%+ | ✅ (follow graph) |
| `async/dispatcher.go` | 95%+ | ✅ YES (goroutines) |
| `ui/model.go` | 90%+ | |
| **Overall** | **85%** | **Critical 95%+** |

### Test Categories

**Unit Tests (70%):**
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

**Race Conditions:**
```bash
# ALWAYS run with -race
go test -race ./...
```

---

## 9. [Common Tasks](#9-common-tasks)

### Task 1: Create a New Package
```bash
# 1. Create directory
mkdir pkgname

# 2. Create files
touch pkgname/pkgname.go
touch pkgname/pkgname_test.go

# 3. Add imports and basic structure
```

### Task 2: Add a New Bubble Tea View
```go
func (m *SoClModel) View() string {
    return lipgloss.JoinVertical(
        lipgloss.Left,
        m.renderHeader(),
        m.renderFeed(),
        m.renderComposer(),
    )
}
```

### Task 3: Add Async Operation
```go
func (m *SoClModel) publishPost(text string) tea.Cmd {
    return func() tea.Msg {
        return PostPublishedMsg{Ref: msgRef}
    }
}
```

---

## 10. [API Quick Reference](#10-api-quick-reference)

### scuttlego Commands (Write Operations)

```go
// Publish a post
content := map[string]interface{}{
    "type": "post",
    "text": text,
}
contentJSON, _ := json.Marshal(content)
rawContent, _ := message.NewRawContent(contentJSON)
cmd, _ := commands.NewPublishRaw(rawContent.Bytes())
msgRef, err := svc.App.Commands.PublishRaw.Handle(cmd)
```

// Follow a peer
peerIdentity, err := refs.NewIdentity(feedRef)
cmd := commands.Follow{Target: peerIdentity}
err = svc.App.Commands.Follow.Handle(cmd)
```

### scuttlego Queries (Read Operations)

```go
// Get recent messages
startSeq, _ := common.NewReceiveLogSequence(0)
query, _ := queries.NewReceiveLog(startSeq, limit)
messages, err := svc.App.Queries.ReceiveLog.Handle(query)
```

---

## 11. [Decision Matrix](#11-decision-matrix)

### When to Code Autonomously ✅

| Scenario | Action | Example |
|----------|--------|---------|
| **UI rendering** | Code it | Add colors, format posts |
| **Validation logic** | Code it | Check 280-char limit |
| **Error handling** | Code it | Wrap errors, log with context |
| **Performance optimization** | Code it | Cache, goroutines for I/O |
| **Testing** | Code it | Write unit/integration tests |
| **Documentation** | Code it | Godoc, comments, README |

### When to Ask for Approval ⏸️️

| Scenario | Why Ask |
|----------|--------|---------|
| **Protocol Changes** | Modify EBT, handshake, message format | Affects all peers, breaking change |
| **Schema Changes** | Add/modify scuttlego internal schema | Breaks backwards compatibility |
| **Security Decisions** | Rate limiting, access control, crypto params | Human intent required |
| **New Dependencies** | Add external Go library | Supply chain + reproducibility |
| **Performance Tradeoffs** | Cache entire feed, disable security | Long-term impact |
| **scuttlego Config** | Change replication scheduler, EBT params | Complex internals |

---

## 12. [Debugging Guidelines](#12-debugging-guidelines)

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

**Issue: High memory usage (>200MB)**
**Diagnosis:**
```bash
# Check goroutine count
fmt.Println("Goroutines:", runtime.NumGoroutine())

# Check BadgerDB cache size
# (scuttlego manages this internally)

# Check feed size
fmt.Println("Feed size:", len(m.feed))
```

**Fix:**
- Implement pagination (max 100 posts in memory)
- Reduce BadgerDB cache size (scuttlego config)
- Clear old feed entries

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

## 13. [Pause & Continue](#13-pause--continue)

### How to Pause

**1. Save current state:**
```bash
# Commit current work
git add .
git commit -m "WIP: [task name]"

# Create checkpoint branch
git branch checkpoint-$(date +%Y%m%d-%H%M%S)
```

**2. Update agent.md:**
- Update "Session State" (above)
- Note blockers
- Note next priority

**3. Clean up:**
```bash
# Stop running processes
pkill so_cl

# Clear temporary files
rm -rf /tmp/so_cl-*
```

### How to Continue

**1. Read agent.md** (this file) - Refresh context
2. **Check TODO list** - Pick next task
3. **Review recent commits** - Understand state

**2. Review branch status:**
```bash
# Check current branch
git status

# Review recent commits
git log --oneline -10

# Review uncommitted changes
git diff HEAD~1
```

---

## 14. [Success Metrics](#14-success-metrics)

### Code Quality

| Metric | Target | How to Measure |
|--------|---------------|---------|
| **Test Coverage** | ≥85% (critical 95%+) | `go test -cover ./...` |
| **Race Conditions** | 0 | `go test -race ./...` |
| **Lint Errors** | 0 | `golangci-lint run` |
| **Build Success** | 100% | `CGO_ENABLED=0 go build` |
| **Binary Size** | ≤50MB | `ls -lh so_cl` |

### User Experience

| Metric | Target | How to Measure |
|--------|---------------|---------|
| **UI Responsiveness** | ≥60fps | Manual testing |
| **Startup Time** | ≤5s | `time ./so_cl` |
| **Post Publish Time** | ≤1s | Manual testing |

---

## 15. [Quick Reference Card](#15-quick-reference)

### Essential Commands

```bash
# Development
go mod tidy                    # Update dependencies
go test -race -cover ./...    # Test with coverage
go vet ./...                       # Static analysis
golangci-lint run                # Linting
CGO_ENABLED=0 go build             # Static build
go test -race -cover ./...    # Race detection
```

### Key Files

| File | Purpose |
|------|---------|
| `agent.md` | THIS FILE - AI agent reference |
| `production.md` | Full product specification |
| `go.mod` | Go dependencies |
| `scuttlego/service.go` | scuttlego service wrapper |
| `ui/model.go` | Bubble Tea model |
| `indexes/hashtags.go` | Hashtag counting & search |
| `indexes/follows.go` | Follow graph tracking |
| `core/types.go` | SoClProfile, settings types |

---

## **EMERGENCY CONTACT**

If you encounter critical issues:

1. **Check logs** (`~/.so_cl/log.txt`)
2. **Review this file** (agent.md)
3. **Read production.md** (full spec)
4. **Check scuttlego docs** (github.com/planetary-social/scuttlego)
5. **Pause and document** (Section 13)

**Remember**: This is a **living document**. Update it as you learn and progress. Good luck! 🚀

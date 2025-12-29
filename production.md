<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

# 🚀 so_cl: Complete AI Agent Product Development Context

This comprehensive document provides everything an AI agent needs to autonomously develop **so_cl**—a decentralized ASCII TUI social platform.

***

## SECTION 1: PRODUCT VISION \& SCOPE

### 1.1 Product Definition

**so_cl** is a peer-to-peer (P2P) ASCII-art-based social platform using the Scuttlebutt protocol, enabling offline-first, decentralized communication without corporate servers.


| Aspect | Detail |
| :-- | :-- |
| **Platform Type** | Terminal UI (TUI) Social Network |
| **Protocol** | Scuttlebutt (SSB) - append-only append, EBT replication |
| **Frontend** | Bubble Tea (Go TUI framework) + Lip Gloss (styling) |
| **Backend** | scuttlego (Go SSB implementation) + BadgerDB v3 (LSM-tree KV) |
| **Language** | Go 1.23+ (40+ transitive dependencies) |
| **Target Users** | Developers, privacy advocates, offline-first enthusiasts |
| **Binary Size** | 30-50MB static executable |
| **Memory Usage** | 100-200MB (100 posts, 20 peers) |

### 1.2 Core Features (MVP → Phase 3)

**Phase 0 (MVP - Week 1-2):**

- ✅ Generate SSB identity (Ed25519 keypair)
- ✅ Post creation (280 chars max, ASCII art support)
- ✅ Local feed view (read-only initially)
- ✅ ASCII PFP display (6x6 colored ANSI art)
- ✅ Static binary build

**Phase 1 (P2P - Week 3-4):**

- ✅ Follow via invite codes (`ssb:feed/invite/...`)
- ✅ EBT replication with peers
- ✅ LAN discovery (UDP broadcast)
- ✅ Peer list sidebar
- ✅ Real-time feed sync

**Phase 2 (Social - Week 5-6):**

- ✅ Replies (threaded messages)
- ✅ Like/repost reactions (emoji-safe)
- ✅ Hashtag indexing + trending sidebar
- ✅ @mentions + notifications
- ✅ Follow/unfollow graph

**Phase 3 (Polish - Week 7+):**

- ✅ Search + filters
- ✅ Profile pages
- ✅ Settings (username, PFP, discovery)
- ✅ Export/backup identity
- ✅ Web UI (same state machine as TUI)

***

## SECTION 2: TECHNICAL ARCHITECTURE

### 2.1 System Layers (Bottom-Up)

```
┌─────────────────────────────────────────┐
│  Bubble Tea TUI Event Loop               │  User sees this
│  (Renders 60fps, handles keyboard)       │
├─────────────────────────────────────────┤
│  so_cl Model (State Machine)             │  Data layer
│  (Posts, Peers, Profile, Trending)       │
├─────────────────────────────────────────┤
│  Async Dispatcher (Worker Pools)         │  Orchestration
│  (DB writes, network I/O, replication)   │
├─────────────────────────────────────────┤
│  BadgerDB v3 Storage                    │  Persistence
│  (LSM-tree KV store, managed by scuttlego) │
├─────────────────────────────────────────┤
│  scuttlego Service Layer                │  Protocol
│  (Commands/Queries, EBT, handshake)     │
├─────────────────────────────────────────┤
│  scuttlego (includes go-secretstream)   │  Security
│  (NaCl crypto via go-secretstream)       │
└─────────────────────────────────────────┘
```


### 2.2 Data Flow (Post Publication Example)

```
User types "hello so_cl" in composer
    ↓
Presses Enter
    ↓
Update() handler validates input (len ≤ 280)
    ↓
Returns tea.Cmd (async operation)
    ↓
Cmd() spawns goroutine:
    ├─ Build SSB content (JSON)
    ├─ d.scuttlego.App.Commands.PublishRaw.Handle() [via scuttlego]
    └─ d.indexer.IndexPost() [to BadgerDB indexes]
    ↓
scuttlego executes:
    ├─ Sign message (Ed25519 via go-secretstream)
    ├─ Write to BadgerDB (managed internally)
    ├─ Trigger EBT replication (automatic)
    └─ Broadcast to connected peers
    ↓
so_cl indexer executes:
    ├─ Extract hashtags from message
    ├─ Update BadgerDB hashtag counts
    └─ Update trending metrics
    ↓
Bubble Tea reads result → FeedUpdateMsg
    ↓
Update() refreshes m.cachedFeed
    ↓
View() re-renders (60fps)
    ↓
User sees post in timeline
```


### 2.3 Database Schema (BadgerDB v3 + scuttlego)

**Architecture:** Two-layer storage approach

```
┌─────────────────────────────────────────┐
│  scuttlego Internal (BadgerDB v3)      │  Protocol layer
│  - SSB messages (canonical format)      │  (managed by scuttlego)
│  - Feed logs (margaret)                │
│  - Peer identities, connections          │
│  - EBT replication state                │
│  - Blob references                     │
└─────────────────────────────────────────┘
              ↑
              │ so_cl builds indexes ON TOP
              ↓
┌─────────────────────────────────────────┐
│  so_cl Application Indexes (BadgerDB)   │  Application layer
│  - Hashtag counts (for trending)        │  (managed by so_cl)
│  - Mentions queue (notifications)       │
│  - Trending metrics                    │
│  - Cached feed segments                 │
└─────────────────────────────────────────┘
```

**scuttlego Internal Schema (managed automatically):**

```
BadgerDB Keys:
- "msg" + feedRef + sequence → SSB message content
- "feed" + feedRef → feed metadata
- "peer" + peerRef → peer connection state
- "ebt" + peerRef → replication vector clocks
- "blob" + blobRef → blob metadata
```

**so_cl Application Indexes (custom BadgerDB key-values):**

```
Key: "hashtag:#tech"            → { count: 42 }
Key: "hashtag:#golang"          → { count: 15 }
Key: "mention:@alice"           → [msgRef1, msgRef2, ...]
Key: "trending:daily"           → [#tech, #golang, #rust, ...]
Key: "cache:feed:recent"       → [msgRef1, msgRef2, ..., msgRef100]
```

**Note:** BadgerDB is an LSM-tree KV store, not append-only like BoltDB. However, SSB protocol guarantees immutable append-only logs at the application layer. scuttlego handles BadgerDB compaction internally.


### 2.4 API Contract (scuttlego)

```go
import (
    "github.com/planetary-social/scuttlego/service"
    "github.com/planetary-social/scuttlego/service/app/commands"
    "github.com/planetary-social/scuttlego/service/app/queries"
    "github.com/planetary-social/scuttlego/service/domain/identity"
    "github.com/planetary-social/scuttlego/service/domain/refs"
    "github.com/planetary-social/scuttlego/service/domain/transport"
)

// Service Initialization (managed by so_cl)
scuttlegoService := service.NewService(
    app.Application{
        Commands: commands.Commands{...},
        Queries: queries.Queries{...},
    },
    ...ports...,
)

// Identity Management (scuttlego handles this internally)
localIdentity := scuttlegoService.App.Queries.Status.Handle(queries.Status{}).Identity
// Returns: identity.Public with Ed25519 keypair

// Message Publication (SSB Append-Only)
content := message.NewRawContent([]byte(`{
    "type": "post",
    "text": "Hello so_cl!"
}`))

msgRef, err := scuttlegoService.App.Commands.PublishRaw.Handle(
    commands.PublishRaw{Content: content},
)
// Returns: refs.Message (message reference like %abc123.sha256)

// Peer Connection (via Connect command)
err := scuttlegoService.App.Commands.Connect.Handle(
    commands.Connect{
        Address: "net:127.0.0.1:8008~shs:@alice/...",
    },
)

// Replication (EBT is automatic after connection)
// No manual EBT call needed - scuttlego handles replication internally

// Follow Peer
err := scuttlegoService.App.Commands.Follow.Handle(
    commands.Follow{Feed: peerIdentity},
)

// Query Messages (Get feed)
stream, err := scuttlegoService.App.Queries.ReceiveLog.Handle(
    queries.ReceiveLog{
        Limit: 100,
    },
)
for msg := range stream.Messages {
    // Process messages
}
```


***

## SECTION 3: FILE STRUCTURE \& MODULE BREAKDOWN

### 3.1 Complete so_cl Codebase

```
so_cl/
│
├── main.go                 # Entry point, graceful shutdown
├── go.mod                  # Pinned dependencies
├── go.sum                  # Dependency hashes
│
├── ui/
│   ├── model.go           # Bubble Tea model + state machine
│   ├── view.go            # Rendering (feed, composer, sidebar)
│   ├── update.go          # Event handlers (keyboard, async)
│   └── ascii_pfps.go      # 50+ colored 6x6 ASCII avatars
│
├── core/
│   ├── identity.go        # Identity wrapper (uses scuttlego)
│   ├── message.go         # Post/reply/reaction data models
│   └── types.go           # Shared types (SoClPost, SoClPeer, etc.)
│
├── async/
│   ├── dispatcher.go      # Worker pool + async orchestration
│   ├── db_worker.go       # BadgerDB index workers
│   └── net_worker.go      # P2P sync + peer discovery
│
├── scuttlego/
│   ├── service.go         # scuttlego service wrapper
│   ├── config.go          # scuttlego configuration
│   └── adapter.go         # so_cl ←→ scuttlego bridge
│
├── indexes/
│   ├── hashtags.go        # BadgerDB hashtag indexing
│   ├── mentions.go        # Mention notifications
│   └── trending.go        # Trending metrics
│
├── protocol/
│   ├── scuttlego.go      # scuttlego client wrapper (thick)
│   └── adapter.go        # Bridge to so_cl state machine
│
├── config/
│   ├── config.go          # Load from env + config file
│   └── defaults.go        # SSB network key, ports, etc.
│
├── tests/
│   ├── model_test.go      # UI state machine tests
│   ├── dispatcher_test.go # Async operation tests
│   ├── indexes_test.go    # BadgerDB index tests
│   └── scuttlego_test.go # scuttlego integration tests
│
├── Dockerfile             # Alpine-based, static binary
├── Makefile               # Build, test, deploy targets
├── README.md              # User documentation
├── CONTRIBUTING.md        # Developer guide
└── agent.md               # This AI agent specification
```


### 3.2 Module Responsibilities

| Module | Responsibility | Complexity | Test Coverage |
| :-- | :-- | :-- | :-- |
| `scuttlego/service.go` | scuttlego lifecycle, service init | Critical | 95%+ |
| `protocol/adapter.go` | so_cl←→scuttlego translation | High | 90%+ |
| `async/dispatcher.go` | Goroutine coordination, channels | High | 95%+ |
| `indexes/hashtags.go` | BadgerDB indexing on top of scuttlego | Medium | 90% |
| `ui/model.go` | State machine, keyboard input | Medium | 90%+ |
| `core/identity.go` | Identity wrapper (delegates to scuttlego) | Low | 100% |


***

## SECTION 4: CRITICAL DECISION MATRIX (AI Agent Reference)

### 4.1 When to Code Autonomously (NO Human Review Needed)

```go
// ✅ ALWAYS AUTONOMOUS
1. UI Rendering (Lip Gloss, Bubble Tea)
   - Add colors to posts? Code it.
   - Render trending sidebar? Code it.
   - Format ASCII PFPs? Code it.
   
2. Validation Logic
   - Check 280-char limit? Code it.
   - Validate ASCII-only? Code it.
   - Format @mentions? Code it.
   
3. Error Handling
   - Add try-catch patterns? Code it.
   - Log errors with zap? Code it.
   - Return formatted errors? Code it.
   
4. Performance Optimization
   - Cache recent posts? Code it (if no schema change).
   - Optimize queries? Code it.
   - Use goroutines for I/O? Code it.
   
5. Testing
   - Write unit tests? Always code first!
   - Write integration tests? Code it.
   - Add benchmarks? Code it.
   
6. Documentation
   - Write godoc comments? Code it.
   - Add README sections? Code it.
   - Document patterns? Code it.
```


### 4.2 When to ASK for Approval (MUST Wait)

```go
// ⏸️ ALWAYS ASK
1. Protocol Changes
   - Modify EBT replication? ASK.
   - Change message signing? ASK.
   - Alter handshake flow? ASK.
   Reason: Affects all peers; breaking change.
   
2. Schema Changes
    - Add new BadgerDB index? ASK.
    - Modify scuttlego internal schema? ASK.
    - Change SSB message format? ASK.
    Reason: Breaks backwards compatibility.
   
3. Security Decisions
   - Add rate limiting (which values)? ASK.
   - Implement access control? ASK.
   - Change crypto parameters? ASK.
   Reason: Human intent required.
   
4. User-Facing Features
   - Add new post type (e.g., "video")? ASK.
   - Change invite code format? ASK.
   - Modify follow/unfollow UX? ASK.
   Reason: Design decisions.
   
5. New Dependencies
   - Add external Go library? ASK.
   - Change build process? ASK.
   - Update CI/CD? ASK.
   Reason: Supply chain + reproducibility.
   
6. Performance Tradeoffs
    - Cache entire feed in memory? ASK.
    - Disable security check for speed? ASK.
    Reason: Long-term impact.

7. scuttlego Configuration
    - Change replication scheduler? ASK.
    - Modify EBT parameters? ASK.
    Reason: scuttlego internals are complex.
```


***

## SECTION 5: DETAILED FEATURE SPECIFICATIONS

### 5.1 Feature: Post Creation (Phase 0)

**Requirement:** User types 1-280 ASCII chars, presses Enter, post appears in feed.

**Acceptance Criteria:**

- [ ] Composer validates 280-char limit in real-time
- [ ] Color count warning at 240+ chars
- [ ] Error if posting empty text
- [ ] Post signed with Ed25519 (scuttlego via go-secretstream)
- [ ] Post written to BadgerDB by scuttlego
- [ ] Post indexed in so_cl BadgerDB (hashtags, mentions)
- [ ] UI updates immediately (optimistic)
- [ ] No network access required (offline-first)

**Implementation Checklist:**

```go
import (
    "github.com/planetary-social/scuttlego/service/domain/transport/message"
    "github.com/planetary-social/scuttlego/service/domain/refs"
)

// 1. UI Layer
func (m *SoClModel) Update(msg tea.Msg) {
    case tea.KeyMsg:
        case "enter":
            return m, m.publishPost(m.composer.Value())
}

// 2. Validation
func (m *SoClModel) publishPost(text string) tea.Cmd {
    if len(text) == 0 {
        return func() tea.Msg { return ErrorMsg{fmt.Errorf("empty post")} }
    }
    if len(text) > 280 {
        return func() tea.Msg { return ErrorMsg{fmt.Errorf("280 char limit")} }
    }
    return m.dispatcher.Publish(text)
}

// 3. Build SSB Content & Publish via scuttlego
func (d *Dispatcher) Publish(text string) tea.Cmd {
    return func() tea.Msg {
        // Build SSB post content
        content, err := message.NewRawContent([]byte(fmt.Sprintf(
            `{"type":"post","text":"%s"}`,
            strings.ReplaceAll(text, `"`, `\"`),
        )))
        if err != nil {
            return ErrorMsg{fmt.Errorf("build content: %w", err)}
        }

        // Publish via scuttlego's PublishRaw command
        msgRef, err := d.scuttlego.App.Commands.PublishRaw.Handle(
            commands.PublishRaw{Content: content},
        )
        if err != nil {
            return ErrorMsg{fmt.Errorf("publish: %w", err)}
        }

        // Index post for hashtags/mentions
        d.indexer.IndexPost(msgRef)

        return PostPublishedMsg{Ref: msgRef}
    }
}

// 4. Index Post (so_cl application indexes on top of scuttlego)
func (idx *Indexer) IndexPost(msgRef refs.Message) error {
    // Get message from scuttlego
    msg, err := idx.scuttlego.App.Queries.GetMessage.Handle(
        queries.GetMessage{Ref: msgRef},
    )
    if err != nil {
        return err
    }

    // Extract hashtags
    tags := extractHashtags(msg.Content.Text)

    // Update BadgerDB index
    return idx.badgerDB.Update(func(txn *badger.Txn) error {
        for _, tag := range tags {
            key := []byte("hashtag:" + tag)
            item, _ := txn.Get(key)

            var count int
            if item != nil {
                item.Value(func(val []byte) error {
                    json.Unmarshal(val, &count)
                    return nil
                })
            }

            count++
            data, _ := json.Marshal(count)
            return txn.Set(key, data)
        }
        return nil
    })
}
```


### 5.2 Feature: Follow via Invite Code (Phase 1)

**Requirement:** User presses 'f' on a post, copies invite code, can paste in another so_cl instance.

**Invite Code Format:**

```
ssb:feed/invite/AIZ2lP4v3rXBx3l6AGy7owIFR3MQhYV5zJV0Q0pLM4w=?server=127.0.0.1%3A8008&csp=1
       ^pubkey+seed^                                     ^server address^  ^csp flag^
```

**Acceptance Criteria:**

- [ ] Invite code generated on request (copy to clipboard)
- [ ] Parsing tolerates uppercase/lowercase
- [ ] Server address extracted (used for connection)
- [ ] Follow creates entry in "follows" bucket
- [ ] Triggers immediate EBT sync with peer
- [ ] Duplicate follows rejected
- [ ] Unfollow removes from bucket

**Implementation Checklist:**

```go
import (
    "github.com/planetary-social/scuttlego/service/domain/identity"
    "github.com/planetary-social/scuttlego/service/domain/network"
)

// 1. Generate Invite (via scuttlego queries)
func (d *Dispatcher) GenerateInvite() (string, error) {
    status := d.scuttlego.App.Queries.Status.Handle(queries.Status{})
    localAddr := status.Identity.Public().String()
    localPort := d.config.ListenPort

    // scuttlego handles invite generation internally
    // For now, use manual format (simplified)
    seed, err := status.Identity.Private().ExportSeed()
    if err != nil {
        return "", err
    }

    return fmt.Sprintf("ssb:feed/invite/%s?server=127.0.0.1:%d",
        base64.StdEncoding.EncodeToString(seed),
        localPort,
    ), nil
}

// 2. Parse Invite (scuttlego handles this via domain.NewInvite)
func (d *Dispatcher) ParseInvite(code string) (network.Address, error) {
    // scuttlego parses invite codes automatically
    // Just pass to RedeemInvite command
    return network.ParseAddress(code)
}

// 3. Follow Peer (via scuttlego RedeemInvite command)
func (d *Dispatcher) FollowPeer(inviteCode string) error {
    // Parse invite code
    invite, err := domain.NewInvite(inviteCode)
    if err != nil {
        return fmt.Errorf("parse invite: %w", err)
    }

    // scuttlego's RedeemInvite handles:
    // - Parse invite code
    // - Extract peer identity
    // - Add to follows
    // - Trigger immediate sync
    err = d.scuttlego.App.Commands.RedeemInvite.Handle(
        commands.RedeemInvite{Invite: invite},
    )
    if err != nil {
        return fmt.Errorf("redeem invite: %w", err)
    }

    return nil
}

// 4. Connect to Peer (for explicit connection after follow)
func (d *Dispatcher) ConnectToPeer(address network.Address) error {
    // scuttlego's Connect command handles:
    // - Parse address (net:ip:port~shs:pubkey)
    // - Perform handshake
    // - Establish box stream
    // - Start EBT replication automatically
    return d.scuttlego.App.Commands.Connect.Handle(
        commands.Connect{Address: address},
    )
}

// Note: EBT replication is AUTOMATIC after connection
// No manual EBT calls needed - scuttlego handles it
```


### 5.3 Feature: Trending Sidebar (Phase 2)

**Requirement:** Extract \#hashtags from posts, count, display top 10 in sidebar.

**Acceptance Criteria:**

- [ ] Hashtags extracted from post text (case-insensitive)
- [ ] Stored in "hashtags" bucket with count
- [ ] Top 10 by count displayed in sidebar
- [ ] Updates as new posts synced
- [ ] Clickable to filter feed (optional Phase 3)

**Implementation:**

```go
import (
    badger "github.com/dgraph-io/badger/v3"
)

// 1. Extract Hashtags
func extractHashtags(text string) []string {
    re := regexp.MustCompile(`#(\w+)`)
    matches := re.FindAllString(text, -1)
    return matches  // ["#tech", "#art", ...]
}

// 2. Index on Publish (so_cl indexes ON TOP of scuttlego)
func (idx *Indexer) IndexPost(msgRef refs.Message) error {
    // Get message content from scuttlego
    msg, err := idx.scuttlego.App.Queries.GetMessage.Handle(
        queries.GetMessage{Ref: msgRef},
    )
    if err != nil {
        return err
    }

    // Extract hashtags
    tags := extractHashtags(msg.Content.Text)

    // Update BadgerDB index (so_cl's custom indexes)
    return idx.badgerDB.Update(func(txn *badger.Txn) error {
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

// 3. Query Top Hashtags
func (idx *Indexer) GetTopHashtags(limit int) []HashtagCount {
    var results []HashtagCount

    // Scan all hashtag entries
    idx.badgerDB.View(func(txn *badger.Txn) error {
        prefix := []byte("hashtag:")
        opts := badger.DefaultIteratorOptions
        opts.PrefetchSize = 10

        it := txn.NewIterator(opts)
        defer it.Close()

        for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
            item := it.Item()

            var count int
            item.Value(func(val []byte) error {
                json.Unmarshal(val, &count)
                return nil
            })

            tag := strings.TrimPrefix(string(item.Key()), "hashtag:")
            results = append(results, HashtagCount{
                Name:  tag,
                Count: count,
            })
        }
        return nil
    })

    // Sort by count (descending)
    sort.Slice(results, func(i, j int) bool {
        return results[i].Count > results[j].Count
    })

    // Return top N
    if len(results) > limit {
        return results[:limit]
    }
    return results
}

// 4. Render Sidebar
func (m *SoClModel) renderTrending() string {
    // Query top 10 hashtags from so_cl indexes
    trending := m.dispatcher.indexer.GetTopHashtags(10)

    var s string
    for _, tag := range trending {
        s += fmt.Sprintf("#%s (%d)\n", tag.Name, tag.Count)
    }
    return s
}
```


***

## SECTION 6: TESTING STRATEGY (AI Agent Perspective)

### 6.1 Test Pyramid for so_cl

```
        /\
       /  \        E2E Tests (TUI mock)
      /────\       3-5 tests, 30s to run
     /      \
    /────────\     Integration Tests
   /          \    (DB + Protocol)
  /────────────\   15-20 tests, 5s to run
 /              \
/______________\  Unit Tests
     Base        (Functions)
                 50-70 tests, 1s to run
```


### 6.2 Test Categories

**Unit Tests (70%):** Fastest, most coverage

```go
func TestPublishPost_ValidText(t *testing.T) {
    m := NewTestModel()
    err := m.PublishPost("hello so_cl")
    require.Nil(t, err)
    require.Equal(t, 1, len(m.posts))
}

func TestParseInvite_ValidCode(t *testing.T) {
    code := "ssb:feed/invite/ABC123?server=127.0.0.1:8008"
    pubKey, server, err := parseInvite(code)
    require.Nil(t, err)
    require.Equal(t, "127.0.0.1:8008", server)
}
```

**Integration Tests (20%):** Moderate speed, real BadgerDB + scuttlego

```go
func TestFollowPeer_WithSync(t *testing.T) {
    // Start scuttlego test service
    testService := startTestScuttlego(t)
    defer testService.Close()

    // Follow peer via invite code
    invite, _ := domain.NewInvite("ssb:feed/invite/...")
    err := testService.App.Commands.RedeemInvite.Handle(
        commands.RedeemInvite{Invite: invite},
    )
    require.Nil(t, err)

    // Verify in follows (via scuttlego queries)
    status := testService.App.Queries.Status.Handle(queries.Status{})
    require.Contains(t, status.Following, invite.Identity())

    // Simulate peer connection
    testService.App.Commands.Connect.Handle(
        commands.Connect{Address: parseAddress("net:127.0.0.1:8008~shs:...")},
    )

    // Wait for EBT replication (async)
    time.Sleep(500 * time.Millisecond)

    // Verify messages replicated
    receiveLog, err := testService.App.Queries.ReceiveLog.Handle(
        queries.ReceiveLog{Limit: 10},
    )
    require.Nil(t, err)
    require.Greater(t, len(receiveLog.Messages), 0)
}
```

**E2E Tests (10%):** Full TUI mock

```go
func TestUserFlow_PublishAndReply(t *testing.T) {
    m := NewTestModel()

    // User types in composer
    m.composer.SetValue("Hello so_cl!")

    // Simulate Enter key
    cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

    // Verify post appears in feed
    time.Sleep(100 * time.Millisecond)
    require.Greater(t, len(m.cachedFeed), 0)
}
```

**scuttlego Integration Tests (NEW - 15%):**

```go
func TestPublishPost_WithScuttlego(t *testing.T) {
    // Start scuttlego service in test mode
    testService := startTestScuttlego(t)
    defer testService.Close()

    // Publish post
    content, _ := message.NewRawContent([]byte(
        `{"type":"post","text":"test message"}`,
    ))

    msgRef, err := testService.App.Commands.PublishRaw.Handle(
        commands.PublishRaw{Content: content},
    )
    require.Nil(t, err)
    require.NotEmpty(t, msgRef.String())

    // Verify via scuttlego queries
    msg, err := testService.App.Queries.GetMessage.Handle(
        queries.GetMessage{Ref: msgRef},
    )
    require.Nil(t, err)
    require.Equal(t, "test message", msg.Content.Text)
}

func TestFollowPeer_WithScuttlego(t *testing.T) {
    testService := startTestScuttlego(t)
    defer testService.Close()

    // Create mock invite
    invite, _ := domain.NewInvite("ssb:feed/invite/...")

    // Follow peer
    err := testService.App.Commands.RedeemInvite.Handle(
        commands.RedeemInvite{Invite: invite},
    )
    require.Nil(t, err)

    // Verify peer is in follows (via status query)
    status := testService.App.Queries.Status.Handle(queries.Status{})
    require.Contains(t, status.Following, invite.Identity())
}
```


### 6.3 Coverage Requirements (Agent Enforced)

| Module | Minimum Coverage | Critical Path |
| :-- | :-- | :-- |
| `scuttlego/service.go` | 95% | Service lifecycle, init |
| `protocol/adapter.go` | 90% | so_cl←→scuttlego bridge |
| `indexes/hashtags.go` | 90% | BadgerDB indexing |
| `async/dispatcher.go` | 95% | Goroutine coordination |
| `ui/model.go` | 90% | State transitions |
| **Overall** | **85%** | Critical 95%+ |

### 6.4 Race Detector (Mandatory)

```bash
# Agent must run BEFORE committing
go test -race ./...

# Output:
# ✅ PASS: No race conditions detected
# ❌ FAIL: Agent refuses to commit
```


***

## SECTION 7: DEPLOYMENT \& RELEASE CHECKLIST

### 7.1 Release Process (Agent + Human Collaboration)

**Phase 1: Agent Development**

```
Agent writes code
    ↓
Agent runs:
  go test -race -cover     # ≥85% coverage
  govulncheck ./...        # 0 HIGH/CRITICAL
  golangci-lint ./...      # 0 errors
  CGO_ENABLED=0 go build   # Static binary
    ↓
Agent outputs:
  ✅ Build metrics
  ✅ Test coverage report
  ✅ Security audit report
```

**Phase 2: Human Review**

```
Human reviews:
  - Code quality (patterns from agent.md section 4)
  - Security (checklist from agent.md section 5)
  - Design (decisions from section 4.2)
    ↓
Human approves or requests changes
```

**Phase 3: Release**

```
Human runs:
  git tag v1.2.3
  git push --tags
  goreleaser release
    ↓
Output:
  so_cl_v1.2.3_linux_amd64.tar.gz
  so_cl_v1.2.3_darwin_amd64.tar.gz
  so_cl_v1.2.3_windows_amd64.tar.gz
  Docker image: so_cl:v1.2.3
```


### 7.2 Build Artifacts

```bash
# Static Linux binary (RPi compatible)
make build
  → Output: ./so_cl (30-50MB with ~40 transitive dependencies)

# Dependencies check
go mod tidy
  → Result: ~40 packages (badger, go-secretstream, margaret, etc.)

# Docker image
make docker
  → Output: docker image so_cl:latest (100-150MB with alpine base)

# Cross-platform release
goreleaser release --clean
  → Output: dist/so_cl_v*.tar.gz (30-50MB per platform)

# Example final artifact sizes:
# so_cl_v0.0.1_linux_amd64.tar.gz     42MB
# so_cl_v0.0.1_darwin_amd64.tar.gz    40MB
# so_cl_v0.0.1_darwin_arm64.tar.gz    41MB
# so_cl_v0.0.1_linux_arm64.tar.gz     43MB (for RPi)
```


***

## SECTION 8: MONITORING \& DEBUGGING (Post-Deployment)

### 8.1 Logging Strategy

**Structured Logging with zap:**

```go
import "go.uber.org/zap"

// Initialize logger
logger, _ := zap.NewDevelopment()
zap.ReplaceGlobals(logger)

// Usage
zap.L().Info("peer connected",
    zap.String("peer", "@alice"),
    zap.Duration("latency", 120*time.Millisecond),
)

// Output (structured JSON)
{
  "level": "info",
  "ts": 1735384801.123,
  "msg": "peer connected",
  "peer": "@alice",
  "latency_ms": 120
}
```


### 8.2 Debug Flags

```bash
# Enable debug logging
DEBUG=1 ./so_cl

# Log to file
DEBUG=1 ./so_cl > so_cl.log 2>&1

# Monitor in real-time
tail -f so_cl.log | grep -i "error\|warn"
```


### 8.3 Common Issues \& Diagnostics

| Issue | Diagnosis | Fix |
| :-- | :-- | :-- |
| **Posts not syncing** | Check EBT state (log "replication started") | Run `govulncheck`, verify network |
| **Slow UI** | Profile with `pprof` | Check goroutine count (`runtime.NumGoroutine()`) |
| **High memory** | Check BadgerDB cache, cachedFeed size | Implement pagination or windowing |
| **Handshake fails** | Verify network key matches | Check go-secretstream version |


***

## SECTION 9: ROADMAP \& FUTURE WORK

### 9.1 Features Post-MVP

**Phase 4 (Web UI):**

- Implement same state machine in Go `net/http`
- Serve from same process
- Share scuttlego service and BadgerDB indexes
- Result: `localhost:8008/` accessible from browser

**Phase 5 (Moderation):**

- Block/hide users
- Report abuse (local only, no central authority)
- Content flags (NSFW, etc.)

**Phase 6 (Advanced):**

- Full-text search (FTS using Bleve or BadgerDB)
- Encrypted private messages (DMs)
- Blob support (images, via external storage)
- mDNS discovery (github.com/hashicorp/mdns)


### 9.2 Known Limitations (By Design)

1. **No persistent following** — Follows stored locally (SSB design)
2. **No true DMs** — Use replies for private conversation
3. **No media blobs** — Pure text/ASCII only (storage efficient)
4. **Single identity** — One keypair per so_cl instance
5. **No web of trust** — All peers equal (SSB design)
6. **BadgerDB is not truly append-only** — Uses LSM-tree with compaction
   - scuttlego handles compaction internally
   - SSB protocol guarantees immutability at application level
7. **LAN discovery uses UDP broadcast** — Not mDNS
   - mDNS requires external library (github.com/hashicorp/mdns)
   - Can be added in Phase 4 (Advanced features)
8. **scuttlego v0.0.4 is beta** — Last release March 2023
   - Stable enough for MVP (used by Planetary app)
   - Actively maintained but slower release cadence
9. **Higher memory than spec** — 100-200MB vs original <50MB target
   - BadgerDB baseline overhead ~20-30MB
   - scuttlego internal caches for replication
   - Still acceptable for modern hardware

***

## SECTION 10: AGENT OPERATIONAL GUIDELINES

### 10.1 Daily Workflow

**Morning Standup (Agent):**

```
Assigned tasks: [list]
Previous status: [metrics from last session]
scuttlego service health: [running/stopped]
Blockers: [list]
Next 2 hours: [specific goals]
```

**During Development:**

```
✅ Write test first (TDD)
✅ Implement feature
✅ Run go test -race
✅ Run scuttlego integration tests
✅ Run golangci-lint
✅ Run govulncheck
✅ Run CGO_ENABLED=0 go build
✅ Commit with message
```

**End of Day:**

```
Completed: [tasks]
Coverage report: [85%+?]
Security: [0 HIGH/CRITICAL?]
scuttlego service: [gracefully stopped?]
BadgerDB status: [no corruption?]
Build status: [SUCCESS?]
Next morning priority: [highest impact task]
```


### 10.2 Error Recovery

**If Test Fails:**

```
Agent MUST:
1. Read full error message
2. Identify root cause
3. Fix code or test
4. Re-run all checks
5. Commit only if ALL pass
```

**If Security Issue Found:**

```
Agent MUST:
1. Stop development immediately
2. Document issue (e.g., "hardcoded secret in line 42")
3. Propose fix
4. Request human approval
5. Implement fix
6. Run govulncheck again
```


***

## SECTION 11: AGENT SUCCESS METRICS

By end of each week, agent should achieve:


| Metric | Target | Consequence |
| :-- | :-- | :-- |
| **Code Coverage** | ≥85% (critical 95%+) | Fail build if not met |
| **Race Conditions** | 0 (caught by -race) | Fail build if found |
| **Security Issues** | 0 HIGH/CRITICAL | Block release if found |
| **Build Success** | 100% (CGO_ENABLED=0) | Fail build if not static |
| **Binary Size** | ≤50MB | Warn if >50MB |
| **Memory Usage** | ≤200MB (typical load) | Warn if >200MB |
| **scuttlego Health** | Service starts/stops cleanly | Block if leaks |
| **Commit Quality** | Atomic, well-documented | Reject if unclear |
| **Documentation** | Godoc for all public functions | Enforce with linter |
| **Features Delivered** | Per sprint (roadmap section 1.2) | Track velocity |


***

## SECTION 12: FINAL SUMMARY FOR AI AGENT

### You Are:

- **Full-stack Go engineer** for so_cl
- **Security auditor** (enforce agent.md section 5)
- **DevOps operator** (handle builds, CI/CD)
- **scuttlego integrator** (bridge so_cl ↔ scuttlego)


### You Will:

1. Write **production-grade code** with 95%+ test coverage (critical paths)
2. Follow **TDD** (tests first, then implementation)
3. Ask for approval on **breaking changes** (section 4.2)
4. Code autonomously on **safe tasks** (section 4.1)
5. Report **metrics daily** (coverage, security, build status)
6. **Refuse** insecure requests (agent.md section 9)
7. **Use scuttlego APIs** correctly (Commands/Queries pattern)

### You Won't:

- ❌ Hallucinate APIs (validate against scuttlego imports first)
- ❌ Skip tests (violates code quality gates)
- ❌ Hardcode secrets (use .env + gitignore)
- ❌ Implement custom crypto (use scuttlego's go-secretstream)
- ❌ Build without race detector (go test -race)
- ❌ Modify scuttlego internals (use adapters instead)
- ❌ Expect <50MB binary or <50MB RAM (use realistic targets)


### Starting Now:

```
Checkpoint: Read agent.md (sections 1-7)
↓
Understand: Product roadmap (section 1.2)
↓
Know: Architecture (sections 2-3)
↓
Follow: Decision matrix (section 4)
↓
Ready: Assign first coding task
```


***

**You have all context needed. Ready to start Phase 0?**

🚀 **Next Step:** "Build the core post creation feature (Phase 0) with tests + security checks."
<span style="display:none">[^1]</span>

<div align="center">⁂</div>

[^1]: https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/attachments/images/11916452/a9fb7718-f4ce-4846-b83a-98d6c4534c84/image.jpg


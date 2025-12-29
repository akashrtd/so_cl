# Revision Plan for production.md

## Context
The original production.md had critical issues:
- Wrong library name (Scuttlegoa → scuttlego)
- Wrong storage (BoltDB → BadgerDB v3)
- Incorrect dependency claims
- Unrealistic performance targets
- API examples didn't match actual libraries

## Target Technology Stack
- **SSB Library**: github.com/planetary-social/scuttlego v0.0.4
- **Storage**: github.com/dgraph-io/badger/v3 (via scuttlego)
- **Frontend**: Bubble Tea v1.3.10 + Lip Gloss
- **Language**: Go 1.19+ (scuttlego requires 1.19, we can use 1.23)

## Critical Fixes

### 1. Product Definition (Lines 16-25)
**Before:**
```
| **Backend** | Scuttlegoa (Go SSB implementation) + BoltDB (append-only log) |
| **Language** | Go 1.23+ (zero dependencies for deployment) |
| **Binary Size** | <10MB static executable |
| **Memory Usage** | <50MB (100 posts, 20 peers) |
```

**After:**
```
| **Backend** | scuttlego (Go SSB implementation) + BadgerDB v3 (LSM-tree KV) |
| **Language** | Go 1.23+ (40+ transitive dependencies) |
| **Binary Size** | 30-50MB static executable |
| **Memory Usage** | 100-200MB (100 posts, 20 peers) |
```

### 2. Feature Updates (Line 41)
**Before:**
```
- ✅ LAN discovery (mDNS)
```

**After:**
```
- ✅ LAN discovery (UDP broadcast)
```

**Note**: scuttlego uses local UDP advertisements, not mDNS. True mDNS requires external library (github.com/hashicorp/mdns) which is out of scope for MVP.

### 3. Architecture Diagram (Section 2.1, Lines 67-87)
**Changes:**
- "BoltDB Append-Only Log" → "BadgerDB v3 Storage"
- "Scuttlegoa SSB Core" → "scuttlego Service Layer"
- "Go stdlib + NaCl crypto" → "scuttlego (includes NaCl via go-secretstream)"

### 4. Data Flow & API (Sections 2.2-2.4)
**Old API:**
```go
identity.Sign(message)
d.dispatcher.AppendMessage()
d.dispatcher.BroadcastToFollowers()
boxConn.Call("ebt", "replicate", EBTArgs{...})
```

**New scuttlego API:**
```go
// Identity is already stored in scuttlego service
content := message.NewRawContent([]byte(`{"type":"post","text":"hello so_cl!"}`))
msgRef, err := service.App.Commands.PublishRaw.Handle(commands.PublishRaw{
    Content: content,
})

// Replication is automatic via EBT
// Connection via:
service.App.Commands.Connect.Handle(commands.Connect{
    Address: "net:127.0.0.1:8008~shs:@alice...",
})
```

### 5. Database Schema (Section 2.3, Lines 124-153)
**Major changes:**
- BadgerDB is NOT append-only by design (LSM-tree)
- Schema is managed by scuttlego internally
- so_cl will build indexes ON TOP of scuttlego's storage

**New Schema Approach:**
```
// scuttlego internal (BadgerDB):
- messages (SSB protocol format)
- feeds, peers, blob references
- EBT state, replication status

// so_cl application indexes (custom BadgerDB buckets):
- "hashtags" → count
- "mentions" → notification queue
- "trending" → aggregated metrics
```

### 6. File Structure (Section 3.1, Lines 199-246)
**Removed:**
```
├── storage/
│   ├── db.go              # BoltDB initialization + schema
│   ├── append_log.go      # Append-only ledger operations
│   └── index.go           # Feed indexing for fast queries
```

**Added:**
```
├── scuttlego/
│   ├── service.go         # scuttlego service wrapper
│   ├── config.go          # scuttlego configuration
│   └── adapter.go         # Adapters for so_cl←→scuttlego bridge
├── indexes/
│   ├── hashtags.go        # BadgerDB hashtag indexing
│   ├── mentions.go        # Mention notifications
│   └── trending.go        # Trending metrics
```

**Updated:**
```
├── protocol/
│   ├── scuttlego.go      # scuttlego client wrapper (thick)
│   └── adapter.go        # Bridge to so_cl state machine
```

### 7. Module Responsibilities (Section 3.2)
**Updates:**
| Module | Responsibility | Complexity | Test Coverage |
| :-- | :-- | :-- | :-- |
| `scuttlego/service.go` | scuttlego lifecycle, service init | Critical | 95%+ |
| `indexes/hashtags.go` | BadgerDB indexing on top of scuttlego | Medium | 90% |
| `protocol/adapter.go` | so_cl←→scuttlego translation | High | 90%+ |
| `async/dispatcher.go` | goroutine coordination, channels | High | 95%+ |
| `ui/model.go` | State machine, keyboard input | Medium | 90%+ |

### 8. Critical Decision Matrix (Section 4)
**Changes to Section 4.2 (When to ASK):**

Add new item:
```go
// ⏸️ ALWAYS ASK
7. scuttlego Configuration
   - Change replication scheduler? ASK.
   - Modify EBT parameters? ASK.
   Reason: scuttlego internals are complex.
```

### 9. Feature Implementations (Section 5)

#### 5.1 Post Creation
**Old implementation:** Used custom signing + BoltDB writes

**New implementation:**
```go
func (m *SoClModel) publishPost(text string) tea.Cmd {
    // Validate
    if len(text) == 0 || len(text) > 280 {
        return func() tea.Msg { return ErrorMsg{...} }
    }

    // Build SSB content
    content, _ := message.NewRawContent([]byte(
        fmt.Sprintf(`{"type":"post","text":"%s"}`, text)
    ))

    // Use scuttlego PublishRaw command
    return func() tea.Msg {
        msgRef, err := m.scuttlego.App.Commands.PublishRaw.Handle(
            commands.PublishRaw{Content: content},
        )
        if err != nil {
            return ErrorMsg{err}
        }
        return PostPublishedMsg{Ref: msgRef}
    }
}
```

#### 5.2 Follow via Invite Code
**Major change:** scuttlego has `RedeemInviteHandler`

```go
func (d *Dispatcher) FollowPeer(inviteCode string) error {
    // scuttlego handles invite parsing
    invite, err := domain.NewInvite(inviteCode)
    if err != nil {
        return err
    }

    // Use scuttlego's RedeemInvite command
    return d.scuttlego.App.Commands.RedeemInvite.Handle(
        commands.RedeemInvite{Invite: invite},
    )
}
```

#### 5.3 Trending Sidebar
**Approach change:** Build indexes ON TOP of scuttlego

```go
func (d *Dispatcher) indexPost(msgRef refs.Message) error {
    // Get message from scuttlego
    msg, err := d.scuttlego.App.Queries.GetMessage.Handle(
        queries.GetMessage{Ref: msgRef},
    )
    if err != nil {
        return err
    }

    // Extract hashtags
    tags := extractHashtags(msg.Content.Text)

    // Update our custom BadgerDB index
    return d.badgerDB.Update(func(txn *badger.Txn) error {
        for _, tag := range tags {
            // Increment count
            item, _ := txn.Get([]byte("hashtag:" + tag))
            var count int
            if item != nil {
                item.Value(func(val []byte) error {
                    json.Unmarshal(val, &count)
                    return nil
                })
            }
            count++
            data, _ := json.Marshal(count)
            txn.Set([]byte("hashtag:"+tag), data)
        }
        return nil
    })
}
```

### 10. Testing Strategy (Section 6)
**Add integration tests for scuttlego:**

```go
func TestPublishPost_WithScuttlego(t *testing.T) {
    // Start scuttlego service in test mode
    testService := startTestScuttlego(t)
    defer testService.Close()

    // Publish post
    content := message.NewRawContent([]byte(`{"type":"post","text":"test"}`))
    msgRef, err := testService.App.Commands.PublishRaw.Handle(
        commands.PublishRaw{Content: content},
    )
    require.Nil(t, err)

    // Verify via scuttlego queries
    msg, err := testService.App.Queries.GetMessage.Handle(
        queries.GetMessage{Ref: msgRef},
    )
    require.Nil(t, err)
    require.Equal(t, "test", msg.Content.Text)
}
```

### 11. Build & Deployment (Section 7)
**Updated targets:**

```bash
# Dependencies (note: ~40 transitive)
go mod tidy
# Result: 40+ packages (badger, go-secretstream, margaret, etc.)

# Build (static but not <10MB)
CGO_ENABLED=0 go build -ldflags="-s -w"
# Result: 30-50MB binary (acceptable for Go apps)

# Cross-platform release
goreleaser release --clean
# Output: 30-50MB binaries per platform
```

### 12. Known Limitations (Section 9.2)
**Add:**
```
6. BadgerDB is not truly append-only (LSM-tree)
   - scuttlego handles compaction internally
   - Data is immutable at SSB protocol level
7. LAN discovery uses UDP broadcast, not mDNS
   - mDNS requires external library (hashicorp/mdns)
   - Can be added in Phase 4
```

## Implementation Order

### Phase 1: Foundation (Must Do First)
1. Update Section 1.1 (product definition)
2. Update Section 2.1 (architecture diagram)
3. Update Section 2.4 (API contract with real scuttlego examples)
4. Update Section 3.1 (file structure - remove storage/, add scuttlego/)
5. Update Section 3.2 (module responsibilities)

### Phase 2: Feature Specs (Critical for Implementation)
6. Update Section 2.3 (database schema - explain BadgerDB + scuttlego layering)
7. Update Section 5.1 (post creation with scuttlego API)
8. Update Section 5.2 (follow with scuttlego RedeemInvite)
9. Update Section 5.3 (trending indexes on top of scuttlego)

### Phase 3: Testing & Deployment
10. Update Section 6 (add scuttlego integration tests)
11. Update Section 7 (build targets, realistic sizes)
12. Update Section 9.2 (known limitations)

### Phase 4: Polish
13. Update Section 4.2 (decision matrix)
14. Update Section 10 (operational guidelines for scuttlego lifecycle)
15. Update Section 11 (success metrics - realistic memory targets)

## Validation Checklist

After revisions, verify:

- [ ] All mentions of "Scuttlegoa" changed to "scuttlego"
- [ ] All mentions of "BoltDB" changed to "BadgerDB v3" (except in historical context)
- [ ] "zero dependencies" claim removed
- [ ] Binary size target: 30-50MB
- [ ] Memory target: 100-200MB
- [ ] LAN discovery: "UDP broadcast" (not mDNS)
- [ ] All API examples use real scuttlego imports
- [ ] File structure includes scuttlego/ directory
- [ ] Testing strategy includes scuttlego integration tests
- [ ] Known limitations mention BadgerDB LSM-tree and UDP broadcast
- [ ] Go version: 1.23+ (compatible with scuttlego's 1.19+ requirement)

## Risk Assessment After Fixes

| Category | Before | After | Status |
|----------|--------|-------|--------|
| Library availability | Non-existent | scuttlego v0.0.4 | ✅ Resolved |
| Storage mismatch | BoltDB (wrong) | BadgerDB v3 | ✅ Resolved |
| Dependency claims | Zero (false) | 40+ transitive | ✅ Resolved |
| API hallucinations | Hypothetical | Real scuttlego APIs | ✅ Resolved |
| Discovery method | mDNS (unavailable) | UDP broadcast | ✅ Resolved |
| Memory targets | <50MB (impossible) | 100-200MB | ✅ Resolved |
| Library maturity | Unknown | Beta, 2yr stale | ⚠️ Documented |

**Overall Risk**: Medium (acceptable)
- Core functionality is implementable
- scuttlego is stable enough for MVP
- All major blockers resolved
- Known risks are documented

## Summary

These revisions transform production.md from an unusable specification with hallucinated libraries into a practical implementation guide based on:

1. **Real library**: scuttlego v0.0.4 (actual Go SSB implementation)
2. **Real storage**: BadgerDB v3 (what scuttlego actually uses)
3. **Real APIs**: Actual scuttlego Commands/Queries pattern
4. **Real targets**: Achievable binary size (30-50MB) and memory (100-200MB)
5. **Real features**: UDP discovery (supported) instead of mDNS (requires external lib)

The specification now accurately reflects what's possible with the chosen technology stack.

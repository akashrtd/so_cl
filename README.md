# so_cl

![alt text](so_clls.jpeg)

> **A decentralized social platform for your terminal**

<div align="center">

so_cl (pronounced "social") is a peer-to-peer (P2P) social platform running entirely in the terminal using the [Secure Scuttlebutt](https://scuttlebutt.nz/) (SSB) protocol.

[![Version](https://img.shields.io/badge/version-0.1.5-blue)](https://github.com/yourusername/so_cl/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/so_cl)](https://goreportcard.com/report/github.com/yourusername/so_cl)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/dl/)

</div>

---

## 📋 **Table of Contents**

- [About](#-about)
- [Features](#-features)
- [Technology Stack](#-technology-stack)
- [Getting Started](#-getting-started)
- [Development](#-development)
- [Contributing](#-contributing)
- [Testing](#-testing)
- [Architecture](#-architecture)
- [License](#-license)
- [Acknowledgments](#-acknowledgments)

---

## 🌟 **About**

**so_cl** is a terminal-based social network that prioritizes:

- **🔒 Privacy**: All data stored locally, encrypted connections
- **🌐 Decentralization**: No central servers, P2P communication
- **💻 Simplicity**: Pure text/ASCII, minimal resource usage
- **📱 Offline-First**: Create and read posts without internet
- **🎨 Creativity**: ASCII art profile pictures (6x6 colored ANSI)

### Why so_cl?

- Perfect for developers who love terminal tools
- Great for privacy advocates and offline-first enthusiasts
- Built on Secure Scuttlebutt (battle-tested P2P protocol)
- Works entirely offline when needed
- No corporate surveillance, no ads, no algorithms

---

## ✨ **Features**

### Current (Phase 0 - MVP) ✅ Complete
- ✅ SSB identity generation (Ed25519 keypair)
- ✅ Post creation (280 character limit)
- ✅ ASCII art profile pictures (6x6 colored ANSI)
- ✅ Local feed view with pagination (100 posts max)
- ✅ Offline-first (works without internet)
- ✅ Hashtag/mention indexing
- ✅ Optimistic UI updates
- ✅ Comprehensive test coverage

### Current (Phase 1 - P2P) ✅ Complete
- ✅ Follow via invite codes (`ssb:feed/invite/...`)
- ✅ EBT replication with peers
- ✅ LAN discovery (UDP broadcast)
- ✅ Peer list sidebar (F1 toggle)
- ✅ Real-time feed sync
- ✅ Keyboard shortcuts: F1 (peers), F2 (invite), F3 (follow)

### Current (Phase 2 - Social) ✅ Complete
- ✅ Replies (threaded messages with root/branch references)
- ✅ Like reactions (emoji-safe, vote type messages)
- ✅ Hashtag indexing + trending sidebar (F6 key)
- ✅ @mentions + notifications (F7 key)
- ✅ Follow/unfollow graph (indexed in BadgerDB)
- ✅ 5 TUI pages: Home, Discover, Peers, Profile, Settings

### Current (Phase 3 - Polish) 🔄 In Progress
- ✅ Search + filters (text, hashtag, author) - F8 key
- ✅ Profile pages (ASCII PFP, stats) - F9 key
- ✅ Settings (username, LAN discovery, PFP) - F10 key
- ⏳ Export/backup identity
- ⏳ Web UI (same state machine as TUI)

### TUI Pages
- **Home**: Social feed with posts, replies, and reactions
- **Discover**: Trending hashtags, popular posts, new peers to follow
- **Peers**: Connected peers, network statistics
- **Profile**: User profile card, post/followers/following counts
- **Settings**: Username, LAN discovery, regenerate PFP

---

## 🛠️ **Technology Stack**

| Component | Technology | Version |
|-----------|-----------|---------|
| **Language** | Go | 1.24+ |
| **SSB Library** | scuttlego | v0.0.4 |
| **Storage** | BadgerDB v3 | v3.2103.5 |
| **Frontend** | Bubble Tea | v1.3.10 |
| **Styling** | Lip Gloss | v1.1.0 |
| **Logging** | Zap | v1.27.1 |
| **Testing** | testify | v1.9.0 |

### Key Dependencies

- **scuttlego**: Go implementation of Secure Scuttlebutt protocol
- **BadgerDB**: Fast, embeddable key-value store (LSM-tree)
- **Bubble Tea**: Powerful TUI framework for Go
- **Lip Gloss**: Style definitions for nice terminal layouts
- **Zap**: Blazing fast, structured, leveled logging

---

## 🚀 **Getting Started**

### Prerequisites

- **Go 1.24+**: [Install Go](https://go.dev/dl/)
- **Git**: [Install Git](https://git-scm.com/downloads)
- **Terminal**: Any terminal emulator (iTerm2, Alacritty, Terminal.app, etc.)

### Installation

#### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/so_cl.git
cd so_cl

# Build the binary
go build -o so_cl

# Run so_cl
./so_cl
```

#### From Release (Coming Soon)

```bash
# Download the latest release for your platform
wget https://github.com/yourusername/so_cl/releases/latest/download/so_cl-linux-amd64.tar.gz

# Extract
tar -xzf so_cl-linux-amd64.tar.gz

# Run
./so_cl
```

#### Using Go Install

```bash
go install github.com/yourusername/so_cl@latest
```

### First Run

When you first run `so_cl`, it will:

1. **Generate an SSB identity** (Ed25519 keypair)
2. **Create a data directory** (`~/.so_cl/` with subdirectories for SSB data and indexes)
3. **Generate an ASCII profile picture** (6x6 colored ANSI, deterministic from feed ref)
4. **Start the TUI** with empty feed and navigation menu

### Basic Usage

```bash
# Start so_cl
./so_cl

# Keyboard shortcuts (TUI):
# ↑/↓         - Navigate feed/menu
# Enter       - Start typing/post
# F1          - Toggle peer list sidebar
# F2          - Toggle invite code input
# F3          - Toggle follow peer input
# F4          - Reply to selected post
# F5          - Like selected post
# F6          - Toggle trending sidebar
# F7          - Toggle mentions list
# F8          - Toggle search input
# F9          - Toggle profile view
# F10         - Toggle settings view
# F11         - Toggle follow graph sidebar
# Ctrl+C/Esc  - Quit

# CLI options:
./so_cl --help
./so_cl
# Environment variables:
#   SO_CL_DATA_DIR=~/.so_cl           # Custom data directory
#   SO_CL_PORT=8008                   # Custom SSB port
#   SO_CL_ENABLE_LAN_DISCOVERY=true   # Enable LAN discovery
#   SO_CL_DEBUG=true                  # Enable debug logging
```

---

## 💻 **Development**

### Setup Development Environment

```bash
# Clone the repository
git clone https://github.com/yourusername/so_cl.git
cd so_cl

# Install dependencies
go mod download

# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/goreleaser/goreleaser@latest

# Verify installation
go version      # Should be 1.23+
golangci-lint version
```

### Project Structure

```
so_cl/
├── main.go                    # Entry point, init logger & Bubble Tea
├── go.mod                     # Go dependencies
├── go.sum                     # Dependency checksums
├── Makefile                   # Build automation
├── .goreleaser.yml           # Release automation
│
├── ui/                       # Presentation layer (Bubble Tea TUI)
│   ├── model.go              # TUI model & state machine
│   └── model_test.go         # Model tests
│
├── core/                     # Domain models & utilities
│   ├── identity.go           # SSB identity (Ed25519)
│   ├── identity_test.go      # Identity tests
│   ├── message.go           # Post/reply types
│   ├── message_test.go      # Message tests
│   ├── types.go             # Common types (SoClPost, SoClPeer, etc.)
│   ├── types_test.go        # Type tests
│   ├── asciipfp.go         # ASCII profile picture generator
│   └── asciipfp_test.go    # PFP tests (determinism)
│
├── scuttlego/                # SSB protocol wrapper
│   ├── service.go            # Service: publish, follow, connect, search
│   └── service_test.go      # Integration tests
│
├── indexes/                  # BadgerDB custom indexes
│   ├── hashtags.go           # Hashtag counting, mentions, search
│   ├── hashtags_test.go      # Index tests
│   ├── follows.go           # Follow graph tracking
│   ├── follows_test.go      # Follow relationship tests
│   └── test_helpers.go      # Test utilities
│
└── config/                   # Configuration
    ├── config.go             # Config loader (env vars)
    └── config_test.go        # Config tests
```

### Building

```bash
# Build for current platform
go build -o so_cl

# Build static binary (CGO_ENABLED=0)
CGO_ENABLED=0 go build -ldflags="-s -w" -o so_cl

# Build for multiple platforms
make build-all  # Requires Makefile

# Or use goreleaser
goreleaser build --snapshot --clean
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -race -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./ui/ -v

# Run specific test
go test -v ./ui/ -run TestModelUpdate
```

### Linting

```bash
# Run golangci-lint
golangci-lint run

# Run go vet
go vet ./...

# Run staticcheck (if installed)
staticcheck ./...
```

---

## 🤝 **Contributing**

We welcome contributions from everyone! Whether you're fixing a bug, adding a feature, improving documentation, or reporting an issue, your help is valuable.

### How to Contribute

#### 1. Find an Issue

- Look for [good first issue](https://github.com/yourusername/so_cl/labels/good%20first%20issue) labels
- Check [open issues](https://github.com/yourusername/so_cl/issues)
- Or create a new issue if you have an idea

#### 2. Fork and Clone

```bash
# Fork the repository on GitHub
# Clone your fork
git clone https://github.com/YOUR_USERNAME/so_cl.git
cd so_cl

# Add upstream remote
git remote add upstream https://github.com/yourusername/so_cl.git
```

#### 3. Create a Branch

```bash
# Create a new branch for your feature/fix
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

#### 4. Make Changes

```bash
# Follow these guidelines:
# - Write tests FIRST (TDD)
# - Run tests with race detector: go test -race ./...
# - Run linter: golangci-lint run
# - Build: CGO_ENABLED=0 go build
# - Document your code (godoc comments)

# Example workflow:
# 1. Write test
# 2. Run test (should fail)
# 3. Implement feature
# 4. Run test (should pass)
# 5. Run all tests: go test -race ./...
# 6. Run linter: golangci-lint run
# 7. Commit changes
```

#### 5. Commit Changes

```bash
# Stage changes
git add .

# Commit with descriptive message
git commit -m "feat: add hashtag trending sidebar

- Extract hashtags from post text
- Store counts in BadgerDB
- Display top 10 in sidebar

Fixes #123"

# Commit message format:
# feat: new feature
# fix: bug fix
# docs: documentation changes
# style: formatting, missing semi colons, etc.
# refactor: refactoring production code
# test: adding missing tests
# chore: updating build tasks, etc.
```

#### 6. Push and Create Pull Request

```bash
# Push to your fork
git push origin feature/your-feature-name

# Go to GitHub and create a pull request
# Link to any related issues
# Describe your changes
```

### Development Guidelines

#### Code Style

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Run `gofmt` (enforced by linter)
- Use meaningful variable/function names
- Keep functions small and focused
- Add comments for complex logic

#### Testing

- **Test Coverage**: Maintain ≥85% (95%+ for critical paths)
- **Race Detection**: Always run `go test -race ./...`
- **TDD**: Write tests before implementation
- **Test Types**:
  - Unit tests: 70% (fast, isolated)
  - Integration tests: 20% (real DB/scuttlego)
  - E2E tests: 10% (TUI mock)

#### Test Coverage

| Package | Test Files | Tests |
|---------|-----------|-------|
| `scuttlego/` | 1 | 20+ |
| `core/` | 4 | 10+ |
| `ui/` | 1 | 20+ |
| `indexes/` | 2 | 15+ |
| `config/` | 1 | 5+ |
| **Total** | **9** | **70+** |

#### Documentation

- Add godoc comments for all exported functions
- Update `README.md` for user-facing changes
- Update `agent.md` for AI agent context
- Update `production.md` for spec changes

#### Precautions

- 🚫 **NEVER** skip tests
- 🚫 **NEVER** modify scuttlego internals (use adapters)
- 🚫 **NEVER** hardcode secrets
- 🚫 **NEVER** implement custom crypto (use scuttlego's go-secretstream)

### Contributing to AI Agent Guide

Since `agent.md` is used by AI agents, please:

- Update TODO list when completing tasks
- Update Session State after work
- Document new patterns in "Common Tasks"
- Add debugging guides for new issues

See [agent.md](agent.md) for details.

---

## 🧪 **Testing**

### Running Tests

```bash
# All tests
go test ./...

# With race detector (ALWAYS use this)
go test -race ./...

# With coverage
go test -race -cover ./...

# Generate HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Structure

```
core/
├── identity_test.go       # SSB identity generation (Ed25519)
├── asciipfp_test.go      # ASCII PFP tests (determinism, rendering)
├── message_test.go       # Post/reply type tests
└── types_test.go        # Common type tests

scuttlego/
└── service_test.go       # scuttlego service tests (publish, follow, connect, messages)

ui/
└── model_test.go         # Bubble Tea model tests (feed, publish, optimistic UI, pages)

indexes/
├── hashtags_test.go      # Hashtag/mention indexing tests
├── follows_test.go       # Follow graph tests
└── test_helpers.go      # Test utilities (test DB setup)

config/
└── config_test.go       # Config loading tests (env vars, defaults)
```

### Test Files

| Module | Test Files | Key Test Coverage |
|--------|------------|------------------|
| `scuttlego/` | 1 | publish, follow, connect, messages |
| `core/` | 4 | identity, message, types, ascii pfp |
| `ui/` | 1 | model updates, page rendering, keyboard |
| `indexes/` | 2 | hashtags, follows, search |
| `config/` | 1 | env vars, defaults |

### Coverage Goals

| Target | Goal |
|--------|-------|
| Unit tests (70%) | ✅ All critical paths |
| Integration tests (20%) | ✅ scuttlego service tests |
| E2E tests (10%) | 📝 TODO |
| Overall 85% | 📝 In progress |
| Critical paths 95% | 📝 In progress |

### Example Test

```go
package ui_test

import (
    "testing"
    "github.com/stretchr/testify/require"
    tea "github.com/charmbracelet/bubbletea"
)

func TestModelUpdate_PublishPost(t *testing.T) {
    // Arrange
    m := NewTestModel()
    m.composer.SetValue("hello so_cl")

    // Act
    model, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

    // Assert
    require.NotNil(t, cmd)
    require.IsType(t, &SoClModel{}, model)
}
```

---

## 🏗️ **Architecture**

### System Layers

```
┌────────────────────────────────────────────────────────────────┐
│  Bubble Tea TUI (ui/model.go)                           │  User interface
│  - 5 pages: Home, Discover, Peers, Profile, Settings      │  - State machine
│  - Keyboard navigation (F1-F11, ↑/↓, Enter)            │  - Optimistic UI
├────────────────────────────────────────────────────────────────┤
│  scuttlego Service Layer (scuttlego/service.go)          │  Protocol wrapper
│  - Publish, Reply, React (vote)                          │  - SSB operations
│  - Follow, Unfollow                                     │  - P2P networking
│  - Connect, RedeemInvite                                │  - LAN discovery
├────────────────────────────────────────────────────────────────┤
│  BadgerDB v3 Storage                                    │  Persistence
│  ├── SSB data (messages, feed, contacts)                │  - scuttlego managed
│  └── so_cl indexes                                      │  - Custom indexes
│      ├── Hashtags (counting)                            │  - hashtags.go
│      ├── Mentions (notification queue)                     │  - hashtags.go
│      ├── Follow graph (following/followers)              │  - follows.go
│      └── Full-text search                               │  - hashtags.go
├────────────────────────────────────────────────────────────────┤
│  scuttlego (SSB Protocol)                              │  Core protocol
│  - Ed25519 identity                                    │  - Cryptography
│  - EBT replication                                     │  - P2P sync
│  - go-secretstream (encryption)                          │  - Security
└────────────────────────────────────────────────────────────────┘
```

### Key Concepts

- **scuttlego**: Go implementation of SSB protocol (manages all P2P logic)
- **BadgerDB**: LSM-tree KV store (used by scuttlego and for so_cl indexes)
- **Bubble Tea**: TUI framework (event loop, state machine)
- **Lip Gloss**: Terminal styling (colors, layouts)

### Documentation

- **[README.md](README.md)** - This file (user guide)
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - Contribution guidelines
- **[agent.md](agent.md)** - AI agent reference guide
- **[production.md](production.md)** - Complete product specification

---

## 📜 **License**

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### TL;DR

- ✅ Free to use
- ✅ Free to modify
- ✅ Free to distribute
- ✅ Free to use commercially
- ⚠️ Must include license and copyright notice
- ⚠️ Software provided "as is" without warranty

---

## 🙏 **Acknowledgments

### Libraries

- **[scuttlego](https://github.com/planetary-social/scuttlego)** - Go SSB implementation
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - TUI framework
- **[Lip Gloss](https://github.com/charmbracelet/lipgloss)** - Terminal styling
- **[BadgerDB](https://github.com/dgraph-io/badger)** - Embeddable KV store
- **[Zap](https://github.com/uber-go/zap)** - Structured logging

### Protocols

- **[Secure Scuttlebutt](https://scuttlebutt.nz/)** - P2P social protocol
- **[Epidemic Broadcast Trees (EBT)](https://github.com/dominictarr/epidemic-broadcast-trees)** - Replication algorithm

### Inspiration

- **[Planetary](https://www.planetary.social/)** - SSB client (iOS/Android)
- **[Patchwork](https://patchwork.foo/)** - SSB client (Electron)
- **[Manyverse](https://manyver.se/)** - SSB client (Android)
- **[ssb-server](https://github.com/ssbc/ssb-server)** - Original SSB implementation

---

## 📞 **Support & Community**

### Getting Help

- **Issues**: [GitHub Issues](https://github.com/yourusername/so_cl/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/so_cl/discussions)
- **SSB Channel**: `#so_cl` on Secure Scuttlebutt
- **Matrix**: `#so_cl:matrix.org` (coming soon)

### Reporting Bugs

When reporting bugs, please include:

1. **so_cl version**: `./so_cl --version`
2. **Go version**: `go version`
3. **OS**: `uname -a` (Linux/macOS) or OS version (Windows)
4. **Steps to reproduce**: What you did before the bug
5. **Expected behavior**: What should have happened
6. **Actual behavior**: What actually happened
7. **Logs**: Run with `--debug` and include logs

### Feature Requests

We welcome feature requests! When suggesting a feature:

1. **Describe the problem**: What problem are you trying to solve?
2. **Proposed solution**: How do you think it should work?
3. **Alternatives**: Have you considered any alternatives?
4. **Additional context**: Any other relevant information?

---

## 🗺️ **Roadmap**

### Phase 0: MVP (Week 1-2) ✅ Complete
- [x] Project structure
- [x] Identity generation
- [x] Post creation
- [x] Feed view
- [x] Static binary build
- [x] Hashtag/mention indexing
- [x] Comprehensive tests
- [x] Release automation (goreleaser)

### Phase 1: P2P (Week 3-4) ✅ Complete
- [x] Follow via invite codes
- [x] Peer connections
- [x] EBT replication
- [x] LAN discovery (UDP broadcast)
- [x] Real-time feed sync

### Phase 2: Social (Week 5-6) ✅ Complete
- [x] Replies (threaded messages)
- [x] Like/repost reactions
- [x] Hashtag indexing
- [x] Trending sidebar
- [x] @mentions + notifications
- [x] Follow/unfollow graph

### Phase 3: Polish (Week 7+) 🔄 In Progress
- [x] Search + filters (text, hashtag, author)
- [x] Profile pages (PFP, bio, stats)
- [x] Settings (username, LAN discovery, PFP regeneration)
- [x] 5 TUI pages: Home, Discover, Peers, Profile, Settings
- [ ] Export/backup identity
- [ ] Web UI (same state machine as TUI)
- [ ] Docker support
- [ ] End-to-end tests

---

<div align="center">

**Made with ❤️ for the terminal**

[⬆ Back to Top](#so_cl)

</div>

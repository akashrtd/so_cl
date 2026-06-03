# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

so_cl (pronounced "social") is a decentralized, terminal-based social network built on the Secure Scuttlebutt (SSB) P2P protocol. It uses Bubble Tea for the TUI, scuttlego for SSB protocol operations, and BadgerDB for persistence and custom indexes.

## Build & Run Commands

```bash
# Build
go build -o so_cl                          # Debug build
CGO_ENABLED=0 go build -ldflags="-s -w" -o so_cl  # Static release build

# Run
./so_cl
SO_CL_DEBUG=true ./so_cl                   # Verbose logging

# Test
go test ./...                              # All tests
go test -race -v ./...                     # With race detector
go test -race -v ./ui/ -run TestModelUpdate  # Single test

# Lint
golangci-lint run
go vet ./...

# Format
go fmt ./...

# Cross-compile
make build-all
```

## Architecture

```
main.go → config.Load() → scuttlego.NewService() → ui.NewSoClModel() → tea.NewProgram()
```

**Four layers, each in its own package:**

1. **`config/`** — Env-var-driven configuration (`SO_CL_DATA_DIR`, `SO_CL_PORT`, `SO_CL_DEBUG`, etc.). Single `Load()` function, no file-based config.

2. **`scuttlego/`** — Wraps the `planetary-social/scuttlego` library. `Service` struct owns the scuttlego service instance, BadgerDB, and the custom `indexes.Indexer`. Provides all SSB operations: `Publish`, `Reply`, `React`, `Follow`, `Unfollow`, `Connect`, `RedeemInvite`, `GetMessages`, `GetFeed`, `Search`. Do NOT modify scuttlego internals — use its public API only.

3. **`indexes/`** — Custom BadgerDB indexes on top of scuttlego's storage. `Indexer` struct handles hashtag counting, mention tracking, trending metrics, follow graph, and full-text search. All stored as Badger key-value pairs with prefixed keys (e.g., `following:`, `followers:`, `hashtag:`).

4. **`ui/`** — Bubble Tea TUI model (`SoClModel`). Five pages: Home, Discover, Peers, Profile, Settings. State machine driven by `tea.Msg` updates. Uses Lip Gloss for neon-themed terminal styling. Keyboard shortcuts: F1–F11 for features, arrows/Enter for navigation, Ctrl+C/Esc to quit.

**`core/`** — Shared domain types (`SoClPost`, `SoClPeer`, `SoClProfile`, `Vote`, `Notification`, `TrendingHashtag`), `Post`/`Reply` constructors, `Identity` (Ed25519 keypair), and ASCII profile picture generator (`asciipfp.go`).

## Key Conventions

- Module path is `github.com/yourusername/so_cl` (not renamed yet)
- All storage goes through BadgerDB — no separate databases
- scuttlego manages its own data under `~/.so_cl/data/`; custom indexes share the same BadgerDB instance
- The TUI uses optimistic updates — posts appear immediately in the feed before SSB confirmation
- Posts are capped at 280 characters and 100 posts in memory
- Identity is Ed25519, feed refs follow SSB format (`@hex.ed25519`)
- Version and build time are injected via ldflags at build time

## Environment Variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `SO_CL_DATA_DIR` | `~/.so_cl` | Data root directory |
| `SO_CL_PORT` | `8008` | SSB listen port |
| `SO_CL_NETWORK_KEY` | empty (mainnet) | SSB network key |
| `SO_CL_ENABLE_LAN_DISCOVERY` | `true` | UDP broadcast peer discovery |
| `SO_CL_LOG_LEVEL` | `info` | Log level |
| `SO_CL_DEBUG` | `false` | Verbose logging |

## Testing

- Use `testify/require` and `testify/assert` for assertions
- `indexes/test_helpers.go` provides test BadgerDB setup helpers
- Tests should work offline — scuttlego integration tests mock network where needed
- Run with `-race` flag always


<!-- BEGIN BEADS INTEGRATION v:1 profile:minimal hash:6cd5cc61 -->
## Beads Issue Tracker

This project uses **bd (beads)** for issue tracking. Run `bd prime` to see full workflow context and commands.

### Quick Reference

```bash
bd ready              # Find available work
bd show <id>          # View issue details
bd update <id> --claim  # Claim work
bd close <id>         # Complete work
```

### Rules

- Use `bd` for ALL task tracking — do NOT use TodoWrite, TaskCreate, or markdown TODO lists
- Run `bd prime` for detailed command reference and session close protocol
- Use `bd remember` for persistent knowledge — do NOT use MEMORY.md files

**Architecture in one line:** issues live in a local Dolt DB; sync uses `refs/dolt/data` on your git remote; `.beads/issues.jsonl` is a passive export. See https://github.com/gastownhall/beads/blob/main/docs/SYNC_CONCEPTS.md for details and anti-patterns.

## Agent Context Profiles

The managed Beads block is task-tracking guidance, not permission to override repository, user, or orchestrator instructions.

- **Conservative (default)**: Use `bd` for task tracking. Do not run git commits, git pushes, or Dolt remote sync unless explicitly asked. At handoff, report changed files, validation, and suggested next commands.
- **Minimal**: Keep tool instruction files as pointers to `bd prime`; use the same conservative git policy unless active instructions say otherwise.
- **Team-maintainer**: Only when the repository explicitly opts in, agents may close beads, run quality gates, commit, and push as part of session close. A current "do not commit" or "do not push" instruction still wins.

## Session Completion

This protocol applies when ending a Beads implementation workflow. It is subordinate to explicit user, repository, and orchestrator instructions.

1. **File issues for remaining work** - Create beads for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **Handle git/sync by active profile**:
   ```bash
   # Conservative/minimal/default: report status and proposed commands; wait for approval.
   git status

   # Team-maintainer opt-in only, unless current instructions forbid it:
   git pull --rebase
   git push
   git status
   ```
5. **Hand off** - Summarize changes, validation, issue status, and any blocked sync/commit/push step

**Critical rules:**
- Explicit user or orchestrator instructions override this Beads block.
- Do not commit or push without clear authority from the active profile or the current user request.
- If a required sync or push is blocked, stop and report the exact command and error.
<!-- END BEADS INTEGRATION -->

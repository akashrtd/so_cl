# agent.md Creation Summary

**Date**: December 29, 2025
**File Created**: `/Users/akashrathod/Desktop/projects/so_cl/agent.md`
**Total Lines**: 1,330 lines
**Sections**: 14 major sections

---

## Overview

`agent.md` is a **comprehensive reference guide and memory bank** for AI agents working on the so_cl project. It ensures continuity across sessions, maintains high coding standards, and provides all necessary context for autonomous development.

---

## Key Features

### 1. **Memory Bank** 🧠
- **Session State**: Editable template for tracking current work
- **Recent Work Log**: Table to document completed tasks
- **Current Session Goals**: Checkbox list for immediate priorities
- **Pause & Continue**: Clear guidelines for stopping/resuming work

### 2. **Complete Project Context** 📋
- Project overview (what is so_cl, goals, target users)
- Technology stack (Go, scuttlego, BadgerDB, Bubble Tea)
- Architecture diagrams and design patterns
- File structure and module responsibilities

### 3. **AI Agent Workflow** 🤖
- **Daily Workflow**: Morning standup, development steps, end-of-day
- **Best Practices**: Code style, testing, memory management, logging
- **Precautions & Constraints**:
  - 🚫 NEVER DO (critical failures to avoid)
  - ⚠️ ALWAYS ASK BEFORE (requires approval)
  - ✅ ALWAYS DO (autonomous actions)

### 4. **Editable TODO List** ✅
Organized by project phase with 21+ checkboxes:

**Phase 0 (MVP - Week 1-2)**: 5 tasks
- Initialize project structure
- Identity generation
- Post creation
- Feed view
- Static binary build

**Phase 1 (P2P - Week 3-4)**: 5 tasks
- Follow via invite codes
- Peer connections
- EBT replication
- LAN discovery (UDP broadcast)
- Real-time feed sync

**Phase 2 (Social - Week 5-6)**: 6 tasks
- Replies (threaded messages)
- Like/repost reactions
- Hashtag indexing
- Trending sidebar
- @mentions + notifications
- Follow/unfollow graph

**Phase 3 (Polish - Week 7+)**: 5 tasks
- Search + filters
- Profile pages
- Settings (username, PFP, discovery)
- Export/backup identity
- Web UI

### 5. **Testing Strategy** 🧪
- **Test Pyramid**: Unit (70%), Integration (20%), E2E (10%)
- **Coverage Targets**: 85% overall, 95%+ for critical paths
- **Test Categories**: Unit, Integration, E2E examples
- **Running Tests**: Commands for race detection, coverage

### 6. **Common Tasks** ⚙️
Code templates for frequent operations:
- Task 1: Create a New Package
- Task 2: Add a New Bubble Tea View
- Task 3: Add Async Operation
- Task 4: Add BadgerDB Index
- Task 5: Add scuttlego Query

### 7. **API Quick Reference** 📚
Real code examples for:
- **scuttlego Commands**: PublishRaw, Follow, Connect, RedeemInvite
- **scuttlego Queries**: GetMessage, ReceiveLog, Status
- **Bubble Tea Events**: KeyMsg, custom messages
- **BadgerDB Operations**: Update, View, transactions

### 8. **Decision Matrix** 🎯
Clear guidance on:
- **When to code autonomously** (UI, validation, error handling, etc.)
- **When to ask for approval** (protocol changes, schema changes, security, etc.)

### 9. **Debugging Guidelines** 🐛
Solutions for common issues:
- Posts not syncing between peers
- High memory usage (>200MB)
- Slow UI (<30fps)
- scuttlego service won't start

Debugging tools:
- Logging (DEBUG mode, log files)
- Profiling (pprof)
- Memory monitoring

### 10. **Pause & Continue** ⏸️▶️
Clear procedures:
- **How to Pause**: Commit work, update agent.md, clean up
- **How to Continue**: Read agent.md, check TODO, review commits
- **Session State Template**: Editable template for tracking

### 11. **Success Metrics** 📊
Quality and delivery targets:
- **Code Quality**: Coverage (85%+), Race conditions (0), Lint errors (0)
- **Feature Delivery**: Phase deadlines (Week 2, 4, 6, 7+)
- **User Experience**: UI responsiveness (≥60fps), startup time (≤5s)

### 12. **Quick Reference Card** ⚡
Essential commands, key files, and imports at a glance.

---

## Usage Guide

### For AI Agents

**First Session:**
1. Read `agent.md` completely (14 sections)
2. Understand project context and architecture
3. Review technology stack and API reference
4. Start with first TODO item in Phase 0

**Daily Workflow:**
```
1. Read agent.md → Refresh context
2. Check Session State → Understand current work
3. Pick next TODO task → Highest priority
4. Write test FIRST (TDD) → Ensure quality
5. Implement feature → Follow best practices
6. Run tests with race detector → Verify no bugs
7. Update TODO list → Mark items complete
8. Update Session State → Document progress
```

**Pausing Work:**
```
1. Update Session State (current task, blockers)
2. Mark completed TODO items
3. Commit changes with descriptive message
4. Note blockers in Session State
```

**Resuming Work:**
```
1. Read agent.md → Refresh context
2. Check Session State → See where you left off
3. Review recent commits → Understand changes
4. Continue with next TODO item
```

### For Human Reviewers

**Checking Progress:**
- Look at **Session State** → Current phase, task, blockers
- Check **TODO list** → Completed vs. remaining tasks
- Review **Recent Work Log** → What was done recently

**Verifying Quality:**
- Run tests: `go test -race -cover ./...`
- Check coverage: ≥85% (critical 95%+)
- Check lint: `golangci-lint run`
- Verify build: `CGO_ENABLED=0 go build`

**Blockers & Issues:**
- Review **Blockers** field in Session State
- Check **Recent Work Log** notes
- Review **Debugging Guidelines** if needed

---

## File Structure

```
agent.md (1,330 lines)
├── 📋 TABLE OF CONTENTS
├── 1. PROJECT CONTEXT
│   ├── What is so_cl?
│   ├── Project Goals (4 phases)
│   └── Target Users
├── 2. TECHNOLOGY STACK
│   ├── Core Dependencies (table)
│   ├── Transitive Dependencies (~40+)
│   └── Build Targets
├── 3. ARCHITECTURE OVERVIEW
│   ├── System Layers (diagram)
│   ├── Key Design Patterns (3 patterns)
│   └── File Structure (tree)
├── 4. AI AGENT WORKFLOW
│   ├── Daily Workflow (3 steps)
│   └── Session State (editable)
├── 5. PRECAUTIONS & CONSTRAINTS
│   ├── 🚫 NEVER DO (5 critical failures)
│   ├── ⚠️ ALWAYS ASK BEFORE (6 scenarios)
│   └── ✅ ALWAYS DO (7 autonomous actions)
├── 6. BEST PRACTICES
│   ├── Code Style (Go conventions)
│   ├── Testing (unit/integration/E2E)
│   ├── Memory Management (BadgerDB, goroutines)
│   └── Logging (structured with zap)
├── 7. EDITABLE TODO LIST
│   ├── 📋 Phase 0: MVP (5 tasks)
│   ├── 📋 Phase 1: P2P (5 tasks)
│   ├── 📋 Phase 2: Social (6 tasks)
│   └── 📋 Phase 3: Polish (5 tasks)
├── 8. TESTING STRATEGY
│   ├── Test Pyramid (diagram)
│   ├── Coverage Targets (table)
│   ├── Test Categories (3 types)
│   ├── Running Tests (commands)
│   └── Race Conditions (examples)
├── 9. COMMON TASKS
│   ├── Task 1: Create a New Package
│   ├── Task 2: Add a New Bubble Tea View
│   ├── Task 3: Add Async Operation
│   ├── Task 4: Add BadgerDB Index
│   └── Task 5: Add scuttlego Query
├── 10. API QUICK REFERENCE
│   ├── scuttlego Commands (write operations)
│   ├── scuttlego Queries (read operations)
│   ├── Bubble Tea Events
│   └── BadgerDB Operations
├── 11. DECISION MATRIX
│   ├── When to Code Autonomously ✅
│   └── When to Ask for Approval ⏸️
├── 12. DEBUGGING GUIDELINES
│   ├── Posts not syncing (diagnosis + fix)
│   ├── High memory usage (diagnosis + fix)
│   ├── Slow UI (diagnosis + fix)
│   ├── scuttlego service won't start (diagnosis + fix)
│   └── Debugging Tools (logging, profiling, memory)
├── 13. PAUSE & CONTINUE
│   ├── How to Pause (5 steps)
│   ├── How to Continue (4 steps)
│   └── Session State Template
├── 14. SUCCESS METRICS
│   ├── Code Quality (6 metrics)
│   ├── Feature Delivery (4 phases)
│   └── User Experience (3 metrics)
└── QUICK REFERENCE CARD
    ├── Essential Commands
    ├── Key Files
    └── Key Imports
```

---

## Key Benefits

### 1. **Continuity Across Sessions** 🔄
- **Session State**: Always know where you left off
- **Recent Work Log**: Track what was done
- **Editable TODO**: Never lose progress

### 2. **Autonomous Development** 🚀
- **Complete Context**: All project info in one file
- **Clear Guidelines**: What to do, what to avoid
- **Code Templates**: Copy-paste for common tasks

### 3. **Quality Assurance** ✅
- **Best Practices**: Code style, testing, logging
- **Precautions**: What NOT to do
- **Success Metrics**: Targets for quality and delivery

### 4. **Easy Reference** 📚
- **API Quick Reference**: Real code examples
- **Debugging Guidelines**: Common issues and solutions
- **Quick Reference Card**: Essential commands at a glance

### 5. **Collaboration Friendly** 👥
- **Editable TODO**: Humans can check progress
- **Session State**: Reviewers understand current work
- **Recent Work Log**: Transparency on what was done

---

## Integration with Project Files

`agent.md` is designed to work with:

1. **production.md** (Full product specification)
   - agent.md provides workflow and context
   - production.md provides detailed specs

2. **go.mod** (Go dependencies)
   - agent.md lists all dependencies
   - go.mod manages actual versions

3. **Source code** (Actual implementation)
   - agent.md provides patterns and examples
   - Source code implements features

---

## Maintenance

### When to Update agent.md

**After Completing a Task:**
- Mark checkbox in TODO list
- Update Session State
- Add entry to Recent Work Log

**When Adding a New Pattern:**
- Document in "Common Tasks" section
- Add example to "API Quick Reference"

**When Encountering New Issues:**
- Document in "Debugging Guidelines"
- Provide diagnosis and fix

**When Changing Technology:**
- Update "Technology Stack" section
- Update "API Quick Reference"

### Version History

| Date | Version | Changes |
|------|---------|---------|
| 2025-12-29 | v1.0 | Initial creation (1,330 lines, 14 sections) |

---

## Summary

`agent.md` is a **comprehensive, living document** that serves as:

✅ **Memory Bank** - Tracks progress, session state, TODOs
✅ **Reference Guide** - APIs, patterns, best practices
✅ **Quality Enforcer** - Precautions, metrics, standards
✅ **Workflow Guide** - Daily tasks, pause/continue procedures
✅ **Debugging Aid** - Common issues, solutions, tools

**Purpose**: Enable high-quality, autonomous AI agent development with continuity across sessions.

**Status**: ✅ **READY TO USE**

---

## Next Steps for AI Agent

1. **Read agent.md** completely (14 sections)
2. **Update Session State** with current information
3. **Start first TODO item**: "Initialize project structure"
4. **Follow best practices** from Section 6
5. **Update TODO list** as tasks complete
6. **Maintain Session State** throughout development

**Remember**: This is a **living document**. Keep it updated as you progress! 🚀

# Contributing to so_cl

Thank you for considering contributing to so_cl! We appreciate your help in making this decentralized social platform better.

This document provides guidelines and instructions for contributing to the project.

---

## 📋 **Table of Contents**

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Commit Messages](#commit-messages)
- [Pull Request Process](#pull-request-process)
- [Documentation](#documentation)
- [Reporting Issues](#reporting-issues)
- [Feature Requests](#feature-requests)

---

## 🤝 **Code of Conduct**

### Our Pledge

We are committed to providing a welcoming and inclusive environment for all contributors. We expect everyone to:

- Be respectful and considerate
- Use inclusive language
- Focus on constructive feedback
- Help others learn and grow
- Respect differing viewpoints and experiences

### Unacceptable Behavior

Harassment, discrimination, or disrespectful behavior is not tolerated. This includes:

- Offensive comments related to gender, gender identity and expression, sexual orientation, disability, mental illness, physical appearance, body size, age, race, or religion
- Unwelcome sexual attention or advances
- Threatening, intimidating, or harassing behavior
- Trolling or insulting/derogatory comments

### Reporting Issues

If you experience or witness unacceptable behavior, please contact us at:
- Email: [your-email@example.com](mailto:your-email@example.com)
- SSB: [@youridentity](ssb:@youridentity)

---

## 🚀 **Getting Started**

### Prerequisites

Before contributing, ensure you have:

- **Go 1.23+**: [Install Go](https://go.dev/dl/)
- **Git**: [Install Git](https://git-scm.com/downloads)
- **Terminal**: Any modern terminal emulator
- **Text Editor**: VS Code, GoLand, vim, etc.

### Setup Development Environment

```bash
# 1. Fork and clone the repository
git clone https://github.com/YOUR_USERNAME/so_cl.git
cd so_cl

# 2. Add upstream remote
git remote add upstream https://github.com/yourusername/so_cl.git

# 3. Install dependencies
go mod download

# 4. Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/goreleaser/goreleaser@latest

# 5. Verify setup
go version      # Should be 1.23+
golangci-lint version

# 6. Build the project
go build -o so_cl

# 7. Run tests
go test ./...
```

### Understanding the Codebase

Before making changes, please read:

1. **[README.md](README.md)** - Project overview and architecture
2. **[production.md](production.md)** - Detailed product specification
3. **[agent.md](agent.md)** - AI agent reference guide (contains code patterns)

Key concepts:
- **scuttlego**: SSB protocol library (don't modify internals)
- **BadgerDB**: Storage engine (LSM-tree KV store)
- **Bubble Tea**: TUI framework (event loop, state machine)
- **Lip Gloss**: Terminal styling (colors, layouts)

---

## 🔄 **Development Workflow**

### 1. Choose an Issue

Find a task to work on:

- **Good First Issues**: Look for `good first issue` label
- **Bug Fixes**: Check `bug` label
- **Features**: Check `enhancement` label
- **Documentation**: Check `documentation` label
- **Or create a new issue** for your idea

### 2. Create a Branch

```bash
# Sync with upstream
git fetch upstream
git checkout main
git merge upstream/main

# Create a new branch (use conventional prefixes)
git checkout -b feature/your-feature-name      # New feature
git checkout -b fix/your-bug-fix                # Bug fix
git checkout -b docs/update-readme              # Documentation
git checkout -b refactor/optimize-dispatcher    # Refactoring
git checkout -b test/add-unit-tests            # Tests
```

### 3. Write Tests First (TDD)

We follow Test-Driven Development:

```go
// 1. Write the test FIRST (it will fail)
func TestPublishPost_ValidText(t *testing.T) {
    m := NewTestModel()
    err := m.PublishPost("hello so_cl")
    require.Nil(t, err)  // This will fail initially
}

// 2. Run the test (it should fail)
go test ./ui/ -run TestPublishPost_ValidText -v

// 3. Implement the feature
func (m *SoClModel) PublishPost(text string) error {
    // Implementation here
}

// 4. Run the test again (it should pass)
go test ./ui/ -run TestPublishPost_ValidText -v
```

### 4. Implement the Feature

Follow these guidelines:

```bash
# 5. Run all tests with race detector
go test -race ./...

# 6. Run linter
golangci-lint run

# 7. Build the project
CGO_ENABLED=0 go build -ldflags="-s -w" -o so_cl

# 8. Test manually (if applicable)
./so_cl
```

### 5. Commit Your Changes

```bash
# Stage changes
git add .

# Commit with descriptive message (see Commit Messages section)
git commit -m "feat: add hashtag trending sidebar

- Extract hashtags from post text
- Store counts in BadgerDB
- Display top 10 in sidebar

Fixes #123"
```

### 6. Push and Create Pull Request

```bash
# Push to your fork
git push origin feature/your-feature-name

# Go to GitHub and create a PR
# Link to any related issues
# Describe your changes
```

---

## 📝 **Coding Standards**

### Go Conventions

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Run `gofmt` (enforced by CI)
- Use `golint` and `staticcheck`
- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Naming Conventions

```go
// ✅ Good: Clear, descriptive names
func (m *SoClModel) PublishPost(text string) error

// ❌ Bad: Vague names
func (m *SoClModel) Do(s string) error

// ✅ Good: Constants use UPPER_SNAKE_CASE
const MAX_POST_LENGTH = 280

// ✅ Good: Exported types use PascalCase
type SoClModel struct {}

// ✅ Good: Private types use camelCase
type messageBuffer struct {}
```

### Error Handling

```go
// ✅ Good: Wrap errors with context
return fmt.Errorf("failed to publish post: %w", err)

// ✅ Good: Use structured logging
zap.L().Error("publish failed",
    zap.String("text", text),
    zap.Error(err),
)

// ❌ Bad: Generic errors
return errors.New("error")

// ❌ Bad: Silent error ignoring
_ = someFunction()  // Don't ignore errors!
```

### Documentation

```go
// Package-level comment
// Package scuttlego provides a wrapper around the scuttlego library
// for use in so_cl.
package scuttlego

// Function comment (godoc format)
// PublishPost publishes a new post to the SSB feed.
// The text must be 1-280 characters.
// Returns an error if the post cannot be published.
func (m *SoClModel) PublishPost(text string) error {
    // ...
}

// Complex logic comments
// Check if the message is too long and truncate if needed.
// This is necessary because SSB has a maximum message size limit.
if len(text) > MAX_POST_LENGTH {
    text = text[:MAX_POST_LENGTH]
}
```

### Code Organization

```
// Exported symbols at the top
const MaxPostLength = 280
var DefaultConfig Config

// Types
type SoClModel struct {}
type Post struct {}

// Constructors
func NewSoClModel() *SoClModel
func NewPost(text string) *Post

// Methods
func (m *SoClModel) PublishPost(text string) error
func (p *Post) String() string

// Private functions
func validateText(text string) error
func extractHashtags(text string) []string
```

### Precautions

- 🚫 **NEVER** skip tests
- 🚫 **NEVER** modify scuttlego internals (use adapters instead)
- 🚫 **NEVER** hardcode secrets (use environment variables)
- 🚫 **NEVER** implement custom crypto (use scuttlego's go-secretstream)
- 🚫 **NEVER** break SSB protocol (don't modify message format/signing)

---

## 🧪 **Testing Guidelines**

### Test Structure

```
tests/
├── model_test.go           # Bubble Tea model tests
├── dispatcher_test.go      # Async orchestration tests
├── indexes_test.go         # BadgerDB index tests
├── scuttlego_test.go      # scuttlego integration tests
└── e2e_test.go           # End-to-end tests
```

### Test Types

#### Unit Tests (70%)

```go
func TestPublishPost_ValidText(t *testing.T) {
    // Arrange
    m := NewTestModel()
    text := "hello world"

    // Act
    err := m.PublishPost(text)

    // Assert
    require.Nil(t, err)
    assert.Equal(t, 1, len(m.posts))
    assert.Equal(t, text, m.posts[0].Text)
}
```

#### Integration Tests (20%)

```go
func TestPublishPost_WithScuttlego(t *testing.T) {
    // Start test service
    svc := startTestScuttlego(t)
    defer svc.Close()

    // Publish post
    content, _ := message.NewRawContent([]byte(`{"type":"post","text":"test"}`))
    msgRef, err := svc.App.Commands.PublishRaw.Handle(
        commands.PublishRaw{Content: content},
    )
    require.Nil(t, err)

    // Verify
    msg, err := svc.App.Queries.GetMessage.Handle(
        queries.GetMessage{Ref: msgRef},
    )
    require.Nil(t, err)
    assert.Equal(t, "test", msg.Content.Text)
}
```

#### End-to-End Tests (10%)

```go
func TestUserFlow_PublishAndReply(t *testing.T) {
    m := NewTestModel()

    // User types in composer
    m.composer.SetValue("Hello so_cl!")

    // Simulate Enter key
    model, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

    // Verify post appears in feed
    time.Sleep(100 * time.Millisecond)
    assert.Greater(t, len(m.cachedFeed), 0)
}
```

### Running Tests

```bash
# All tests
go test ./...

# With race detector (ALWAYS use this)
go test -race ./...

# With coverage
go test -race -cover ./...

# Specific package
go test ./ui/ -v

# Specific test
go test -v ./ui/ -run TestModelUpdate

# Generate HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Coverage Targets

| Module | Target | Critical Path |
|--------|--------|---------------|
| `scuttlego/service.go` | 95%+ | ✅ YES |
| `protocol/adapter.go` | 90%+ | ✅ YES |
| `indexes/hashtags.go` | 90% | - |
| `async/dispatcher.go` | 95%+ | ✅ YES |
| `ui/model.go` | 90% | - |
| **Overall** | **85%** | **Critical 95%+** |

### Test Best Practices

- ✅ **Test FIRST** (TDD)
- ✅ **ALWAYS** run `go test -race ./...`
- ✅ Use **table-driven tests** for multiple scenarios
- ✅ Use **testify** for assertions (`require` vs `assert`)
- ✅ Mock external dependencies (scuttlego)
- ✅ Clean up resources (`defer svc.Close()`)

---

## 💬 **Commit Messages**

We follow [Conventional Commits](https://www.conventionalcommits.org/):

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation changes |
| `style` | Formatting, missing semicolons, etc. |
| `refactor` | Refactoring production code |
| `test` | Adding missing tests |
| `chore` | Updating build tasks, package manager configs, etc. |

### Examples

```bash
# Feature
feat(ui): add hashtag trending sidebar

- Extract hashtags from post text
- Store counts in BadgerDB
- Display top 10 in sidebar

Closes #123

# Bug fix
fix(p2p): resolve connection timeout issue

The connection timeout was too aggressive, causing
connections to be dropped prematurely. Increased
timeout from 5s to 30s.

Fixes #456

# Documentation
docs(readme): update installation instructions

Added instructions for building from source and
installing via go install.

# Refactor
 refactor(async): optimize dispatcher worker pool

Improved load balancing across workers using
work-stealing algorithm. Improves throughput
by 15%.

# Test
test(ui): add unit tests for composer validation

Added tests for:
- 280 character limit
- Empty text validation
- ASCII-only validation

Related to #789
```

---

## 🔀 **Pull Request Process**

### Before Submitting

1. **Check the checklist**:
   - [ ] Tests pass (`go test -race ./...`)
   - [ ] Coverage ≥85% (critical 95%+)
   - [ ] Linter passes (`golangci-lint run`)
   - [ ] Builds successfully (`CGO_ENABLED=0 go build`)
   - [ ] Documentation updated
   - [ ] Commit messages follow format

2. **Sync with main**:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

3. **Resolve conflicts** (if any):
   ```bash
   # Rebase onto main
   git rebase upstream/main

   # Fix conflicts, then:
   git add .
   git rebase --continue
   ```

### Creating a Pull Request

1. **Go to GitHub**: https://github.com/yourusername/so_cl/compare
2. **Select your branch**: feature/your-feature-name
3. **Fill out the PR template**:
   - Title: Use conventional commit format (e.g., "feat: add hashtag trending")
   - Description:
     - What changes were made?
     - Why were they made?
     - How were they tested?
   - Linked issues: `Fixes #123`, `Closes #456`
   - Screenshots: For UI changes

### PR Review Process

1. **Automated Checks**: CI runs tests, linter, build
2. **Code Review**: Maintainers review the code
3. **Feedback**: Address review comments
4. **Approval**: At least one maintainer approval
5. **Merge**: Squash merge into main branch

### Review Guidelines

- Be constructive and respectful
- Focus on the code, not the person
- Provide clear, actionable feedback
- Ask questions if something is unclear
- Acknowledge and thank the contributor

---

## 📚 **Documentation**

### What to Document

- **Public APIs**: Add godoc comments to all exported functions
- **Complex logic**: Add inline comments explaining the "why"
- **User-facing changes**: Update README.md
- **Spec changes**: Update production.md
- **AI agent context**: Update agent.md

### Documentation Files

| File | Purpose | Audience |
|------|----------|-----------|
| `README.md` | User guide | Users |
| `CONTRIBUTING.md` | Contribution guidelines | Contributors |
| `agent.md` | AI agent reference | AI agents |
| `production.md` | Product specification | Developers/AI agents |

### Documentation Style

```go
// Function comment (godoc format)
// PublishPost publishes a new post to the SSB feed.
//
// The text must be 1-280 characters and contain only ASCII characters.
// Returns an error if the post cannot be published or if validation fails.
//
// Example:
//
//   err := model.PublishPost("hello so_cl")
//   if err != nil {
//       log.Fatal(err)
//   }
func (m *SoClModel) PublishPost(text string) error
```

---

## 🐛 **Reporting Issues**

### Bug Reports

When reporting bugs, please include:

1. **so_cl version**: `./so_cl --version`
2. **Go version**: `go version`
3. **OS**: `uname -a` (Linux/macOS) or Windows version
4. **Steps to reproduce**:
   ```bash
   1. Run: ./so_cl
   2. Type: "test"
   3. Press: Enter
   4. Error appears
   ```
5. **Expected behavior**: What should have happened
6. **Actual behavior**: What actually happened (include error message)
7. **Logs**: Run with `--debug` and include logs

### Issue Template

```markdown
## Description
Brief description of the bug.

## Steps to Reproduce
1. Run...
2. Type...
3. Press...
4. Error...

## Expected Behavior
What should have happened

## Actual Behavior
What actually happened

## Environment
- so_cl version: X.Y.Z
- Go version: 1.23.X
- OS: Ubuntu 22.04

## Logs
```

---

## 💡 **Feature Requests**

### Proposing a Feature

When suggesting a feature:

1. **Describe the problem**: What problem are you trying to solve?
2. **Proposed solution**: How do you think it should work?
3. **Alternatives**: Have you considered any alternatives?
4. **Additional context**: Any other relevant information (mockups, examples)

### Feature Request Template

```markdown
## Problem
Describe the problem you're facing or the gap in functionality.

## Proposed Solution
Describe how you think this feature should work.

## Alternatives
Describe any alternative solutions or features you considered.

## Additional Context
Add any other context, screenshots, or examples here.
```

---

## 🎓 **Learning Resources**

### Go

- [A Tour of Go](https://go.dev/tour/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)

### Testing

- [Testing in Go](https://go.dev/doc/tutorial/add-a-test)
- [Table-Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Testify](https://github.com/stretchr/testify)

### SSB Protocol

- [Secure Scuttlebutt Website](https://scuttlebutt.nz/)
- [SSB Protocol Guide](https://ssbc.github.io/scuttlebutt-protocol-guide/)
- [scuttlego Documentation](https://pkg.go.dev/github.com/planetary-social/scuttlego)

### TUI Development

- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- [Termenv](https://github.com/muesli/termenv)

---

## 🏆 **Recognition**

Contributors will be recognized in:

- **README.md**: Contributors section
- **CHANGELOG.md**: Listed with each release
- **GitHub Contributors**: Automatic recognition

### Top Contributors

<!-- TODO: Add top contributors here -->

---

## ❓ **Questions?

- **GitHub Issues**: [Ask a question](https://github.com/yourusername/so_cl/issues/new?labels=question)
- **Discussions**: [Start a discussion](https://github.com/yourusername/so_cl/discussions)
- **SSB**: `#so_cl` channel on Secure Scuttlebutt
- **Email**: [your-email@example.com](mailto:your-email@example.com)

---

## 📜 **Code of Conduct Summary**

- Be respectful and inclusive
- Focus on constructive feedback
- Help others learn and grow
- Report unacceptable behavior

For full details, see the Code of Conduct section above.

---

**Thank you for contributing to so_cl!** 🚀

<div align="center">

Made with ❤️ by the so_cl community

[⬆ Back to Top](#contributing-to-so_cl)

</div>

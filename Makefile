.PHONY: build clean test test-race lint run help all

# Build variables
BINARY_NAME=so_cl
BUILD_DIR=./build
GO=go
GOFLAGS=-v
LDFLAGS=-ldflags="-s -w"

# Directories
CMD_DIR=./
DATA_DIR=~/.so_cl/data

# Build for current platform
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build static binary (cross-platform)
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)

	# Linux AMD64
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./$(CMD_DIR)

	# Linux ARM64 (RPi)
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./$(CMD_DIR)

	# macOS AMD64
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./$(CMD_DIR)

	# macOS ARM64 (M1/M2)
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./$(CMD_DIR)

	# Windows AMD64
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./$(CMD_DIR)

	@echo "Build complete for all platforms"
	@ls -lh $(BUILD_DIR)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Run all tests
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	$(GO) test -race -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -race -cover ./...
	@$(GO) tool cover -html=coverage.out
	@echo "Coverage report generated: coverage.out"
	@echo "HTML report: coverage.html"

# Run linter
lint:
	@echo "Running golangci-lint..."
	golangci-lint run
	@echo "Linting complete"

# Run go vet
vet:
	@echo "Running go vet..."
	$(GO) vet ./...
	@echo "Vetting complete"

# Run application
run: build
	@echo "Running $(BINARY_NAME)..."
	$(BUILD_DIR)/$(BINARY_NAME)

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "Dependencies updated"

# Create data directory
data-dir:
	@mkdir -p $(DATA_DIR)

# Docker build (for testing)
docker-build:
	@echo "Building Docker image..."
	docker build -t so_cl:latest .

# Docker run
docker-run:
	@echo "Running Docker container..."
	docker run --rm -it so_cl:latest

# Install development tools
install-tools:
	@echo "Installing development tools..."
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install github.com/goreleaser/goreleaser@latest
	@echo "Tools installed"

# Help target
help:
	@echo "Available targets:"
	@echo "  make build       - Build for current platform"
	@echo "  make build-all   - Build for all platforms"
	@echo "  make clean       - Remove build artifacts"
	@echo "  make test        - Run all tests"
	@echo "  make test-race   - Run tests with race detector"
	@echo "  make test-coverage - Run tests with coverage"
	@echo "  make lint        - Run golangci-lint"
	@echo "  make vet         - Run go vet"
	@echo "  make run         - Build and run application"
	@echo "  make fmt         - Format code"
	@echo "  make deps        - Download/update dependencies"
	@echo "  make data-dir    - Create data directory"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run  - Run Docker container"
	@echo "  make install-tools - Install development tools"
	@echo "  make help        - Show this help message"

# All target (build + test + lint)
all: build test lint
	@echo "All targets completed successfully!"

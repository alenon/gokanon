# Makefile for GoKanon

# Variables
BINARY_NAME=gokanon
VERSION?=dev
BUILD_DIR=./bin
CMD_DIR=.
COVERAGE_DIR=./coverage

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean

# Build flags
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X github.com/alenon/gokanon/internal/cli.Version=$(VERSION) -X github.com/alenon/gokanon/internal/cli.GitCommit=$(GIT_COMMIT) -X github.com/alenon/gokanon/internal/cli.BuildDate=$(BUILD_DATE) -s -w"
BUILD_FLAGS=-trimpath

# Test flags
TEST_FLAGS=-v -race
COVERAGE_FLAGS=-coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic

# Detect OS for binary naming
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

.PHONY: all build test clean install uninstall fmt vet lint coverage help

# Default target
all: clean fmt vet test build

## help: Display this help message
help:
	@echo "GoKanon Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /'

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "Binary built: $(BUILD_DIR)/$(BINARY_NAME)"

## build-all: Build binaries for all platforms
build-all:
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(CMD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)
	@echo "All binaries built in $(BUILD_DIR)/"

## install: Install the binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GOINSTALL) $(LDFLAGS) $(CMD_DIR)
	@echo "Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

## uninstall: Remove the installed binary
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(shell go env GOPATH)/bin/$(BINARY_NAME)
	@echo "Uninstalled $(BINARY_NAME)"

## test: Run tests
test:
	@echo "Running tests..."
	$(GOTEST) $(TEST_FLAGS) ./...

## test-short: Run tests in short mode (skip long tests)
test-short:
	@echo "Running tests in short mode..."
	$(GOTEST) -short $(TEST_FLAGS) ./...

## test-verbose: Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	$(GOTEST) -v -race ./...

## test-dashboard: Run dashboard package tests only
test-dashboard:
	@echo "Running dashboard tests..."
	$(GOTEST) $(TEST_FLAGS) ./internal/dashboard/...

## bench: Run benchmarks
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

## bench-dashboard: Run dashboard benchmarks only
bench-dashboard:
	@echo "Running dashboard benchmarks..."
	$(GOTEST) -bench=. -benchmem ./internal/dashboard/...

## coverage: Generate test coverage report
coverage:
	@echo "Generating coverage report..."
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) ./... $(COVERAGE_FLAGS)
	$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report: $(COVERAGE_DIR)/coverage.html"
	@$(GOCMD) tool cover -func=$(COVERAGE_DIR)/coverage.out | grep total | awk '{print "Total coverage: " $$3}'

## coverage-dashboard: Generate coverage for dashboard package
coverage-dashboard:
	@echo "Generating dashboard coverage..."
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) ./internal/dashboard/... -coverprofile=$(COVERAGE_DIR)/dashboard-coverage.out -covermode=atomic
	$(GOCMD) tool cover -html=$(COVERAGE_DIR)/dashboard-coverage.out -o $(COVERAGE_DIR)/dashboard-coverage.html
	@echo "Dashboard coverage: $(COVERAGE_DIR)/dashboard-coverage.html"
	@$(GOCMD) tool cover -func=$(COVERAGE_DIR)/dashboard-coverage.out | grep total | awk '{print "Dashboard coverage: " $$3}'

## fmt: Format Go code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...
	@echo "Code formatted"

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...
	@echo "Vet complete"

## lint: Run golangci-lint (requires golangci-lint installed)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install it with:"; \
		echo "  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin"; \
	fi

## mod-download: Download dependencies
mod-download:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	@echo "Dependencies downloaded"

## mod-tidy: Tidy up dependencies
mod-tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy
	@echo "Dependencies tidied"

## mod-verify: Verify dependencies
mod-verify:
	@echo "Verifying dependencies..."
	$(GOMOD) verify
	@echo "Dependencies verified"

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -rf $(COVERAGE_DIR)
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out
	@echo "Cleaned"

## clean-test-cache: Clean test cache
clean-test-cache:
	@echo "Cleaning test cache..."
	$(GOCMD) clean -testcache
	@echo "Test cache cleaned"

## run: Build and run the binary
run: build
	@echo "Running $(BINARY_NAME)..."
	$(BUILD_DIR)/$(BINARY_NAME)

## run-serve: Build and run the dashboard server
run-serve: build
	@echo "Starting dashboard server..."
	$(BUILD_DIR)/$(BINARY_NAME) serve

## run-example: Run example benchmarks
run-example: build
	@echo "Running example benchmarks..."
	$(BUILD_DIR)/$(BINARY_NAME) run -pkg=./examples

## dev: Development mode - build and run tests
dev: clean fmt vet test build
	@echo "Development build complete"

## ci: CI mode - full validation
ci: clean mod-verify fmt vet test coverage
	@echo "CI validation complete"

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t gokanon:$(VERSION) .
	@echo "Docker image built: gokanon:$(VERSION)"

## check: Quick validation (fmt, vet, test)
check: fmt vet test-short
	@echo "Quick check complete"

## release: Prepare release build
release: clean test build-all
	@echo "Release build complete"
	@echo "Binaries in $(BUILD_DIR)/"
	@ls -lh $(BUILD_DIR)/

## info: Display build information
info:
	@echo "GoKanon Build Information"
	@echo "========================="
	@echo "Binary name:  $(BINARY_NAME)"
	@echo "Version:      $(VERSION)"
	@echo "Build dir:    $(BUILD_DIR)"
	@echo "GOOS:         $(GOOS)"
	@echo "GOARCH:       $(GOARCH)"
	@echo "Go version:   $$(go version)"
	@echo ""
	@echo "Targets:"
	@echo "  make build          - Build binary for current platform"
	@echo "  make test           - Run all tests"
	@echo "  make coverage       - Generate coverage report"
	@echo "  make install        - Install to GOPATH/bin"
	@echo "  make help           - Show all targets"

## deps: Install development dependencies
deps:
	@echo "Installing development dependencies..."
	@if ! command -v golangci-lint > /dev/null; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin; \
	fi
	@echo "Dependencies installed"

## upgrade: Upgrade dependencies
upgrade:
	@echo "Upgrading dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy
	@echo "Dependencies upgraded"

## watch: Watch for changes and rebuild (requires entr)
watch:
	@if command -v entr > /dev/null; then \
		echo "Watching for changes..."; \
		find . -name '*.go' | entr -c make check; \
	else \
		echo "entr not installed. Install it with:"; \
		echo "  apt-get install entr  (Debian/Ubuntu)"; \
		echo "  brew install entr     (macOS)"; \
	fi

## stats: Show code statistics
stats:
	@echo "Code Statistics"
	@echo "==============="
	@echo "Go files:     $$(find . -name '*.go' | wc -l)"
	@echo "Lines of code: $$(find . -name '*.go' -exec cat {} \; | wc -l)"
	@echo "Test files:   $$(find . -name '*_test.go' | wc -l)"
	@echo ""
	@echo "Package breakdown:"
	@for dir in ./internal/*/; do \
		pkg=$$(basename $$dir); \
		count=$$(find $$dir -name "*.go" | wc -l); \
		echo "  $$pkg: $$count files"; \
	done

.DEFAULT_GOAL := help

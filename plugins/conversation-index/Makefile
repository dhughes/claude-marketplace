.PHONY: all build test clean install

# Binary output locations
INDEXER_BIN=bin/cidx-index
SEARCH_BIN=bin/cidx-search

# Go source files
GO_FILES=$(shell find . -name '*.go' -type f)

# Use local cache to avoid sandbox permission issues
export GOCACHE=$(shell pwd)/.go-cache

# Build both binaries for current architecture
all: build

build: $(INDEXER_BIN) $(SEARCH_BIN)

$(INDEXER_BIN): $(GO_FILES)
	@mkdir -p bin
	go build -o $(INDEXER_BIN) ./cmd/cidx-index

$(SEARCH_BIN): $(GO_FILES)
	@mkdir -p bin
	go build -o $(SEARCH_BIN) ./cmd/cidx-search

# Build universal binaries for macOS (arm64 + amd64)
universal: clean
	@echo "Building for macOS arm64..."
	@mkdir -p .build
	GOOS=darwin GOARCH=arm64 go build -o .build/cidx-index-arm64 ./cmd/cidx-index
	GOOS=darwin GOARCH=amd64 go build -o .build/cidx-index-amd64 ./cmd/cidx-index
	GOOS=darwin GOARCH=arm64 go build -o .build/cidx-search-arm64 ./cmd/cidx-search
	GOOS=darwin GOARCH=amd64 go build -o .build/cidx-search-amd64 ./cmd/cidx-search

	@echo "Creating universal binaries..."
	@mkdir -p bin
	lipo -create -output $(INDEXER_BIN) .build/cidx-index-arm64 .build/cidx-index-amd64
	lipo -create -output $(SEARCH_BIN) .build/cidx-search-arm64 .build/cidx-search-amd64
	chmod +x $(INDEXER_BIN) $(SEARCH_BIN)

	@echo "Universal binaries created:"
	@file $(INDEXER_BIN)
	@file $(SEARCH_BIN)

# Run tests
test:
	go test -v -race ./...

# Run tests with coverage
coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Install binaries to plugin scripts directory (for local testing)
install: build
	@echo "Installing binaries to scripts directories..."
	cp $(INDEXER_BIN) scripts/
	cp $(SEARCH_BIN) skills/conversation-search/scripts/
	chmod +x scripts/cidx-index skills/conversation-search/scripts/cidx-search
	@echo "Installed successfully"

# Clean build artifacts
clean:
	rm -rf bin/ .build/
	rm -f coverage.out coverage.html

# Development: build and run indexer
run-index: build
	$(INDEXER_BIN) --full-reindex

# Development: build and run search
run-search: build
	$(SEARCH_BIN) --json "test"

# Check Go mod dependencies
deps:
	go mod tidy
	go mod verify

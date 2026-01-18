.PHONY: all clean darwin-amd64 darwin-arm64 linux-amd64 linux-arm64 test

ORCHESTRATOR=statusline
TOOLS=status-directory
SRC_DIR=scripts/src
BIN_DIR=scripts/bin
PLATFORMS=darwin-amd64 darwin-arm64 linux-amd64 linux-arm64

all: darwin-amd64 darwin-arm64 linux-amd64 linux-arm64

darwin-amd64:
	@echo "Building orchestrator for darwin-amd64..."
	@cd $(SRC_DIR) && GOOS=darwin GOARCH=amd64 go build -o ../../$(BIN_DIR)/$(ORCHESTRATOR)-darwin-amd64 .
	@echo "Building tools for darwin-amd64..."
	@cd $(SRC_DIR) && GOOS=darwin GOARCH=amd64 go build -o ../../$(BIN_DIR)/status-directory-darwin-amd64 ./cmd/tools/directory
	@cd $(SRC_DIR) && GOOS=darwin GOARCH=amd64 go build -o ../../$(BIN_DIR)/status-git-branch-darwin-amd64 ./cmd/tools/git-branch
	@cd $(SRC_DIR) && GOOS=darwin GOARCH=amd64 go build -o ../../$(BIN_DIR)/status-model-darwin-amd64 ./cmd/tools/model

darwin-arm64:
	@echo "Building orchestrator for darwin-arm64..."
	@cd $(SRC_DIR) && GOOS=darwin GOARCH=arm64 go build -o ../../$(BIN_DIR)/$(ORCHESTRATOR)-darwin-arm64 .
	@echo "Building tools for darwin-arm64..."
	@cd $(SRC_DIR) && GOOS=darwin GOARCH=arm64 go build -o ../../$(BIN_DIR)/status-directory-darwin-arm64 ./cmd/tools/directory
	@cd $(SRC_DIR) && GOOS=darwin GOARCH=arm64 go build -o ../../$(BIN_DIR)/status-git-branch-darwin-arm64 ./cmd/tools/git-branch
	@cd $(SRC_DIR) && GOOS=darwin GOARCH=arm64 go build -o ../../$(BIN_DIR)/status-model-darwin-arm64 ./cmd/tools/model

linux-amd64:
	@echo "Building orchestrator for linux-amd64..."
	@cd $(SRC_DIR) && GOOS=linux GOARCH=amd64 go build -o ../../$(BIN_DIR)/$(ORCHESTRATOR)-linux-amd64 .
	@echo "Building tools for linux-amd64..."
	@cd $(SRC_DIR) && GOOS=linux GOARCH=amd64 go build -o ../../$(BIN_DIR)/status-directory-linux-amd64 ./cmd/tools/directory
	@cd $(SRC_DIR) && GOOS=linux GOARCH=amd64 go build -o ../../$(BIN_DIR)/status-git-branch-linux-amd64 ./cmd/tools/git-branch
	@cd $(SRC_DIR) && GOOS=linux GOARCH=amd64 go build -o ../../$(BIN_DIR)/status-model-linux-amd64 ./cmd/tools/model

linux-arm64:
	@echo "Building orchestrator for linux-arm64..."
	@cd $(SRC_DIR) && GOOS=linux GOARCH=arm64 go build -o ../../$(BIN_DIR)/$(ORCHESTRATOR)-linux-arm64 .
	@echo "Building tools for linux-arm64..."
	@cd $(SRC_DIR) && GOOS=linux GOARCH=arm64 go build -o ../../$(BIN_DIR)/status-directory-linux-arm64 ./cmd/tools/directory
	@cd $(SRC_DIR) && GOOS=linux GOARCH=arm64 go build -o ../../$(BIN_DIR)/status-git-branch-linux-arm64 ./cmd/tools/git-branch
	@cd $(SRC_DIR) && GOOS=linux GOARCH=arm64 go build -o ../../$(BIN_DIR)/status-model-linux-arm64 ./cmd/tools/model

clean:
	@echo "Cleaning binaries..."
	@rm -f $(BIN_DIR)/$(ORCHESTRATOR)-*
	@rm -f $(BIN_DIR)/status-*

test-go:
	@echo "Running Go unit tests..."
	@cd $(SRC_DIR) && go test -v ./...

test-integration:
	@echo "Testing statusline script..."
	@echo '{}' | ./scripts/statusline.sh

test: test-go test-integration

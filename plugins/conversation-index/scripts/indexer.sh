#!/bin/bash
# Go-based indexer hook script
# Called by PostToolUse hook to incrementally index conversations

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INDEXER_BIN="${SCRIPT_DIR}/cidx-index"
LOCK_FILE="${HOME}/.claude/conversation-index.lock"

# Check if binary exists
if [ ! -x "$INDEXER_BIN" ]; then
    echo "Error: Indexer binary not found at $INDEXER_BIN" >&2
    echo "Run 'make install' to build and install binaries" >&2
    exit 1
fi

# Use flock to ensure only one indexer runs at a time
# Exit immediately if lock can't be acquired (another instance is running)
exec 200>"$LOCK_FILE"
if ! flock -n 200; then
    # Another indexer is already running, exit silently
    exit 0
fi

# Run indexer (pass through all arguments)
"$INDEXER_BIN" "$@"
EXIT_CODE=$?

# Release lock
flock -u 200

exit $EXIT_CODE

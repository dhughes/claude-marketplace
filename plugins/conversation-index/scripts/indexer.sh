#!/bin/bash
# Go-based indexer hook script
# Called by PostToolUse hook to incrementally index conversations

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INDEXER_BIN="${SCRIPT_DIR}/cidx-index"

# Check if binary exists
if [ ! -x "$INDEXER_BIN" ]; then
    echo "Error: Indexer binary not found at $INDEXER_BIN" >&2
    echo "Run 'make install' to build and install binaries" >&2
    exit 1
fi

# Run indexer (pass through all arguments)
exec "$INDEXER_BIN" "$@"

#!/bin/bash
# Go-based search script
# Called by conversation-search skill

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SEARCH_BIN="${SCRIPT_DIR}/cidx-search"

# Check if binary exists
if [ ! -x "$SEARCH_BIN" ]; then
    echo "Error: Search binary not found at $SEARCH_BIN" >&2
    echo "Run 'make install' to build and install binaries" >&2
    exit 1
fi

# Run search (pass through all arguments)
exec "$SEARCH_BIN" "$@"

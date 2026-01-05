#!/bin/bash
# Initialization script that runs on SessionStart
# Builds Go binaries and creates index if needed

PLUGIN_ROOT="${CLAUDE_PLUGIN_ROOT:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}"
DB_PATH="${HOME}/.claude/conversation-index.db"

cd "$PLUGIN_ROOT"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21+ from https://golang.org/dl/" >&2
    echo "  macOS: brew install go" >&2
    exit 1
fi

# Check if binaries need to be built
NEED_BUILD=false

if [ ! -x "scripts/cidx-index" ] || [ ! -x "skills/conversation-search/scripts/cidx-search" ]; then
    NEED_BUILD=true
fi

# Build binaries if needed
if [ "$NEED_BUILD" = true ]; then
    echo "[1/2] Building conversation-index binaries..." >&2

    if ! make install > /dev/null 2>&1; then
        echo "Error: Failed to build Go binaries" >&2
        echo "Try running 'make install' manually in $PLUGIN_ROOT" >&2
        exit 1
    fi

    echo "[2/2] Building conversation index..." >&2
    INDEX_OUTPUT=$(scripts/cidx-index --full-reindex 2>&1)

    # Return JSON response that Claude can see
    cat << EOF
{
  "hookSpecificOutput": {
    "hookEventName": "SessionStart",
    "additionalContext": "IMPORTANT: Display this to the user:\n\nconversation-index plugin initialized (Go version):\n✓ Built native Go binaries\n✓ ${INDEX_OUTPUT}\n\nYou can now search conversation history by asking questions like:\n- 'when did we first discuss X?'\n- 'find conversations about Y'\n- 'show me conversations mentioning Z'"
  }
}
EOF
    exit 0
fi

# Check if database needs initialization
if [ ! -f "$DB_PATH" ]; then
    INDEX_OUTPUT=$(scripts/cidx-index --full-reindex 2>&1)
    cat << EOF
{
  "hookSpecificOutput": {
    "hookEventName": "SessionStart",
    "additionalContext": "NOTE: Display this to the user:\n\nRebuilt conversation search index:\n${INDEX_OUTPUT}\n\nconversation-index plugin ready."
  }
}
EOF
fi

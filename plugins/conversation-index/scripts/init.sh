#!/bin/bash

# Initialization script that runs on SessionStart
# Installs dependencies and creates index if needed

PLUGIN_ROOT="${CLAUDE_PLUGIN_ROOT:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}"
DB_PATH="${HOME}/.claude/conversation-index.db"

cd "$PLUGIN_ROOT"

# Check if dependencies are installed
if [ ! -d "node_modules/better-sqlite3" ]; then
  # Run installation (stderr will show in verbose mode)
  echo "[1/3] Creating package.json..." >&2
  npm init -y > /dev/null 2>&1

  echo "[2/3] Installing dependencies (better-sqlite3)..." >&2
  npm install better-sqlite3 --save --loglevel=error > /dev/null 2>&1

  echo "[3/3] Building conversation index..." >&2
  INDEX_OUTPUT=$(./scripts/indexer.sh --full-reindex 2>&1)

  # Return JSON response that Claude can see
  cat << EOF
{
  "hookSpecificOutput": {
    "hookEventName": "SessionStart",
    "additionalContext": "conversation-index plugin initialized on first run:\n- Installed better-sqlite3 dependency\n- Built conversation index\n- Ready to search conversation history\n\n${INDEX_OUTPUT}\n\nUsers can now ask: 'when did we first discuss X?' or 'find conversations about Y'"
  }
}
EOF
  exit 0
fi

# Check if database needs initialization
if [ ! -f "$DB_PATH" ]; then
  INDEX_OUTPUT=$(./scripts/indexer.sh --full-reindex 2>&1)
  cat << EOF
{
  "hookSpecificOutput": {
    "hookEventName": "SessionStart",
    "additionalContext": "conversation-index: Rebuilt index on session start.\n\n${INDEX_OUTPUT}"
  }
}
EOF
fi

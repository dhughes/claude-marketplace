#!/bin/bash

# Silent initialization script that runs on SessionStart
# Installs dependencies and creates index if needed

PLUGIN_ROOT="${CLAUDE_PLUGIN_ROOT:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}"
DB_PATH="${HOME}/.claude/conversation-index.db"

cd "$PLUGIN_ROOT"

# Check if dependencies are installed
if [ ! -d "node_modules/better-sqlite3" ]; then
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo "  Initializing conversation-index plugin (first run)"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo ""

  if [ ! -f "package.json" ]; then
    echo "[1/3] Creating package.json..."
    npm init -y > /dev/null 2>&1
  fi

  echo "[2/3] Installing dependencies (better-sqlite3)..."
  npm install better-sqlite3 --save --loglevel=error 2>&1 | grep -v "^npm"

  echo "[3/3] Building conversation index..."
  echo "       This may take 10-30 seconds..."
  ./scripts/indexer.sh --full-reindex 2>&1 | tail -1

  echo ""
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo "  ✓ Plugin ready. Try: 'when did we first discuss X?'"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo ""
  exit 0
fi

# Check if database needs initialization
if [ ! -f "$DB_PATH" ]; then
  echo "Initializing conversation index (first use)..."
  ./scripts/indexer.sh --full-reindex 2>&1 | tail -1
fi

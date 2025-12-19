#!/bin/bash

set -e

PLUGIN_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  Installing conversation-index plugin"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Install npm dependencies
cd "$PLUGIN_DIR"
if [ ! -f "package.json" ]; then
  echo "[1/3] Creating package.json..."
  npm init -y > /dev/null
fi

echo "[2/3] Installing dependencies (better-sqlite3)..."
npm install better-sqlite3 --save --loglevel=error

# Initialize database
echo "[3/3] Building initial conversation index..."
echo "       This may take 10-30 seconds depending on conversation history..."
./scripts/indexer.sh --full-reindex

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  ✓ Conversation index plugin installed successfully"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "The index will update automatically after each interaction."
echo "Try asking: 'when did we first discuss X?'"
echo ""

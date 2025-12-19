#!/bin/bash

set -e

PLUGIN_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "Installing conversation-index plugin..."

# Install npm dependencies
cd "$PLUGIN_DIR"
if [ ! -f "package.json" ]; then
  npm init -y > /dev/null
fi

echo "Installing dependencies..."
npm install better-sqlite3 --save

# Initialize database
echo "Initializing conversation index database..."
node scripts/indexer.sh --full-reindex

echo "✓ Conversation index plugin installed successfully"
echo ""
echo "The indexer will run automatically after each tool call via hooks."
echo "You can manually rebuild the index with: node scripts/indexer.sh --full-reindex"

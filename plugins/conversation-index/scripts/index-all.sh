#!/bin/bash

# Conversation Index - Bulk Historical Indexing Script
# Indexes all historical conversations from ~/.claude/projects/
# Run this once after installing the plugin to index your history

set -e

echo "🔍 Bulk indexing historical conversations..."

# Check dependencies
if ! command -v sqlite3 &> /dev/null; then
    echo "Error: sqlite3 not found" >&2
    exit 1
fi

if ! command -v jq &> /dev/null; then
    echo "Error: jq not found. Install: brew install jq" >&2
    exit 1
fi

PROJECTS_DIR="$HOME/.claude/projects"
DB_PATH="$HOME/.claude/conversation_index.db"

# Count total files to process
TOTAL_FILES=$(find "$PROJECTS_DIR" -name "*.jsonl" ! -name "agent-*.jsonl" -type f | wc -l | tr -d ' ')
echo "Found $TOTAL_FILES conversation files to index"

PROCESSED=0

# Process each project directory
find "$PROJECTS_DIR" -name "*.jsonl" ! -name "agent-*.jsonl" -type f | while read -r transcript_path; do
    SESSION_ID=$(basename "$transcript_path" .jsonl)

    PROCESSED=$((PROCESSED + 1))
    echo "[$PROCESSED/$TOTAL_FILES] Indexing session: $SESSION_ID"

    # Call the index script for this session
    echo "{\"transcript_path\":\"$transcript_path\",\"session_id\":\"$SESSION_ID\"}" | \
        "$(dirname "$0")/index.sh" 2>/dev/null || echo "  ⚠️  Failed to index $SESSION_ID"
done

# Show summary
TOTAL_MESSAGES=$(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM messages" 2>/dev/null)
TOTAL_SESSIONS=$(sqlite3 "$DB_PATH" "SELECT COUNT(DISTINCT session_id) FROM messages" 2>/dev/null)

echo ""
echo "✅ Indexing complete!"
echo "   Total messages indexed: $TOTAL_MESSAGES"
echo "   Total sessions indexed: $TOTAL_SESSIONS"
echo ""
echo "You can now search with:"
echo "  $(dirname "$0")/search.sh \"your search term\""

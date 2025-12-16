#!/bin/bash

# Conversation Index - Search Script
# Query the SQLite FTS5 index for conversations matching a search term
# Usage: search.sh <search_term> [limit]

set -e

SEARCH_TERM="$1"
LIMIT="${2:-20}"

if [ -z "$SEARCH_TERM" ]; then
    echo "Usage: $0 <search_term> [limit]" >&2
    echo "Example: $0 \"zeebe workers\" 10" >&2
    exit 1
fi

# Check dependencies
if ! command -v sqlite3 &> /dev/null; then
    echo "Error: sqlite3 not found" >&2
    exit 1
fi

DB_PATH="$HOME/.claude/conversation_index.db"

# Check if database exists
if [ ! -f "$DB_PATH" ]; then
    echo "No conversation index found. The index is built automatically as you use Claude Code." >&2
    exit 1
fi

# Perform FTS5 search
# Group by session_id and show:
# - Session ID
# - Earliest timestamp
# - Number of matching messages
# - Preview of first matching content

sqlite3 "$DB_PATH" -column -header <<EOF
SELECT
    session_id,
    MIN(timestamp) as first_timestamp,
    COUNT(*) as match_count,
    SUBSTR(
        (SELECT content FROM messages_fts
         WHERE session_id = m.session_id
           AND messages_fts MATCH '$SEARCH_TERM'
         LIMIT 1),
        1, 200
    ) as preview
FROM messages_fts m
WHERE messages_fts MATCH '$SEARCH_TERM'
GROUP BY session_id
ORDER BY MIN(timestamp) DESC
LIMIT $LIMIT;
EOF

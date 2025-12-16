#!/bin/bash

# Conversation Index - Indexing Script
# Called by SessionStart, UserPromptSubmit, Stop, and SessionEnd hooks
# Idempotently indexes new messages from conversation transcripts into SQLite FTS5

set -e

# Read hook input from stdin
HOOK_INPUT=$(cat)

# Check dependencies - fail gracefully
if ! command -v sqlite3 &> /dev/null; then
    echo "⚠️  conversation-index: sqlite3 not found. Indexing disabled." >&2
    echo "   Install: Pre-installed on macOS/Linux (check your installation)" >&2
    exit 0  # Exit 0 so we don't block the hook chain
fi

if ! command -v jq &> /dev/null; then
    echo "⚠️  conversation-index: jq not found. Indexing disabled." >&2
    echo "   Install: brew install jq (macOS) or apt-get install jq (Linux)" >&2
    exit 0  # Exit 0 so we don't block the hook chain
fi

# Parse hook input
TRANSCRIPT_PATH=$(echo "$HOOK_INPUT" | jq -r '.transcript_path // empty')
SESSION_ID=$(echo "$HOOK_INPUT" | jq -r '.session_id // empty')

# If we don't have a transcript path or session ID, nothing to index
if [ -z "$TRANSCRIPT_PATH" ] || [ -z "$SESSION_ID" ]; then
    exit 0
fi

# Database location
DB_PATH="$HOME/.claude/conversation_index.db"

# Initialize database if it doesn't exist
sqlite3 "$DB_PATH" <<'EOF'
-- Main messages table
CREATE TABLE IF NOT EXISTS messages (
    uuid TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    message_type TEXT NOT NULL,
    timestamp TEXT,
    content TEXT NOT NULL
);

-- FTS5 virtual table for full-text search
CREATE VIRTUAL TABLE IF NOT EXISTS messages_fts USING fts5(
    uuid,
    session_id,
    message_type,
    timestamp,
    content,
    content='messages',
    content_rowid='rowid'
);

-- Triggers to keep FTS index in sync
CREATE TRIGGER IF NOT EXISTS messages_ai AFTER INSERT ON messages BEGIN
  INSERT INTO messages_fts(rowid, uuid, session_id, message_type, timestamp, content)
  VALUES (new.rowid, new.uuid, new.session_id, new.message_type, new.timestamp, new.content);
END;

CREATE TRIGGER IF NOT EXISTS messages_ad AFTER DELETE ON messages BEGIN
  INSERT INTO messages_fts(messages_fts, rowid, uuid, session_id, message_type, timestamp, content)
  VALUES('delete', old.rowid, old.uuid, old.session_id, old.message_type, old.timestamp, old.content);
END;

-- Track last indexed message per session
CREATE TABLE IF NOT EXISTS index_state (
    session_id TEXT PRIMARY KEY,
    last_message_uuid TEXT NOT NULL,
    last_indexed_at TEXT DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster lookups
CREATE INDEX IF NOT EXISTS idx_messages_session ON messages(session_id);
CREATE INDEX IF NOT EXISTS idx_messages_timestamp ON messages(timestamp);
EOF

# Process transcript and generate SQL with proper escaping
# Use a single transaction for speed, let SQLite PRIMARY KEY handle duplicates

TEMP_SQL=$(mktemp)

{
    echo "BEGIN TRANSACTION;"

    # Use jq to parse JSONL and generate SQL - much faster and safer than bash loops
    # Only index actual text content (skip tool_results and other non-text content)
    jq -r --arg session "$SESSION_ID" '
        select(.type == "user" or .type == "assistant") |
        {
            uuid: .uuid,
            session_id: $session,
            type: .type,
            timestamp: .timestamp,
            content: (
                if .type == "user" then
                    (if (.message.content | type) == "string"
                     then .message.content
                     else ""
                     end)
                else
                    ([.message.content[]? | select(.type == "text") | .text] | join(" "))
                end
            )
        } |
        select((.content | length) > 0) |
        "INSERT OR IGNORE INTO messages (uuid, session_id, message_type, timestamp, content) VALUES (" +
            (.uuid | @json) + ", " +
            (.session_id | @json) + ", " +
            (.type | @json) + ", " +
            (.timestamp | @json) + ", " +
            (.content | @json) +
        ");"
    ' "$TRANSCRIPT_PATH" 2>/dev/null

    echo "COMMIT;"
} > "$TEMP_SQL"

# Execute in single sqlite3 invocation - very fast!
sqlite3 "$DB_PATH" < "$TEMP_SQL" 2>&1

# Cleanup
rm -f "$TEMP_SQL"

exit 0

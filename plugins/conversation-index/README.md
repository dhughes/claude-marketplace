# Conversation Index Plugin

Fast full-text search across all your Claude Code conversation history using SQLite FTS5.

## Features

- **Automatic Indexing**: Hooks into Claude Code events to index conversations in real-time
- **Fast Search**: FTS5-powered queries return results in milliseconds
- **Idempotent**: Safe to run multiple times, only indexes new messages
- **Cross-Project**: Searches across all your Claude Code projects
- **Rich Results**: Returns session IDs, timestamps, match counts, and content previews

## Requirements

- **sqlite3** - Pre-installed on macOS/Linux
- **jq** - JSON parser for processing conversation transcripts
  - macOS: `brew install jq`
  - Linux: `apt-get install jq` or `yum install jq`

## Installation

```bash
# Add the marketplace
/plugin marketplace add doughughes/claude-marketplace

# Install the plugin
/plugin install conversation-index@doug-marketplace
```

## Usage

### Ask Claude to Search

Just ask Claude naturally:

```
"When did we discuss auto compensation?"
"Find conversations about zeebe workers"
"Search for phase 2 tickets"
```

The `search-conversations` skill will automatically activate and search the index.

### Direct Search

You can also use the search script directly:

```bash
# Basic search
~/.claude/plugins/conversation-index/scripts/search.sh "search term"

# Limit results
~/.claude/plugins/conversation-index/scripts/search.sh "search term" 10

# Advanced FTS5 queries
~/.claude/plugins/conversation-index/scripts/search.sh "zeebe AND worker"
~/.claude/plugins/conversation-index/scripts/search.sh '"exact phrase"'
```

## How It Works

### Automatic Indexing

The plugin uses hooks to index conversations automatically:

- **SessionStart**: Catches up on any missed indexing when you start Claude
- **UserPromptSubmit**: Indexes after each prompt you send
- **Stop**: Indexes after Claude finishes responding
- **SessionEnd**: Final cleanup when you quit or `/clear`

All hooks call the same `index.sh` script, which is idempotent and only indexes new messages.

### Database Schema

Located at `~/.claude/conversation_index.db`:

```sql
-- Main messages table
messages (
  uuid TEXT PRIMARY KEY,
  session_id TEXT,
  message_type TEXT,  -- 'user' or 'assistant'
  timestamp TEXT,
  content TEXT
)

-- FTS5 virtual table for full-text search
messages_fts USING fts5(uuid, session_id, message_type, timestamp, content)

-- Tracks indexing progress per session
index_state (
  session_id TEXT PRIMARY KEY,
  last_message_uuid TEXT,
  last_indexed_at TEXT
)
```

### Search Performance

- **Indexing**: ~5ms per message (incremental)
- **Search**: <10ms even with thousands of conversations
- **Storage**: ~1-2KB per message

## FTS5 Query Syntax

The search supports advanced SQLite FTS5 query syntax:

| Query | Meaning |
|-------|---------|
| `zeebe worker` | Both terms (implicit AND) |
| `zeebe OR worker` | Either term |
| `zeebe NOT worker` | First but not second |
| `"phase 2"` | Exact phrase |
| `zeeb*` | Prefix match |
| `NEAR(term1 term2, 5)` | Terms within 5 words |

## Troubleshooting

### "sqlite3 not found" or "jq not found"

Install missing dependencies:

```bash
# macOS
brew install jq

# Linux (Debian/Ubuntu)
sudo apt-get install jq

# Linux (RHEL/CentOS)
sudo yum install jq
```

`sqlite3` should be pre-installed on macOS and most Linux distributions.

### "No conversation index found"

The index is built automatically as you use Claude Code. If you just installed the plugin:

1. Start a new Claude session
2. Send a few messages
3. The index will be created at `~/.claude/conversation_index.db`

### Search returns no results

- Try alternative search terms
- Very recent messages (within the last second) might not be indexed yet
- Check if the conversation happened in a different project

### Check index status

```bash
# See if database exists
ls -lh ~/.claude/conversation_index.db

# Count indexed messages
sqlite3 ~/.claude/conversation_index.db "SELECT COUNT(*) FROM messages"

# List indexed sessions
sqlite3 ~/.claude/conversation_index.db "SELECT DISTINCT session_id FROM messages"
```

## Development

### File Structure

```
conversation-index/
├── .claude-plugin/
│   └── plugin.json           # Plugin manifest
├── hooks/
│   └── hooks.json            # Hook definitions
├── scripts/
│   ├── index.sh              # Indexing script (called by hooks)
│   └── search.sh             # Search script (called by skill)
├── skills/
│   └── search-conversations/
│       └── SKILL.md          # Skill definition
└── README.md
```

### Testing Changes

```bash
# Reinstall after changes
/plugin uninstall conversation-index
/plugin install conversation-index@doug-marketplace --scope local

# Or just restart Claude if only scripts changed
```

### Manual Indexing

If you want to manually trigger indexing:

```bash
# Simulate a hook call
echo '{"transcript_path":"path/to/file.jsonl","session_id":"uuid"}' | \
  ~/.claude/plugins/conversation-index/scripts/index.sh
```

## License

MIT

## Author

Doug Hughes (doug.hughes@ezcater.com)

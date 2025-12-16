---
name: search-conversations
description: Search through indexed Claude Code conversation history using full-text search. Activates when user asks to find past conversations, search conversation history, or locate when topics were discussed. Returns session IDs with timestamps and previews.
---

# Search Conversations

Fast full-text search across all indexed Claude Code conversations using SQLite FTS5. This skill searches the conversation index that is automatically built by the conversation-index plugin.

## Activation

Use this skill when the user asks questions like:
- "When did we discuss [topic]?"
- "Find conversations about [topic]"
- "Search my conversation history for [topic]"
- "Which session covered [topic]?"
- "Show me conversations where we talked about [topic]"

## How It Works

The conversation-index plugin automatically indexes conversations in the background using hooks:
- **SessionStart**: Catches up on any missed indexing
- **UserPromptSubmit**: Indexes after each user prompt
- **Stop**: Indexes after Claude finishes responding
- **SessionEnd**: Final cleanup when session ends

All messages are indexed into a SQLite FTS5 database at `~/.claude/conversation_index.db`.

## Instructions

### Step 1: Run the Search

Use the pre-approved search script to query the FTS index:

```bash
~/.claude/plugins/conversation-index/scripts/search.sh "<search_term>" [limit]
```

**Parameters:**
- `search_term`: The topic to search for (supports FTS5 query syntax)
- `limit`: Optional, number of results to return (default: 20)

**Examples:**
```bash
# Basic search
~/.claude/plugins/conversation-index/scripts/search.sh "zeebe workers"

# Search with custom limit
~/.claude/plugins/conversation-index/scripts/search.sh "auto compensation" 5

# FTS5 query syntax (AND, OR, NOT, phrases)
~/.claude/plugins/conversation-index/scripts/search.sh "zeebe AND worker"
~/.claude/plugins/conversation-index/scripts/search.sh '"phase 2" OR "phase two"'
```

**Important Notes:**
- This script is PRE-APPROVED and executes without user confirmation
- FTS5 queries are very fast (milliseconds even across thousands of messages)
- Results are sorted by date (newest first)

### Step 2: Parse and Present Results

The search script returns results in column format:

```
session_id                            first_timestamp      match_count  preview
d6e5bd36-51d2-4010-b194-5370758b3992  2025-12-15 19:07:12  24          I need you to review how Zeebe workers are created...
190d48e9-b342-43ee-adc6-3741237ad51d  2025-12-12 14:25:09  7           I just turned on the FX::ExcludeInvalid...
```

Format these for the user in a clear, readable way:

```
Found 2 conversations about "zeebe workers":

1. Session: d6e5bd36-51d2-4010-b194-5370758b3992
   Date: 2025-12-15 19:07
   Matches: 24 messages
   Preview: I need you to review how Zeebe workers are created...

2. Session: 190d48e9-b342-43ee-adc6-3741237ad51d
   Date: 2025-12-12 14:25
   Matches: 7 messages
   Preview: I just turned on the FX::ExcludeInvalid...
```

### Step 3: Provide Session File Paths

If the user wants to view a conversation, tell them the file location:

```
The conversation is stored at:
~/.claude/projects/<project-dir>/<session-id>.jsonl

To convert the current working directory to project-dir format:
echo "$CWD" | sed 's|^/|-|' | tr '/' '-'

Example:
/Users/doughughes/code/ezcater/ez-rails
→ -Users-doughughes-code-ezcater-ez-rails
```

### Step 4: Handle No Results

If the search returns no results:
- Suggest alternative search terms
- Remind that only indexed conversations are searchable
- Note that very recent messages might not be indexed yet (though hooks make this rare)

## FTS5 Query Syntax

Users can use advanced FTS5 syntax:

| Query | Meaning |
|-------|---------|
| `zeebe worker` | Messages containing both "zeebe" AND "worker" |
| `zeebe OR worker` | Messages containing either term |
| `zeebe NOT worker` | Messages with "zeebe" but not "worker" |
| `"phase 2"` | Exact phrase match |
| `zeeb*` | Prefix match (zeebe, zeebes, etc.) |
| `NEAR(zeebe worker, 5)` | Terms within 5 tokens of each other |

## Examples

### Example 1: Basic Search
```
User: "When did we discuss auto compensation?"

Claude:
1. Runs: ~/.claude/plugins/conversation-index/scripts/search.sh "auto compensation"
2. Parses results
3. Presents formatted list of matching sessions
```

### Example 2: No Results
```
User: "Find conversations about kubernetes pods"

Claude: Runs search, gets no results

Response: "I didn't find any indexed conversations about 'kubernetes pods'.
Try alternative terms like 'k8s', 'containers', or 'deployments'."
```

### Example 3: Advanced Query
```
User: "Search for conversations about phase 1 or phase 2"

Claude: Runs: ~/.claude/plugins/conversation-index/scripts/search.sh '"phase 1" OR "phase 2"'
```

## Troubleshooting

### Database Not Found
If the search script reports no index:
- The plugin needs time to index conversations
- Indexing happens automatically via hooks
- Check `~/.claude/conversation_index.db` exists

### Missing Dependencies
If indexing fails:
- Requires `sqlite3` (pre-installed on macOS/Linux)
- Requires `jq` (install: `brew install jq`)
- Errors are printed to stderr during hook execution

### Empty Results
- Very recent messages might not be indexed yet
- Try alternative search terms
- Check if the conversation was in the current or a different project

## Performance

- **Index building**: Fast, incremental (~5ms per message)
- **Search queries**: Very fast (<10ms even with thousands of conversations)
- **Storage**: Minimal (~1-2KB per message with compression)

## Technical Details

**Database Schema:**
- `messages` table: UUID, session_id, message_type, timestamp, content
- `messages_fts` FTS5 virtual table: Full-text index over content
- `index_state` table: Tracks last indexed message per session

**Hook Triggers:**
- All hooks call the same idempotent `index.sh` script
- Script tracks last indexed UUID per session
- Only indexes new messages since last run
- Safe to call multiple times (idempotent)

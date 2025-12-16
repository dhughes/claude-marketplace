# Doug's Claude Code Marketplace

Personal collection of Claude Code plugins for productivity and workflow enhancement.

## Installation

Add this marketplace to your Claude Code installation:

```bash
# From GitHub (when published)
/plugin marketplace add doughughes/claude-marketplace

# From local path (for development)
/plugin marketplace add ~/code/doug/claude-marketplace
```

## Available Plugins

### conversation-index

Fast full-text search across all your Claude Code conversation history using SQLite FTS5.

**Features:**
- Automatic background indexing via hooks
- Fast FTS5-powered search (milliseconds even with thousands of conversations)
- Idempotent indexing (safe to run multiple times)
- Works across all projects
- Session IDs with timestamps and previews

**Install:**
```bash
/plugin install conversation-index@doug-marketplace
```

**Usage:**
Just ask Claude to search:
- "When did we discuss auto compensation?"
- "Find conversations about zeebe workers"
- "Search for phase 2 tickets"

Or use the search script directly:
```bash
~/.claude/plugins/conversation-index/scripts/search.sh "search term"
```

**Requirements:**
- `sqlite3` (pre-installed on macOS/Linux)
- `jq` (install: `brew install jq`)

## Development

### Structure

```
claude-marketplace/
├── .claude-plugin/
│   └── marketplace.json       # Marketplace manifest
├── plugins/
│   └── conversation-index/    # Individual plugins
└── README.md
```

### Adding Plugins

1. Create plugin directory under `plugins/`
2. Add plugin manifest at `plugins/<name>/.claude-plugin/plugin.json`
3. Update `marketplace.json` to reference the new plugin
4. Commit and push

### Testing Locally

```bash
# Add marketplace from local path
/plugin marketplace add ~/code/doug/claude-marketplace

# Install plugin for testing
/plugin install conversation-index@doug-marketplace --scope local

# Remove when done
/plugin uninstall conversation-index
```

## License

MIT

## Author

Doug Hughes (doug.hughes@ezcater.com)

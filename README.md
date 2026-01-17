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

### conversation-loader

Load the full content of previous Claude Code conversations into current context.

**Features:**
- Resume previous conversations with full context
- Works with conversation IDs from conversation-index
- Loads complete conversation history including code and artifacts

**Install:**
```bash
/plugin install conversation-loader@doug-marketplace
```

**Usage:**
```
You: "Load conversation abc-123-def"
```

### ezcater-research

Comprehensive ezCater research toolkit for investigating codebases, architectural decisions, and project histories.

**Features:**
- Deep research agent (Opus) for complex multi-source investigations
- Lightweight research skill for straightforward searches
- Integrates with Glean, Atlassian (Jira/Confluence), GitHub, and Git
- Comprehensive documentation with working examples
- 4 detailed reference guides (Glean, atl CLI, gh CLI, git analysis)
- 3 complete research workflow examples

**Install:**
```bash
/plugin install ezcater-research@doug-marketplace
```

**Usage:**
For deep research:
- "Do deep research on why Liberty uses Omnichannel for feature flags"
- "I need deep research on how seeders work in DM-rails"

For lighter research:
- "Find recent PRs about order management in ez-rails"
- "What's the history of the authentication refactor?"
- "How does the order state machine work?"

**Requirements:**
- Glean MCP server (for ezCater internal docs)
- `atl` CLI (for Jira/Confluence access)
- `gh` CLI (for GitHub operations)
- `git` (for code history analysis)

### tdd-feature-dev

Augments feature development workflows with test-driven development practices. Ensures tests are written early and explicitly planned as separate TODO items.

**Features:**
- Automatically composes with feature-dev workflows
- Tests as explicit, separate TODO items (not buried in implementation)
- TDD-ish approach: tests written alongside code
- Strong emphasis on testing behavior, not implementation
- Never test private functions
- 80% test coverage target
- Project-agnostic (works with any testing framework)

**Install:**
```bash
/plugin install tdd-feature-dev@doug-marketplace
```

**Usage:**
The skill automatically loads when you're implementing features:
```
/feature-dev Add user authentication
```

The TDD skill will compose with feature-dev, ensuring:
- Tests are added as explicit TODO items
- Tests are written before/alongside implementation
- Tests focus on behavior and public APIs
- Coverage targets are met

**Trigger phrases:**
- "implement a feature"
- "build new functionality"
- "add a feature"
- Any use of `/feature-dev`

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

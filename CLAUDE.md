# Claude Code Marketplace - Plugin Development Guide

**This is a MARKETPLACE repository containing multiple Claude Code plugins.** All plugins in this marketplace are stored in the `plugins/` directory at the repository root. Each plugin is a separate subdirectory within `plugins/`.

## Marketplace Structure

```
claude-marketplace/
├── .claude-plugin/
│   └── marketplace.json     # Marketplace registry
├── plugins/                 # ALL PLUGINS GO HERE
│   ├── plugin-name-1/       # Individual plugin directory
│   │   ├── .claude-plugin/
│   │   │   └── plugin.json
│   │   ├── commands/
│   │   ├── agents/
│   │   └── skills/
│   ├── plugin-name-2/       # Another plugin
│   │   └── ...
│   └── plugin-name-3/       # Yet another plugin
│       └── ...
└── CLAUDE.md                # This file
```

**CRITICAL: When creating new plugins, always create them as subdirectories within `plugins/`.**

This project is a marketplace of Claude Code plugins designed to extend Claude's capabilities with specialized tools, agents, skills, and automation.

## Plugin Development Workflow

**CRITICAL REQUIREMENTS:**
1. **Always use the `plugin-dev` plugin for creating new plugins.**
2. **Always create plugins in the `plugins/` directory.**
3. **Always maintain `.claude-plugin/marketplace.json` when creating or updating plugins** (add entries, bump versions, update descriptions).

The `plugin-dev` plugin provides comprehensive support for plugin development:

- **Commands**: `/create-plugin` - Guided plugin creation workflow
- **Agents**:
  - `agent-creator` - Design and implement custom agents
  - `skill-reviewer` - Review and improve skill implementations
  - `plugin-validator` - Validate plugin structure and configuration
- **Skills**:
  - `/command-development` - Create slash commands
  - `/skill-development` - Build new skills
  - `/plugin-settings` - Configure plugin settings and state
  - `/plugin-structure` - Understand plugin architecture
  - `/hook-development` - Create event-driven hooks
  - `/mcp-integration` - Integrate Model Context Protocol servers
  - `/agent-development` - Build autonomous agents

To start creating a plugin, run `/create-plugin` and follow the guided workflow.

## Official Documentation

### Core Plugin Documentation
- **Plugin Development Guide**: https://code.claude.com/docs/en/plugins.md
- **Plugins Reference**: https://code.claude.com/docs/en/plugins-reference.md
- **Plugin Marketplaces**: https://code.claude.com/docs/en/plugin-marketplaces.md
- **Discover & Install Plugins**: https://code.claude.com/docs/en/discover-plugins.md

### Plugin Components
- **Commands**: https://code.claude.com/docs/en/slash-commands.md (includes plugin commands)
- **Skills**: https://code.claude.com/docs/en/skills.md
- **Agents (Subagents)**: https://code.claude.com/docs/en/sub-agents.md
- **Hooks**: https://code.claude.com/docs/en/hooks.md (full reference) and https://code.claude.com/docs/en/hooks-guide.md (quickstart)
- **MCP Integration**: https://code.claude.com/docs/en/mcp.md

### Additional Resources
- **Documentation Map**: https://code.claude.com/docs/en/claude_code_docs_map.md
- **Changelog**: https://github.com/anthropics/claude-code/blob/main/CHANGELOG.md

**IMPORTANT**: This documentation may change. Always verify current documentation URLs and check for updates before implementing features. Use the claude-code-guide agent to look up current documentation when needed.

---

## Individual Plugin Structure

**Note**: The directory layout below shows the structure of a single plugin within the marketplace. Each plugin lives in its own directory under `plugins/` (e.g., `plugins/my-plugin/`).

### Directory Layout (for a single plugin)
```
plugins/my-plugin/           # Plugin lives inside plugins/ directory
├── .claude-plugin/
│   └── plugin.json          # Required manifest (ONLY file in .claude-plugin/)
├── commands/                # Slash commands (optional)
│   └── hello.md
├── agents/                  # Subagents (optional)
│   └── reviewer.md
├── skills/                  # Agent Skills (optional)
│   └── code-review/
│       └── SKILL.md
├── hooks/                   # Event handlers (optional)
│   └── hooks.json
├── .mcp.json                # MCP servers (optional)
├── .lsp.json                # LSP servers (optional)
└── scripts/                 # Supporting files
    └── validate.sh
```

**Critical rules**:
- All plugins MUST be created as subdirectories within `plugins/` at the marketplace root
- Only `plugin.json` goes inside `.claude-plugin/`. All other directories must be at the plugin root.

### Environment Variables
- `${CLAUDE_PLUGIN_ROOT}`: Absolute path to plugin directory (use in scripts, hooks, MCP configs)

---

## Component Details

### 1. Commands (Slash Commands)

**What they are**: Reusable prompts that users invoke with `/plugin-name:command-name`

**Key capabilities**:
- Namespaced to plugin (prevents conflicts)
- Support user arguments via `$ARGUMENTS`, `$1`, `$2`, etc.
- Can execute bash commands
- Can reference files
- Can define inline hooks (PreToolUse, PostToolUse, Stop)

**File format**: Markdown files with YAML frontmatter in `commands/` directory

**Example**:
```markdown
---
description: Greet a user by name with personalized message
---

# Hello Command

Greet the user named "$ARGUMENTS" warmly and ask how you can help them today.
```

**Documentation**: https://code.claude.com/docs/en/slash-commands.md#plugin-commands

---

### 2. Skills (Agent Skills)

**What they are**: Markdown files that teach Claude specialized knowledge; Claude automatically invokes them based on task context (model-invoked, not user-invoked)

**Key capabilities**:
- Auto-discovery when plugin installed
- Triggered by Claude based on description match (no explicit invocation needed)
- Support progressive disclosure (split into SKILL.md + supporting files)
- Restrict tool access with `allowed-tools`
- Run in forked context with `context: fork`
- Define hooks (PreToolUse, PostToolUse, Stop)
- Control visibility (user-invocable, disable-model-invocation)

**File structure**:
```
skills/
└── code-review/
    ├── SKILL.md (required)
    ├── reference.md (optional)
    └── scripts/ (optional)
```

**SKILL.md metadata** (YAML frontmatter):
- `name` (required): Lowercase, hyphens, max 64 chars
- `description` (required): Claude uses this to decide when to apply. Max 1024 chars. Include keywords users would say.
- `allowed-tools`: Comma-separated or YAML list (Read, Grep, Glob, Bash, etc.)
- `context`: Set to `fork` for isolated sub-agent context
- `agent`: Specify agent type when `context: fork` (Explore, Plan, general-purpose, custom)
- `model`: Override model (claude-sonnet-4, claude-opus-4, claude-haiku-4, inherit)
- `hooks`: Define event handlers scoped to Skill
- `user-invocable`: true (default) shows in slash menu, false hides it
- `disable-model-invocation`: true blocks programmatic invocation

**Documentation**: https://code.claude.com/docs/en/skills.md

---

### 3. Agents (Subagents)

**What they are**: Specialized AI assistants that run in isolated contexts with custom system prompts, specific tools, and custom permissions

**Key capabilities**:
- Run in own context window (preserves main conversation context)
- Custom system prompts
- Tool restriction (allowlist/denylist)
- Permission modes (default, acceptEdits, dontAsk, bypassPermissions, plan)
- Model selection (sonnet, opus, haiku, inherit)
- Lifecycle hooks (PreToolUse, PostToolUse, SubagentStart, SubagentStop)
- Load specific Skills
- Claude auto-delegates based on task type

**Metadata fields** (YAML frontmatter):
- `name` (required): Lowercase, hyphens
- `description` (required): When Claude should delegate to this agent
- `tools`: Allowed tools (Read, Grep, Glob, Bash, Write, Edit, WebFetch, WebSearch, etc.)
- `disallowedTools`: Explicitly deny tools
- `model`: sonnet (default), opus, haiku, inherit
- `permissionMode`: default, acceptEdits, dontAsk, bypassPermissions, plan
- `skills`: Comma-separated skill names to load
- `hooks`: Lifecycle hooks (PreToolUse, PostToolUse, SubagentStart, SubagentStop)
- `color`: Visual identification color

**Available tools** (can be restricted):
- **Read-only**: Read, Grep, Glob, Bash (for read-only commands)
- **File ops**: Write, Edit, MultiEdit
- **Web**: WebFetch, WebSearch
- **Task delegation**: Task (for subagents)
- **MCP tools**: Custom tools from connected MCP servers

**Documentation**: https://code.claude.com/docs/en/sub-agents.md

---

### 4. Hooks (Event-Driven Automation)

**What they are**: Event handlers that run automatically in response to Claude Code events (file edits, tool use, session start/end, etc.)

**Key capabilities**:
- Execute bash commands or LLM-based prompts
- Pre-tool validation (PreToolUse) - can block operations
- Post-tool processing (PostToolUse)
- Permission request handling (PermissionRequest)
- User prompt validation (UserPromptSubmit)
- Session lifecycle (SessionStart, SessionEnd)
- Subagent lifecycle (SubagentStart, SubagentStop)
- Stop decision control (Stop, SubagentStop)
- Tool input modification
- Conditional logic with exit codes or JSON output

**Supported hook events**:
- `PreToolUse`: Before Claude uses a tool
- `PostToolUse`: After tool execution succeeds
- `PermissionRequest`: When permission dialog shown
- `UserPromptSubmit`: When user submits prompt
- `Notification`: When notifications sent
- `Stop`: When main agent stops
- `SubagentStop`: When subagent stops
- `SubagentStart`: When subagent starts
- `SessionStart`: Session begins
- `SessionEnd`: Session ends
- `PreCompact`: Before context compaction

**Hook types**:
- `command`: Execute bash script
- `prompt`: LLM-based evaluation (Haiku model, returns JSON)

**Configuration location**: `hooks/hooks.json` or inline in `plugin.json`

**MCP tool naming** for hooks:
- MCP tools follow pattern: `mcp__servername__toolname`
- Example matcher: `"mcp__memory__.*"` to match all memory server tools

**Documentation**:
- Full reference: https://code.claude.com/docs/en/hooks.md
- Quickstart: https://code.claude.com/docs/en/hooks-guide.md

---

### 5. MCP Servers (Model Context Protocol)

**What they are**: External tool and service integrations that provide Claude with access to APIs, databases, and custom tools

**Key capabilities**:
- Connect to 100+ MCP servers (GitHub, Sentry, databases, APIs, etc.)
- Three transport types: HTTP, SSE (deprecated), stdio (local)
- Environmental variable substitution (`${VARIABLE}`, `${VARIABLE:-default}`)
- Three scopes: local (project-specific), project (team via version control), user (across projects)
- OAuth 2.0 authentication
- Dynamic tool updates
- Tool search for large tool sets
- Resource references via @ mentions
- Use as slash commands via MCP prompts

**Configuration location**: `.mcp.json` in plugin root or inline in `plugin.json`

**Transport types**:
- `stdio`: Local executables
- `http`: Remote HTTP endpoints
- `sse`: Server-sent events (deprecated)

**MCP tool naming** (in hooks, permissions):
- Format: `mcp__servername__toolname`
- Example: `mcp__github__search_repositories`

**Documentation**: https://code.claude.com/docs/en/mcp.md

---

## Plugin Development Guidelines

### Planning
1. Start by clearly defining the plugin's purpose and target use cases
2. Identify what components you need: commands, skills, agents, hooks, or MCP servers
3. Consider whether existing marketplace plugins could be extended instead of creating new ones
4. Use the `plugin-dev` skills to understand component architecture before building

### Structure
1. **CRITICAL: Create all new plugins in the `plugins/` directory** at the marketplace root
2. Each plugin must be its own subdirectory (e.g., `plugins/my-plugin/`)
3. Follow the standard plugin directory structure within each plugin directory (see Individual Plugin Structure above)
4. Use `${CLAUDE_PLUGIN_ROOT}` for relative paths in scripts
5. Follow auto-discovery conventions for file naming
6. Only `plugin.json` goes in `.claude-plugin/` directory

### Implementation
1. Keep components focused and single-purpose
2. Write clear, descriptive names that indicate functionality
3. Provide comprehensive examples in agent/skill system prompts
4. Test each component independently before integration
5. Use YAML frontmatter for metadata and configuration
6. **After creating or updating a plugin, immediately update `.claude-plugin/marketplace.json`**:
   - Add new plugin entries when creating plugins
   - Bump version numbers when updating plugins
   - Update descriptions when functionality changes
   - Maintain metadata synchronization

### Quality Standards
1. Validate plugins using the `plugin-validator` agent
2. Write clear documentation in plugin.json descriptions
3. Include example usage in skill/command descriptions
4. Test with real-world scenarios before publishing
5. Follow the principle of progressive disclosure in skills
6. **Verify `.claude-plugin/marketplace.json` is valid JSON after every change**

### Integration
1. **CRITICAL: Immediately register new plugins in `.claude-plugin/marketplace.json`**
2. **CRITICAL: Update marketplace.json whenever plugin versions or descriptions change**
3. Include proper version numbers (semantic versioning)
4. Provide clear installation instructions if needed
5. Document any dependencies or prerequisites

## Best Practices

### Commands
- Use YAML frontmatter for metadata
- Support both interactive and non-interactive modes
- Provide clear error messages and validation

### Skills
- Start with problem statement, not solution
- Use progressive disclosure (reveal details as needed)
- Include concrete examples
- Focus on teaching Claude, not directly solving user problems
- Write descriptions that include keywords users would say

### Agents
- Define clear triggering conditions in description
- Specify available tools explicitly
- Use distinct colors for visual identification
- Provide comprehensive system prompts with examples
- Consider permission modes carefully

### Hooks
- Use prompt-based hooks for complex logic
- Keep PreToolUse hooks fast and focused
- Document hook behavior in plugin.json
- Test hook interactions with other plugins
- Use exit code 2 to block operations

### MCP Servers
- Use environment variables for configuration
- Document required environment variables
- Test with various scopes (local, project, user)
- Handle authentication properly

## Plugin Lifecycle

1. **Ideation** - Define purpose and components
2. **Planning** - Use plugin-dev skills to research and plan
3. **Creation** - Run `/create-plugin` or manually create structure **in the `plugins/` directory**
4. **Development** - Implement components using plugin-dev agents
5. **Validation** - Test with `plugin-validator` agent
6. **Registration** - Add to marketplace.json at the marketplace root (`.claude-plugin/marketplace.json`)
7. **Iteration** - Refine based on real-world usage

## Marketplace Registration & Maintenance

**CRITICAL: Claude must actively maintain the `.claude-plugin/marketplace.json` file whenever plugins are created, updated, or modified.**

### When Creating a New Plugin:
1. Ensure the plugin is located in `plugins/<plugin-name>/`
2. Ensure plugin.json has complete metadata (name, version, description, author)
3. **Immediately add an entry to `.claude-plugin/marketplace.json`** at the marketplace root
4. Test plugin installation and functionality
5. Document any special configuration or setup requirements

### When Updating an Existing Plugin:
1. **Update the version number** in the plugin's `plugin.json` (follow semantic versioning)
2. **Update the version number** in `.claude-plugin/marketplace.json` to match
3. **Update the description** in `.claude-plugin/marketplace.json` if the plugin's functionality has changed
4. Update any other relevant metadata (keywords, dependencies, etc.)
5. Test the updated plugin thoroughly

### Maintenance Responsibilities:
- **Version Bumping**: Always increment version numbers appropriately:
  - **Major** (x.0.0): Breaking changes or major new features
  - **Minor** (0.x.0): New features, backward-compatible
  - **Patch** (0.0.x): Bug fixes, minor tweaks
- **Description Updates**: Keep descriptions accurate and reflective of current functionality
- **Metadata Sync**: Ensure `.claude-plugin/marketplace.json` stays in sync with individual plugin.json files
- **Registry Integrity**: Verify marketplace.json remains valid JSON after every change

## Plugin Manifest (plugin.json)

**Required fields**:
- `name` (string, required): Unique kebab-case identifier, used for namespacing commands
- `version` (string): Semantic versioning (e.g., "2.1.0")
- `description` (string): Brief explanation
- `author` (object): `{ "name": "...", "email": "..." }`

**Optional fields**:
- `homepage`, `repository`, `license`, `keywords`
- `commands` (string|array): Additional command files/directories
- `agents` (string|array): Additional agent files
- `skills` (string|array): Additional skill directories
- `hooks` (string|object): Hook config path or inline
- `mcpServers` (string|object): MCP config
- `lspServers` (string|object): LSP server configs
- `outputStyles` (string|array): Output style files

## Installation & Management Commands

```bash
# Install
claude plugin install my-plugin@marketplace-name --scope user|project|local

# Uninstall
claude plugin uninstall my-plugin@marketplace-name --scope user|project|local

# Enable/Disable
claude plugin enable my-plugin@marketplace-name
claude plugin disable my-plugin@marketplace-name

# Update
claude plugin update my-plugin@marketplace-name

# Validate
claude plugin validate .

# List
claude plugin list

# Marketplace management
claude plugin marketplace add owner/repo          # GitHub
claude plugin marketplace add https://url.git     # Git
claude plugin marketplace add ./local/path        # Local
claude plugin marketplace update
claude plugin marketplace list
```

## Notes

- The marketplace and plugin system are actively evolving - expect changes
- Always verify documentation URLs before implementing features (they may change)
- Use the `claude-code-guide` agent to look up current documentation
- Focus on creating high-quality, focused plugins over feature-rich ones
- Leverage existing plugins as examples, but verify current patterns
- Check the changelog regularly for breaking changes: https://github.com/anthropics/claude-code/blob/main/CHANGELOG.md

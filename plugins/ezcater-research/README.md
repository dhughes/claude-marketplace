# ezCater Research Plugin

Comprehensive research toolkit for investigating ezCater codebases, architectural decisions, and project histories.

## Overview

This plugin provides both autonomous deep research capabilities (agent) and lightweight research methodology (skill) for understanding ezCater systems. It integrates with multiple data sources including Glean, Atlassian (Jira/Confluence), GitHub, and Git to provide comprehensive context about code, decisions, and project evolution.

## Features

- **Deep Research Agent**: Autonomous Opus-powered agent for complex, multi-source investigations
- **Research Skill**: Lightweight methodology for common research tasks using main Claude
- **Multi-Source Integration**: Glean, atl CLI, gh CLI, and git history analysis
- **Comprehensive Methodology**: Proven patterns for architectural decisions, code evolution, and process understanding

## Prerequisites

### Required

1. **Glean MCP Server**: Must be configured for ezCater instance
   - Provides access to internal documentation, Slack, Google Docs, etc.
   - Cannot access Google Docs directly - must use Glean

2. **atl CLI**: Atlassian command-line tool
   - Install: Follow Atlassian documentation
   - Configure with ezCater credentials
   - Used for all Jira and Confluence access
   - Cannot access ezcater.atlassian.net directly - must use atl

3. **gh CLI**: GitHub command-line tool
   - Install: `brew install gh` (macOS) or see GitHub docs
   - Authenticate: `gh auth login`
   - Used for PR search, issue lookup, etc.

4. **git**: Standard git installation
   - Used for history analysis, blame, log, etc.

## Installation

### From Marketplace

```bash
claude install ezcater-research
```

### Local Development

```bash
cd /path/to/claude-marketplace/plugins/ezcater-research
claude --plugin-dir .
```

## Usage

### Deep Research (Agent)

For complex investigations requiring synthesis across multiple sources:

```
You: "Can you do some deep research on why Liberty uses Omnichannel for feature flags?"
```

The agent will:
- Search Glean for design docs and discussions
- Query Jira for related tickets using atl CLI
- Find relevant PRs using gh CLI
- Analyze git history for implementation details
- Synthesize findings into comprehensive report

**When to use deep research:**
- Architectural decision investigations
- Code evolution analysis requiring context
- Multi-system integration understanding
- Historical context for current implementations

### Lightweight Research (Skill)

For straightforward research tasks:

```
You: "Find recent PRs about order management in ez-rails"
```

Main Claude follows research methodology:
- Uses appropriate tools (Glean, atl, gh, git)
- Applies search strategies
- Stays in current conversation context

**When to use regular research:**
- Specific PR/ticket lookups
- Recent code changes
- Single-source investigations
- Follow-up questions on existing research

## Components

### Agent: ezcater-research-analyst

Autonomous research specialist using Opus model. Triggers on "deep research" keyword.

**Tools**: Bash, Glob, Grep, Read, WebFetch, TodoWrite, Glean MCP, atl CLI, gh CLI

**Best for**: Complex multi-source investigations requiring synthesis

### Skill: ezcater-research

Research methodology and techniques for ezCater systems.

**Includes**:
- Glean search strategies
- Atlassian CLI usage patterns
- GitHub/Git analysis techniques
- Information synthesis and reporting

**Best for**: Lighter research where main Claude follows methodology

## Research Patterns

### Architectural Decisions

"Why was X chosen over Y?"

1. Search Glean for design docs, RFCs, ADRs
2. Find Jira epic/tickets using atl CLI
3. Locate implementation PRs with gh CLI
4. Review git history for context

### Code Evolution

"How did feature X evolve over time?"

1. Use git log --follow for file history
2. Find related PRs with gh CLI
3. Search Glean for related discussions
4. Check Jira for feature requests/bugs using atl

### Process Understanding

"How do I get an order to state X?"

1. Search Glean for process documentation
2. Find related code with Grep/Glob
3. Check Confluence docs using atl CLI
4. Locate workflow diagrams/specs

## Examples

See `skills/ezcater-research/examples/` for detailed research scenarios:

- `architectural-decision.md`: Investigating technology choices
- `code-evolution.md`: Tracing feature implementation
- `process-understanding.md`: Understanding workflows

## Limitations

- **No direct Google Docs access**: Must use Glean MCP
- **No direct Atlassian web access**: Must use atl CLI
- **Requires ezCater context**: Not useful outside ezCater
- **Depends on external tools**: Glean, atl, gh must be configured

## Troubleshooting

**"Permission denied" errors**:
- Check atl CLI authentication: `atl me`
- Check gh CLI authentication: `gh auth status`
- Verify Glean MCP server is running

**Agent not triggering**:
- Use explicit "deep research" keyword
- Check agent is enabled in Claude Code settings

**Skill not providing enough detail**:
- Use deep research agent for complex investigations
- Check references/ directory for detailed guides

## Contributing

This plugin is specific to ezCater infrastructure. For internal improvements:

1. Test changes locally with `claude --plugin-dir .`
2. Update version in plugin.json
3. Document changes in this README
4. Submit PR to marketplace repository

## License

MIT

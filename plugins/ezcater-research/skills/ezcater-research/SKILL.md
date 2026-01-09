---
name: ezcater-research
description: This skill should be used when the user asks to "investigate ezCater systems", "research architectural decisions", "understand code evolution", "find historical context", "analyze project history", "search internal documentation", or needs to understand why technical decisions were made at ezCater. Use this for lighter research where the user doesn't explicitly request "deep research" (which triggers the agent).
version: 0.1.0
---

# ezCater Research Methodology

Investigate ezCater codebases, architectural decisions, and project histories using multiple data sources including Glean, Atlassian (Jira/Confluence), GitHub, and Git.

## Overview

ezCater research requires synthesizing information from multiple internal sources. This skill provides proven methodologies for:
- Finding why architectural decisions were made
- Understanding code evolution over time
- Discovering historical context for current implementations
- Researching multi-system interactions and workflows

**When to use this skill vs deep research agent:**
- **Use this skill**: Straightforward investigations, specific PR/ticket lookups, single-source searches
- **Use deep research agent**: Complex multi-source synthesis, historical timelines, architectural decision investigations requiring extensive context

## CRITICAL: Access Requirements

**Cannot access resources directly. Must use appropriate tools:**

### Google Docs / Google Drive
- ❌ **CANNOT** access Google Docs URLs directly
- ✅ **MUST** use Glean MCP tools (`mcp__glean__company_search`, `mcp__glean__chat`)
- Glean indexes Google Docs and retrieves content

### Atlassian (Jira / Confluence)
- ❌ **CANNOT** access ezcater.atlassian.net URLs directly
- ✅ **MUST** use `atl` CLI for all Jira and Confluence access
- Commands: `atl jira issue view FX-123`, `atl confluence search "keywords"`

### GitHub
- ✅ Use `gh` CLI for PR/issue searches
- ✅ Use git commands for repository analysis
- ✅ Use Read tool for files in checked out repositories

## Research Strategy by Question Type

### Architectural Decisions

"Why was X chosen over Y?"

1. **Search Glean** for design docs, RFCs, ADRs using keywords like:
   - "[feature] decision", "[component] architecture", "why [technology]"
2. **Find epics/tickets** using `atl jira search "project = <PROJECT> AND text ~ 'keywords'"`
3. **Locate PRs** with `gh pr list --search "keywords" --repo owner/repo`
4. **Review git history** for implementation context using `git log --grep="pattern"`

See `references/glean-search-guide.md` for detailed Glean querying strategies.
See `examples/architectural-decision.md` for a complete walkthrough.

### Code Evolution

"How did feature X evolve over time?"

1. **Trace file history** with `git log --follow -p -- <file>`
2. **Find related PRs** using `gh pr list --search "feature-name"`
3. **Search Glean** for discussions using feature name keywords
4. **Check Jira** for feature requests: `atl jira search "project = FX AND summary ~ 'feature'"`

See `references/github-git-analysis.md` for git history techniques.
See `examples/code-evolution.md` for a complete example.

### Process Understanding

"How do I get an order to state X?" or "What triggers workflow Y?"

1. **Search Glean** for process documentation using keywords like "order lifecycle", "state machine"
2. **Find relevant code** using Grep for state transitions, workflow definitions
3. **Check Confluence** using `atl confluence search "order process"` for diagrams and specs
4. **Locate test fixtures** that demonstrate the process

See `references/atlassian-cli-guide.md` for atl CLI patterns.
See `examples/process-understanding.md` for a workflow investigation example.

### Multi-System Interactions

"How does system A interact with system B?"

1. **Search Glean** for integration documentation
2. **Find API contracts** in code using Grep for client/service definitions
3. **Check Confluence** for architecture diagrams: `atl confluence search "architecture"`
4. **Review recent PRs** affecting both systems using `gh pr list --search "systemA systemB"`

## Tool Usage Patterns

### Glean Search

Use Glean MCP for internal information:

```
Query: "order management migration decision"
Datasources: ["confluence", "slack", "gdrive"]
```

**Effective queries:**
- Be specific: "feature flag eppo provider" not just "feature flags"
- Include context: "liberty authentication refactor" not just "authentication"
- Try variations: "DM rails seeders", "delivery management backfills"

See `references/glean-search-guide.md` for complete search strategies.

### Atlassian CLI

**Search Jira tickets:**
```bash
atl jira search "project = FX AND text ~ 'delivery tracking'"
atl jira issue view FX-4623
```

**Search Confluence:**
```bash
atl confluence search "order state machine"
atl confluence page view 12345678
```

See `references/atlassian-cli-guide.md` for full command reference and examples.

### GitHub CLI

**Search PRs:**
```bash
gh pr list --repo ezcater/ez-rails --search "authentication" --limit 20
gh pr view 1234 --repo ezcater/ez-rails --json title,body,comments
```

**Search issues:**
```bash
gh issue list --repo ezcater/delivery-management-rails --search "seeder"
```

### Git History

**File evolution:**
```bash
git log --follow -p -- path/to/file.rb
git blame path/to/file.rb
```

**Commit searches:**
```bash
git log --grep="feature flag"
git log --all --grep="migration"
```

See `references/github-git-analysis.md` for advanced git techniques.

## Information Synthesis

After gathering information from multiple sources:

1. **Create timeline** - Order findings chronologically to show evolution
2. **Identify stakeholders** - Note who made decisions and why
3. **Understand business context** - Connect technical decisions to business needs
4. **Recognize trade-offs** - Explain what was gained and what was sacrificed
5. **Cross-reference sources** - Confirm findings across Glean, Jira, GitHub, git

See `references/synthesis-reporting.md` for detailed synthesis techniques.

## Response Format

Structure research findings clearly:

**Executive Summary**: Direct answer (2-3 sentences)

**Key Findings**:
- Finding 1 with evidence
- Finding 2 with evidence
- Finding 3 with evidence

**Historical Context** (if relevant):
- Timeline of decisions and changes
- Stakeholders involved
- Business drivers

**References**:
- Jira: FX-1234
- PR: ezcater/ez-rails#5678
- Commit: abc123def
- Confluence: [Title](page-id)

## Quality Standards

- **Be specific**: Include ticket numbers, PR links, commit SHAs, dates
- **Cite sources**: Always reference where information was found
- **Acknowledge gaps**: State explicitly what couldn't be found
- **Distinguish facts from speculation**: Make clear what's confirmed vs. inferred
- **Anticipate follow-ups**: Address obvious related questions

## Common Research Scenarios

### Finding Original Implementation

1. Use `gh pr list --search "initial [feature]"` to find first PR
2. Check PR description for original requirements
3. Find linked Jira ticket for business context
4. Review commit messages for implementation decisions

### Understanding Recent Changes

1. Use `git log --since="1 month ago" --grep="[feature]"`
2. Check recent PRs with `gh pr list --search "is:merged [feature]"`
3. Search Glean for Slack discussions in recent timeframe
4. Check Jira for recent bugs or improvements

### Investigating Production Issues

1. Search Jira for related tickets: `atl jira search "project = FX AND text ~ 'error'"`
2. Find relevant code with Grep
3. Check git blame for recent changes
4. Search Glean for related incidents or postmortems

## Special Considerations

### EngPortal Documentation

EngPortal content comes from GitHub files:
- URL: `https://engportal.ezcater.com/docs/.../component/repo-name/path`
- Source: `https://github.com/ezcater/repo-name/blob/main/path.md`

Access via:
- Read file directly if repo checked out under `/Users/doughughes/code/ezcater`
- Ask Glean for engportal content
- Use `gh` CLI to read from GitHub

### Checked Out Repositories

If repository exists under `/Users/doughughes/code/ezcater`:
- Read files directly with Read tool
- Run git commands in repository directory
- Use grep/glob for code searches

### Multiple Projects

ezCater has many interconnected projects:
- ez-rails: Main monolith
- delivery-management-rails: Delivery tracking
- omnichannel-rails: Multi-channel features
- liberty: Mobile app backend

When researching cross-project features, check all relevant repositories.

## Additional Resources

### Reference Files

For detailed techniques and patterns, consult:
- **`references/glean-search-guide.md`** - Glean search strategies and datasource selection
- **`references/atlassian-cli-guide.md`** - Complete atl CLI command reference
- **`references/github-git-analysis.md`** - Advanced gh CLI and git history techniques
- **`references/synthesis-reporting.md`** - Connecting dots and reporting findings

### Example Files

Working examples of research investigations in `examples/`:
- **`examples/architectural-decision.md`** - Researching why a technology was chosen
- **`examples/code-evolution.md`** - Tracing feature implementation over time
- **`examples/process-understanding.md`** - Understanding workflows and state transitions

## When to Escalate to Deep Research Agent

Use the deep research agent (by saying "do deep research" or "deep research needed") when:
- Investigation requires synthesizing 5+ sources
- Need comprehensive historical timeline
- Architectural decision investigation spanning multiple systems
- Research involves connecting dots across quarters/years
- Question requires extensive Opus-level reasoning

For straightforward searches, PR lookups, or single-system investigations, this skill provides efficient guidance without the overhead of the autonomous agent.

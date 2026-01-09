# Glean Search Guide for ezCater Research

Comprehensive guide for effective Glean searches when researching ezCater systems.

## Overview

Glean is ezCater's internal search engine that indexes content across multiple platforms including Slack, Google Docs, Confluence, Jira, Gmail, and more. Use Glean MCP tools to access this content - direct Google Docs or web access will not work.

## Available Tools

### company_search

Search across ezCater's internal datasources:

```
mcp__glean__company_search
- query: "your search terms"
- datasources: ["confluence", "slack", "gdrive", "jira"] (optional)
```

**When to use**: Finding documents, discussions, or information across multiple sources

### chat

Interactive chat with Glean Assistant using RAG:

```
mcp__glean__chat
- message: "your question"
- context: ["previous message 1", "previous message 2"] (optional)
```

**When to use**: Natural language questions requiring context and synthesis

### people_profile_search

Search for people by name, role, department:

```
mcp__glean__people_profile_search
- query: "person name or criteria"
- filters: {"department": "Engineering", "title": "Senior Engineer"}
```

**When to use**: Finding stakeholders, team members, or decision makers

## Effective Query Strategies

### Be Specific

❌ **Bad**: "authentication"
✅ **Good**: "liberty authentication refactor 2024"

❌ **Bad**: "feature flags"
✅ **Good**: "feature flag eppo provider migration"

**Why**: Specific queries return focused, relevant results instead of overwhelming general information.

### Include Context

Add project names, system names, or timeframes:

- "ez-rails order state machine"
- "delivery management rails seeders Q4 2024"
- "omnichannel feature flag implementation"

### Try Multiple Phrasings

If first search doesn't yield results, try variations:

1. "DM rails seeders"
2. "delivery management backfills"
3. "delivery tracking seeder implementation"
4. "backfill tasks delivery management"

### Use Quotes for Exact Phrases

Search for specific phrases:

- "order management migration"
- "why did we choose Omnichannel"
- "authentication system refactor"

## Datasource Selection

### Confluence

Technical documentation, RFCs, design docs, ADRs:

```
datasources: ["confluence"]
query: "order state machine RFC"
```

**Best for**:
- Design documents
- Architecture decision records
- Project requirements
- Technical specifications
- Team documentation

**Search strategies**:
- Include "RFC", "ADR", "design doc" in queries
- Search for component or system names
- Look for "architecture", "design", "requirements"

### Slack

Discussions, decisions, quick context:

```
datasources: ["slack"]
query: "feature flag provider discussion"
```

**Best for**:
- Recent discussions
- Decision rationale
- Quick context
- Team communication
- Incident discussions

**Search strategies**:
- Include dates/timeframes for recent discussions
- Search channel names if known (#eng-fulfillment, #dev-experience)
- Look for key participants' names
- Try "discussion", "decision", "why"

### Google Drive (gdrive)

Docs, sheets, slides, meeting notes:

```
datasources: ["gdrive"]
query: "Q4 planning order management"
```

**Best for**:
- Meeting notes
- Planning documents
- Presentations
- Spreadsheets
- Collaborative docs

**Search strategies**:
- Include meeting types: "standup", "planning", "retro"
- Search for quarters: "Q4 2024", "2024 planning"
- Look for "notes", "agenda", "slides"

### Jira

Tickets, epics, stories:

```
datasources: ["jira"]
query: "delivery tracking fulfillment issue"
```

**Best for**:
- Feature requirements
- Bug reports
- Implementation tasks
- Project tracking

**Note**: For detailed Jira research, prefer `atl` CLI which provides more control. Use Glean for quick searches when you don't know ticket numbers.

### Multiple Datasources

Search across multiple sources for comprehensive results:

```
datasources: ["confluence", "slack", "gdrive"]
query: "order management migration"
```

**When to use**: Initial broad search to identify where information exists

## Query Patterns by Research Type

### Architectural Decisions

**Goal**: Understand why a technical choice was made

**Queries**:
- "[feature] decision rationale"
- "why [technology] instead of [alternative]"
- "[component] architecture RFC"
- "[technology] evaluation comparison"

**Datasources**: confluence, slack, gdrive

**Example**:
```
query: "why Omnichannel feature flag provider instead of Eppo"
datasources: ["confluence", "slack"]
```

### Code Evolution

**Goal**: Understand how code changed over time

**Queries**:
- "[feature] implementation history"
- "[component] refactor motivation"
- "[system] migration planning"

**Datasources**: confluence, slack, jira

**Example**:
```
query: "authentication system refactor ez-rails 2023"
datasources: ["confluence", "jira", "slack"]
```

### Process Understanding

**Goal**: Learn how business processes work

**Queries**:
- "[process] workflow documentation"
- "how [state transition] works"
- "[system] state machine design"
- "[process] lifecycle diagram"

**Datasources**: confluence, gdrive

**Example**:
```
query: "order lifecycle state transitions documentation"
datasources: ["confluence", "gdrive"]
```

### Multi-System Integration

**Goal**: Understand how systems communicate

**Queries**:
- "[systemA] [systemB] integration"
- "[service] API contract"
- "[component] communication architecture"

**Datasources**: confluence, gdrive

**Example**:
```
query: "DM rails ez-rails integration architecture"
datasources: ["confluence"]
```

### Historical Context

**Goal**: Find when and why something was built

**Queries**:
- "[feature] original requirements"
- "[project] kickoff planning"
- "[initiative] business justification"

**Datasources**: confluence, gdrive, jira

**Example**:
```
query: "delivery tracking project kickoff Q3 2023"
datasources: ["confluence", "gdrive"]
```

## Using Glean Chat

For natural language questions requiring synthesis:

```
mcp__glean__chat
message: "What systems are involved in order fulfillment tracking?"
```

**Good questions**:
- "How does the order state machine work?"
- "What are the main components of delivery tracking?"
- "Who owns the authentication system in ez-rails?"
- "What's the current status of the mobile app migration?"

**Follow-up with context**:
```
mcp__glean__chat
message: "What changes were made in Q4 2024?"
context: ["What systems are involved in order fulfillment tracking?"]
```

## Interpreting Results

### Result Quality Indicators

**High quality results**:
- Recent dates (within last year for current info)
- From authoritative sources (Confluence, design docs)
- Contain technical details
- Include decision rationale

**Low quality results**:
- Very old (>2 years unless historical research)
- Brief Slack messages without context
- Duplicate information
- Vague or high-level only

### When Results Are Insufficient

If Glean doesn't return useful results:

1. **Try alternative queries** with different keywords
2. **Narrow datasources** to specific platforms
3. **Search for people** who might know: `people_profile_search`
4. **Check other tools**: atl CLI, gh CLI, git history
5. **Ask in Glean chat** for synthesis

### Common Issues

**Too many results**: Add more specific terms, narrow datasources
**Too few results**: Broaden query, try synonyms, check spelling
**Wrong context**: Add project/system names to query
**Outdated results**: Include timeframes like "2024", "recent"

## Advanced Techniques

### Stakeholder Discovery

Find who made decisions or owns systems:

1. Search for technical content: `"authentication refactor"`
2. Note author names in results
3. Use `people_profile_search` to find their current role
4. Search for their name + topic for more context

### Timeline Construction

Build chronological understanding:

1. Search with older timeframe: "feature 2022"
2. Search with recent timeframe: "feature 2024"
3. Compare results to see evolution
4. Fill gaps with Jira and git history

### Cross-Platform Verification

Confirm findings across multiple sources:

1. Find design doc in Confluence
2. Find discussion in Slack
3. Find implementation ticket in Jira
4. Verify dates and details match

## Integration with Other Tools

Glean works best in combination with other research tools:

**Glean → Jira**: Find ticket numbers in Glean, get details with atl CLI
**Glean → GitHub**: Find PR numbers in Glean, review with gh CLI
**Glean → Git**: Find context in Glean, trace code with git
**Glean → Confluence**: Quick search with Glean, deep read with atl CLI

## Examples

### Example 1: Feature Flag Decision

**Question**: "Why does Liberty use Omnichannel for feature flags?"

**Search**:
```
query: "Liberty Omnichannel feature flag decision"
datasources: ["confluence", "slack"]
```

**Results interpretation**:
- Check Confluence for design docs or ADRs
- Check Slack for decision discussions
- Note dates and participants
- Follow up with atl/gh for implementation details

### Example 2: Order State Understanding

**Question**: "How do orders transition to completed state?"

**Search**:
```
query: "order state machine completed transition"
datasources: ["confluence", "gdrive"]
```

**Results interpretation**:
- Look for state machine diagrams
- Find workflow documentation
- Identify code locations from docs
- Follow up with code search using Grep

### Example 3: Migration Status

**Question**: "What's the status of order management migration?"

**Chat**:
```
message: "What's the current status of the order management system migration?"
```

**Results interpretation**:
- Glean chat synthesizes multiple sources
- Note mentioned tickets and docs
- Follow up on specific tickets with atl CLI
- Check recent PRs with gh CLI

## Best Practices

✅ **DO**:
- Start with broad searches, then narrow
- Try multiple query variations
- Use datasource filtering when appropriate
- Follow up Glean results with specific tools (atl, gh, git)
- Note key people mentioned in results
- Cross-reference findings across sources

❌ **DON'T**:
- Rely solely on Glean for detailed information
- Accept first results without verification
- Skip reading actual documents (don't just rely on snippets)
- Ignore timestamps on results
- Forget to check multiple datasources

## Quick Reference

| Research Type | Glean Datasources | Key Query Terms |
|--------------|-------------------|-----------------|
| Architecture | confluence, gdrive | RFC, ADR, design, architecture |
| Discussions | slack | discussion, decision, why |
| Requirements | confluence, jira | requirements, PRD, spec |
| Planning | gdrive, confluence | planning, kickoff, roadmap |
| Historical | confluence, slack, gdrive | original, history, evolution |
| People | (use people_profile_search) | name, role, team |

## Summary

Glean is the starting point for ezCater research, providing access to internal documentation and discussions. Use it to:
- Find design documents and decisions
- Locate relevant discussions
- Identify stakeholders
- Get initial context

Then deepen research with:
- `atl` CLI for detailed Jira/Confluence reading
- `gh` CLI for PR and issue investigation
- `git` commands for code history
- Direct file reading for checked out code

Effective Glean use requires specific queries, appropriate datasource selection, and integration with other research tools for comprehensive investigations.

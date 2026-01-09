# Atlassian CLI Guide for ezCater Research

Complete reference for using the `atl` CLI tool to research Jira tickets and Confluence pages at ezCater.

## Overview

The `atl` CLI is the **ONLY** way to access ezcater.atlassian.net content. Direct web access will not work. Use `atl` for all Jira and Confluence operations.

## Prerequisites

- `atl` CLI must be installed and configured
- Authenticated with ezCater Atlassian credentials
- Verify with: `atl me`

## Jira Commands

### Search Tickets

Find tickets using JQL (Jira Query Language):

```bash
atl jira search "JQL query"
```

**Common searches**:

```bash
# Search by project and keywords
atl jira search "project = FX AND text ~ 'delivery tracking'"

# Search by status
atl jira search "project = FX AND status = 'In Progress'"

# Search by assignee
atl jira search "project = FX AND assignee = currentUser()"

# Search created recently
atl jira search "project = FX AND created >= -30d"

# Search by type
atl jira search "project = FX AND type = Epic"

# Complex query
atl jira search "project = FX AND text ~ 'order' AND status IN ('In Progress', 'To Do')"
```

**Output options**:

```bash
# Table format (default)
atl jira search "query"

# JSON format
atl jira search "query" --output json

# Limit results
atl jira search "query" --limit 50
```

### View Ticket Details

Get complete ticket information:

```bash
atl jira issue view TICKET-ID
```

**Examples**:

```bash
# View ticket
atl jira issue view FX-4623

# View with JSON output
atl jira issue view FX-4623 --output json

# View multiple tickets
atl jira issue view FX-4623 FX-4625
```

**Output includes**:
- Summary and description
- Status and priority
- Assignee and reporter
- Comments
- Linked issues
- Custom fields
- History/changelog

### View Ticket Comments

Get just the comments:

```bash
atl jira issue view TICKET-ID --comments
```

**Use case**: When ticket has many comments and you want focused view

### List Projects

Find project keys:

```bash
atl jira project list
```

**Common ezCater projects**:
- FX - Fulfillment Experience
- DM - Delivery Management
- MOB - Mobile
- PLAT - Platform

## Confluence Commands

### Search Pages

Find pages by keywords:

```bash
atl confluence search "search terms"
```

**Examples**:

```bash
# Basic search
atl confluence search "order state machine"

# Search in specific space
atl confluence search "authentication" --space ENG

# Search with filters
atl confluence search "RFC" --type page

# Limit results
atl confluence search "migration" --limit 20
```

### View Page

Read complete page content:

```bash
atl confluence page view PAGE-ID
```

**Examples**:

```bash
# View by page ID
atl confluence page view 123456789

# View by page title
atl confluence page view --title "Order Management RFC"

# View with specific space
atl confluence page view --title "Authentication" --space ENG

# Output as markdown
atl confluence page view 123456789 --output markdown
```

**Note**: Page IDs are long numeric values. Get them from search results or page URLs.

### List Spaces

Find Confluence spaces:

```bash
atl confluence space list
```

**Common ezCater spaces**:
- ENG - Engineering
- PROD - Product
- ~USERNAME - Personal spaces

## JQL Query Patterns

JQL (Jira Query Language) powers Jira searches. Master these patterns for effective research.

### Basic Syntax

```
field operator value
```

**Combine with**:
- `AND` - both conditions must be true
- `OR` - either condition must be true
- `NOT` - condition must be false

### Common Fields

| Field | Description | Example |
|-------|-------------|---------|
| project | Project key | `project = FX` |
| text | Full-text search | `text ~ 'keywords'` |
| summary | Issue title | `summary ~ 'bug'` |
| description | Issue body | `description ~ 'error'` |
| status | Current status | `status = 'In Progress'` |
| assignee | Assigned to | `assignee = 'doug.hughes'` |
| reporter | Created by | `reporter = currentUser()` |
| created | Creation date | `created >= -30d` |
| updated | Last updated | `updated >= -7d` |
| type | Issue type | `type = Epic` |
| priority | Priority level | `priority = High` |
| labels | Issue labels | `labels = backend` |

### Operators

| Operator | Meaning | Example |
|----------|---------|---------|
| = | Equals | `status = 'Done'` |
| != | Not equals | `status != 'Done'` |
| ~ | Contains | `text ~ 'order'` |
| !~ | Not contains | `text !~ 'test'` |
| IN | In list | `status IN ('To Do', 'In Progress')` |
| NOT IN | Not in list | `type NOT IN (Sub-task)` |
| >= | Greater/equal | `created >= -30d` |
| <= | Less/equal | `priority <= Medium` |
| IS | Is value | `assignee IS EMPTY` |
| IS NOT | Is not value | `resolution IS NOT EMPTY` |

### Date/Time

**Relative dates**:
- `-1d` - 1 day ago
- `-1w` - 1 week ago
- `-30d` - 30 days ago
- `-3m` - 3 months ago

**Functions**:
- `now()` - Current time
- `startOfDay()` - Start of today
- `startOfWeek()` - Start of this week
- `endOfMonth()` - End of this month

**Examples**:
```jql
created >= -30d
updated >= startOfDay()
created >= '2024-01-01' AND created <= '2024-12-31'
```

### Example Queries

**Find recent work**:
```jql
project = FX AND updated >= -7d ORDER BY updated DESC
```

**Find my open tickets**:
```jql
project = FX AND assignee = currentUser() AND status != Done
```

**Find epics about orders**:
```jql
project = FX AND type = Epic AND text ~ 'order'
```

**Find bugs created this month**:
```jql
project = FX AND type = Bug AND created >= startOfMonth()
```

**Find unassigned high priority**:
```jql
project = FX AND assignee IS EMPTY AND priority = High
```

**Find tickets with specific label**:
```jql
project = FX AND labels = 'delivery-tracking'
```

**Complex multi-condition**:
```jql
project = FX AND
type IN (Story, Task) AND
status IN ('In Progress', 'To Do') AND
text ~ 'authentication' AND
updated >= -14d
ORDER BY updated DESC
```

## Research Workflows

### Investigating a Feature

**Goal**: Understand feature history and implementation

1. **Find epic**:
   ```bash
   atl jira search "project = FX AND type = Epic AND text ~ 'feature name'"
   ```

2. **View epic details**:
   ```bash
   atl jira issue view FX-1234
   ```

3. **Find related stories**:
   ```bash
   atl jira search "'Epic Link' = FX-1234"
   ```

4. **Check design docs** (note page IDs from epic):
   ```bash
   atl confluence page view 123456789
   ```

### Understanding Technical Decision

**Goal**: Find why a choice was made

1. **Search Confluence for RFC/ADR**:
   ```bash
   atl confluence search "feature name RFC"
   atl confluence search "feature name ADR"
   atl confluence search "feature name design doc"
   ```

2. **Read design documents**:
   ```bash
   atl confluence page view PAGE-ID --output markdown
   ```

3. **Find discussion tickets**:
   ```bash
   atl jira search "project = FX AND text ~ 'feature discussion'"
   ```

4. **Review comments** for decision rationale:
   ```bash
   atl jira issue view FX-1234 --comments
   ```

### Tracing Bug Resolution

**Goal**: Understand how a bug was fixed

1. **Find bug ticket**:
   ```bash
   atl jira search "project = FX AND type = Bug AND text ~ 'error description'"
   ```

2. **View full ticket**:
   ```bash
   atl jira issue view FX-5678
   ```

3. **Check linked issues**:
   - Note linked PRs in ticket
   - Note related bugs or stories
   - Note parent epic if present

4. **Find follow-up work**:
   ```bash
   atl jira search "project = FX AND text ~ 'bug-5678'"
   ```

### Understanding Process

**Goal**: Learn how business process works

1. **Search Confluence for process docs**:
   ```bash
   atl confluence search "order lifecycle"
   atl confluence search "state machine"
   atl confluence search "workflow diagram"
   ```

2. **Find implementation tickets**:
   ```bash
   atl jira search "project = FX AND text ~ 'order state'"
   ```

3. **Read process documentation**:
   ```bash
   atl confluence page view PAGE-ID --output markdown
   ```

## Output Formats

### Table Format (Default)

Human-readable tables:

```bash
atl jira search "project = FX" --limit 10
```

**Good for**: Quick scanning, terminal viewing

### JSON Format

Machine-parseable data:

```bash
atl jira issue view FX-1234 --output json
```

**Good for**: Parsing with jq, detailed analysis, automation

**Parse with jq**:
```bash
atl jira issue view FX-1234 --output json | jq '.fields.description'
atl jira search "project = FX" --output json | jq '.[].key'
```

### Markdown Format (Confluence only)

Formatted markdown:

```bash
atl confluence page view 123456789 --output markdown
```

**Good for**: Reading documentation, including in reports

## Common Issues and Solutions

### Authentication Errors

**Problem**: "Authentication failed" or "Not authorized"

**Solutions**:
1. Check authentication: `atl me`
2. Re-authenticate: `atl auth login`
3. Verify ezCater instance configured

### Too Many Results

**Problem**: Search returns hundreds of tickets

**Solutions**:
1. Add more specific keywords
2. Filter by project: `project = FX`
3. Limit timeframe: `created >= -30d`
4. Use limit flag: `--limit 20`
5. Sort results: `ORDER BY updated DESC`

### Can't Find Page

**Problem**: Confluence page not found by title

**Solutions**:
1. Search first to get page ID
2. Try different title variations
3. Specify space: `--space ENG`
4. Use page ID directly instead of title

### JQL Syntax Errors

**Problem**: "Invalid JQL query"

**Solutions**:
1. Check quotes around values with spaces
2. Use `~` for text search, not `=`
3. Use parentheses for complex conditions
4. Check field names are correct

## Integration with Other Tools

### Atl → GitHub

Use ticket IDs to find related PRs:

```bash
# Find ticket
atl jira issue view FX-1234

# Note PR numbers mentioned
# Then use gh CLI:
gh pr view 5678 --repo ezcater/ez-rails
```

### Atl → Git

Use information from tickets to guide git searches:

```bash
# Find ticket with component name
atl jira issue view FX-1234

# Search git history:
git log --grep="FX-1234"
git log --all --grep="component-name"
```

### Atl → Glean

Use atl for detailed reading, Glean for broad searching:

```bash
# Broad Glean search to find ticket numbers
# Then detailed atl viewing:
atl jira issue view FX-1234
```

## Best Practices

✅ **DO**:
- Use JQL for complex searches
- Start broad, then narrow with filters
- Save common queries in notes
- Use JSON output for parsing
- Check page IDs in search results
- Order results by relevant fields
- Limit results to reduce noise

❌ **DON'T**:
- Try to access ezcater.atlassian.net URLs directly
- Use vague search terms without filters
- Ignore linked issues in tickets
- Skip reading full ticket descriptions
- Forget to check comments for context
- Assume first result is most relevant

## Quick Reference

### Most Common Commands

```bash
# Search Jira
atl jira search "project = FX AND text ~ 'keywords'"

# View ticket
atl jira issue view FX-1234

# Search Confluence
atl confluence search "keywords"

# View page
atl confluence page view 123456789

# Check authentication
atl me
```

### Most Useful JQL Patterns

```jql
# Recent updates
project = FX AND updated >= -7d

# My open work
project = FX AND assignee = currentUser() AND status != Done

# Find by keyword
project = FX AND text ~ 'keyword'

# Find epics
project = FX AND type = Epic

# Find by date range
project = FX AND created >= '2024-01-01' AND created <= '2024-12-31'
```

## Summary

The `atl` CLI is essential for ezCater research, providing the only access to Jira and Confluence. Use it to:
- Search tickets and pages with JQL and keywords
- Read full content with detailed views
- Follow links between tickets and documentation
- Understand project history and decisions

Master JQL for powerful searches, use JSON output for parsing, and integrate with other tools (gh, git, Glean) for comprehensive research.

Remember: **Always use atl CLI** - never attempt direct web access to ezcater.atlassian.net.

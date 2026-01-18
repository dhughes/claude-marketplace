---
name: atlassian-usage
description: This skill should be used when the user asks to "search jira", "find tickets", "look up an issue", "search confluence", "find pages", "read a document", "create a ticket", "update an issue", "add a comment", mentions JQL, CQL, Atlassian, Jira issues, Confluence pages, or provides an Atlassian URL (*.atlassian.net). Provides guidance for using the atl CLI to interact with Atlassian products.
---

# Atlassian CLI Usage

This skill provides guidance for using the `atl` CLI tool to interact with Jira and Confluence.

## CRITICAL: Handling Atlassian URLs

**NEVER attempt to fetch Atlassian URLs directly via web tools.** Atlassian Cloud requires authentication that web fetch tools cannot provide. Instead, extract identifiers from URLs and use the `atl` CLI.

### Jira URL Patterns

When a user provides a Jira URL, extract the relevant identifier and use `atl`:

| URL Pattern | Extract | Command |
|-------------|---------|---------|
| `https://*.atlassian.net/jira/software/projects/PROJ/boards/123` | `PROJ` | `atl jira search-jql "project = PROJ"` |
| `https://*.atlassian.net/jira/core/projects/PROJ/board` | `PROJ` | `atl jira search-jql "project = PROJ"` |

**Examples:**
- `https://example.atlassian.net/jira/software/projects/PROJ/boards/5` → `atl jira search-jql "project = PROJ"`

### Confluence URL Patterns

When a user provides a Confluence URL, extract the page ID and use `atl`:

| URL Pattern | Extract | Command |
|-------------|---------|---------|
| `https://*.atlassian.net/wiki/spaces/SPACE/pages/123456789/Page+Title` | `123456789` | `atl confluence get-page 123456789` |
| `https://*.atlassian.net/wiki/spaces/SPACE/overview` | `SPACE` | `atl confluence get-pages-in-space SPACE` |

**Examples:**
- `https://example.atlassian.net/wiki/spaces/ENG/pages/987654321/Architecture` → `atl confluence get-page 987654321`
- `https://example.atlassian.net/wiki/spaces/TEAM/pages/123456789` → `atl confluence get-page 123456789`

### URL Recognition Rules

1. Any URL containing `.atlassian.net` should trigger this skill
2. NEVER use WebFetch or similar tools on Atlassian URLs
3. Extract the identifier (issue key or page ID) from the URL
4. Use the appropriate `atl` command to retrieve the content

## Overview

The `atl` CLI provides direct access to Jira and Confluence through their REST APIs. Use it for searching, reading, creating, and updating content in Atlassian products.

## Discovering Commands

The `atl` CLI is actively developed and may have commands not covered in this skill. Always use `--help` to discover available commands and options:

```bash
atl --help
atl jira --help
atl confluence --help
atl jira <command> --help
```

## Jira Operations

### Searching Issues (JQL)

Search using Jira Query Language:

```bash
atl jira search-jql "project = PROJ AND status = 'In Progress'"
atl jira search-jql "assignee = currentUser() AND created >= -7d"
atl jira search-jql "text ~ 'keyword'" --max-results 20 --json
```

For detailed JQL syntax, see `references/jql-reference.md`.

### Reading Issues

```bash
atl jira get-issue PROJ-123
atl jira get-issue PROJ-123 --json
atl jira get-issue PROJ-123 --fields summary,status,assignee
```

### Creating Issues

```bash
atl jira create-issue --project PROJ --type Task --summary "Issue title"
atl jira create-issue --project PROJ --type Bug --summary "Bug title" \
  --description "**Details:** markdown supported"
```

The `--description` flag supports markdown formatting including headings, bold, italic, code blocks, and lists.

### Updating Issues

```bash
atl jira edit-issue PROJ-123 --summary "New title"
atl jira edit-issue PROJ-123 --description "Updated description"
atl jira edit-issue PROJ-123 --assignee <account-id>
```

To find account IDs:

```bash
atl jira lookup-account-id "username or email"
```

### Comments

```bash
atl jira add-comment PROJ-123 "Comment text with **markdown**"
```

### Transitions

```bash
atl jira get-transitions PROJ-123
atl jira transition-issue PROJ-123 --transition "In Progress"
```

### Issue Links

```bash
atl jira get-link-types
atl jira create-issue-link --inward PROJ-123 --outward PROJ-456 --type "Blocks"
atl jira get-issue-links PROJ-123
```

### Project Information

```bash
atl jira get-projects
atl jira get-project-issue-types PROJ
atl jira get-create-meta --project PROJ --issue-type Task
```

## Confluence Operations

### Searching Pages (CQL)

Search using Confluence Query Language:

```bash
atl confluence search-cql "space = TEAM AND type = page"
atl confluence search-cql "title ~ 'Onboarding'"
atl confluence search-cql "text ~ 'documentation'" --limit 20 --json
```

For detailed CQL syntax, see `references/cql-reference.md`.

### Reading Pages

Page IDs are numeric values found in page URLs or search results:

```bash
atl confluence get-page 123456789
atl confluence get-page 123456789 --json
```

### Creating Pages

```bash
atl confluence create-page --space TEAM --title "Page Title" --body "Content"
```

### Updating Pages

```bash
atl confluence update-page 123456789 --title "New Title"
atl confluence update-page 123456789 --body "Updated content"
```

### Comments

```bash
atl confluence get-page-comments 123456789
atl confluence add-comment 123456789 "Comment text"
atl confluence create-inline-comment 123456789 --body "Inline comment"
```

### Navigation

```bash
atl confluence get-spaces
atl confluence get-pages-in-space TEAM
atl confluence get-page-ancestors 123456789
atl confluence get-page-descendants 123456789
```

## Output Formats

Most commands support `--json` for machine-readable output:

```bash
atl jira get-issue PROJ-123 --json | jq '.fields.status.name'
atl jira search-jql "project = PROJ" --json | jq '.[].key'
```

## Common Workflows

### Research a Topic

1. Search Jira for relevant tickets:
   ```bash
   atl jira search-jql "project = PROJ AND text ~ 'topic'"
   ```

2. Read promising tickets:
   ```bash
   atl jira get-issue PROJ-123
   ```

3. Search Confluence for documentation:
   ```bash
   atl confluence search-cql "text ~ 'topic'"
   ```

4. Read relevant pages:
   ```bash
   atl confluence get-page 123456789
   ```

### Create and Track Work

1. Create the issue:
   ```bash
   atl jira create-issue --project PROJ --type Task --summary "Task name" \
     --description "Details here"
   ```

2. Transition to in progress:
   ```bash
   atl jira transition-issue PROJ-123 --transition "In Progress"
   ```

3. Add updates as comments:
   ```bash
   atl jira add-comment PROJ-123 "Progress update"
   ```

## Authentication

Verify authentication status:

```bash
atl meta user-info
```

If not authenticated or authentication has expired:

```bash
atl auth
```

## Additional Resources

### Reference Files

For detailed query syntax, consult:
- **`references/jql-reference.md`** - Complete JQL syntax, fields, operators, and functions
- **`references/cql-reference.md`** - Complete CQL syntax for Confluence searches

### Discovering New Features

The CLI is actively developed. To discover commands not covered here:

```bash
atl --help
atl jira --help
atl confluence --help
atl <product> <command> --help
```

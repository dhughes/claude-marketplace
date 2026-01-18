# JQL Reference

Complete reference for Jira Query Language (JQL) used with `atl jira search-jql`.

## Basic Syntax

```
field operator value [AND|OR] field operator value
```

## Common Fields

| Field | Description | Example |
|-------|-------------|---------|
| `project` | Project key | `project = PROJ` |
| `text` | Full-text search | `text ~ 'keywords'` |
| `summary` | Issue title | `summary ~ 'bug'` |
| `description` | Issue body | `description ~ 'error'` |
| `status` | Current status | `status = 'In Progress'` |
| `assignee` | Assigned to | `assignee = 'user.name'` |
| `reporter` | Created by | `reporter = currentUser()` |
| `created` | Creation date | `created >= -30d` |
| `updated` | Last updated | `updated >= -7d` |
| `resolved` | Resolution date | `resolved >= -7d` |
| `type` | Issue type | `type = Epic` |
| `priority` | Priority level | `priority = High` |
| `labels` | Issue labels | `labels = backend` |
| `component` | Component | `component = "API"` |
| `fixVersion` | Fix version | `fixVersion = "1.0"` |
| `affectedVersion` | Affected version | `affectedVersion = "0.9"` |
| `resolution` | Resolution status | `resolution = Done` |
| `parent` | Parent issue | `parent = PROJ-100` |
| `"Epic Link"` | Epic parent | `"Epic Link" = PROJ-50` |

## Operators

| Operator | Meaning | Example |
|----------|---------|---------|
| `=` | Equals | `status = 'Done'` |
| `!=` | Not equals | `status != 'Done'` |
| `~` | Contains (text) | `text ~ 'order'` |
| `!~` | Not contains | `text !~ 'test'` |
| `>` | Greater than | `created > '2024-01-01'` |
| `>=` | Greater or equal | `created >= -30d` |
| `<` | Less than | `priority < High` |
| `<=` | Less or equal | `priority <= Medium` |
| `IN` | In list | `status IN ('To Do', 'In Progress')` |
| `NOT IN` | Not in list | `type NOT IN (Sub-task)` |
| `IS` | Is value | `assignee IS EMPTY` |
| `IS NOT` | Is not value | `resolution IS NOT EMPTY` |
| `WAS` | Historical value | `status WAS 'In Progress'` |
| `WAS NOT` | Was not value | `status WAS NOT 'Done'` |
| `CHANGED` | Value changed | `status CHANGED` |

## Logical Operators

| Operator | Description |
|----------|-------------|
| `AND` | Both conditions must be true |
| `OR` | Either condition must be true |
| `NOT` | Condition must be false |
| `()` | Group conditions |

Example with grouping:

```jql
project = PROJ AND (status = 'To Do' OR status = 'In Progress')
```

## Date/Time Values

### Relative Dates

| Value | Meaning |
|-------|---------|
| `-1d` | 1 day ago |
| `-1w` | 1 week ago |
| `-30d` | 30 days ago |
| `-3m` | 3 months ago |
| `-1y` | 1 year ago |

### Date Functions

| Function | Description |
|----------|-------------|
| `now()` | Current time |
| `currentLogin()` | Last login time |
| `startOfDay()` | Start of today |
| `startOfWeek()` | Start of this week |
| `startOfMonth()` | Start of this month |
| `startOfYear()` | Start of this year |
| `endOfDay()` | End of today |
| `endOfWeek()` | End of this week |
| `endOfMonth()` | End of this month |
| `endOfYear()` | End of this year |

### Absolute Dates

```jql
created >= '2024-01-01'
created >= '2024-01-01' AND created <= '2024-12-31'
updated >= '2024-06-15 09:00'
```

## User Functions

| Function | Description |
|----------|-------------|
| `currentUser()` | Logged-in user |
| `membersOf("group")` | Members of group |

## Ordering Results

```jql
project = PROJ ORDER BY created DESC
project = PROJ ORDER BY priority ASC, created DESC
```

Sort directions:
- `ASC` - Ascending (oldest/lowest first)
- `DESC` - Descending (newest/highest first)

## Example Queries

### Find Recent Work

```jql
project = PROJ AND updated >= -7d ORDER BY updated DESC
```

### My Open Tickets

```jql
project = PROJ AND assignee = currentUser() AND status != Done
```

### Unassigned High Priority

```jql
project = PROJ AND assignee IS EMPTY AND priority = High
```

### Epics About a Topic

```jql
project = PROJ AND type = Epic AND text ~ 'authentication'
```

### Bugs Created This Month

```jql
project = PROJ AND type = Bug AND created >= startOfMonth()
```

### Tickets with Specific Label

```jql
project = PROJ AND labels = 'backend'
```

### Stories in an Epic

```jql
"Epic Link" = PROJ-100
```

### Recently Resolved

```jql
project = PROJ AND resolved >= -7d ORDER BY resolved DESC
```

### Complex Multi-Condition

```jql
project = PROJ AND
type IN (Story, Task) AND
status IN ('In Progress', 'To Do') AND
text ~ 'api' AND
updated >= -14d
ORDER BY updated DESC
```

### Status Changed Recently

```jql
project = PROJ AND status CHANGED AFTER -1d
```

### Overdue Issues

```jql
project = PROJ AND duedate < now() AND resolution IS EMPTY
```

## Text Search Tips

The `~` operator performs full-text search:

```jql
text ~ 'exact phrase'
text ~ 'word1 word2'  -- matches if both words appear
text ~ 'word1 OR word2'  -- matches either word
text ~ 'word*'  -- wildcard
```

Search specific fields:

```jql
summary ~ 'keyword'  -- title only
description ~ 'keyword'  -- description only
comment ~ 'keyword'  -- comments only
text ~ 'keyword'  -- all text fields
```

## Escaping Special Characters

Quote values with spaces or special characters:

```jql
status = 'In Progress'
summary ~ 'can\'t reproduce'
project = "My Project"
```

## Performance Tips

1. Always include `project = X` to limit scope
2. Use `created >= -Nd` to limit date range
3. Avoid `text ~` on large result sets
4. Use `--max-results` flag to limit results
5. Order by indexed fields (created, updated, priority)

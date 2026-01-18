# CQL Reference

Complete reference for Confluence Query Language (CQL) used with `atl confluence search-cql`.

## Basic Syntax

```
field operator value [AND|OR] field operator value
```

## Common Fields

| Field | Description | Example |
|-------|-------------|---------|
| `type` | Content type | `type = page` |
| `space` | Space key | `space = TEAM` |
| `title` | Page title | `title ~ 'Onboarding'` |
| `text` | Full content | `text ~ 'documentation'` |
| `creator` | Created by | `creator = 'user.name'` |
| `contributor` | Any contributor | `contributor = 'user.name'` |
| `created` | Creation date | `created >= '2024-01-01'` |
| `lastModified` | Last modified | `lastModified >= -7d` |
| `label` | Content labels | `label = 'architecture'` |
| `parent` | Parent page ID | `parent = 123456789` |
| `ancestor` | Any ancestor | `ancestor = 123456789` |
| `id` | Content ID | `id = 123456789` |

## Content Types

| Type | Description |
|------|-------------|
| `page` | Wiki pages |
| `blogpost` | Blog posts |
| `comment` | Comments |
| `attachment` | Attachments |

Example:

```cql
type = page AND space = TEAM
type IN (page, blogpost) AND text ~ 'release'
```

## Operators

| Operator | Meaning | Example |
|----------|---------|---------|
| `=` | Equals | `space = 'TEAM'` |
| `!=` | Not equals | `space != 'ARCHIVE'` |
| `~` | Contains | `title ~ 'guide'` |
| `!~` | Not contains | `title !~ 'draft'` |
| `>` | After (dates) | `created > '2024-01-01'` |
| `>=` | On or after | `created >= -30d` |
| `<` | Before (dates) | `lastModified < '2024-01-01'` |
| `<=` | On or before | `lastModified <= -7d` |
| `IN` | In list | `space IN ('TEAM', 'ENG')` |
| `NOT IN` | Not in list | `space NOT IN ('ARCHIVE')` |

## Logical Operators

| Operator | Description |
|----------|-------------|
| `AND` | Both conditions |
| `OR` | Either condition |
| `NOT` | Negation |
| `()` | Grouping |

Example:

```cql
type = page AND (space = TEAM OR space = ENG)
```

## Date Values

### Relative Dates

| Value | Meaning |
|-------|---------|
| `-1d` | 1 day ago |
| `-1w` | 1 week ago |
| `-30d` | 30 days ago |
| `-1m` | 1 month ago |
| `-1y` | 1 year ago |

### Absolute Dates

```cql
created >= '2024-01-01'
lastModified >= '2024-06-15'
```

## User References

Reference users by account ID or username:

```cql
creator = 'user.name'
contributor = currentUser()
```

## Ordering Results

```cql
type = page AND space = TEAM ORDER BY created DESC
type = page ORDER BY lastModified DESC, title ASC
```

Sort fields:
- `created` - Creation date
- `lastModified` - Last modified date
- `title` - Alphabetical title

## Example Queries

### Pages in a Space

```cql
type = page AND space = TEAM
```

### Search by Title

```cql
type = page AND title ~ 'Onboarding'
```

### Full-Text Search

```cql
type = page AND text ~ 'authentication'
```

### Recently Modified

```cql
type = page AND lastModified >= -7d ORDER BY lastModified DESC
```

### Pages with Label

```cql
type = page AND label = 'architecture'
```

### Multiple Labels

```cql
type = page AND label IN ('rfc', 'approved')
```

### Pages by Creator

```cql
type = page AND creator = 'user.name' AND space = TEAM
```

### Child Pages

```cql
type = page AND parent = 123456789
```

### Descendant Pages

```cql
type = page AND ancestor = 123456789
```

### Combined Search

```cql
type = page AND
space = ENG AND
text ~ 'api' AND
lastModified >= -30d
ORDER BY lastModified DESC
```

### Exclude Archived

```cql
type = page AND space NOT IN ('ARCHIVE', 'OLD')
```

### Blog Posts

```cql
type = blogpost AND created >= -30d ORDER BY created DESC
```

### Find RFCs or Design Docs

```cql
type = page AND (title ~ 'RFC' OR title ~ 'Design Doc')
```

## Text Search Tips

The `~` operator searches content:

```cql
text ~ 'exact phrase'
title ~ 'keyword'
text ~ 'word1 word2'  -- both words must appear
text ~ 'word1 OR word2'  -- either word
```

## Combining with Space

Always scope searches to improve performance:

```cql
space = TEAM AND text ~ 'keyword'
space IN ('ENG', 'PROD') AND title ~ 'architecture'
```

## Escaping Values

Quote values with spaces:

```cql
space = 'My Space'
title ~ 'Getting Started Guide'
```

## Performance Tips

1. Include `space = X` to limit scope
2. Include `type = page` to exclude comments/attachments
3. Use date filters (`lastModified >= -Nd`) for large spaces
4. Use `--limit` flag to control result count
5. Search by `title ~` before `text ~` when possible (faster)

## Differences from JQL

| Aspect | JQL (Jira) | CQL (Confluence) |
|--------|------------|------------------|
| Primary field | `project` | `space` |
| Content type | `type` (Bug, Task) | `type` (page, blogpost) |
| Full text | `text ~` | `text ~` |
| Title | `summary ~` | `title ~` |
| Date fields | `created`, `updated` | `created`, `lastModified` |
| Hierarchy | `parent`, `"Epic Link"` | `parent`, `ancestor` |

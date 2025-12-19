---
name: conversation-search
description: Search indexed conversation history to find past conversations, locate when topics were discussed, or identify which project conversations occurred in. Use when user asks questions like "when did we discuss X?", "find conversations about Y", "in which project did we talk about Z?", or "show me conversations mentioning W". Works across all projects or within current project.
allowed-tools: Bash, Read
---

# Conversation Search Skill

Search through indexed Claude Code conversation history using SQLite full-text search.

## How to Search Conversations

### 1. Parse the User's Query

Extract:
- **Search terms**: The topic/keywords to search for
- **Scope**: Current project (default) or all projects
- **Result type**: First match only, or all matches

User mentions "all projects", "across projects", "everywhere", or "in which project" → use `all_projects` scope

### 2. Execute the Search Script

The search script is located in this skill's `scripts/` directory. Execute it:

```bash
cd <skill-base-directory>/scripts && ./search.sh --json --scope <scope> --project <current-working-directory> "<search-query>"
```

**Scope options:**
- `current_project` - Search only in current project (default)
- `all_projects` - Search across all projects

**Note:** The `<skill-base-directory>` is automatically provided by Claude Code as the base directory for this skill.

### 3. Parse and Present Results

The script returns JSON with matches:
```json
{
  "query": "...",
  "scope": "all_projects",
  "total_matches": 5,
  "matches": [
    {
      "uuid": "abc-123",
      "project_path": "/Users/.../project",
      "created_at": "2025-12-19T...",
      "message_count": 42,
      "summary": "Brief summary...",
      "relevance_score": 1.23
    }
  ]
}
```

### 4. Format Output

Present results based on user query:

**For "when did we first..."**: Show only the earliest match (last in results array)

**For "find all..." or "show me..."**: List all matches

**Format:**
```
Found N conversation(s) about [query]:

1. Conversation ID: abc-123
   Project: /path/to/project (only if all_projects)
   Date: Dec 19, 2025 at 10:30 AM
   Messages: 42
   Summary: Brief description...

[If many results] ...and X more conversations
```

**Important formatting notes:**
- Use "Conversation ID:" instead of "UUID:" for better readability
- Include Project only when searching all_projects
- Format date as human-readable (e.g., "Dec 19, 2025 at 10:30 AM")

### 5. Help User Resume Conversations

Remind users they can use `/continue <conversation-id>` to resume any conversation.

## Examples

**User**: "when did we first discuss authentication?"
→ Search for "authentication" in current project, show earliest match

**User**: "in which project did we talk about delivery tracking?"
→ Search "delivery tracking" with `all_projects` scope, group by project

**User**: "find all conversations mentioning bug fixes"
→ Search "bug fixes" in current project, show all matches

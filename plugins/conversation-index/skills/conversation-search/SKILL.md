---
name: conversation-search
description: Search indexed conversation history to find past conversations, locate when topics were discussed, or identify which project conversations occurred in. Use when user asks questions like "when did we discuss X?", "find conversations about Y", "in which project did we talk about Z?", or "show me conversations mentioning W". Works across all projects or within current project.
allowed-tools: Task
---

# Conversation Search Skill

## Purpose

This skill helps users find past conversations by searching through an indexed database of conversation history.

## When to Use

Activate when the user asks questions like:
- "when did we discuss X?"
- "find conversation about X"
- "which conversation had X?"
- "show me conversations about X"
- "what did we talk about regarding X?"
- "search for conversations where we did X"

## How It Works

1. Launch the search-agent with the user's query
2. The agent will:
   - Execute the search script against the conversation index
   - Determine scope (current project vs all projects)
   - Return matching conversations with UUIDs, dates, and summaries
3. Present results to the user

## Examples

**User**: "when did we first discuss authentication?"
**Action**: Launch search-agent with query "authentication", show earliest match

**User**: "find all conversations about the marketplace plugin"
**Action**: Launch search-agent with query "marketplace plugin", show all matches

**User**: "show me conversations about bug fixes across all projects"
**Action**: Launch search-agent with scope=all_projects

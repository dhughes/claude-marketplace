# Conversation Loader Plugin

Load the full content of previous Claude Code conversations into your current context.

## Overview

This plugin provides a skill that allows Claude to load and display complete conversation transcripts from previous sessions. This is particularly useful when you want to continue work from a previous conversation or reference detailed discussions.

## Features

- **Full Transcript Loading**: Loads complete user-assistant dialogue from previous conversations
- **Clean Format**: Displays only conversational content, filtering out tool execution details
- **File References**: Shows a summary of all files that were read or edited during the conversation
- **Large Context Support**: Designed for Sonnet models with 1M token context windows

## Installation

### From Doug's Marketplace

```bash
/plugin marketplace add doughughes/claude-marketplace
/plugin install conversation-loader@doug-marketplace
```

### Local Development

```bash
claude --plugin-dir /path/to/claude-marketplace/plugins/conversation-loader
```

## Usage

### Finding Conversations

First, use the `conversation-search` skill to find the conversation you want to load:

```
Find conversations about retry_before_failing
```

This will return conversation IDs like: `e06a3702-af08-41d7-a425-403622c2f266`

### Loading Conversations

Once you have a conversation ID, ask Claude to load it:

```
Load conversation e06a3702-af08-41d7-a425-403622c2f266
```

The skill will automatically activate and load the full transcript.

## Output Format

The loaded conversation includes:

1. **Header**: Metadata about the conversation
   - Conversation ID
   - Project path
   - Date created
   - Message count

2. **Transcript**: Complete dialogue
   - All user messages
   - All assistant text responses
   - Excludes thinking blocks and tool execution details

3. **Files Accessed**: Summary of files involved
   - List of files read
   - List of files edited

## Example Output

```
================================================================================
CONVERSATION
================================================================================
ID: e06a3702-af08-41d7-a425-403622c2f266
Project: /Users/you/Projects/my-project
Date: December 19, 2025 at 11:47 AM
Messages: 573

================================================================================
TRANSCRIPT
================================================================================

[USER]:
I want to refactor the authentication system...

[ASSISTANT]:
I'll help you refactor the authentication system...

...

================================================================================
FILES ACCESSED
================================================================================

Read:
  - /path/to/file1.rb
  - /path/to/file2.rb

Edited:
  - /path/to/file3.rb
```

## How It Works

The plugin:

1. Searches for the conversation file in `~/.claude/projects/`
2. Parses the JSONL file line by line
3. Extracts user and assistant messages
4. Tracks Read and Edit/Write tool uses
5. Formats everything into a readable transcript

## Use Cases

- **Continue Previous Work**: Resume complex refactoring or feature implementations
- **Reference Past Decisions**: Review architectural discussions and decisions
- **Debug Issues**: See the full context of how a bug was investigated
- **Learn from History**: Study how complex problems were solved

## Technical Details

- **File Location**: Conversations are stored at `~/.claude/projects/<encoded-path>/<uuid>.jsonl`
- **Format**: Each line in the JSONL file is a separate JSON document
- **Filtering**: Only user messages and assistant text responses are included
- **Tool Tracking**: Monitors Read, Edit, and Write tool calls for file references

## Requirements

- Claude Code CLI
- Python 3.x (for JSONL parsing)
- Bash shell

## License

MIT

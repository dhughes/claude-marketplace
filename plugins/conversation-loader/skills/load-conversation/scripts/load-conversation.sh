#!/usr/bin/env bash

# Load Conversation Script
# Parses a Claude Code conversation JSONL file and outputs a formatted transcript

set -euo pipefail

CONVERSATION_UUID="${1:-}"

if [ -z "$CONVERSATION_UUID" ]; then
    echo "Error: Conversation UUID is required" >&2
    echo "Usage: $0 <conversation-uuid>" >&2
    exit 1
fi

# Find the conversation file
CONVERSATION_FILE=$(find ~/.claude/projects -name "${CONVERSATION_UUID}.jsonl" -type f 2>/dev/null | head -1)

if [ -z "$CONVERSATION_FILE" ]; then
    echo "Error: Conversation file not found for UUID: $CONVERSATION_UUID" >&2
    echo "Searched in: ~/.claude/projects/*/${CONVERSATION_UUID}.jsonl" >&2
    exit 1
fi

# Export variables for Python script
export CONVERSATION_FILE
export CONVERSATION_UUID

# Use Python to parse the JSONL and format the output
python3 <<'PYTHON_SCRIPT'
import json
import sys
import os
from datetime import datetime
from pathlib import Path

conversation_file = os.environ.get('CONVERSATION_FILE')
conversation_uuid = os.environ.get('CONVERSATION_UUID')

if not conversation_file or not Path(conversation_file).exists():
    print(f"Error: Conversation file not found", file=sys.stderr)
    sys.exit(1)

# Parse the JSONL file
messages = []
files_read = set()
files_edited = set()
metadata = {
    'project': None,
    'created_at': None,
    'message_count': 0
}

with open(conversation_file, 'r') as f:
    for line in f:
        try:
            entry = json.loads(line.strip())
            entry_type = entry.get('type')

            # Extract metadata
            if metadata['project'] is None and 'cwd' in entry:
                metadata['project'] = entry['cwd']
            if metadata['created_at'] is None and 'timestamp' in entry:
                metadata['created_at'] = entry['timestamp']

            # Process user messages
            if entry_type == 'user':
                message = entry.get('message', {})
                content = message.get('content', '')

                # Skip tool results - they're not actual user messages
                # Tool results are arrays of dicts with 'tool_use_id' and 'type': 'tool_result'
                if isinstance(content, str) and content:
                    messages.append({
                        'role': 'user',
                        'content': content,
                        'timestamp': entry.get('timestamp')
                    })
                    metadata['message_count'] += 1

            # Process assistant messages
            elif entry_type == 'assistant':
                message = entry.get('message', {})
                content_blocks = message.get('content', [])

                text_parts = []
                for block in content_blocks:
                    block_type = block.get('type')

                    # Extract text content
                    if block_type == 'text':
                        text_parts.append(block.get('text', ''))

                    # Track file operations
                    elif block_type == 'tool_use':
                        tool_name = block.get('name', '')
                        tool_input = block.get('input', {})

                        # Track Read operations
                        if tool_name == 'Read':
                            file_path = tool_input.get('file_path')
                            if file_path:
                                files_read.add(file_path)

                        # Track Edit/Write operations
                        elif tool_name in ['Edit', 'Write']:
                            file_path = tool_input.get('file_path')
                            if file_path:
                                files_edited.add(file_path)

                # Add assistant message if there's text content
                if text_parts:
                    messages.append({
                        'role': 'assistant',
                        'content': '\n\n'.join(text_parts),
                        'timestamp': entry.get('timestamp')
                    })
                    metadata['message_count'] += 1

        except json.JSONDecodeError:
            continue

# Format and output the transcript
print("=" * 80)
print("CONVERSATION")
print("=" * 80)
print(f"ID: {conversation_uuid}")
print(f"Project: {metadata['project'] or 'Unknown'}")
if metadata['created_at']:
    dt = datetime.fromisoformat(metadata['created_at'].replace('Z', '+00:00'))
    print(f"Date: {dt.strftime('%B %d, %Y at %I:%M %p')}")
print(f"Messages: {metadata['message_count']}")
print()

print("=" * 80)
print("TRANSCRIPT")
print("=" * 80)
print()

for msg in messages:
    role = msg['role'].upper()
    content = msg['content']
    print(f"[{role}]:")
    print(content)
    print()

# Output files accessed
if files_read or files_edited:
    print("=" * 80)
    print("FILES ACCESSED")
    print("=" * 80)
    print()

    if files_read:
        print("Read:")
        for file_path in sorted(files_read):
            print(f"  - {file_path}")
        print()

    if files_edited:
        print("Edited:")
        for file_path in sorted(files_edited):
            print(f"  - {file_path}")
        print()

PYTHON_SCRIPT

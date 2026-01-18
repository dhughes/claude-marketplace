#!/bin/bash
set -euo pipefail

input=$(cat)
url=$(echo "$input" | jq -r '.tool_input.url // ""')

if [[ "$url" != *".atlassian.net"* ]]; then
  exit 0
fi

if [[ "$url" == *"/jira"* ]]; then
  cat >&2 << 'EOF'
{
  "decision": "block",
  "reason": "Cannot fetch Atlassian Jira URLs directly - authentication required.\n\nUse the atl CLI instead:\n  atl jira --help\n\nCommon commands:\n  atl jira get-issue PROJ-123\n  atl jira search-jql \"project = PROJ\""
}
EOF
  exit 2
fi

if [[ "$url" == *"/wiki"* ]]; then
  cat >&2 << 'EOF'
{
  "decision": "block",
  "reason": "Cannot fetch Atlassian Confluence URLs directly - authentication required.\n\nUse the atl CLI instead:\n  atl confluence --help\n\nCommon commands:\n  atl confluence get-page PAGE_ID\n  atl confluence search-cql \"type = page AND space = SPACE\""
}
EOF
  exit 2
fi

cat >&2 << 'EOF'
{
  "decision": "block",
  "reason": "Cannot fetch Atlassian URLs directly - authentication required.\n\nUse the atl CLI to interact with Atlassian:\n  atl --help"
}
EOF
exit 2

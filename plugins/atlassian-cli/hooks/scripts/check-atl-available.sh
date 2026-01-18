#!/bin/bash
set -euo pipefail

input=$(cat)
command_str=$(echo "$input" | jq -r '.tool_input.command // ""')

if [[ "$command_str" != *"atl "* ]] && [[ "$command_str" != "atl" ]]; then
  exit 0
fi

if command -v atl &> /dev/null; then
  exit 0
fi

cat >&2 << 'EOF'
{
  "decision": "block",
  "reason": "The `atl` CLI tool is not installed or not in PATH.\n\nTo install:\n1. Visit https://github.com/dhughes/atlassian-cli\n2. Follow the installation instructions\n3. Run `atl auth` to authenticate\n\nOnce installed, try your command again."
}
EOF
exit 2

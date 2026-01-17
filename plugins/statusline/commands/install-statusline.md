---
description: Configure Claude Code to use the statusline plugin
allowed-tools: Read, Write, Bash, AskUserQuestion
---

# Install Statusline Command

Configure Claude Code's settings to use the statusline plugin.

## Steps

1. **Locate settings file**: Find `~/.claude/settings.json`

2. **Check current configuration**: Read the settings file and check if `statusLine` is already configured

3. **If statusLine exists**:
   - Show the user what's currently configured
   - Use AskUserQuestion to ask: "A statusline is already configured. Overwrite with the statusline plugin?"
   - If user says no, exit without changes
   - If user says yes, continue

4. **Update settings**:
   - Set `statusLine.type` to `"command"`
   - Set `statusLine.command` to a bash command that dynamically finds the latest installed version
   - Use this exact command: `bash -c 'exec "$(ls -1d ~/.claude/plugins/cache/doug-marketplace/statusline/*/ 2>/dev/null | sort -V | tail -1)scripts/statusline.sh"'`
   - This ensures the statusline works across plugin updates without needing to re-run this command
   - Set `statusLine.padding` to `0`

5. **Write updated settings**: Save the modified settings.json

6. **Confirm success**: Tell the user: "Statusline installed successfully! Restart Claude Code or start a new session to see the statusline."

## Important Notes

- Do NOT show a preview of the statusline output
- The settings file path is `~/.claude/settings.json` (expand ~ to the user's home directory)
- The command uses a bash one-liner that automatically finds the latest installed version
- This approach ensures the statusline continues working after plugin updates
- If settings.json doesn't exist, create it with just the statusLine configuration
- Preserve all other settings in the file

## Example Settings Structure

```json
{
  "statusLine": {
    "type": "command",
    "command": "bash -c 'exec \"$(ls -1d ~/.claude/plugins/cache/doug-marketplace/statusline/*/ 2>/dev/null | sort -V | tail -1)scripts/statusline.sh\"'",
    "padding": 0
  }
}
```

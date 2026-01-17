# PR Best Practices Plugin

Enforces pull request best practices in Claude Code by automatically applying preferred workflows when creating PRs.

## Features

- **Include ticket numbers**: Automatically detects and includes ticket numbers from Jira or GitHub Issues in PR title and body
- **Always create PRs as drafts**: PRs are created with the `--draft` flag by default, allowing review before marking ready
- **Use PR templates**: Automatically checks for and uses GitHub PR templates from `.github/` directory

## Installation

```bash
# From this marketplace
claude plugin install pr-best-practices@claude-marketplace --scope user
```

## Usage

The plugin works automatically. When you ask Claude to create or open a pull request, the skill activates and ensures:

1. Claude detects ticket numbers from:
   - Conversation context (tickets mentioned in discussion)
   - Branch name (e.g., `MSI-608/zendesk2`, `feature/PROJ-123-description`, `456-fix-bug`)
   - Commit messages
   - Asks you if no ticket is found

2. Ticket is added to PR title: `[MSI-608] Add Zendesk integration`

3. Ticket is linked in PR body (when URL can be determined from context or project configuration)

4. Claude checks for PR templates in standard locations:
   - `.github/pull_request_template.md`
   - `.github/PULL_REQUEST_TEMPLATE.md`
   - `.github/PULL_REQUEST_TEMPLATE/*.md`

5. The PR is created as a draft using `gh pr create --draft`

### Examples

Simply ask Claude:
- "Create a PR for this branch"
- "Open a pull request"
- "Make a PR with these changes"

The plugin will automatically apply the best practices.

## How It Works

The plugin provides a skill that Claude automatically invokes when you request PR creation. The skill teaches Claude to:
- Detect and include ticket numbers from multiple sources (context, branch name, commits)
- Format PR titles with ticket references: `[TICKET] Title`
- Link tickets in PR body when URL can be determined
- Look for PR templates in your repository
- Always use the `--draft` flag when creating PRs
- Follow GitHub best practices for pull requests

### Ticket URL Detection

The plugin determines ticket URLs from:
- Project CLAUDE.md files (Jira base URL, project configuration)
- Git remote URLs (for GitHub issues)
- Conversation context and project knowledge
- Project configuration files

## Version

0.2.0

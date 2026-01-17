---
name: pr-creation
description: This skill should be used when the user asks to "create a PR", "open a PR", "submit a PR", "raise a PR", "file a PR", "make a PR", "create a pull request", "open a pull request", "new PR", or any variation requesting GitHub pull request creation.
allowed-tools: Read, Bash, Grep, Glob, AskUserQuestion
---

# Pull Request Creation Best Practices

## Purpose

This skill enforces three critical best practices when creating GitHub pull requests:

1. **Include ticket numbers** - Associates PRs with tracking tickets (Jira, GitHub Issues)
2. **Use PR templates when available** - Ensures consistent, complete PR descriptions
3. **Create PRs as drafts by default** - Allows review before marking ready for review

## Ticket Detection and Integration

Before creating the pull request, attempt to identify an associated ticket number from Jira or GitHub Issues.

### Detection Sources

Check the following sources for ticket numbers (in order of preference):

1. **Conversation context** - Ticket mentioned in recent messages or conversation history
2. **Branch name** - Parse current branch name for ticket patterns:
   - `PROJ-123/description` (Jira-style prefix)
   - `feature/PROJ-123-description` (Jira after prefix)
   - `123-description` (GitHub issue number)
   - Other common patterns using intuition
3. **Recent commit messages** - Check for ticket references in commit messages
4. **Ask user** - If no ticket found in above sources, use AskUserQuestion to ask if there's an associated ticket

### Ticket Number Patterns

Recognize these common patterns:
- **Jira tickets**: `PROJ-123`, `ABC-456` (uppercase letters, hyphen, numbers)
- **GitHub issues**: `#123`, `456` (numbers, optionally with # prefix)

### PR Title Integration

When ticket number is identified, format the PR title as:

```
[TICKET-NUM] Descriptive PR title
```

Examples:
- `[MSI-608] Add Zendesk integration`
- `[#456] Fix authentication bug`
- `[PLATFORM-123] Refactor API endpoints`

### PR Body Integration

**If PR template exists with ticket field:**
- Look for placeholders like `[TICKET]`, `[ISSUE]`, `Ticket:`, `Issue:`, `JIRA:`, etc.
- Replace placeholder with ticket number (as link if URL can be determined)
- Preserve template formatting

**If template has no specific ticket field:**
- Add ticket reference at the top of the PR body
- Format as: `**Ticket:** [TICKET-NUM](url)` or `**Ticket:** TICKET-NUM` if no URL

**If no template:**
- Include ticket reference at the start of the PR description
- Format with link if available

### Ticket URL Determination

Determine ticket URLs from available context and knowledge:

1. **Check CLAUDE.md or project configuration** - Look for Jira base URL, GitHub repo URL, or ticket system configuration
2. **Parse git remote** - Extract GitHub repository URL from git config for GitHub issues
3. **Use conversation context** - Apply knowledge from earlier in conversation about project's ticket systems
4. **Infer from project structure** - Look for `.jira`, project config files, or documentation mentioning ticket URLs

**URL formats:**
- Jira: `https://company.atlassian.net/browse/PROJ-123`
- GitHub Issues: `https://github.com/owner/repo/issues/123`

If URL cannot be determined, include ticket number without link.

## PR Template Discovery

Before creating the pull request, check for GitHub PR templates in the `.github/` directory using the following search order:

1. `.github/pull_request_template.md` (lowercase, standard GitHub convention)
2. `.github/PULL_REQUEST_TEMPLATE.md` (uppercase variant)
3. `.github/PULL_REQUEST_TEMPLATE/*.md` (multiple templates in subdirectory)

### Template Usage

If a template is found:
- Read the template file to understand the required sections
- Use the template structure when constructing the PR body with `gh pr create --body`
- Fill in each template section appropriately based on the changes being submitted
- Preserve template formatting, headings, and structure

If multiple templates exist in `.github/PULL_REQUEST_TEMPLATE/`:
- List available templates to the user
- Ask which template to use for this PR
- Proceed with the selected template

If no template is found:
- Proceed with standard PR creation
- Generate a comprehensive PR description including:
  - Summary of changes
  - Test plan or verification steps
  - Any relevant context or decisions

## Draft PR Creation

Always create pull requests as drafts using the `--draft` flag:

```bash
gh pr create --draft --title "PR title" --body "$(cat <<'EOF'
PR description here
EOF
)"
```

### Why Draft by Default

Draft PRs provide several benefits:
- Allows CI/CD checks to run before requesting reviews
- Provides opportunity to review the PR diff on GitHub before soliciting feedback
- Prevents accidental notification spam to reviewers for incomplete work
- Enables incremental commits and refinements before marking ready

### Making PR Ready for Review

After creating the draft PR, inform the user they can mark it ready for review when appropriate:

```bash
gh pr ready <pr-number>
```

Or they can mark it ready directly on GitHub.

## Complete Workflow

Follow this sequence when the user requests PR creation:

1. **Detect ticket number**:
   - Check conversation context for ticket references
   - Parse current branch name for ticket patterns
   - Check recent commit messages
   - If no ticket found, use AskUserQuestion to ask user
   - Determine ticket URL if possible from context, CLAUDE.md, or git config

2. **Check for PR templates**:
   - Search `.github/` directory for template files
   - If found, read template content
   - Identify if template has ticket field/placeholder

3. **Gather PR information**:
   - Analyze recent commits using `git log` and `git diff`
   - Generate appropriate PR title with ticket: `[TICKET] Title`
   - Generate PR description:
     - If template with ticket field: fill in ticket (with link if available)
     - If template without ticket field: add ticket at top of body
     - If no template: include ticket at start of description
   - If using template, structure description according to template format

4. **Create draft PR**:
   - Use `gh pr create --draft` with ticket-formatted title and body
   - Include template content if applicable
   - Ensure body is passed via heredoc for proper formatting

5. **Verify creation**:
   - Confirm PR was created successfully
   - Display PR URL to user
   - Remind user that PR is in draft state

6. **Provide next steps**:
   - Inform user how to mark PR ready when appropriate
   - Note any CI/CD checks that will run

## Example Commands

### PR with ticket from branch name

```bash
# Detect ticket from branch (e.g., MSI-608/zendesk2)
BRANCH=$(git branch --show-current)
# Parse ticket: MSI-608

# Create draft PR with ticket in title
gh pr create --draft --title "[MSI-608] Add Zendesk integration" --body "$(cat <<'EOF'
**Ticket:** [MSI-608](https://company.atlassian.net/browse/MSI-608)

## Summary
Added Zendesk integration for customer support tickets

## Changes
- Implemented Zendesk API client
- Added webhook handlers
EOF
)"
```

### PR with template containing ticket placeholder

```bash
# Template has "Ticket: [TICKET]" placeholder
# Replace with actual ticket and link

gh pr create --draft --title "[PLATFORM-123] Refactor API endpoints" --body "$(cat <<'EOF'
## Ticket
[PLATFORM-123](https://jira.company.com/browse/PLATFORM-123)

## Summary
[Filled from template]

## Test Plan
[Filled from template]
EOF
)"
```

### PR with GitHub issue number

```bash
# Detect from branch: 456-fix-auth-bug

gh pr create --draft --title "[#456] Fix authentication bug" --body "$(cat <<'EOF'
**Related Issue:** [#456](https://github.com/owner/repo/issues/456)

## Summary
Fixed authentication token expiration handling
EOF
)"
```

## Important Notes

- **Detect ticket numbers** - Always attempt to find associated ticket from context, branch, or commits before asking user
- **Always use `--draft` flag** - Never create PRs as immediately ready for review
- **Include ticket in title** - Format as `[TICKET] Title` when ticket is identified
- **Link tickets when possible** - Use context knowledge to determine ticket URLs
- **Respect template structure** - When templates exist, follow their format exactly
- **Use heredoc for body** - Ensures proper multiline formatting and prevents shell escaping issues
- **Verify template locations** - Check all three standard locations before proceeding
- **Inform the user** - Always communicate that the PR was created as a draft and how to mark it ready

## Edge Cases

**No ticket found and user declines to provide**: If user indicates there's no associated ticket or chooses not to provide one, proceed with PR creation without ticket in title or body

**Multiple ticket references**: If multiple tickets are found, ask user which is the primary ticket for this PR

**Invalid ticket format**: If detected ticket doesn't match expected patterns, verify with user before using

**Cannot determine ticket URL**: If ticket is identified but URL cannot be determined from context, include ticket number without link

**No git repository**: Verify current directory is a git repository before proceeding

**Not a GitHub repository**: Confirm remote is GitHub before using `gh` commands

**No GitHub CLI**: If `gh` command is not available, inform user to install GitHub CLI

**Branch already has PR**: If a PR already exists for the current branch, inform user rather than creating duplicate

**No commits to push**: Ensure there are commits to include in the PR before creating it

# GitHub and Git Analysis Guide for ezCater Research

Comprehensive guide for using `gh` CLI and `git` commands to research code history, PRs, and implementation details at ezCater.

## Overview

Use GitHub CLI (`gh`) for PRs, issues, and GitHub-hosted content. Use `git` commands for local repository history and blame analysis. Together they provide complete code evolution understanding.

## GitHub CLI (gh)

### Prerequisites

- `gh` CLI installed: `brew install gh` (macOS)
- Authenticated: `gh auth login`
- Verify: `gh auth status`

### Pull Request Commands

#### Search PRs

Find PRs by keywords:

```bash
gh pr list --repo OWNER/REPO --search "keywords"
```

**Common searches**:

```bash
# Search by title/body
gh pr list --repo ezcater/ez-rails --search "authentication"

# Search merged PRs
gh pr list --repo ezcater/ez-rails --search "is:merged order state"

# Search by author
gh pr list --repo ezcater/ez-rails --search "author:username"

# Search by date
gh pr list --repo ezcater/ez-rails --search "created:>=2024-01-01"

# Search open PRs
gh pr list --repo ezcater/ez-rails --search "is:open feature-flag"

# Complex search
gh pr list --repo ezcater/ez-rails --search "is:merged authentication created:>=2024-01-01" --limit 50
```

**Useful flags**:
- `--limit N` - Return N results (default 30)
- `--state {open|closed|merged|all}` - Filter by state
- `--json` - Output as JSON
- `--jq` - Parse JSON with jq expression

#### View PR Details

Get complete PR information:

```bash
gh pr view PR-NUMBER --repo OWNER/REPO
```

**Examples**:

```bash
# View PR
gh pr view 1234 --repo ezcater/ez-rails

# View with comments
gh pr view 1234 --repo ezcater/ez-rails --comments

# View as JSON
gh pr view 1234 --repo ezcater/ez-rails --json title,body,author,comments,files

# View PR from URL
gh pr view https://github.com/ezcater/ez-rails/pull/1234
```

**JSON output fields**:
```bash
gh pr view 1234 --repo ezcater/ez-rails --json \
  title,body,author,createdAt,mergedAt,state,reviews,comments,files,commits
```

#### View PR Diff

See code changes:

```bash
gh pr diff PR-NUMBER --repo OWNER/REPO
```

**Examples**:

```bash
# View diff
gh pr diff 1234 --repo ezcater/ez-rails

# View specific file
gh pr diff 1234 --repo ezcater/ez-rails -- path/to/file.rb
```

### Issue Commands

#### Search Issues

Find issues by keywords:

```bash
gh issue list --repo OWNER/REPO --search "keywords"
```

**Examples**:

```bash
# Search by keywords
gh issue list --repo ezcater/ez-rails --search "bug authentication"

# Search open issues
gh issue list --repo ezcater/ez-rails --state open

# Search by label
gh issue list --repo ezcater/ez-rails --label bug

# Search by assignee
gh issue list --repo ezcater/ez-rails --assignee username
```

#### View Issue

Get issue details:

```bash
gh issue view ISSUE-NUMBER --repo OWNER/REPO
```

**Example**:

```bash
gh issue view 567 --repo ezcater/ez-rails --comments
```

### Repository Commands

#### Clone Repository

```bash
gh repo clone OWNER/REPO [DIRECTORY]
```

**Example**:

```bash
gh repo clone ezcater/ez-rails ~/code/ezcater/ez/rails
```

#### View Repository

```bash
gh repo view OWNER/REPO
```

**Example**:

```bash
gh repo view ezcater/ez-rails --json description,url,defaultBranch
```

## Git History Analysis

### Log Commands

#### Basic Log

View commit history:

```bash
git log
```

**Useful flags**:
- `--oneline` - Compact format
- `--graph` - Show branch graph
- `-n N` - Limit to N commits
- `--since="date"` - Commits after date
- `--until="date"` - Commits before date
- `--author="name"` - Filter by author

**Examples**:

```bash
# Recent commits
git log --oneline -20

# Commits in last month
git log --since="1 month ago"

# Commits by author
git log --author="Doug Hughes"

# Visual branch history
git log --oneline --graph --all -20
```

#### Search Commit Messages

Find commits by message content:

```bash
git log --grep="PATTERN"
```

**Examples**:

```bash
# Find commits mentioning ticket
git log --grep="FX-1234"

# Find commits about feature
git log --grep="authentication"

# Case-insensitive search
git log --grep="bug" -i

# Search all branches
git log --all --grep="migration"

# Multiple patterns (OR)
git log --grep="feature" --grep="enhancement" --regexp-ignore-case

# Multiple patterns (AND)
git log --grep="auth" --grep="refactor" --all-match
```

#### File History

Trace evolution of specific file:

```bash
git log -- path/to/file
```

**Follow renames**:

```bash
git log --follow -- path/to/file
```

**Show changes**:

```bash
git log -p -- path/to/file
git log --follow -p -- path/to/file
```

**Examples**:

```bash
# File history with diffs
git log -p -- app/models/order.rb

# Follow through renames
git log --follow -- app/services/auth_service.rb

# Recent changes to file
git log --since="2024-01-01" -- config/routes.rb

# Who changed file
git log --pretty=format:"%h %an %ad %s" -- app/models/order.rb
```

#### Search Code Changes

Find when specific code was added or removed:

```bash
git log -S "CODE_STRING"
```

**Examples**:

```bash
# Find when function was added/removed
git log -S "def calculate_total"

# Find when string appeared
git log -S "error_message_text"

# With patches
git log -p -S "authentication_token"

# Regex search (more flexible)
git log -G "regex_pattern"
```

### Blame Analysis

#### Basic Blame

See who last modified each line:

```bash
git blame path/to/file
```

**Examples**:

```bash
# Blame file
git blame app/models/order.rb

# Blame specific lines
git blame -L 10,20 app/models/order.rb

# Ignore whitespace changes
git blame -w app/models/order.rb

# Show email addresses
git blame -e app/models/order.rb
```

#### Blame with Context

Understand why line was changed:

```bash
# Get commit SHA from blame
git blame app/models/order.rb

# View full commit
git show COMMIT_SHA

# View commit message only
git log -1 COMMIT_SHA
```

### Show Commands

#### View Commit

See complete commit details:

```bash
git show COMMIT_SHA
```

**Examples**:

```bash
# Full commit with diff
git show abc123def

# Just message and stats
git show --stat abc123def

# Just the diff
git show --no-patch abc123def

# Specific file in commit
git show abc123def:path/to/file
```

### Diff Commands

#### Compare Branches

```bash
git diff branch1...branch2
```

**Examples**:

```bash
# Changes from main
git diff main...feature-branch

# Just file names
git diff --name-only main...feature-branch

# Summary stats
git diff --stat main...feature-branch
```

#### Compare Commits

```bash
git diff commit1 commit2
```

**Example**:

```bash
# Specific file between commits
git diff abc123 def456 -- path/to/file.rb
```

## Common Research Workflows

### Find Original Implementation

**Goal**: Understand how feature was first implemented

1. **Search for initial PR**:
   ```bash
   gh pr list --repo ezcater/ez-rails --search "is:merged initial order" --limit 50
   ```

2. **Find first commit**:
   ```bash
   git log --all --grep="order" --reverse | head -20
   ```

3. **Trace file creation**:
   ```bash
   git log --diff-filter=A -- app/models/order.rb
   ```

4. **View original implementation**:
   ```bash
   git show COMMIT_SHA
   ```

### Understand Code Evolution

**Goal**: See how code changed over time

1. **File history with changes**:
   ```bash
   git log --follow -p -- app/services/order_processor.rb
   ```

2. **Find all PRs touching file**:
   ```bash
   # Get recent commits
   git log --oneline --follow -50 -- app/services/order_processor.rb

   # Search PRs by commit messages
   gh pr list --repo ezcater/ez-rails --search "is:merged order processor"
   ```

3. **Compare old vs new**:
   ```bash
   # Get commit from 6 months ago
   git log --since="6 months ago" --until="5 months ago" -- file.rb

   # Compare
   git diff OLD_SHA HEAD -- file.rb
   ```

### Investigate Bug Fix

**Goal**: Understand how bug was resolved

1. **Find bug ticket in commits**:
   ```bash
   git log --grep="FX-5678"
   ```

2. **Find related PRs**:
   ```bash
   gh pr list --repo ezcater/ez-rails --search "is:merged FX-5678"
   ```

3. **View fix PR**:
   ```bash
   gh pr view 1234 --repo ezcater/ez-rails --comments
   ```

4. **See code changes**:
   ```bash
   gh pr diff 1234 --repo ezcater/ez-rails
   ```

5. **Check if issue recurred**:
   ```bash
   git log --since="PR_MERGE_DATE" --grep="similar error"
   ```

### Trace Refactoring

**Goal**: Understand major code restructuring

1. **Find refactoring PR**:
   ```bash
   gh pr list --repo ezcater/ez-rails --search "is:merged refactor authentication"
   ```

2. **View PR details**:
   ```bash
   gh pr view PR_NUM --repo ezcater/ez-rails --json title,body,reviews,files
   ```

3. **See what was moved**:
   ```bash
   git log --follow -- app/services/new_location.rb
   ```

4. **Compare before/after**:
   ```bash
   # From PR, get base and head commits
   git diff BASE_SHA HEAD_SHA
   ```

### Find Code Author/Owner

**Goal**: Identify who knows about specific code

1. **Recent contributors**:
   ```bash
   git log --since="6 months ago" --pretty=format:"%an" -- path/to/file | sort | uniq -c | sort -rn
   ```

2. **Primary author**:
   ```bash
   git blame path/to/file | cut -d'(' -f2 | cut -d' ' -f1 | sort | uniq -c | sort -rn
   ```

3. **Recent PRs by author**:
   ```bash
   gh pr list --repo ezcater/ez-rails --search "is:merged author:username" --limit 20
   ```

## Advanced Techniques

### Finding When Bug Was Introduced

Use git bisect to find commit that introduced bug:

```bash
# Start bisect
git bisect start

# Mark current as bad
git bisect bad

# Mark old known-good commit
git bisect good COMMIT_SHA

# Git checks out middle commit, test it
# Mark as good or bad
git bisect good   # or git bisect bad

# Repeat until found
# Then reset
git bisect reset
```

### Finding Deleted Code

Find when code was removed:

```bash
# Search for deleted code
git log -S "deleted_function_name" --all

# Show when it was removed
git log -p -S "deleted_function_name" --all

# Find in all branches
git log --all --full-history -- path/to/deleted/file
```

### Comparing Branches

See what's different between branches:

```bash
# Commits in feature not in main
git log main..feature-branch --oneline

# Commits in main not in feature
git log feature-branch..main --oneline

# View diff
git diff main...feature-branch

# Files changed
git diff --name-only main...feature-branch
```

### Searching Across All Branches

Find commits in any branch:

```bash
# Search all branches
git log --all --grep="pattern"

# Find code in any branch
git log --all -S "code_string"

# Show which branch contains commit
git branch --contains COMMIT_SHA
```

## Integration with Other Tools

### gh → git

Use PR information to guide git exploration:

```bash
# Get PR details
gh pr view 1234 --repo ezcater/ez-rails --json commits

# Explore commits locally
git show COMMIT_SHA
git log COMMIT_SHA~5..COMMIT_SHA
```

### git → gh

Use commit information to find PRs:

```bash
# Find commit SHA
git log --grep="feature"

# Find PR containing commit
gh pr list --repo ezcater/ez-rails --search "SHA_FIRST_7_CHARS"
```

### git → atl

Use commit messages to find tickets:

```bash
# Find commits mentioning ticket
git log --grep="FX-1234"

# Then view ticket
atl jira issue view FX-1234
```

## Common Repositories

ezCater's main repositories:

| Repository | Path (if checked out) | Purpose |
|------------|----------------------|---------|
| ez-rails | ~/code/ezcater/ez/rails | Main monolith |
| delivery-management-rails | ~/code/ezcater/delivery-management-rails | Delivery tracking |
| omnichannel-rails | ~/code/ezcater/omnichannel-rails | Multi-channel features |
| liberty | ~/code/ezcater/liberty | Mobile backend |

**Note**: If repository is checked out locally, use `git` commands directly in that directory. Otherwise, use `gh` CLI with `--repo` flag.

## Best Practices

✅ **DO**:
- Use `--follow` with git log for renamed files
- Search commit messages for ticket numbers (FX-1234)
- Use JSON output from gh for parsing
- Check both PR description and comments
- Use git blame to understand why code exists
- Search all branches with `--all` when unsure
- Limit output with `-n` or `--limit` to avoid overwhelming results

❌ **DON'T**:
- Forget that files may have been renamed (use --follow)
- Ignore PR review comments (they contain context)
- Skip checking commit messages (they explain why)
- Search only current branch (may miss history)
- Use git for GitHub operations (use gh CLI instead)
- Assume first result is most relevant

## Quick Reference

### Most Common gh Commands

```bash
# Search PRs
gh pr list --repo ezcater/ez-rails --search "is:merged keywords"

# View PR
gh pr view 1234 --repo ezcater/ez-rails --comments

# Search issues
gh issue list --repo ezcater/ez-rails --search "keywords"

# Clone repo
gh repo clone ezcater/ez-rails ~/code/ezcater/ez/rails
```

### Most Common git Commands

```bash
# Search commits
git log --grep="pattern"

# File history
git log --follow -p -- path/to/file

# Find code changes
git log -S "code_string"

# Blame
git blame path/to/file

# Show commit
git show COMMIT_SHA

# Recent commits
git log --oneline -20
```

### Useful Git Aliases

Add to `.gitconfig` for shortcuts:

```ini
[alias]
  lg = log --oneline --graph --all -20
  find = log --all --grep
  who = shortlog -sn
  changed = log -p --follow
```

## Summary

GitHub CLI (`gh`) and Git provide powerful code archaeology capabilities. Use them to:
- Find PRs and issues related to features
- Trace code evolution through history
- Understand who wrote code and why
- Investigate bugs and their fixes
- Discover deleted or moved code

Master these tools for complete code history understanding, and integrate with Glean and atl CLI for comprehensive ezCater research.

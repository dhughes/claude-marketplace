---
description: Monitor GitHub PR CI checks until completion with audio notification
argument-hint: [pr-number]
allowed-tools: Bash, Read
---

Monitor the CI/CD checks for a GitHub Pull Request until all checks complete, then announce the result.

If a PR number is provided ($ARGUMENTS), monitor that specific PR.
If no PR number is provided, detect the PR from the current git branch.

Use the ci-monitor skill to perform the monitoring. Be tenacious - continue polling every 60 seconds until ALL checks complete (pass or fail). Never give up, even if CI takes an hour or more.

On completion:
- Print a clear summary of the results
- If checks failed, print failure details and the PR URL
- Use the macOS `say` command (if available) to announce the result audibly

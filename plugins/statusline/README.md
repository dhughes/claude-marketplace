# Statusline Plugin

A modular, pluggable statusline system for Claude Code that allows context-aware customization of your status line.

## Overview

This plugin provides a flexible statusline system where you can customize what information appears based on your context (work, side projects, personal projects).

## Installation

Install the plugin from the marketplace:

```bash
claude plugin install statusline@doug-marketplace --scope user
```

## Usage

### Install Statusline

Run the installation command to configure Claude Code to use the statusline:

```bash
/statusline:install-statusline
```

This command will:
1. Check if you have an existing statusline configuration
2. Update your `~/.claude/settings.json` to use the statusline script
3. Confirm successful installation

### Current Status (v0.5.1)

Current features:
- Installs the statusline entry point with version-agnostic configuration
- Displays workspace directory as a clickable hyperlink (opens in Finder/file browser)
  - Shows current directory in parentheses if different from workspace
- Shows current git branch (or "N/A" in muted gray if not in a git repo)
- Modular architecture with standalone tool binaries
- Comprehensive unit test coverage
- Supports macOS (Intel & Apple Silicon) and Linux (amd64 & arm64)

Future versions will add:
- Additional status tools (GitHub PR, Jira tickets, kubectl, tilt, usage stats, etc.)
- Per-project and global configuration
- Tool enable/disable configuration
- Custom status tool support from other plugins

## Requirements

- Claude Code with native statusline support
- Supported platforms: macOS, Linux (Windows support planned)

## Development

### Building from Source

The plugin includes pre-compiled binaries, but you can rebuild if needed:

```bash
cd plugins/statusline
make all
```

Requires Go 1.21 or later.

### Running Tests

Run all tests (unit + integration):

```bash
make test
```

Or run individually:
```bash
make test-go           # Unit tests only
make test-integration  # Integration test only
```

## License

MIT

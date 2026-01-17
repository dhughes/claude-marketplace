#!/bin/bash

# Platform detection wrapper for statusline
# This script detects the platform and executes the appropriate binary

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Detect operating system
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

# Detect architecture
ARCH=$(uname -m)

# Map architecture names to our binary naming convention
case "$ARCH" in
  x86_64)
    ARCH="amd64"
    ;;
  aarch64|arm64)
    ARCH="arm64"
    ;;
  *)
    echo "Error: Unsupported architecture: $ARCH" >&2
    exit 1
    ;;
esac

# Construct binary path
BINARY="${SCRIPT_DIR}/bin/statusline-${OS}-${ARCH}"

# Check if binary exists
if [ ! -f "$BINARY" ]; then
  echo "Error: Binary not found for ${OS}-${ARCH}" >&2
  echo "Expected: $BINARY" >&2
  exit 1
fi

# Make sure binary is executable
chmod +x "$BINARY" 2>/dev/null

# Execute the binary, passing stdin through
exec "$BINARY"

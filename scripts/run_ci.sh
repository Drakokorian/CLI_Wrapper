#!/usr/bin/env bash
set -euo pipefail

# Navigate to repo root
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$REPO_ROOT"

# Build and run the CI helper written in Go
GO_RUN="go run ./cmd/runci"

$GO_RUN

# Build desktop binaries if wails is available
if command -v wails >/dev/null 2>&1; then
  wails build -platform windows/amd64
  wails build -platform darwin/universal
  wails build -platform linux/amd64
  "$SCRIPT_DIR/sign.sh" build/bin/*
else
  echo "wails CLI not found, skipping UI build" >&2
fi

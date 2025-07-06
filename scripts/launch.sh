#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$REPO_ROOT"

# Determine log path
LOG_FILE="${SENTINEL_LOG_PATH:-}"
if [ -z "$LOG_FILE" ]; then
  case "$(uname -s)" in
    CYGWIN*|MINGW*|MSYS*) LOG_FILE="${SystemDrive:-C:}\\Temp\\sentinel.log" ;;
    *) LOG_FILE="/var/log/sentinel.log" ;;
  esac
fi
mkdir -p "$(dirname "$LOG_FILE")"

utc_now() { date -u +"%Y-%m-%dT%H:%M:%SZ"; }
log() { echo "$(utc_now) launch.sh: $1" >> "$LOG_FILE"; }

log "starting launch script"

if [ -z "${LAUNCH_SKIP_INSTALL:-}" ]; then
  log "installing prerequisites"
  "$SCRIPT_DIR/install.sh" >> "$LOG_FILE" 2>&1
fi

if ! command -v wails >/dev/null 2>&1; then
  log "wails CLI not found"
  echo "Please run scripts/install.sh to install prerequisites." >&2
  exit 1
fi

case "$(uname -s)" in
  Linux*)   PLATFORM="linux/amd64";   BIN="./build/bin/ai-cli-ui" ;;
  Darwin*)  PLATFORM="darwin/universal"; BIN="build/bin/ai-cli-ui.app" ;;
  CYGWIN*|MINGW*|MSYS*) PLATFORM="windows/amd64"; BIN="build\\bin\\ai-cli-ui.exe" ;;
  *) echo "Unsupported platform $(uname -s)" >&2; exit 1 ;;
esac

log "building for $PLATFORM"
if ! wails build -platform "$PLATFORM" >> "$LOG_FILE" 2>&1; then
  log "wails build failed"
  exit 1
fi

log "signing $BIN"
if ! "$SCRIPT_DIR/sign.sh" "$BIN" >> "$LOG_FILE" 2>&1; then
  log "signing failed"
  exit 1
fi

if [ -n "${LAUNCH_NOP:-}" ]; then
  log "dry run mode, skipping launch"
  exit 0
fi

log "launching $BIN"
if [ "$(uname -s)" = "Darwin" ]; then
  open "$BIN" >> "$LOG_FILE" 2>&1 &
else
  "$BIN" >> "$LOG_FILE" 2>&1 &
fi

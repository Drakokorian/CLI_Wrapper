# Phase 01: Core App Foundation

## Objective
Establish the base Wails application with cross-platform configuration and initial user interface.

## Detailed Steps
- Scaffold a new Wails v2 project using a Go backend and React frontend.
- Provide a `wails.json` file with configuration for Windows, macOS, and Linux.
- Create directories for logs, configuration, and state on launch:
  - Windows: `%APPDATA%\ai-cli-ui\`
  - macOS: `~/Library/Application Support/ai-cli-ui/`
  - Linux: `~/.config/ai-cli-ui/`
- Detect whether the `openai` or `gemini` CLI is available in the system path.
- Implement the initial UI with a prompt input, an output pane, and a model selector.
- Log each CLI invocation to `logs.txt` in the app directory.

## Expected Output
The application should launch successfully with one model available, create the necessary directories, and log the first CLI call.

## Dependencies
None. This phase creates the foundation for all subsequent work.

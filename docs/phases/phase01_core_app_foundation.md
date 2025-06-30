# Phase 01: Core App Foundation

## Objective
Establish the base Wails application with cross-platform configuration and initial user interface.

## Detailed Steps
- Scaffold a new Wails v2 project using a Go backend and React frontend.
- Install Node.js 18 and the Wails CLI before generating the project.
- Provide a `wails.json` file with configuration for Windows, macOS, and Linux.
- Create directories for logs, configuration, and state on launch:
  - Windows: `%APPDATA%\ai-cli-ui\`
  - macOS: `~/Library/Application Support/ai-cli-ui/`
  - Linux: `~/.config/ai-cli-ui/`
- Detect whether the `openai` or `gemini` CLI is available in the system path.
- Implement the initial UI with a prompt input, an output pane, and a model selector.
- Log each CLI invocation to `sentinel.log` in the external log directory
  (e.g., `/var/log` or `C:\\Temp`).
- Add `toolchain go1.24.x` to `go.mod` for consistent builds across environments.
- Document cross-platform build commands such as `wails build -platform windows/amd64`.

## Expected Output
The application should launch successfully with one model available, create the necessary directories, and log the first CLI call.

## Dependencies
None. This phase creates the foundation for all subsequent work.

## Status
Completed in **Sprint 1**.

## Sprint Plan
- [x] Scaffold the Wails project and React frontend
- [x] Configure cross-platform directories and logging
- [x] Detect available AI CLI tools
- [x] Build the initial prompt and output interface

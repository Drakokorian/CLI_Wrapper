# CLI Wrapper Utilities

This repository contains helper utilities for the CLI Wrapper project.

## Installing prerequisites

Run `scripts/install.sh` to install Go 1.24+, Node.js 18+, and the Wails CLI automatically for your platform.

## `phasefirst` Command

The `phasefirst` command extracts the first bullet from the **Detailed Steps**
section of each phase document under `docs/phases`. It outputs the result as
formatted JSON.

### Building

```
go build ./cmd/phasefirst
```

### Running

```
./phasefirst
```

The command writes operational logs to `sentinel.log` in `/var/log` on Unix or
`C:\Temp` on Windows. Logs are stored in JSON with automatic rotation. Both the
log file path and verbosity level are configurable via `config.json`.

### Testing

Run unit tests with:

```
go vet ./...
go test -race ./...
npm test --prefix frontend
```

You can also run `scripts/run_ci.sh` to execute `go vet`, `go test`, and `go build` with output logged to the sentinel log.

## Wails Application

The repository also contains a minimal Wails v2 application located in the
`frontend` directory with a Go backend.

### Building the desktop app

Use the Wails CLI to compile for your platform:

```
wails build -platform windows/amd64
wails build -platform darwin/universal
wails build -platform linux/amd64
```

Binaries are written to `build/bin`.

### Launching

Run the generated binary on your platform:

- **Windows**: `build\bin\ai-cli-ui.exe`
- **macOS**: `open build/bin/ai-cli-ui.app`
- **Linux**: `./build/bin/ai-cli-ui`

During development you can start a hot reload server with `wails dev`.

## Session Manager

The `internal/app` package now provides a `SessionManager` which limits concurrent
CLI sessions. The limit is stored in `config.json` under the application's config
directory. A default limit of one is used if no configuration exists. Sessions
receive unique UUIDs and can be terminated programmatically.

Resource throttling parameters are also configurable via `config.json`:

```
{
  "cpuThreshold": 35,
  "memoryThreshold": 35,
  "pollInterval": 2,
  "logLevel": "info",
  "logPath": "<path>"
}
```

If a running session exceeds these CPU or memory limits it is first throttled by
reducing process priority. Persistent overages cause the session to be
terminated. All alerts are recorded in `sentinel.log` within `/var/log` on
Unix or `C:\\Temp` on Windows.

## Model Switch Alerts

The Wails backend tracks the last-used model. When a new prompt specifies a different
model, the backend emits a `model:switched` event via the Wails event bus and logs the
change in JSON Lines format. Log files automatically rotate after reaching 20 MB.

## Streaming Session Output

Sessions now stream stdout and stderr lines in real time. The backend exposes a
`/stream` WebSocket endpoint that accepts a session ID and forwards output
chunks as they are produced. React components subscribe to these streams and
append lines dynamically so results appear incrementally.

## CLI Plugins

The application supports extensible CLI integrations via a plugin mechanism. Each plugin must implement the `plugins.Plugin` interface and register itself during initialization using `plugins.Register`. Registered plugins are automatically detected at runtime. See `internal/plugins/openai.go` and `internal/plugins/gemini.go` for reference implementations.

To force the application to use a specific CLI, set the `CLIWRAP_TOOL` environment variable to the plugin name before launching. When set and the selected tool is available, `DetectCLITool` returns that tool instead of auto-detecting.

## `cliwrap` Command

The `cliwrap` tool provides a minimal terminal interface for running a model and managing history.

### Examples

Run a prompt and wait for the result:

```bash
cliwrap run -model sh -prompt "echo hello"
```

Export previous runs to a file:

```bash
cliwrap history export > history.json
```

Import history from standard input:

```bash
cat history.json | cliwrap history import
```

Global flags `-concurrency` and `-workdir` update the same `config.json` used by the desktop app:

```bash
cliwrap -concurrency 2 -workdir /tmp run -model sh -prompt "echo hi"
```

You can override automatic CLI detection by exporting `CLIWRAP_TOOL` with the desired plugin name before running commands.


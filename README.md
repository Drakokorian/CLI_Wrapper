# CLI Wrapper Utilities

This repository contains helper utilities for the CLI Wrapper project.

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
`C:\Temp` on Windows. Logs are stored in JSON with automatic rotation.

### Testing

Run unit tests with:

```
go vet ./...
go test -race ./...
```

## Wails Application

The repository also contains a minimal Wails v2 application located in the
`frontend` directory with a Go backend. Building requires Node.js 18 and the
Wails CLI:

```
wails build -platform windows/amd64
wails build -platform darwin/universal
wails build -platform linux/amd64
```

## Session Manager

The `internal/app` package now provides a `SessionManager` which limits concurrent
CLI sessions. The limit is stored in `config.json` under the application's config
directory. A default limit of one is used if no configuration exists. Sessions
receive unique UUIDs and can be terminated programmatically.


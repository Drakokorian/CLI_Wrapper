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


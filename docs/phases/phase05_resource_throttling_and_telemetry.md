# Phase 05: Resource Throttling and Telemetry

## Objective
Monitor system usage and prevent the application from consuming excessive resources.

## Detailed Steps
- Use the `gopsutil` library to read process memory and CPU statistics.
- Alert the user if consumption exceeds 35% of system memory or CPU.
- If high usage persists, automatically throttle or cancel the affected session.
- Record resource alerts and throttling actions in `sentinel.log` located
  in the external log directory (e.g., `/var/log` or `C:\\Temp`).
- Poll resource metrics every two seconds to remain lightweight.
- Throttle by reducing process priority with platform-specific system calls.
- Allow thresholds and polling interval to be configured in `config.json`.

## Expected Output
The application adapts to system load, keeps resource usage under control, and logs significant events for review.

## Dependencies
Builds on the session management from Phase 02 and logging from Phase 04.

## Status
Completed in **Sprint 5**.

## Sprint Plan
- [x] Read CPU and memory metrics with gopsutil
- [x] Throttle or cancel overused sessions
- [x] Persist thresholds in configuration

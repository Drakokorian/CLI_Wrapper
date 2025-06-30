# Phase 04: Model Change Alerts and Logging

## Objective
Track model selections per session and make changes visible through the UI and logs.

## Detailed Steps
- Persist the last-used model for each session.
- When a model switch occurs, highlight the UI with a message like “⚠️ Model switched from gpt-4 to gemini-1.5-pro.”
- Append the change event, CLI command, prompt, and timestamp to
  `sentinel.log` in the external log directory (e.g., `/var/log` or
  `C:\\Temp`).
- Ensure all logged timestamps use UTC.
- Emit a runtime event to display alerts using the Wails event bus.
- Write logs in JSON Lines format and rotate files when they exceed 20 MB.
- Offer an option in settings to silence model-switch notifications.

## Expected Output
Users can easily identify which model was used for any prompt, both in the interface and in the log file.

## Dependencies
Requires theme and layout groundwork from Phase 03.

## Status
Completed in **Sprint 4**.

## Sprint Plan
- [x] Persist model selections per session
- [x] Emit alerts and log model switches
- [x] Rotate logs automatically

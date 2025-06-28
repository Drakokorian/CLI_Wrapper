# Phase 04: Model Change Alerts and Logging

## Objective
Track model selections per session and make changes visible through the UI and logs.

## Detailed Steps
- Persist the last-used model for each session.
- When a model switch occurs, highlight the UI with a message like “⚠️ Model switched from gpt-4 to gemini-1.5-pro.”
- Append the change event, CLI command, prompt, and timestamp to `logs.txt`.
- Ensure all logged timestamps use UTC.

## Expected Output
Users can easily identify which model was used for any prompt, both in the interface and in the log file.

## Dependencies
Requires theme and layout groundwork from Phase 03.

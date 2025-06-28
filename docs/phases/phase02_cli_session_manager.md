# Phase 02: CLI Session Manager

## Objective
Allow multiple concurrent CLI sessions with user control over execution.

## Detailed Steps
- On application start, prompt the user to specify the allowed number of concurrent AI CLI instances (1 to 5).
- Maintain an execution queue and allocate sessions dynamically based on the user's choice.
- Assign a UUID to each running session and display it in a sidebar list.
- Provide controls for terminating any session safely without leaving orphan processes.

## Expected Output
Users can queue or run several prompts at once, manage sessions from the UI, and gracefully stop any session when needed.

## Dependencies
Requires the core application foundation from Phase 01.

# Phase 07: Direct Billing Links and Session Summary

## Objective
Provide convenient links to billing information and optional usage summaries.

## Detailed Steps
- Detect which CLI tool is active—Gemini or OpenAI.
- Display a “🧾 View Billing” link that opens the correct provider page:
  - OpenAI: <https://platform.openai.com/account/billing>
  - Gemini: <https://makersuite.google.com/app/apikey>
- If the CLI exposes usage statistics, show the latest numbers in the session summary view.
- Open links using the Wails runtime so the default browser is used on every platform.
- Cache the last usage totals locally and refresh them each time the application starts.

## Expected Output
Users can view their billing page with a single click and optionally monitor usage data.

## Dependencies
Builds on previous phases to identify the CLI and store history data from Phase 06.

## Status
Completed in **Sprint 7**.

## Sprint Plan
- [x] Link to provider billing portals
- [x] Show usage statistics when available
- [x] Cache latest totals on startup

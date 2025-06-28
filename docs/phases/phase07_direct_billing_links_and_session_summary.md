# Phase 07: Direct Billing Links and Session Summary

## Objective
Provide convenient links to billing information and optional usage summaries.

## Detailed Steps
- Detect which CLI tool is active—Gemini or OpenAI.
- Display a “🧾 View Billing” link that opens the correct provider page:
  - OpenAI: <https://platform.openai.com/account/billing>
  - Gemini: <https://makersuite.google.com/app/apikey>
- If the CLI exposes usage statistics, show the latest numbers in the session summary view.

## Expected Output
Users can view their billing page with a single click and optionally monitor usage data.

## Dependencies
Builds on previous phases to identify the CLI and store history data from Phase 06.

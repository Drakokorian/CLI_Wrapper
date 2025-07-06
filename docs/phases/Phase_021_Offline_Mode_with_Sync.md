---
Phase: 021
Title: Offline Mode with Sync
Created: 2025-07-06
CompatibleWith:
  - Windows 10+ (PowerShell 5.1+)
  - Ubuntu 18.04+ (Bash 3.2+)
  - macOS 10.15+ (zsh/bash)
SprintWindow: Sprint 1 (Day 1–10)
Author: AutoDocSystem
---

# Phase 021: Offline Mode with Sync

## Objective
This phase covers offline mode with sync across all supported platforms.

## Tasks
- **PH021-T01**: Create local offline cache directory
  - TaskDetails: Create a cross-platform directory to store offline cache files. Use environment variables for custom paths. Document default locations.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 1
  - Documentation Block:
    ```bash
echo "create_local_offline_cache_directory"
```

- **PH021-T02**: Add configuration flag for offline mode
  - TaskDetails: Introduce an `offline` boolean in the configuration file. Provide CLI override via `--offline`. Make sure it works on all OSes.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 1
  - Documentation Block:
    ```bash
echo "add_configuration_flag_for_offline_mode"
```

- **PH021-T03**: Determine network connectivity method
  - TaskDetails: Implement a function that checks network reachability using a lightweight ping. Use platform-specific commands with fallbacks.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 1
  - Documentation Block:
    ```bash
echo "determine_network_connectivity_method"
```

- **PH021-T04**: Build CLI command to enter offline mode
  - TaskDetails: Add `cliapp offline start` which forces offline mode even if the network is available. Output the current status to the log.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 1
  - Documentation Block:
    ```bash
echo "build_cli_command_to_enter_offline_mode"
```

- **PH021-T05**: Implement offline queue for commands
  - TaskDetails: Store user commands in a local SQLite database when offline. Include timestamp and parameters for each entry.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 1
  - Documentation Block:
    ```bash
echo "implement_offline_queue_for_commands"
```

- **PH021-T06**: Add sync service to flush queue
  - TaskDetails: Create a background process that sends queued commands to the server when connectivity is restored. Retry on failure.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 2
  - Documentation Block:
    ```bash
echo "add_sync_service_to_flush_queue"
```

- **PH021-T07**: Encrypt cached data
  - TaskDetails: Protect offline cache using AES-256 encryption with a key stored in the user's config directory. Document key rotation.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 2
  - Documentation Block:
    ```bash
echo "encrypt_cached_data"
```

- **PH021-T08**: Write unit tests for offline queue
  - TaskDetails: Ensure commands are queued correctly and preserve order. Test across PowerShell and Bash.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 2
  - Documentation Block:
    ```bash
echo "write_unit_tests_for_offline_queue"
```

- **PH021-T09**: Provide manual sync command
  - TaskDetails: Allow users to run `cliapp sync` to force synchronization. Support `--dry-run` to preview actions.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 2
  - Documentation Block:
    ```bash
echo "provide_manual_sync_command"
```

- **PH021-T10**: Add UI notification for offline state
  - TaskDetails: Display a banner in the UI when offline mode is active. Include a button to trigger manual sync.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 2
  - Documentation Block:
    ```bash
echo "add_ui_notification_for_offline_state"
```

- **PH021-T11**: Handle conflict resolution
  - TaskDetails: Implement last-write-win strategy when local changes conflict with server data. Log conflicts for review.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 3
  - Documentation Block:
    ```bash
echo "handle_conflict_resolution"
```

- **PH021-T12**: Document environment variables
  - TaskDetails: Describe `CLIAPP_OFFLINE_PATH` and `CLIAPP_SYNC_URL` variables. Include platform-specific examples.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 3
  - Documentation Block:
    ```bash
echo "document_environment_variables"
```

- **PH021-T13**: Provide cache purge script
  - TaskDetails: Add `scripts/clear_offline_cache.sh` and `.ps1` variants to delete cached files securely.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 3
  - Documentation Block:
    ```bash
echo "provide_cache_purge_script"
```

- **PH021-T14**: Add exponential backoff to sync retries
  - TaskDetails: Retry sync operations with increasing delays. Cap retries at 5 attempts.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 3
  - Documentation Block:
    ```bash
echo "add_exponential_backoff_to_sync_retries"
```

- **PH021-T15**: Document troubleshooting steps
  - TaskDetails: Create `docs/offline_troubleshooting.md` with common issues and resolutions.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 3
  - Documentation Block:
    ```bash
echo "document_troubleshooting_steps"
```

- **PH021-T16**: Log offline events
  - TaskDetails: Write structured log entries to `sentinel_cli.log` whenever offline mode is toggled.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 4
  - Documentation Block:
    ```bash
echo "log_offline_events"
```

- **PH021-T17**: Secure offline storage
  - TaskDetails: Set directory permissions so only the current user can read cached files. Mention commands for chmod and icacls.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 4
  - Documentation Block:
    ```bash
echo "secure_offline_storage"
```

- **PH021-T18**: Implement incremental sync
  - TaskDetails: Only send changed items during synchronization to reduce bandwidth usage.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 4
  - Documentation Block:
    ```bash
echo "implement_incremental_sync"
```

- **PH021-T19**: Add CLI switch to disable offline mode
  - TaskDetails: `cliapp offline stop` exits offline mode and initiates sync. Return error if already online.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 4
  - Documentation Block:
    ```bash
echo "add_cli_switch_to_disable_offline_mode"
```

- **PH021-T20**: Integration test for offline workflow
  - TaskDetails: Write tests that simulate network loss and recovery to verify data integrity.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 4
  - Documentation Block:
    ```bash
echo "integration_test_for_offline_workflow"
```

- **PH021-T21**: Check data size before caching
  - TaskDetails: Warn the user if the cache exceeds 100 MB. Provide a setting to override the limit.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 5
  - Documentation Block:
    ```bash
echo "check_data_size_before_caching"
```

- **PH021-T22**: Show sync status command
  - TaskDetails: Implement `cliapp offline status` to display queue length and last sync time.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 5
  - Documentation Block:
    ```bash
echo "show_sync_status_command"
```

- **PH021-T23**: Add progress indicator during sync
  - TaskDetails: Display a percentage progress in the terminal while queued items are uploading.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 5
  - Documentation Block:
    ```bash
echo "add_progress_indicator_during_sync"
```

- **PH021-T24**: Schedule background sync
  - TaskDetails: Use `cron` on Linux/macOS and Task Scheduler on Windows to run sync every hour.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 5
  - Documentation Block:
    ```bash
echo "schedule_background_sync"
```

- **PH021-T25**: API endpoint for offline status
  - TaskDetails: Expose `/offline-status` so remote tools can query whether the CLI is offline.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 5
  - Documentation Block:
    ```bash
echo "api_endpoint_for_offline_status"
```

- **PH021-T26**: Sample config file
  - TaskDetails: Add `configs/offline.example.yaml` showing how to enable offline mode and custom paths.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 6
  - Documentation Block:
    ```bash
echo "sample_config_file"
```

- **PH021-T27**: Export offline queue command
  - TaskDetails: Provide `cliapp offline export <file>` to dump the pending queue to JSON.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 6
  - Documentation Block:
    ```bash
echo "export_offline_queue_command"
```

- **PH021-T28**: Compress cache for backup
  - TaskDetails: Offer `cliapp offline compress` that zips the cache folder for manual backups.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 6
  - Documentation Block:
    ```bash
echo "compress_cache_for_backup"
```

- **PH021-T29**: Telemetry for offline transitions
  - TaskDetails: Record metrics each time the app switches between offline and online.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 6
  - Documentation Block:
    ```bash
echo "telemetry_for_offline_transitions"
```

- **PH021-T30**: Warning for large cache
  - TaskDetails: Alert the user when the cache grows over 500 MB with instructions to prune.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 6
  - Documentation Block:
    ```bash
echo "warning_for_large_cache"
```

- **PH021-T31**: Migrate offline data
  - TaskDetails: Create scripts to move cached data to a new location, preserving permissions.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 7
  - Documentation Block:
    ```bash
echo "migrate_offline_data"
```

- **PH021-T32**: Override sync server URL
  - TaskDetails: Allow `--sync-url` flag to change the destination server. Validate HTTPS.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 7
  - Documentation Block:
    ```bash
echo "override_sync_server_url"
```

- **PH021-T33**: CI testing guidance
  - TaskDetails: Document how to simulate offline mode in continuous integration pipelines.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 7
  - Documentation Block:
    ```bash
echo "ci_testing_guidance"
```

- **PH021-T34**: Prevent queue overflow
  - TaskDetails: Cap the number of queued items and reject new ones once the limit is hit.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 7
  - Documentation Block:
    ```bash
echo "prevent_queue_overflow"
```

- **PH021-T35**: UI toggle for settings page
  - TaskDetails: Add a switch in Settings that lets users enable or disable offline mode.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 7
  - Documentation Block:
    ```bash
echo "ui_toggle_for_settings_page"
```

- **PH021-T36**: Example PowerShell network check
  - TaskDetails: Provide a snippet using `Test-Connection` to verify connectivity.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 8
  - Documentation Block:
    ```bash
echo "example_powershell_network_check"
```

- **PH021-T37**: Pre-flight offline check
  - TaskDetails: Run diagnostics before enabling offline mode to ensure storage is writable.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 8
  - Documentation Block:
    ```bash
echo "pre-flight_offline_check"
```

- **PH021-T38**: File locking for database
  - TaskDetails: Use platform-native file locks to prevent multiple instances from corrupting the offline DB.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 8
  - Documentation Block:
    ```bash
echo "file_locking_for_database"
```

- **PH021-T39**: Timestamped sync logs
  - TaskDetails: Include UTC timestamps in every sync attempt log entry.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 8
  - Documentation Block:
    ```bash
echo "timestamped_sync_logs"
```

- **PH021-T40**: Metrics for sync duration
  - TaskDetails: Collect statistics on how long each sync cycle takes and log them.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 8
  - Documentation Block:
    ```bash
echo "metrics_for_sync_duration"
```

- **PH021-T41**: Reset offline mode steps
  - TaskDetails: Document how to clear the queue and disable offline mode if errors persist.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 9
  - Documentation Block:
    ```bash
echo "reset_offline_mode_steps"
```

- **PH021-T42**: Cleanup orphaned files
  - TaskDetails: Scan cache for files not referenced by the queue and delete them.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 9
  - Documentation Block:
    ```bash
echo "cleanup_orphaned_files"
```

- **PH021-T43**: Admin audit script
  - TaskDetails: Provide a script that administrators can run to audit offline actions by user.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 9
  - Documentation Block:
    ```bash
echo "admin_audit_script"
```

- **PH021-T44**: Handle partial sync errors
  - TaskDetails: Ensure the queue resumes correctly after a failed upload without duplicating commands.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 9
  - Documentation Block:
    ```bash
echo "handle_partial_sync_errors"
```

- **PH021-T45**: Dry-run queue flush
  - TaskDetails: Support `cliapp sync --dry-run` to show what would be uploaded without sending data.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 9
  - Documentation Block:
    ```bash
echo "dry-run_queue_flush"
```

- **PH021-T46**: Offline mode user guide
  - TaskDetails: Finalize documentation detailing all offline features with examples for each OS.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 10
  - Documentation Block:
    ```bash
echo "offline_mode_user_guide"
```

- **PH021-T47**: Offline mode metrics dashboard
  - TaskDetails: Display metrics showing how often offline mode is used and queue sizes over time.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 10
  - Documentation Block:
    ```bash
    echo "offline_mode_metrics_dashboard"
    ```

- **PH021-T48**: Notifier when sync completes
  - TaskDetails: Show a desktop notification when a background sync completes successfully.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 10
  - Documentation Block:
    ```bash
    echo "notifier_when_sync_completes"
    ```

- **PH021-T49**: Limit offline retention days
  - TaskDetails: Allow configuration of how many days cached data is kept before purging.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 10
  - Documentation Block:
    ```bash
    echo "limit_offline_retention_days"
    ```

- **PH021-T50**: Deployment checklist
  - TaskDetails: Provide a checklist covering configuration, security, and testing before deployment.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 10
  - Documentation Block:
    ```bash
    echo "deployment_checklist"
    ```

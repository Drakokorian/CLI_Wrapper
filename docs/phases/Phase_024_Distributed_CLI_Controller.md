---
Phase: 024
Title: Distributed CLI Controller
Created: 2025-07-06
CompatibleWith:
  - Windows 10+ (PowerShell 5.1+)
  - Ubuntu 18.04+ (Bash 3.2+)
  - macOS 10.15+ (zsh/bash)
SprintWindow: Sprint 1 (Day 1–10)
Author: AutoDocSystem
---

# Phase 024: Distributed CLI Controller

## Objective
This phase covers distributed cli controller across all supported platforms.

## Tasks
- **PH024-T01**: Design distributed controller architecture
  - TaskDetails: Create diagrams outlining controller, agents, and communication channels.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 1
  - Documentation Block:
    ```bash
echo "design_distributed_controller_architecture"
```

- **PH024-T02**: Generate RSA keys script
  - TaskDetails: Write `scripts/gen_keys.sh` to create 4096-bit keys for TLS.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 1
  - Documentation Block:
    ```bash
echo "generate_rsa_keys_script"
```

- **PH024-T03**: Implement central controller
  - TaskDetails: Develop a Go service that queues and dispatches CLI commands to agents.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 1
  - Documentation Block:
    ```bash
echo "implement_central_controller"
```

- **PH024-T04**: Document controller API
  - TaskDetails: Describe REST endpoints for registering agents and submitting jobs.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 1
  - Documentation Block:
    ```bash
echo "document_controller_api"
```

- **PH024-T05**: Agent registration subcommand
  - TaskDetails: Add `cliapp agent register --controller <url>` to enroll a machine.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 1
  - Documentation Block:
    ```bash
echo "agent_registration_subcommand"
```

- **PH024-T06**: Windows agent service script
  - TaskDetails: Provide PowerShell instructions to run the agent as a Windows service.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 2
  - Documentation Block:
    ```bash
echo "windows_agent_service_script"
```

- **PH024-T07**: Linux agent service script
  - TaskDetails: Include systemd unit example for auto-start on boot.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 2
  - Documentation Block:
    ```bash
echo "linux_agent_service_script"
```

- **PH024-T08**: macOS agent launchd script
  - TaskDetails: Offer a plist example so the agent can run in the background.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 2
  - Documentation Block:
    ```bash
echo "macos_agent_launchd_script"
```

- **PH024-T09**: Heartbeat protocol
  - TaskDetails: Agents send a heartbeat every 30 seconds to indicate they are alive.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 2
  - Documentation Block:
    ```bash
echo "heartbeat_protocol"
```

- **PH024-T10**: List registered agents
  - TaskDetails: Command `cliapp agent list` returns all agents with status.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 2
  - Documentation Block:
    ```bash
echo "list_registered_agents"
```

- **PH024-T11**: Remote job queue
  - TaskDetails: Implement a persistent queue in the controller for submitted commands.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 3
  - Documentation Block:
    ```bash
echo "remote_job_queue"
```

- **PH024-T12**: Agent authentication model
  - TaskDetails: Use mutual TLS and signed certificates to verify each agent.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 3
  - Documentation Block:
    ```bash
echo "agent_authentication_model"
```

- **PH024-T13**: Controller URL environment variable
  - TaskDetails: Allow `CLIAPP_CONTROLLER_URL` to set the target controller.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 3
  - Documentation Block:
    ```bash
echo "controller_url_environment_variable"
```

- **PH024-T14**: Systemd unit example
  - TaskDetails: Provide `examples/agent.service` with recommended settings.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 3
  - Documentation Block:
    ```bash
echo "systemd_unit_example"
```

- **PH024-T15**: Launchd plist example
  - TaskDetails: Provide `examples/agent.plist` for macOS installations.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 3
  - Documentation Block:
    ```bash
echo "launchd_plist_example"
```

- **PH024-T16**: PowerShell service config
  - TaskDetails: Explain how to create a Windows service using `New-Service`.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 4
  - Documentation Block:
    ```bash
echo "powershell_service_config"
```

- **PH024-T17**: Fetch agent logs command
  - TaskDetails: Enable `cliapp agent logs <id>` to pull recent logs from an agent.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 4
  - Documentation Block:
    ```bash
echo "fetch_agent_logs_command"
```

- **PH024-T18**: Retry logic for dispatch
  - TaskDetails: Controller retries failed commands with exponential backoff.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 4
  - Documentation Block:
    ```bash
echo "retry_logic_for_dispatch"
```

- **PH024-T19**: Dashboard showing activity
  - TaskDetails: Offer a web UI summarizing active agents and queued jobs.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 4
  - Documentation Block:
    ```bash
echo "dashboard_showing_activity"
```

- **PH024-T20**: Remove agent script
  - TaskDetails: `cliapp agent remove <id>` deregisters an agent and cleans its credentials.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 4
  - Documentation Block:
    ```bash
echo "remove_agent_script"
```

- **PH024-T21**: Encrypt network traffic
  - TaskDetails: Enforce TLS 1.2+ on all connections and store certificates securely.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 5
  - Documentation Block:
    ```bash
echo "encrypt_network_traffic"
```

- **PH024-T22**: Query agent capabilities
  - TaskDetails: Expose an API to list features or available resources of each agent.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 5
  - Documentation Block:
    ```bash
echo "query_agent_capabilities"
```

- **PH024-T23**: Health check script
  - TaskDetails: Provide cross-platform script that verifies connectivity and required ports.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 5
  - Documentation Block:
    ```bash
echo "health_check_script"
```

- **PH024-T24**: Interactive session flag
  - TaskDetails: Support `cliapp run --controller <url> --interactive` for remote shells.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 5
  - Documentation Block:
    ```bash
echo "interactive_session_flag"
```

- **PH024-T25**: Local execution fallback
  - TaskDetails: If the controller is unreachable, commands fall back to local execution.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 5
  - Documentation Block:
    ```bash
echo "local_execution_fallback"
```

- **PH024-T26**: Agent registration tests
  - TaskDetails: Write tests covering registration flow including error cases.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 6
  - Documentation Block:
    ```bash
echo "agent_registration_tests"
```

- **PH024-T27**: Sample agent config file
  - TaskDetails: Add `configs/agent.yaml` with defaults for polling interval and log path.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 6
  - Documentation Block:
    ```bash
echo "sample_agent_config_file"
```

- **PH024-T28**: Custom agent plugins guide
  - TaskDetails: Document how to extend agents with additional command handlers.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 6
  - Documentation Block:
    ```bash
echo "custom_agent_plugins_guide"
```

- **PH024-T29**: Telemetry for distributed commands
  - TaskDetails: Record execution time and status for each dispatched job.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 6
  - Documentation Block:
    ```bash
echo "telemetry_for_distributed_commands"
```

- **PH024-T30**: Batch job execution example
  - TaskDetails: Illustrate how to run a batch script on multiple agents simultaneously.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 6
  - Documentation Block:
    ```bash
echo "batch_job_execution_example"
```

- **PH024-T31**: Role-based access control
  - TaskDetails: Only authorized users can dispatch or cancel remote jobs.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 7
  - Documentation Block:
    ```bash
echo "role-based_access_control"
```

- **PH024-T32**: Scale controllers horizontally
  - TaskDetails: Describe how to run multiple controllers behind a load balancer.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 7
  - Documentation Block:
    ```bash
echo "scale_controllers_horizontally"
```

- **PH024-T33**: Backup controller database
  - TaskDetails: Provide script to export the controller's state to a backup file.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 7
  - Documentation Block:
    ```bash
echo "backup_controller_database"
```

- **PH024-T34**: Cancel remote job command
  - TaskDetails: `cliapp job cancel <id>` stops a job in progress.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 7
  - Documentation Block:
    ```bash
echo "cancel_remote_job_command"
```

- **PH024-T35**: Offline agent queue behavior
  - TaskDetails: Agents store jobs locally until connectivity is restored.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 7
  - Documentation Block:
    ```bash
echo "offline_agent_queue_behavior"
```

- **PH024-T36**: Execution latency metrics
  - TaskDetails: Measure how long each command takes from dispatch to completion.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 8
  - Documentation Block:
    ```bash
echo "execution_latency_metrics"
```

- **PH024-T37**: Agent upgrade script
  - TaskDetails: Provide cross-platform instructions to update agents to the latest version.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 8
  - Documentation Block:
    ```bash
echo "agent_upgrade_script"
```

- **PH024-T38**: Audit remote command history
  - TaskDetails: Store details of who ran what command and when.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 8
  - Documentation Block:
    ```bash
echo "audit_remote_command_history"
```

- **PH024-T39**: Agent update channel variable
  - TaskDetails: Use `CLIAPP_AGENT_CHANNEL` to switch between stable and beta releases.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 8
  - Documentation Block:
    ```bash
echo "agent_update_channel_variable"
```

- **PH024-T40**: Integration tests for operations
  - TaskDetails: Automate tests that spin up controller and agents in containers.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 8
  - Documentation Block:
    ```bash
echo "integration_tests_for_operations"
```

- **PH024-T41**: Secure controller behind proxy
  - TaskDetails: Recommend placing the controller behind an authenticating reverse proxy.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 9
  - Documentation Block:
    ```bash
echo "secure_controller_behind_proxy"
```

- **PH024-T42**: Stream agent logs API
  - TaskDetails: Allow clients to stream logs in real time using WebSockets.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 9
  - Documentation Block:
    ```bash
echo "stream_agent_logs_api"
```

- **PH024-T43**: Filter jobs by tag
  - TaskDetails: Support labeling jobs and retrieving them by tag via API.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 9
  - Documentation Block:
    ```bash
echo "filter_jobs_by_tag"
```

- **PH024-T44**: Version mismatch handling
  - TaskDetails: Warn when agent version differs from controller version.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 9
  - Documentation Block:
    ```bash
echo "version_mismatch_handling"
```

- **PH024-T45**: Pause agent activity
  - TaskDetails: `cliapp agent pause <id>` temporarily stops new job execution.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 9
  - Documentation Block:
    ```bash
echo "pause_agent_activity"
```

- **PH024-T46**: Auto-restart failed agents
  - TaskDetails: Provide a watchdog script that restarts the agent if it crashes.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 10
  - Documentation Block:
    ```bash
echo "auto-restart_failed_agents"
```

- **PH024-T47**: Recommended timeouts
  - TaskDetails: Document default and maximum timeouts for remote commands.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 10
  - Documentation Block:
    ```bash
echo "recommended_timeouts"
```

- **PH024-T48**: Cleanup stale artifacts
  - TaskDetails: Remove old job logs and temp files after completion.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 10
  - Documentation Block:
    ```bash
echo "cleanup_stale_artifacts"
```

- **PH024-T49**: Tail job logs example
  - TaskDetails: `cliapp job logs --follow <id>` streams output like `tail -f`.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 10
  - Documentation Block:
    ```bash
echo "tail_job_logs_example"
```

- **PH024-T50**: Deployment guide
  - TaskDetails: Publish step-by-step instructions for deploying the controller in production.
  - Platform Compatibility Notes: PowerShell 5.1+, Bash 3.2+, and zsh compatible.
  - Security Considerations: Validate inputs, sanitize paths, and restrict permissions.
  - Estimated Sprint Day: Day 10
  - Documentation Block:
    ```bash
echo "deployment_guide"
```


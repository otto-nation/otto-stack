---
title: CLI Reference
description: Complete command reference for otto-stack CLI
lead: Comprehensive reference for all otto-stack CLI commands and their usage
date: "2025-10-01"
lastmod: "2026-02-27"
draft: false
weight: 50
toc: true
---

<!--
  ⚠️  AUTO-GENERATED FILE - DO NOT EDIT DIRECTLY
  To make changes, edit source files and run: task generate:docs
-->

# otto-stack CLI Reference

A powerful development stack management tool for streamlined local development automation

## Command Categories

### 📁 Project Management

Initialize, validate, and manage project setup

**Commands:** `init`, `validate`, `services`, `deps`, `conflicts`, `doctor`

### 🚀 Service Lifecycle

Start, stop, and manage running services

**Commands:** `up`, `down`, `restart`, `cleanup`

### ⚙️ Operations & Data

Monitor and manage service data

**Commands:** `status`, `logs`

### 🛠️ Utility

Information and development tools

**Commands:** `version`, `help`, `web-interfaces`

## Commands

### `up`

Start development stack services

Start one or more services in the development stack. The command is context-aware:
- **In a project directory**: Starts project services (including shared containers)
- **Outside a project**: Starts only shared containers (requires service names)

When sharing is enabled, containers are registered in ~/.otto-stack/shared/containers.yaml
to track which projects use them.

**Usage:** `otto-stack up [service...]`

**Aliases:** `start`, `run`

**Examples:**

```bash
otto-stack up
```

Start all configured services (in project context)

```bash
otto-stack up postgres redis
```

Start specific services

```bash
cd ~ && otto-stack up redis
```

Start shared containers from global context

```bash
otto-stack up --detach --build
```

Build images and start services in background

**Flags:**

- `--detach` (`bool`): Run services in background (detached mode) (default: `false`)
- `--build` (`bool`): Build images before starting services (default: `false`)
- `--force-recreate` (`bool`): Recreate containers even if config hasn't changed (default: `false`)
- `--no-deps` (`bool`): Don't start linked services (default: `false`)
- `--timeout` (`string`): Timeout for service startup (e.g., 30s, 2m) (default: `30s`)

**Related Commands:** [`down`](#down), [`restart`](#restart), [`status`](#status)

**Tips:**

- Add --build if you've made changes to Dockerfiles
- Use --detach to free up your terminal while services run
- Shared containers persist across projects when sharing is enabled

### `down`

Stop development stack services

Stop one or more services in the development stack. The command is context-aware:
- **In a project directory**: Stops project services, prompts before stopping shared containers
- **Outside a project**: Stops shared containers (requires service names)
- **With --shared flag**: Stops all shared containers from any location
- **With --all flag**: Stops both project and shared containers

When stopping shared containers, you'll be prompted if they're used by other projects.
The registry at ~/.otto-stack/shared/containers.yaml is updated to remove the project.

**Usage:** `otto-stack down [service...]`

**Aliases:** `stop`

**Examples:**

```bash
otto-stack down
```

Stop all running services (prompts for shared containers)

```bash
otto-stack down postgres redis
```

Stop specific services

```bash
otto-stack down --shared
```

Stop all shared containers

```bash
otto-stack down --all
```

Stop both project and shared containers

```bash
otto-stack down --volumes
```

Stop services and remove volumes

```bash
otto-stack down --timeout 5
```

Stop services with custom timeout

**Flags:**

- `--shared` (`bool`): Stop all shared containers (default: `false`)
- `--all` (`bool`): Stop both project and shared containers (default: `false`)
- `--volumes` (`bool`): Remove named volumes and anonymous volumes (default: `false`)
- `--remove-orphans` (`bool`): Remove containers for services not in compose file (default: `false`)
- `--timeout` (`int`): Shutdown timeout in seconds (default: `10`)

**Related Commands:** [`up`](#up), [`cleanup`](#cleanup), [`status`](#status)

**Tips:**

- Use --volumes carefully as it will delete all data
- Shared containers prompt before stopping if used by other projects
- Add --remove-orphans to clean up unused containers

### `restart`

Restart development stack services

Restart one or more services. This is equivalent to running down followed
by up, but more efficient for quick restarts.

**Usage:** `otto-stack restart [service...]`

**Examples:**

```bash
otto-stack restart
```

Restart all services

```bash
otto-stack restart postgres
```

Restart a specific service

```bash
otto-stack restart --timeout 5
```

Restart with custom timeout

**Flags:**

- `--timeout` (`int`): Restart timeout in seconds (default: `10`)
- `--no-deps` (`bool`): Don't restart linked services (default: `false`)

**Related Commands:** [`up`](#up), [`down`](#down), [`status`](#status)

### `status`

Show status of development stack services

Display comprehensive status information for services. The command is context-aware:
- **In a project directory**: Shows project services status
- **Outside a project**: Use --all or --shared flag to see shared containers
- **Specific project**: Use --project flag to see what a project uses

**Usage:** `otto-stack status [service...]`

**Aliases:** `ps`, `ls`

**Examples:**

```bash
otto-stack status
```

Show status of all services (in project context)

```bash
otto-stack status postgres redis
```

Show status of specific services

```bash
cd ~ && otto-stack status --all
```

Show all projects' shared containers (global context)

```bash
otto-stack status --shared
```

Show detailed shared container usage

```bash
otto-stack status --project my-app
```

Show shared containers used by specific project

```bash
otto-stack status --format json
```

Output status in JSON format

**Flags:**

- `--format` (`string`): Output format (table|json|yaml) (default: `table`) (options: `table`, `json`, `yaml`)
- `--all` (`bool`): Show status across all projects (including shared containers) (default: `false`)
- `--shared` (`bool`): Show detailed shared container usage and metrics (default: `false`)
- `--project` (`string`): Show shared containers used by specific project (default: ``)

**Related Commands:** [`logs`](#logs), [`status`](#status)

**Tips:**

- Try --format json for programmatic access
- Use --shared to see detailed container usage metrics
- Use --project to see what a specific project uses

### `logs`

View logs from services

View and follow logs from one or more services. Supports filtering,
timestamps, and real-time following. Logs from multiple services are
color-coded for easy identification.

**Usage:** `otto-stack logs [service...]`

**Examples:**

```bash
otto-stack logs
```

Show the last 100 lines from all services

```bash
otto-stack logs postgres redis
```

Show logs from specific services

```bash
otto-stack logs --follow postgres
```

Stream live logs from postgres

```bash
otto-stack logs --since 30m
```

Show logs from the last 30 minutes

**Flags:**

- `--follow` (`bool`): Follow log output in real-time (default: `false`)
- `--timestamps` (`bool`): Show timestamps (default: `false`)
- `--tail` (`string`): Number of lines to show from the end of the logs (default: `100`)
- `--since` (`string`): Show logs since a relative duration (e.g. 30m, 1h) or timestamp (default: ``)

**Related Commands:** [`status`](#status)

**Tips:**

- Logs from multiple services are color-coded for identification
- Use --follow to stream logs in real-time; interrupt with Ctrl+C

### `doctor`

Diagnose and troubleshoot stack health

Run comprehensive health checks on your development stack. Identifies
common issues, provides troubleshooting suggestions, and validates
service configurations.

**Usage:** `otto-stack doctor [service...]`

**Examples:**

```bash
otto-stack doctor
```

Run health checks on all services

```bash
otto-stack doctor postgres
```

Diagnose a specific service

**Flags:**

- `--format` (`string`): Output format (table|json) (default: `table`) (options: `table`, `json`)

**Related Commands:** [`status`](#status), [`logs`](#logs)

**Tips:**

- Run doctor when services aren't behaving as expected
- Use --format json for programmatic access to health check results

### `cleanup`

Clean up unused resources and data

Clean up unused Docker resources, temporary files, and orphaned data
created by otto-stack services. Helps reclaim disk space and maintain
a clean development environment.

**Usage:** `otto-stack cleanup [options]`

**Examples:**

```bash
otto-stack cleanup
```

Interactive cleanup with confirmations

```bash
otto-stack cleanup --all --force
```

Clean up everything without prompts

```bash
otto-stack cleanup --dry-run
```

Preview what would be cleaned up

**Flags:**

- `--all`, `-a` (`bool`): Clean up all resources (containers, volumes, images) (default: `false`)
- `--volumes` (`bool`): Remove unused volumes (default: `false`)
- `--images` (`bool`): Remove unused images (default: `false`)
- `--networks` (`bool`): Remove unused networks (default: `false`)
- `--force`, `-f` (`bool`): Don't prompt for confirmation (default: `false`)
- `--orphans` (`bool`): Clean up orphaned shared containers (containers with no projects) (default: `false`)
- `--project`, `-p` (`string`): Clean specific project (if not specified, cleans current project) (default: ``)

**Related Commands:** [`down`](#down), [`doctor`](#doctor)

**Tips:**

- Use --dry-run first to see what will be removed
- Be careful with --volumes as it removes all data

### `init`

Initialize a new otto-stack project interactively

Initialize a new otto-stack project in the current directory with an
interactive setup process. Guides you through selecting services,
configuring validation and advanced settings, and creates all
necessary configuration files.

**Usage:** `otto-stack init [flags]`

**Examples:**

```bash
otto-stack init
```

Interactive project initialization (recommended)

```bash
otto-stack init --project-name myproject --services postgres,redis
```

Non-interactive setup with specific project name and services

```bash
otto-stack init --force
```

Overwrite existing configuration

**Flags:**

- `--force`, `-f` (`bool`): Overwrite existing files (default: `false`)
- `--project-name` (`string`): Project name (defaults to current directory name) (default: ``)
- `--services` (`string`): Comma-separated list of services to include (required for non-interactive mode) (default: ``)
- `--no-shared-containers` (`bool`): Disable shared containers for all services (default: `false`)
- `--shared-services` (`string`): Comma-separated list of services to share across projects (e.g., postgres,redis). Overrides global sharing setting for specified services. (default: ``)

**Related Commands:** [`validate`](#validate)

### `web-interfaces`

Show web interfaces for running services

Display web interfaces (dashboards, UIs) for running services.
Shows URLs and availability status for easy access to service
management interfaces.

**Usage:** `otto-stack web-interfaces [service-name] [flags]`

**Examples:**

```bash
otto-stack web-interfaces
```

Show all web interfaces for running services

```bash
otto-stack web-interfaces localstack
```

Show interfaces for specific service

```bash
otto-stack web-interfaces --all
```

Show interfaces for all enabled services

**Flags:**

- `--all` (`bool`): Show interfaces for all enabled services, even if not running (default: `false`)

**Related Commands:** [`status`](#status), [`up`](#up)

### `services`

List available services by category

List all available services organized by category (database, cache,
messaging, observability, cloud). Shows service descriptions and
dependencies for easy discovery and selection.

**Usage:** `otto-stack services [flags]`

**Examples:**

```bash
otto-stack services
```

List all services grouped by category

```bash
otto-stack services --category database
```

List services in database category

```bash
otto-stack services --category cache
```

List cache services

```bash
otto-stack services --format table
```

List services in table format

**Flags:**

- `--category` (`string`): Show services in specific category (default: ``) (options: `database`, `cache`, `messaging`, `observability`, `cloud`)
- `--format` (`string`): Output format (group|table|json|yaml) (default: `group`) (options: `group`, `table`, `json`, `yaml`)

**Related Commands:** [`deps`](#deps), [`conflicts`](#conflicts), [`init`](#init)

### `deps`

Show dependency information for enabled project services

Display full dependency information for services in your project stack.
Shows required dependencies, soft dependencies, declared conflicts, and
provided capabilities for each enabled service. Columns with no data are
hidden automatically. Optionally filter to a specific service by name.

**Usage:** `otto-stack deps [service]`

**Examples:**

```bash
otto-stack deps
```

Show dependency info for all enabled services

```bash
otto-stack deps postgres
```

Show dependency info for a specific service

```bash
otto-stack deps kafka
```

Show what kafka requires and provides

**Related Commands:** [`services`](#services), [`conflicts`](#conflicts), [`up`](#up)

### `conflicts`

Detect service conflicts in the project stack

Analyze the enabled services in your project stack for conflicts.
Checks declared service incompatibilities and shared capability overlaps.
Use --check-ports to also verify that required host ports are available.
Returns exit code 1 when conflicts are found, making it safe to use in scripts.

**Usage:** `otto-stack conflicts [--check-ports]`

**Examples:**

```bash
otto-stack conflicts
```

Check for semantic conflicts in the project stack

```bash
otto-stack conflicts --check-ports
```

Also check if required ports are available on the host

```bash
otto-stack conflicts && otto-stack up
```

Fail fast if conflicts exist before starting services

**Flags:**

- `--check-ports` (`bool`): Also check if required service ports are available on the host (default: `false`)

**Related Commands:** [`services`](#services), [`deps`](#deps), [`up`](#up)

### `validate`

Validate project configuration and service definitions

Validate the otto-stack project configuration, project name, and service
definitions. Checks configuration syntax, validates that all enabled
services exist in the catalog, and resolves their dependencies.

With --strict, also verifies that Docker is available and warns when
the project is not inside a git repository.

**Usage:** `otto-stack validate`

**Examples:**

```bash
otto-stack validate
```

Validate configuration and service definitions

```bash
otto-stack validate --strict
```

Also check Docker availability and git repository

**Flags:**

- `--strict` (`bool`): Also check Docker availability and git repository (default: `false`)

**Related Commands:** [`doctor`](#doctor), [`deps`](#deps)

### `version`

Show version information

Display version information for otto-stack including build details,
Git commit, and platform information. Can check for available updates.

**Usage:** `otto-stack version`

**Examples:**

```bash
otto-stack version
```

Show basic version information

```bash
otto-stack version --full
```

Show detailed build information

```bash
otto-stack version --check-updates
```

Check for available updates

**Flags:**

- `--full` (`bool`): Show detailed version information (default: `false`)
- `--check-updates` (`bool`): Check for available updates (default: `false`)
- `--format` (`string`): Output format (text, json, yaml) (default: `text`)


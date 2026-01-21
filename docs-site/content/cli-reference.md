---
title: CLI Reference
description: Complete command reference for otto-stack CLI
lead: Comprehensive reference for all otto-stack CLI commands and their usage
date: "2025-10-01"
lastmod: "2026-01-21"
draft: false
weight: 50
toc: true
---

<!--
  ⚠️  AUTO-GENERATED FILE - DO NOT EDIT DIRECTLY
  This file is generated from internal/config/commands.yaml
  To make changes, edit the source file and run: task generate:docs
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

Monitor, connect to, and manage service data

**Commands:** `status`, `logs`, `exec`, `connect`

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
otto-stack down --volumes
```

Stop services and remove volumes

```bash
otto-stack down --timeout 5
```

Stop services with custom timeout

**Flags:**

- `--remove` (`bool`): Remove containers (default: false, just stops them) (default: `false`)
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
- **Outside a project**: Use --all flag to see all projects' shared containers

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
otto-stack status --format json
```

Output status in JSON format

**Flags:**

- `--format` (`string`): Output format (table|json|yaml) (default: `table`) (options: `table`, `json`, `yaml`)
- `--all` (`bool`): Show status across all projects (including shared containers) (default: `false`)

**Related Commands:** [`logs`](#logs), [`status`](#status)

**Tips:**

- Try --format json for programmatic access
- Use --all to see shared containers across all projects

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

Show logs from all services

```bash
otto-stack logs postgres redis
```

Show logs from specific services

**Related Commands:** [`status`](#status)

**Tips:**

- Logs from multiple services are color-coded for identification

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

```bash
otto-stack doctor --fix
```

Attempt to fix detected issues

**Flags:**

- `--fix` (`bool`): Attempt to automatically fix issues (default: `false`)
- `--format` (`string`): Output format (table|json) (default: `table`) (options: `table`, `json`)

**Related Commands:** [`status`](#status), [`logs`](#logs)

**Tips:**

- Run doctor when services aren't behaving as expected
- Use --fix to attempt automatic resolution of common issues

### `exec`

Execute commands in running service containers

Execute commands inside running service containers. Useful for database
operations, debugging, and maintenance tasks. Supports interactive and
non-interactive modes.

**Usage:** `otto-stack exec <service> <command> [args...]`

**Examples:**

```bash
otto-stack exec postgres psql -U postgres
```

Connect to PostgreSQL with psql

```bash
otto-stack exec redis redis-cli
```

Connect to Redis CLI

```bash
otto-stack exec postgres bash
```

Open bash shell in postgres container

**Flags:**

- `--user`, `-u` (`string`): Username to execute command as (default: ``)
- `--workdir` (`string`): Working directory for command (default: ``)
- `--interactive`, `-i` (`bool`): Keep STDIN open (interactive mode) (default: `true`)
- `--tty` (`bool`): Allocate a pseudo-TTY (default: `true`)
- `--detach` (`bool`): Run command in background (default: `false`)
- `--env` (`string`): Set environment variables (comma-separated key=value pairs) (default: ``)

**Related Commands:** [`connect`](#connect), [`logs`](#logs)

**Tips:**

- Use for database maintenance and debugging
- Combine with --user to run as specific user

### `connect`

Quick connect to service databases and interfaces

Quickly connect to service databases and management interfaces using
appropriate client tools. Automatically configures connection parameters
based on service configuration.

**Usage:** `otto-stack connect <service>`

**Examples:**

```bash
otto-stack connect postgres
```

Connect to PostgreSQL database

```bash
otto-stack connect redis
```

Connect to Redis CLI

```bash
otto-stack connect mysql
```

Connect to MySQL database

**Flags:**

- `--database` (`string`): Database name to connect to (default: ``)
- `--user`, `-u` (`string`): Username for connection (default: ``)
- `--host`, `-h` (`string`): Host to connect to (default: `localhost`)
- `--port`, `-p` (`int`): Port to connect to (default: `0`)
- `--read-only` (`bool`): Connect in read-only mode (default: `false`)

**Related Commands:** [`exec`](#exec), [`status`](#status)

**Tips:**

- Automatically uses correct client tools for each service
- Use --read-only for safe data exploration

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

Show dependency tree for a service

Display the complete dependency tree for a service, showing all required
dependencies and the resolved start order. Helps understand service
relationships and startup sequences.

**Usage:** `otto-stack deps <service>`

**Examples:**

```bash
otto-stack deps kafka-ui
```

Show dependencies for kafka-ui service

```bash
otto-stack deps postgres
```

Show dependencies for postgres service

**Related Commands:** [`services`](#services), [`conflicts`](#conflicts), [`up`](#up)

### `conflicts`

Check for conflicts between services

Check if the specified services have any conflicts that would prevent
them from running together. Identifies port conflicts, resource conflicts,
and incompatible service combinations.

**Usage:** `otto-stack conflicts <service1> <service2> [service...]`

**Examples:**

```bash
otto-stack conflicts postgres mysql
```

Check if postgres and mysql conflict

```bash
otto-stack conflicts postgres redis kafka-broker
```

Check conflicts between multiple services

**Related Commands:** [`services`](#services), [`deps`](#deps), [`up`](#up)

### `validate`

Validate configurations and manifests

Validate otto-stack configurations, service definitions, and YAML
manifests. Checks for syntax errors, missing dependencies, and
configuration inconsistencies.

**Usage:** `otto-stack validate [file...]`

**Examples:**

```bash
otto-stack validate
```

Validate all configuration files

```bash
otto-stack validate otto-stack-config.yaml
```

Validate specific configuration file

```bash
otto-stack validate --strict
```

Use strict validation rules

**Flags:**

- `--strict` (`bool`): Use strict validation rules (default: `false`)
- `--format` (`string`): Output format (table|json) (default: `table`) (options: `table`, `json`)
- `--fix` (`bool`): Attempt to fix validation errors (default: `false`)

**Related Commands:** [`doctor`](#doctor)

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

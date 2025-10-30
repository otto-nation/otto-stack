---
title: "CLI Reference"
description: "Complete command reference for otto-stack CLI"
lead: "Comprehensive reference for all otto-stack CLI commands and their usage"
date: "2025-10-01"
lastmod: "2025-10-30"
draft: false
weight: 50
toc: true
---

# otto-stack CLI Reference

Development stack management tool

Version: 0.1.0

## Commands

### cleanup

```
Clean up unused Docker resources, temporary files, and orphaned data
created by otto-stack services. Helps reclaim disk space and maintain
a clean development environment.

Usage:
  otto-stack cleanup [options] [flags]

Examples:
  otto-stack cleanup
    Interactive cleanup with confirmations

  otto-stack cleanup --all --force
    Clean up everything without prompts

  otto-stack cleanup --dry-run
    Preview what would be cleaned up



Flags:
  -a, --all        Clean up all resources (containers, volumes, images)
      --dry-run    Show what would be cleaned without doing it
  -f, --force      Don't prompt for confirmation
  -i, --images     Remove unused images
  -n, --networks   Remove unused networks
  -v, --volumes    Remove unused volumes

Global Flags:
  -c, --config string     Config file (default: $HOME/.otto-stack.yaml)
  -h, --help              Show help information
      --json              Output in JSON format (CI-friendly)
      --no-color          Disable colored output (CI-friendly)
      --non-interactive   Run in non-interactive mode (CI-friendly)
  -q, --quiet             Suppress non-essential output (CI-friendly)
  -v, --verbose           Enable verbose output
      --version           Show version information
```

### conflicts

```
Check if the specified services have any conflicts that would prevent
them from running together. Identifies port conflicts, resource conflicts,
and incompatible service combinations.

Usage:
  otto-stack conflicts <service1> <service2> [service...] [flags]

Examples:
  otto-stack conflicts postgres mysql
    Check if postgres and mysql conflict

  otto-stack conflicts postgres redis kafka-broker
    Check conflicts between multiple services



Global Flags:
  -c, --config string     Config file (default: $HOME/.otto-stack.yaml)
  -h, --help              Show help information
      --json              Output in JSON format (CI-friendly)
      --no-color          Disable colored output (CI-friendly)
      --non-interactive   Run in non-interactive mode (CI-friendly)
  -q, --quiet             Suppress non-essential output (CI-friendly)
  -v, --verbose           Enable verbose output
      --version           Show version information
```

### connect

```
Quickly connect to service databases and management interfaces using
appropriate client tools. Automatically configures connection parameters
based on service configuration.

Usage:
  otto-stack connect <service> [flags]

Examples:
  otto-stack connect postgres
    Connect to PostgreSQL database

  otto-stack connect redis
    Connect to Redis CLI

  otto-stack connect mysql
    Connect to MySQL database



Flags:
  -d, --database string   Database name to connect to
  -h, --host string       Host to connect to (default "localhost")
  -p, --port int          Port to connect to
      --read-only         Connect in read-only mode
  -u, --user string       Username for connection

Global Flags:
  -c, --config string     Config file (default: $HOME/.otto-stack.yaml)
  -h, --help              Show help information
      --json              Output in JSON format (CI-friendly)
      --no-color          Disable colored output (CI-friendly)
      --non-interactive   Run in non-interactive mode (CI-friendly)
  -q, --quiet             Suppress non-essential output (CI-friendly)
  -v, --verbose           Enable verbose output
      --version           Show version information
```

### deps

```
Display the complete dependency tree for a service, showing all required
dependencies and the resolved start order. Helps understand service
relationships and startup sequences.

Usage:
  otto-stack deps <service> [flags]

Examples:
  otto-stack deps kafka-ui
    Show dependencies for kafka-ui service

  otto-stack deps postgres
    Show dependencies for postgres service



Global Flags:
  -c, --config string     Config file (default: $HOME/.otto-stack.yaml)
  -h, --help              Show help information
      --json              Output in JSON format (CI-friendly)
      --no-color          Disable colored output (CI-friendly)
      --non-interactive   Run in non-interactive mode (CI-friendly)
  -q, --quiet             Suppress non-essential output (CI-friendly)
  -v, --verbose           Enable verbose output
      --version           Show version information
```

### doctor

```
Run comprehensive health checks on your development stack. Identifies
common issues, provides troubleshooting suggestions, and validates
service configurations.

Usage:
  otto-stack doctor [service...] [flags]

Examples:
  otto-stack doctor
    Run health checks on all services

  otto-stack doctor postgres
    Diagnose a specific service

  otto-stack doctor --fix
    Attempt to fix detected issues



Flags:
      --fix             Attempt to automatically fix issues
  -f, --format string   Output format (table|json) (default "table")
  -v, --verbose         Show detailed diagnostic information

Global Flags:
  -c, --config string     Config file (default: $HOME/.otto-stack.yaml)
  -h, --help              Show help information
      --json              Output in JSON format (CI-friendly)
      --no-color          Disable colored output (CI-friendly)
      --non-interactive   Run in non-interactive mode (CI-friendly)
  -q, --quiet             Suppress non-essential output (CI-friendly)
      --version           Show version information
```

### down

```
Stop one or more services in the development stack. By default, containers
are removed but volumes are preserved. Use --volumes to also remove data.

Usage:
  otto-stack down [service...] [flags]

Aliases:
  down, stop

Examples:
  otto-stack down
    Stop all running services

  otto-stack down postgres redis
    Stop specific services

  otto-stack down --volumes
    Stop services and remove volumes

  otto-stack down --timeout 5
    Stop services with custom timeout



Flags:
      --remove-images string   Remove images (all|local)
      --remove-orphans         Remove containers for services not in compose file
  -t, --timeout int            Shutdown timeout in seconds (default 10)
  -v, --volumes                Remove named volumes and anonymous volumes

Global Flags:
  -c, --config string     Config file (default: $HOME/.otto-stack.yaml)
  -h, --help              Show help information
      --json              Output in JSON format (CI-friendly)
      --no-color          Disable colored output (CI-friendly)
      --non-interactive   Run in non-interactive mode (CI-friendly)
  -q, --quiet             Suppress non-essential output (CI-friendly)
  -v, --verbose           Enable verbose output
      --version           Show version information
```

### exec

```
Execute commands inside running service containers. Useful for database
operations, debugging, and maintenance tasks. Supports interactive and
non-interactive modes.

Usage:
  otto-stack exec <service> <command> [args...] [flags]

Examples:
  otto-stack exec postgres psql -U postgres
    Connect to PostgreSQL with psql

  otto-stack exec redis redis-cli
    Connect to Redis CLI

  otto-stack exec postgres bash
    Open bash shell in postgres container



Flags:
  -d, --detach           Run command in background
  -e, --env string       Set environment variables (comma-separated key=value pairs)
  -i, --interactive      Keep STDIN open (interactive mode) (default true)
  -t, --tty              Allocate a pseudo-TTY (default true)
  -u, --user string      Username to execute command as
  -w, --workdir string   Working directory for command

Global Flags:
  -c, --config string     Config file (default: $HOME/.otto-stack.yaml)
  -h, --help              Show help information
      --json              Output in JSON format (CI-friendly)
      --no-color          Disable colored output (CI-friendly)
      --non-interactive   Run in non-interactive mode (CI-friendly)
  -q, --quiet             Suppress non-essential output (CI-friendly)
  -v, --verbose           Enable verbose output
      --version           Show version information
```

### init

```
Initialize a new otto-stack project in the current directory with an
interactive setup process. Guides you through selecting services,
configuring validation and advanced settings, and creates all
necessary configuration files.

Usage:
  otto-stack init [flags]

Examples:
  otto-stack init
    Interactive project initialization (recommended)

  otto-stack init --name myproject --minimal
    Non-interactive minimal setup

  otto-stack init --force
    Overwrite existing configuration



Flags:
  -f, --force   Overwrite existing files

Global Flags:
  -c, --config string     Config file (default: $HOME/.otto-stack.yaml)
  -h, --help              Show help information
      --json              Output in JSON format (CI-friendly)
      --no-color          Disable colored output (CI-friendly)
      --non-interactive   Run in non-interactive mode (CI-friendly)
  -q, --quiet             Suppress non-essential output (CI-friendly)
  -v, --verbose           Enable verbose output
      --version           Show version information
```

### logs

```
View and follow logs from one or more services. Supports filtering,
timestamps, and real-time following. Logs from multiple services are
color-coded for easy identification.

Usage:
  otto-stack logs [service...] [flags]

Examples:
  otto-stack logs
    Show logs from all services

  otto-stack logs postgres redis
    Show logs from specific services

  otto-stack logs --follow postgres
    Follow logs from postgres in real-time

  otto-stack logs --tail 100 --since 1h
    Show last 100 lines from the past hour



Flags:
  -f, --follow         Follow log output in real-time
      --no-color       Disable colored output
      --no-prefix      Don't show service name prefix
      --since string   Show logs since timestamp or relative time
  -t, --tail string    Number of lines to show from end of logs (default "all")
      --timestamps     Show timestamps in log output

Global Flags:
  -c, --config string     Config file (default: $HOME/.otto-stack.yaml)
  -h, --help              Show help information
      --json              Output in JSON format (CI-friendly)
      --non-interactive   Run in non-interactive mode (CI-friendly)
  -q, --quiet             Suppress non-essential output (CI-friendly)
  -v, --verbose           Enable verbose output
      --version           Show version information
```

### restart

```
Restart one or more services. This is equivalent to running down followed
by up, but more efficient for quick restarts.

Usage:
  otto-stack restart [service...] [flags]

Examples:
  otto-stack restart
    Restart all services

  otto-stack restart postgres
    Restart a specific service

  otto-stack restart --timeout 5
    Restart with custom timeout



Flags:
      --no-deps       Don't restart linked services
  -t, --timeout int   Restart timeout in seconds (default 10)

Global Flags:
  -c, --config string     Config file (default: $HOME/.otto-stack.yaml)
  -h, --help              Show help information
      --json              Output in JSON format (CI-friendly)
      --no-color          Disable colored output (CI-friendly)
      --non-interactive   Run in non-interactive mode (CI-friendly)
  -q, --quiet             Suppress non-essential output (CI-friendly)
  -v, --verbose           Enable verbose output
      --version           Show version information
```

### services

```
List all available services organized by category (database, cache,
messaging, observability, cloud). Shows service descriptions and
dependencies for easy discovery and selection.

Usage:
  otto-stack services [flags]

Examples:
  otto-stack services
    List all services grouped by category

  otto-stack services --category database
    List services in database category

  otto-stack services --category cache
    List cache services



Flags:
  -c, --category string   Show services in specific category

Global Flags:
  -c, --config string     Config file (default: $HOME/.otto-stack.yaml)
  -h, --help              Show help information
      --json              Output in JSON format (CI-friendly)
      --no-color          Disable colored output (CI-friendly)
      --non-interactive   Run in non-interactive mode (CI-friendly)
  -q, --quiet             Suppress non-essential output (CI-friendly)
  -v, --verbose           Enable verbose output
      --version           Show version information
```

### status

```
Display comprehensive status information for services including running
state, health checks, resource usage, and port mappings. Supports multiple
output formats and real-time monitoring.

Usage:
  otto-stack status [service...] [flags]

Aliases:
  status, ps, ls

Examples:
  otto-stack status
    Show status of all services

  otto-stack status postgres redis
    Show status of specific services

  otto-stack status --format json
    Output status in JSON format

  otto-stack status --watch
    Watch for status changes in real-time

  otto-stack status --filter running
    Show only running services



Flags:
      --filter string   Filter services by status
  -f, --format string   Output format (table|json|yaml) (default "table")
      --no-trunc        Don't truncate output
  -q, --quiet           Only show service names and basic status
  -w, --watch           Watch for status changes

Global Flags:
  -c, --config string     Config file (default: $HOME/.otto-stack.yaml)
  -h, --help              Show help information
      --json              Output in JSON format (CI-friendly)
      --no-color          Disable colored output (CI-friendly)
      --non-interactive   Run in non-interactive mode (CI-friendly)
  -v, --verbose           Enable verbose output
      --version           Show version information
```

### up

```
Start one or more services in the development stack. Services are started
with their configured dependencies and health checks. Use profiles to start
predefined service combinations.

Usage:
  otto-stack up [service...] [flags]

Aliases:
  up, start, run

Examples:
  otto-stack up
    Start all configured services

  otto-stack up postgres redis
    Start specific services

  otto-stack up --profile web
    Start services using the 'web' profile

  otto-stack up --detach --build
    Build images and start services in background



Flags:
  -b, --build             Build images before starting services
      --check-conflicts   Check for service conflicts before starting
  -d, --detach            Run services in background (detached mode)
      --force-recreate    Recreate containers even if config hasn't changed
      --no-deps           Don't start linked services
  -p, --profile string    Use a specific service profile
      --resolve-deps      Show dependency resolution tree before starting
  -t, --timeout string    Timeout for service startup (e.g., 30s, 2m) (default "30s")

Global Flags:
  -c, --config string     Config file (default: $HOME/.otto-stack.yaml)
  -h, --help              Show help information
      --json              Output in JSON format (CI-friendly)
      --no-color          Disable colored output (CI-friendly)
      --non-interactive   Run in non-interactive mode (CI-friendly)
  -q, --quiet             Suppress non-essential output (CI-friendly)
  -v, --verbose           Enable verbose output
      --version           Show version information
```

### validate

```
Validate otto-stack configurations, service definitions, and YAML
manifests. Checks for syntax errors, missing dependencies, and
configuration inconsistencies.

Usage:
  otto-stack validate [file...] [flags]

Examples:
  otto-stack validate
    Validate all configuration files

  otto-stack validate otto-stack-config.yaml
    Validate specific configuration file

  otto-stack validate --strict
    Use strict validation rules



Flags:
      --fix             Attempt to fix validation errors
  -f, --format string   Output format (table|json) (default "table")
  -s, --strict          Use strict validation rules

Global Flags:
  -c, --config string     Config file (default: $HOME/.otto-stack.yaml)
  -h, --help              Show help information
      --json              Output in JSON format (CI-friendly)
      --no-color          Disable colored output (CI-friendly)
      --non-interactive   Run in non-interactive mode (CI-friendly)
  -q, --quiet             Suppress non-essential output (CI-friendly)
  -v, --verbose           Enable verbose output
      --version           Show version information
```

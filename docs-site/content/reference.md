---
title: "CLI Reference"
description: "Complete command reference for otto-stack CLI"
lead: "Comprehensive reference for all otto-stack CLI commands and their usage"
date: "2025-10-01"
lastmod: "2025-11-10"
draft: "false"
weight: "50"
toc: "true"
---

# otto-stack CLI Reference

Usage:
otto-stack [command]

## Commands

### conflicts

```
Check if the specified services have any conflicts that would prevent
them from running together. Identifies port conflicts, resource conflicts,
and incompatible service combinations.

Usage:
  otto-stack conflicts [flags]

Global Flags:
      --dry-run           Show what would be done without executing
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
  otto-stack deps [flags]

Global Flags:
      --dry-run           Show what would be done without executing
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
  otto-stack doctor [flags]

Flags:
      --fix             Attempt to automatically fix issues
      --format string   Output format (table|json) (default "table")
  -v, --verbose         Show detailed diagnostic information

Global Flags:
      --dry-run           Show what would be done without executing
  -h, --help              Show help information
      --json              Output in JSON format (CI-friendly)
      --no-color          Disable colored output (CI-friendly)
      --non-interactive   Run in non-interactive mode (CI-friendly)
  -q, --quiet             Suppress non-essential output (CI-friendly)
      --version           Show version information

```

### exec

```
Execute commands inside running service containers. Useful for database
operations, debugging, and maintenance tasks. Supports interactive and
non-interactive modes.

Usage:
  otto-stack exec [flags]

Flags:
  -d, --detach           Run command in background
  -e, --env string       Set environment variables (comma-separated key=value pairs)
  -i, --interactive      Keep STDIN open (interactive mode) (default true)
  -t, --tty              Allocate a pseudo-TTY (default true)
  -u, --user string      Username to execute command as
  -w, --workdir string   Working directory for command

Global Flags:
      --dry-run           Show what would be done without executing
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

Flags:
  -f, --force   Overwrite existing files

Global Flags:
      --dry-run           Show what would be done without executing
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
  otto-stack logs [flags]

Flags:
      --follow         Follow log output in real-time
      --no-color       Disable colored output
      --no-prefix      Don't show service name prefix
      --since string   Show logs since timestamp or relative time
  -t, --tail string    Number of lines to show from end of logs (default "all")
      --timestamps     Show timestamps in log output

Global Flags:
      --dry-run           Show what would be done without executing
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
  otto-stack restart [flags]

Flags:
      --no-deps       Don't restart linked services
  -t, --timeout int   Restart timeout in seconds (default 10)

Global Flags:
      --dry-run           Show what would be done without executing
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

Flags:
      --category string   Show services in specific category
      --format string     Output format (group|table|json|yaml) (default "group")

Global Flags:
      --dry-run           Show what would be done without executing
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
  otto-stack status [flags]

Flags:
      --filter string   Filter services by status (default "<nil>")
      --format string   Output format (table|json|yaml) (default "table")
      --no-trunc        Don't truncate output
  -q, --quiet           Only show service names and basic status
  -w, --watch           Watch for status changes

Global Flags:
      --dry-run           Show what would be done without executing
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
with their configured dependencies and health checks.

Usage:
  otto-stack up [flags]

Flags:
  -b, --build             Build images before starting services
      --check-conflicts   Check for service conflicts before starting
  -d, --detach            Run services in background (detached mode)
      --force-recreate    Recreate containers even if config hasn't changed
      --no-deps           Don't start linked services
      --resolve-deps      Show dependency resolution tree before starting
  -t, --timeout string    Timeout for service startup (e.g., 30s, 2m) (default "30s")

Global Flags:
      --dry-run           Show what would be done without executing
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
  otto-stack validate [flags]

Flags:
      --fix             Attempt to fix validation errors
      --format string   Output format (table|json) (default "table")
  -s, --strict          Use strict validation rules

Global Flags:
      --dry-run           Show what would be done without executing
  -h, --help              Show help information
      --json              Output in JSON format (CI-friendly)
      --no-color          Disable colored output (CI-friendly)
      --non-interactive   Run in non-interactive mode (CI-friendly)
  -q, --quiet             Suppress non-essential output (CI-friendly)
  -v, --verbose           Enable verbose output
      --version           Show version information

```

### version

```
Display version information for otto-stack including build details,
Git commit, and platform information. Can check for available updates.

Usage:
  otto-stack version [flags]

Flags:
      --check-updates   Check for available updates
      --format string   Output format (text, json, yaml) (default "text")
      --full            Show detailed version information

Global Flags:
      --dry-run           Show what would be done without executing
  -h, --help              Show help information
      --json              Output in JSON format (CI-friendly)
      --no-color          Disable colored output (CI-friendly)
      --non-interactive   Run in non-interactive mode (CI-friendly)
  -q, --quiet             Suppress non-essential output (CI-friendly)
  -v, --verbose           Enable verbose output
      --version           Show version information

```

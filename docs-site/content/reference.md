---
title: "CLI Reference"
description: "Complete command reference for otto-stack CLI"
lead: "Comprehensive reference for all otto-stack CLI commands and their usage"
date: "2025-10-01"
lastmod: "2025-11-05"
draft: false
weight: 50
toc: true
---

# otto-stack CLI Reference

otto-stack - Development stack management tool

Usage:

## Commands

### cleanup

```
Clean up unused Docker resources, temporary files, and orphaned data
created by otto-stack services. Helps reclaim disk space and maintain
a clean development environment.

Usage:
  otto-stack cleanup [flags]

Flags:
  -h, --help   help for cleanup

Global Flags:
      --json              Output in JSON format
      --no-color          Disable colored output
      --non-interactive   Run in non-interactive mode
      --quiet             Suppress output
      --strict            Enable strict validation
```

### conflicts

```
Check if the specified services have any conflicts that would prevent
them from running together. Identifies port conflicts, resource conflicts,
and incompatible service combinations.

Usage:
  otto-stack conflicts [flags]

Flags:
  -h, --help   help for conflicts

Global Flags:
      --json              Output in JSON format
      --no-color          Disable colored output
      --non-interactive   Run in non-interactive mode
      --quiet             Suppress output
      --strict            Enable strict validation
```

### connect

```
Quickly connect to service databases and management interfaces using
appropriate client tools. Automatically configures connection parameters
based on service configuration.

Usage:
  otto-stack connect [flags]

Flags:
  -h, --help   help for connect

Global Flags:
      --json              Output in JSON format
      --no-color          Disable colored output
      --non-interactive   Run in non-interactive mode
      --quiet             Suppress output
      --strict            Enable strict validation
```

### deps

```
Display the complete dependency tree for a service, showing all required
dependencies and the resolved start order. Helps understand service
relationships and startup sequences.

Usage:
  otto-stack deps [flags]

Flags:
  -h, --help   help for deps

Global Flags:
      --json              Output in JSON format
      --no-color          Disable colored output
      --non-interactive   Run in non-interactive mode
      --quiet             Suppress output
      --strict            Enable strict validation
```

### doctor

```
Run comprehensive health checks on your development stack. Identifies
common issues, provides troubleshooting suggestions, and validates
service configurations.

Usage:
  otto-stack doctor [flags]

Flags:
  -h, --help   help for doctor

Global Flags:
      --json              Output in JSON format
      --no-color          Disable colored output
      --non-interactive   Run in non-interactive mode
      --quiet             Suppress output
      --strict            Enable strict validation
```

### down

```
Stop one or more services in the development stack. By default, containers
are removed but volumes are preserved. Use --volumes to also remove data.

Usage:
  otto-stack down [flags]

Flags:
  -h, --help   help for down

Global Flags:
      --json              Output in JSON format
      --no-color          Disable colored output
      --non-interactive   Run in non-interactive mode
      --quiet             Suppress output
      --strict            Enable strict validation
```

### exec

```
Execute commands inside running service containers. Useful for database
operations, debugging, and maintenance tasks. Supports interactive and
non-interactive modes.

Usage:
  otto-stack exec [flags]

Flags:
  -h, --help   help for exec

Global Flags:
      --json              Output in JSON format
      --no-color          Disable colored output
      --non-interactive   Run in non-interactive mode
      --quiet             Suppress output
      --strict            Enable strict validation
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
  -h, --help   help for init

Global Flags:
      --json              Output in JSON format
      --no-color          Disable colored output
      --non-interactive   Run in non-interactive mode
      --quiet             Suppress output
      --strict            Enable strict validation
```

### logs

```
View and follow logs from one or more services. Supports filtering,
timestamps, and real-time following. Logs from multiple services are
color-coded for easy identification.

Usage:
  otto-stack logs [flags]

Flags:
  -h, --help   help for logs

Global Flags:
      --json              Output in JSON format
      --no-color          Disable colored output
      --non-interactive   Run in non-interactive mode
      --quiet             Suppress output
      --strict            Enable strict validation
```

### restart

```
Restart one or more services. This is equivalent to running down followed
by up, but more efficient for quick restarts.

Usage:
  otto-stack restart [flags]

Flags:
  -h, --help   help for restart

Global Flags:
      --json              Output in JSON format
      --no-color          Disable colored output
      --non-interactive   Run in non-interactive mode
      --quiet             Suppress output
      --strict            Enable strict validation
```

### services

```
List all available services organized by category (database, cache,
messaging, observability, cloud). Shows service descriptions and
dependencies for easy discovery and selection.

Usage:
  otto-stack services [flags]

Flags:
  -h, --help   help for services

Global Flags:
      --json              Output in JSON format
      --no-color          Disable colored output
      --non-interactive   Run in non-interactive mode
      --quiet             Suppress output
      --strict            Enable strict validation
```

### status

```
Display comprehensive status information for services including running
state, health checks, resource usage, and port mappings. Supports multiple
output formats and real-time monitoring.

Usage:
  otto-stack status [flags]

Flags:
  -h, --help   help for status

Global Flags:
      --json              Output in JSON format
      --no-color          Disable colored output
      --non-interactive   Run in non-interactive mode
      --quiet             Suppress output
      --strict            Enable strict validation
```

### up

```
Start one or more services in the development stack. Services are started
with their configured dependencies and health checks.

Usage:
  otto-stack up [flags]

Flags:
  -h, --help   help for up

Global Flags:
      --json              Output in JSON format
      --no-color          Disable colored output
      --non-interactive   Run in non-interactive mode
      --quiet             Suppress output
      --strict            Enable strict validation
```

### validate

```
Validate otto-stack configurations, service definitions, and YAML
manifests. Checks for syntax errors, missing dependencies, and
configuration inconsistencies.

Usage:
  otto-stack validate [flags]

Flags:
  -h, --help   help for validate

Global Flags:
      --json              Output in JSON format
      --no-color          Disable colored output
      --non-interactive   Run in non-interactive mode
      --quiet             Suppress output
      --strict            Enable strict validation
```

### version

```
Display version information for otto-stack including build details,
Git commit, and platform information. Can check for available updates.

Usage:
  otto-stack version [flags]

Flags:
  -h, --help   help for version

Global Flags:
      --json              Output in JSON format
      --no-color          Disable colored output
      --non-interactive   Run in non-interactive mode
      --quiet             Suppress output
      --strict            Enable strict validation
```

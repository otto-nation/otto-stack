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

### version

```
Display version information for otto-stack including build details,
Git commit, and platform information. Can check for available updates.

Usage:
  otto-stack version [flags]

Examples:
  otto-stack version
    Show basic version information

  otto-stack version --full
    Show detailed build information

  otto-stack version --check-updates
    Check for available updates



Flags:
      --check-updates   Check for available updates
      --format string   Output format (text, json, yaml) (default "text")
      --full            Show detailed version information

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

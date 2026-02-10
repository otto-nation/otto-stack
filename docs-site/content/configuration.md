---
title: Configuration Guide
description: Configure your otto-stack development environment
lead: Learn how to configure your development stack
date: "2025-10-01"
lastmod: "2026-02-10"
draft: false
weight: 25
toc: true
---

<!--
  ⚠️  PARTIALLY GENERATED FILE
  - Sections marked with triple braces are auto-generated from internal/config/schema.yaml
  - Custom content (like "Sharing Configuration Details") is maintained in docs-site/templates/configuration-guide.md
  - To regenerate, run: task generate:docs
-->

# Configuration Guide

Otto-stack uses `.otto-stack/config.yaml` to define your development stack.

## File Structure

After running `otto-stack init`, you'll have:

```
.otto-stack/
├── config.yaml              # Main configuration
├── generated/
│   ├── .env.generated       # Available environment variables
│   └── docker-compose.yml   # Generated Docker Compose
├── services/                # Service metadata
│   ├── postgres.yml
│   └── redis.yml
├── .gitignore
└── README.md
```

## Main Configuration

**`.otto-stack/config.yaml`:**

```yaml
project:
  name: my-app
stack:
  enabled:
    - postgres
    - redis
sharing:
  enabled: false
validation:
  skip_warnings: false
  allow_multiple_databases: false
advanced:
  auto_start: false
  pull_latest_images: false
  cleanup_on_recreate: false
version_config:
  required_version: ">=1.0.0"
```

### Project

Project configuration settings

- **name**: Project name

### Stack

Stack service configuration

- **enabled**: List of enabled services

### Sharing

Container sharing configuration allows services to be shared across multiple projects, reducing resource usage and startup time

- **enabled**: Enable container sharing across projects. When enabled, containers are prefixed with 'otto-stack-' and tracked in ~/.otto-stack/shared/containers.yaml
- **services**: Per-service sharing overrides (service_name: true/false). If empty, all services are shared when enabled is true

### Validation

Validation and safety settings

- **skip_warnings**: Skip validation warnings during startup
- **allow_multiple_databases**: Allow multiple database services

### Advanced

Advanced operational settings

- **auto_start**: Start services automatically after setup
- **pull_latest_images**: Pull latest Docker images
- **cleanup_on_recreate**: Keep data when recreating services

### Version Config

Version management and update settings

- **required_version**: Required otto-stack version

### Sharing Configuration Details

When sharing is enabled:

1. Containers are prefixed with `otto-stack-` (e.g., `otto-stack-redis`)
2. A registry at `~/.otto-stack/shared/containers.yaml` tracks which projects use each shared container
3. The `down` command prompts before stopping shared containers used by other projects
4. Shared containers persist across project switches

**Example configurations:**

```yaml
# Share all services (default)
sharing:
  enabled: true

# Share specific services only
sharing:
  enabled: true
  services:
    postgres: true
    redis: true
    kafka: false  # Not shared

# Disable sharing completely
sharing:
  enabled: false
```

**Registry location:** `~/.otto-stack/shared/containers.yaml`

## Service Configuration

Services are configured through environment variables. Otto-stack generates `.otto-stack/generated/.env.generated` showing all available variables with defaults:

**Example `.env.generated`:**

```bash
# # POSTGRES
DATABASE_URL=postgresql://${POSTGRES_USER:-postgres}:${POSTGRES_PASSWORD:-password}@${POSTGRES_HOST:-localhost}:${POSTGRES_PORT:-5432}/${POSTGRES_DB:-local_dev}
PGHOST=${POSTGRES_HOST:-localhost}
POSTGRES_DB=${POSTGRES_DB:-local_dev}
POSTGRES_HOST=${POSTGRES_HOST:-localhost}
# REDIS
REDIS_HOST=${REDIS_HOST:-localhost}
REDIS_PASSWORD=${REDIS_PASSWORD:-password}
REDIS_PORT=${REDIS_PORT:-6379}
REDIS_URL=redis://:${REDIS_PASSWORD:-password}@${REDIS_HOST:-localhost}:${REDIS_PORT:-6379}
```

### Customizing Services

Create a `.env` file in your project root to override defaults:

```bash
# Postgres
DATABASE_URL=my_custom_value
PGHOST=my_custom_value
# Redis
REDIS_HOST=my_custom_value
REDIS_PASSWORD=my_custom_value
```

These values will be used by Docker Compose when starting services.

## Service Metadata Files

Files in `.otto-stack/services/` contain service metadata:

**`.otto-stack/services/postgres.yml`:**

```yaml
name: postgres
description: Configuration for postgres service
```

These are informational and don't affect service behavior. Configuration happens via environment variables.

## Complete Example

**`.otto-stack/config.yaml`:**

```yaml
project:
  name: my-fullstack-app
  type: docker
stack:
  enabled:
    - postgres
    - redis
    - kafka
validation:
  options:
    config-syntax: true
    docker: true
```

**`.env` (your customizations):**

```bash
# Postgres
DATABASE_URL=production_value
PGHOST=production_value
# Redis
REDIS_HOST=production_value
REDIS_PASSWORD=production_value
```

## Next Steps

- **[Services Guide](/services)** - Available services and environment variables
- **[CLI Reference](/cli-reference)** - Command usage
- **[Troubleshooting](/troubleshooting)** - Common issues

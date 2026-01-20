---
title: Configuration Guide
description: Configure your otto-stack development environment
lead: Learn how to configure your development stack
date: "2025-10-01"
lastmod: "2026-01-20"
draft: false
weight: 25
toc: true
---

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
  name: "my-app"
  type: "docker"

stack:
  enabled:
    - postgres
    - redis

validation:
  options:
    config-syntax: true
    docker: true
    service-definitions: true
```

### Project Section

- **name** (required): Project identifier
- **type** (optional): Project type (docker, web, api, microservice)

### Stack Section

- **enabled**: Array of service names to run

See [Services Guide](services.md) for available services.

### Validation Section

Control validation checks:

- **config-syntax**: Validate YAML syntax
- **docker**: Check Docker availability
- **service-definitions**: Validate service configs

## Service Configuration

Services are configured through environment variables. Otto-stack generates `.otto-stack/generated/.env.generated` showing all available variables with defaults:

**Example `.env.generated`:**

```bash
# POSTGRES
POSTGRES_DB=${POSTGRES_DB:-local_dev}
POSTGRES_USER=${POSTGRES_USER:-postgres}
POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-password}
POSTGRES_PORT=${POSTGRES_PORT:-5432}

# REDIS
REDIS_PASSWORD=${REDIS_PASSWORD:-password}
REDIS_PORT=${REDIS_PORT:-6379}
```

### Customizing Services

Create a `.env` file in your project root to override defaults:

```bash
# .env
POSTGRES_DB=my_custom_db
POSTGRES_PASSWORD=secure_password
POSTGRES_PORT=5433

REDIS_PASSWORD=redis_secure
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
  type: web

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
# Database
POSTGRES_DB=production_db
POSTGRES_PASSWORD=super_secure_password

# Cache
REDIS_PASSWORD=redis_secure_password

# Messaging
KAFKA_HEAP_OPTS=-Xmx512M -Xms256M
```

## Next Steps

- **[Services Guide](services.md)** - Available services and environment variables
- **[CLI Reference](cli-reference.md)** - Command usage
- **[Troubleshooting](troubleshooting.md)** - Common issues

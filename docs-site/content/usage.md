---
title: "Usage"
description: "Daily usage patterns and common workflows for otto-stack"
lead: "Learn how to effectively use otto-stack for your development workflow"
date: "2025-10-01"
lastmod: "2025-11-10"
draft: false
weight: 20
toc: true
---

# Usage Guide

## 🚀 Quick Start

```bash
# Initialize a new project
otto-stack init

# Start your development stack
otto-stack up

# Check service status
otto-stack status

# View logs
otto-stack logs

# Stop services
otto-stack down
```

## 📋 Basic Commands

### Project Management

```bash
# Initialize new project with configuration
otto-stack init

# Start all configured services
otto-stack up

# Stop all services
otto-stack down

# Restart services
otto-stack restart
```

### Service Management

```bash
# Start specific service
otto-stack up postgres

# Stop specific service
otto-stack down redis

# Check service status
otto-stack status postgres

# View service logs
otto-stack logs postgres --follow
```

### Health & Diagnostics

```bash
# Check overall system health
otto-stack doctor

# View detailed status
otto-stack status --verbose

# List available services
otto-stack services list
```

## 🔧 Common Workflows

### Database Development

```bash
# Start database services
otto-stack up postgres redis

# Connect to database
otto-stack exec postgres psql -U postgres

# View database logs
otto-stack logs postgres
```

### Full Stack Development

```bash
# Start complete development environment
otto-stack up

# Monitor all services
otto-stack status

# View aggregated logs
otto-stack logs --all
```

### Service Testing

```bash
# Start specific services for testing
otto-stack up postgres redis kafka

# Run health checks
otto-stack doctor

# Clean restart
otto-stack down && otto-stack up
```

## 📁 Configuration

Otto-stack uses `otto-stack-config.yaml` for configuration:

```yaml
# Example configuration
services:
  - postgres
  - redis
  - kafka

service_configuration:
  postgres:
    database: my_app_dev
    password: secure_password
  redis:
    password: redis_password
```

See [Configuration Guide](configuration.md) for detailed options.

## 🔍 Monitoring

### Check Service Health

```bash
# Overall health check
otto-stack doctor

# Individual service status
otto-stack status postgres

# Resource usage
docker stats
```

### View Logs

```bash
# All services
otto-stack logs

# Specific service
otto-stack logs postgres

# Follow logs in real-time
otto-stack logs postgres --follow

# Last 100 lines
otto-stack logs postgres --tail 100
```

## 🧹 Cleanup

```bash
# Stop all services
otto-stack down

# Remove containers and volumes
otto-stack down --volumes

# Clean up Docker resources
docker system prune
```

## 📚 Next Steps

- **[Services Guide](services.md)** - Available services and configuration
- **[CLI Reference](cli-reference.md)** - Complete command reference
- **[Configuration](configuration.md)** - Detailed configuration options
- **[Troubleshooting](troubleshooting.md)** - Common issues and solutions

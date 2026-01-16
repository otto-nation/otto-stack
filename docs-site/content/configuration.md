---
title: Configuration Guide
description: Complete guide to configuring otto-stack services
lead: Learn how to configure services for your specific needs
date: "2025-10-01"
lastmod: "2026-01-16"
draft: false
weight: 25
toc: true
---

# Configuration Guide

Otto-stack uses a single configuration file `otto-stack-config.yaml` to define your entire development stack.

## Configuration File Structure

```yaml
project:
  name: "my-app"

stack:
  enabled:
    - postgres
    - redis
    - kafka

service_configuration:
  postgres:
    database: "my_app_db"
    password: "secure_password"
  redis:
    password: "redis_password"
    max_memory: "512m"
```

## Configuration Sections

### Project Configuration

The `project` section defines basic project metadata:

- **name** (required): Your project name

### Stack Configuration

The `stack` section defines which services to enable:

- **enabled**: Array of service names to include in your stack

## Available Services

Otto-stack supports the following services:

- **jaeger** - Jaeger distributed tracing system for monitoring and troubleshooting microservices

- **kafka** - Complete Apache Kafka messaging platform with UI and topic management

- **kafka-broker** - Apache Kafka broker for event streaming and messaging

- **kafka-ui** - Web UI for Kafka cluster management and topic browsing

- **localstack-dynamodb** - LocalStack DynamoDB NoSQL database emulation

- **localstack-s3** - LocalStack S3 (Simple Storage Service) emulation

- **localstack-sns** - LocalStack SNS (Simple Notification Service) emulation

- **localstack-sqs** - LocalStack SQS (Simple Queue Service) emulation

- **mysql** - MySQL relational database for persistent data storage

- **postgres** - PostgreSQL relational database for persistent data storage

- **prometheus-service** - Prometheus metrics collection and monitoring system

- **redis** - Redis in-memory data store for caching and session storage

- **zookeeper** - Apache Zookeeper coordination service for distributed systems

For detailed information about each service, including configuration options and examples, see the [Services Guide](services.md).

## Service Configuration

Services can be configured in the `service_configuration` section. Each service has its own configuration schema with specific options for customization.

For complete service configuration details, examples, and available options, refer to the [Services Guide](services.md).

## Complete Example

Here's a complete configuration example with multiple services:

```yaml
project:
  name: my-fullstack-app
  type: web
stack:
  enabled:
    - postgres
    - redis
    - kafka
service_configuration:
  postgres:
    database: my_app_db
    password: secure_password
    port: 5432
  redis:
    password: redis_password
    max_memory: 512m
  kafka:
    topics:
      - events
      - notifications
```

For more configuration examples and service-specific options, see the [Services Guide](services.md).

---
title: "Configuration Guide"
description: "Complete guide to configuring otto-stack services"
lead: "Learn how to configure services for your specific needs"
date: "2025-10-01"
lastmod: "2025-11-10"
draft: "false"
weight: "25"
toc: "true"
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

- **jaeger** - Distributed tracing system
- **kafka** - Apache Kafka messaging platform
- **localstack-dynamodb** - DynamoDB emulation
- **localstack-s3** - S3 storage emulation
- **localstack-sns** - SNS notification service emulation
- **localstack-sqs** - SQS queue service emulation
- **mysql** - MySQL relational database
- **postgres** - PostgreSQL relational database
- **prometheus-service** - Metrics collection and monitoring
- **redis** - In-memory data store

For detailed information about each service, see the [Services Guide](services.md).

## Service Configuration Details

### jaeger

Jaeger distributed tracing system for monitoring and troubleshooting microservices

#### sampling

Jaeger sampling configuration

- Type: `object`

**Properties:**

- **default_strategy** (`string`) = `probabilistic`: Default sampling strategy
- **max_traces_per_second** (`integer`) = `100`: Maximum traces per second

#### storage

Storage configuration

- Type: `object`

**Properties:**

- **type** (`string`) = `memory`: Storage backend type

##### Example Configuration

```yaml
sampling:
  default_strategy: probabilistic
  max_traces_per_second: 100
storage:
  type: memory
```

### kafka

Complete Apache Kafka messaging platform with UI and topic management

#### topics

Kafka topics to create

- Type: `array`

**Items:**

- **name** (`string`) _required_: Topic name
- **partitions** (`integer`) = `3`: Number of partitions
- **replication_factor** (`integer`) = `1`: Replication factor

##### Example Configuration

```yaml
topics:
  - name: example-name
    partitions: 3
    replication_factor: 1
```

### localstack-dynamodb

LocalStack DynamoDB NoSQL database emulation

#### tables

DynamoDB tables to create

- Type: `array`

**Items:**

- **name** (`string`) _required_: Table name
- **hash_key** (`string`) _required_: Partition key
- **range_key** (`string`): Sort key
- **read_capacity** (`integer`) = `5`: Read capacity units
- **write_capacity** (`integer`) = `5`: Write capacity units

##### Example Configuration

```yaml
tables:
  - name: example-name
    hash_key: example-hash_key
    range_key: example-range_key
    read_capacity: 5
    write_capacity: 5
```

### localstack-s3

LocalStack S3 (Simple Storage Service) emulation

#### buckets

S3 buckets to create

- Type: `array`

**Items:**

- **name** (`string`) _required_: Bucket name
- **versioning** (`boolean`) = `false`: Enable versioning
- **public_read** (`boolean`) = `false`: Allow public read access

##### Example Configuration

```yaml
buckets:
  - name: example-name
    versioning: false
    public_read: false
```

### localstack-sns

LocalStack SNS (Simple Notification Service) emulation

#### topics

SNS topics to create

- Type: `array`

**Items:**

- **name** (`string`) _required_: Topic name
- **subscriptions** (`array`): Topic subscriptions

##### Example Configuration

```yaml
topics:
  - name: example-name
```

### localstack-sqs

LocalStack SQS (Simple Queue Service) emulation

#### queues

SQS queues to create

- Type: `array`

**Items:**

- **name** (`string`) _required_: Queue name
- **visibility_timeout** (`integer`) = `30`: Message visibility timeout in seconds
- **dead_letter_queue** (`string`): Dead letter queue name
- **max_receive_count** (`integer`) = `3`: Max receive count before moving to DLQ

##### Example Configuration

```yaml
queues:
  - name: example-name
    visibility_timeout: 30
    dead_letter_queue: example-dead_letter_queue
    max_receive_count: 3
```

### mysql

MySQL relational database for persistent data storage

#### database

Default database name

- Type: `string`
- Default: `local_dev`

#### password

Root password

- Type: `string`
- Default: `password`

#### user

Database user

- Type: `string`
- Default: `root`

##### Example Configuration

```yaml
database: local_dev
password: password
user: root
```

### postgres

PostgreSQL relational database for persistent data storage

#### database

Default database name

- Type: `string`
- Default: `local_dev`

#### password

Database password

- Type: `string`
- Default: `password`

#### user

Database user

- Type: `string`
- Default: `postgres`

##### Example Configuration

```yaml
database: local_dev
password: password
user: postgres
```

### prometheus-service

Prometheus metrics collection and monitoring system

#### scrape_configs

Prometheus scrape configurations

- Type: `array`

**Items:**

- **job_name** (`string`) _required_: Job name
- **static_configs** (`array`): Static target configurations
- **scrape_interval** (`string`) = `15s`: Scrape interval

##### Example Configuration

```yaml
scrape_configs:
  - job_name: example-job_name
    scrape_interval: 15s
```

### redis

Redis in-memory data store for caching and session storage

#### password

Redis password

- Type: `string`
- Default: `password`

#### max_memory

Maximum memory limit

- Type: `string`
- Default: `256m`

#### databases

Number of databases

- Type: `integer`
- Default: `16`

##### Example Configuration

```yaml
password: password
max_memory: 256m
databases: 16
```

## Complete Example

```yaml
project:
  name: my-app
stack:
  enabled:
    - postgres
    - redis
    - kafka
    - localstack-sqs
service_configuration:
  redis:
    password: password
    max_memory: 256m
    databases: 16
  localstack-dynamodb:
    tables:
      - name: example-name
        hash_key: example-hash_key
        range_key: example-range_key
        read_capacity: 5
        write_capacity: 5
  localstack-s3:
    buckets:
      - name: example-name
        versioning: false
        public_read: false
  localstack-sns:
    topics:
      - name: example-name
  localstack-sqs:
    queues:
      - name: example-name
        visibility_timeout: 30
        dead_letter_queue: example-dead_letter_queue
        max_receive_count: 3
  mysql:
    database: local_dev
    password: password
    user: root
  postgres:
    database: local_dev
    password: password
    user: postgres
  kafka:
    topics:
      - name: example-name
        partitions: 3
        replication_factor: 1
  jaeger:
    sampling:
      default_strategy: probabilistic
      max_traces_per_second: 100
    storage:
      type: memory
  prometheus-service:
    scrape_configs:
      - job_name: example-job_name
        scrape_interval: 15s
```

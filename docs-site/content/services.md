---
title: Services
description: Available services and configuration options
lead: Explore all the services you can use with otto-stack
date: "2025-10-01"
lastmod: "2026-01-21"
draft: false
weight: 30
toc: true
---

# Available Services

13 services available for your development stack.

Each service can be configured through the `service_configuration` section in your `otto-stack-config.yaml` file. For detailed configuration instructions, see the [Configuration Guide](configuration.md).

## 🗄️ Database

### mysql

MySQL relational database for persistent data storage

#### Configuration Options

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

#### databases

Additional databases to create

- Type: `array`

**Items:**

- **name** (`string`)

#### users

Additional users to create

- Type: `array`

**Items:**

- **name** (`string`)

- **password** (`string`)

- **database** (`string`)

##### Example Configuration

```yaml
database: local_dev
password: password
user: root
databases:
  - name: example-name
users:
  - name: example-name
    password: example-password
    database: example-database
```

#### Use Cases

- Primary application database

- Transactional data storage

- Relational data modeling

---

### postgres

PostgreSQL relational database for persistent data storage

#### Configuration Options

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

#### Use Cases

- Primary application database

- Transactional data storage

- Relational data modeling

- ACID compliance requirements

---

## ⚡ Cache

### redis

Redis in-memory data store for caching and session storage

#### Configuration Options

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

#### Use Cases

- Session storage

- Application caching

- Rate limiting

- Pub/Sub messaging

- Temporary data storage

---

## 📨 Messaging

### kafka

Complete Apache Kafka messaging platform with UI and topic management

#### Configuration Options

#### topics

Kafka topics to create

- Type: `array`

**Items:**

- **name** (`string`): Topic name

- **partitions** (`integer`) = `3`: Number of partitions

- **replication_factor** (`integer`) = `1`: Replication factor

##### Example Configuration

```yaml
topics:
  - name: example-name
    partitions: 3
    replication_factor: 1
```

#### Use Cases

- Event streaming and processing

- Message queuing and pub/sub

- Real-time data pipelines

- Microservices communication

- Log aggregation and analytics

---

### kafka-broker

Apache Kafka broker for event streaming and messaging

#### Configuration Options

#### topics

Kafka topics to create on startup

- Type: `array`

**Items:**

- **name** (`string`): Topic name

- **partitions** (`integer`) = `3`: Number of partitions

- **replication_factor** (`integer`) = `1`: Replication factor

##### Example Configuration

```yaml
topics:
  - name: example-name
    partitions: 3
    replication_factor: 1
```

#### Use Cases

- Event-driven microservices architecture

- Real-time data streaming and processing

- Message queuing between services

- Data pipeline and ETL processes

---

### kafka-ui

Web UI for Kafka cluster management and topic browsing

#### Use Cases

- Topic management and browsing

- Message inspection and debugging

- Cluster monitoring and health checks

- Consumer group management

---

### zookeeper

Apache Zookeeper coordination service for distributed systems

#### Use Cases

- Kafka cluster coordination

- Distributed configuration management

- Service discovery

- Leader election

---

## ☁️ Cloud

### localstack-dynamodb

LocalStack DynamoDB NoSQL database emulation

#### Configuration Options

#### tables

DynamoDB tables to create

- Type: `array`

**Items:**

- **name** (`string`): Table name

- **hash_key** (`string`): Partition key

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

#### Use Cases

- NoSQL database testing

- High-performance data storage

- Session management

- Event sourcing patterns

---

### localstack-s3

LocalStack S3 (Simple Storage Service) emulation

#### Configuration Options

#### buckets

S3 buckets to create

- Type: `array`

**Items:**

- **name** (`string`): Bucket name

- **versioning** (`boolean`): Enable versioning

- **public_read** (`boolean`): Allow public read access

##### Example Configuration

```yaml
buckets:
  - name: example-name
    versioning: false
    public_read: false
```

#### Use Cases

- File storage and retrieval testing

- Static asset hosting

- Data backup and archival

- Content distribution testing

---

### localstack-sns

LocalStack SNS (Simple Notification Service) emulation

#### Configuration Options

#### topics

SNS topics to create

- Type: `array`

**Items:**

- **name** (`string`): Topic name

- **subscriptions** (`array`): Topic subscriptions

##### Example Configuration

```yaml
topics:
  - name: example-name
```

#### Use Cases

- Pub/sub messaging patterns

- Event notifications

- Fan-out message distribution

- Integration with SQS subscriptions

---

### localstack-sqs

LocalStack SQS (Simple Queue Service) emulation

#### Configuration Options

#### queues

SQS queues to create

- Type: `array`

**Items:**

- **name** (`string`): Queue name

- **visibility_timeout** (`integer`) = `30`: Message visibility timeout in seconds

- **dead_letter_queue** (`string`): Dead letter queue name

- **max_receive_count** (`integer`) = `3`: Max receive count before moving to DLQ

- **subscription** (`object`) _required_

##### Example Configuration

```yaml
queues:
  - name: example-name
    visibility_timeout: 30
    dead_letter_queue: example-dead_letter_queue
    max_receive_count: 3
```

#### Use Cases

- Message queue testing

- Asynchronous processing development

- Event-driven architecture testing

---

## 🔍 Observability

### jaeger

Jaeger distributed tracing system for monitoring and troubleshooting microservices

#### Configuration Options

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

#### Use Cases

- Distributed tracing

- Performance monitoring

- Service dependency analysis

- Request flow visualization

---

### prometheus-service

Prometheus metrics collection and monitoring system

#### Configuration Options

#### scrape_configs

Prometheus scrape configurations

- Type: `array`

**Items:**

- **job_name** (`string`): Job name

- **static_configs** (`array`): Static target configurations

- **scrape_interval** (`string`) = `15s`: Scrape interval

##### Example Configuration

```yaml
scrape_configs:
  - job_name: example-job_name
    scrape_interval: 15s
```

#### Use Cases

- Application metrics collection

- Performance monitoring

- Alerting and notifications

---

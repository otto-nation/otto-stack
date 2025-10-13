---
title: "Service Profiles"
description: "Predefined service combinations for common development scenarios"
lead: "Quickly start with predefined service combinations"
date: "2025-10-01"
lastmod: "2025-10-11"
draft: false
weight: 35
toc: true
---

<!-- AUTO-GENERATED-START -->

# Service Profiles

Service profiles are predefined combinations of services for common development scenarios. Use them to quickly start your development environment with the right services.

## Using Profiles

```bash
# Start services using a profile
otto-stack up --profile <profile-name>

# List available profiles
otto-stack up --profile <TAB>
```

## Available Profiles

### API Development

Services for API development and testing

**Services included:**

- postgres

- redis

- prometheus

**Quick start:**

```bash
otto-stack up --profile api development
```

---

### Data Engineering

Services for data processing and analytics

**Services included:**

- postgres

- redis

- kafka

- localstack

**Quick start:**

```bash
otto-stack up --profile data engineering
```

---

### Microservices

Full microservices development stack

**Services included:**

- postgres

- redis

- kafka

- jaeger

- prometheus

**Quick start:**

```bash
otto-stack up --profile microservices
```

---

### Minimal Stack

Minimal services for basic development

**Services included:**

- postgres

**Quick start:**

```bash
otto-stack up --profile minimal stack
```

---

### Web Development

Services for web application development

**Services included:**

- postgres

- redis

- jaeger

**Quick start:**

```bash
otto-stack up --profile web development
```

---

## Creating Custom Profiles

You can define custom profiles in your `otto-stack-config.yaml` file:

```yaml
profiles:
  my-profile:
    name: "My Custom Profile"
    description: "Custom services for my project"
    services:
      - postgres
      - redis
      - my-service
```

<!-- AUTO-GENERATED-END -->

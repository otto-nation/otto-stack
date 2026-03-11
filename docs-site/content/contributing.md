---
title: Contributing
description: Guide for contributing to otto-stack development
lead: Learn how to contribute to otto-stack development
date: "2025-10-01"
lastmod: "2026-03-11"
draft: false
weight: 60
toc: true
---

<!--
  ⚠️  AUTO-GENERATED FILE - DO NOT EDIT DIRECTLY
  To make changes, edit source files and run: task generate:docs
-->

# Contributing to otto-stack

Thank you for contributing! This guide covers the essentials to get you started.

## Prerequisites

| Tool | Version | Purpose |
| --- | --- | --- |
| [Go](https://go.dev/dl/) | 1.25.7 (see `.go-version`) | Build and run the project |
| [Task](https://taskfile.dev/installation/) | any | Build automation |
| [Docker](https://docs.docker.com/get-started/get-docker/) | 24+ | Service management and E2E tests |
| [Node.js](https://nodejs.org/) | >=18 | Documentation generation |
| [Hugo](https://gohugo.io/installation/) | >=0.148.1 extended | Documentation site |

**Go — install via [goenv](https://github.com/go-nv/goenv) (recommended):**

```bash
# macOS
brew install goenv

# then task setup will install and activate the correct version automatically
```

Or download directly from [go.dev/dl](https://go.dev/dl/) if you don't use a version manager.

**Task — [taskfile.dev](https://taskfile.dev/installation/):**

```bash
# macOS
brew install go-task

# Linux / Windows: see https://taskfile.dev/installation/
```

> **Auto-installed by task:** `golangci-lint` and `goimports` are installed automatically
> when you run `task lint` or `task fmt` — no manual installation needed.

## Quick Start

```bash
# 1. Fork and clone
git clone https://github.com/your-username/otto-stack.git
cd otto-stack

# 2. Check prerequisites, install dependencies, and configure Git hooks
task setup

# 3. Build
task build

# 4. Run tests
task test

# 5. Verify
./build/otto-stack version
```

## Making Changes

### Code Changes

1. Create a feature branch
2. Make your changes
3. Run `task test` and `task lint`
4. Submit a PR with clear description

### Adding Services

1. Create YAML file in `internal/config/services/{category}/`
2. Follow existing service structure (see `internal/config/services/database/postgres.yaml`)
3. Run `task docs` to regenerate documentation
4. Test with `task docs-serve`

**Service categories:**
- `database/` - Databases (postgres, mysql)
- `cache/` - Caching (redis)
- `messaging/` - Message queues (kafka)
- `observability/` - Monitoring (prometheus, jaeger)
- `cloud/` - Cloud emulation (localstack)

### Documentation Changes

**Auto-generated** (regenerated from source, don't edit):
- `docs-site/content/cli-reference.md` - From CLI commands
- `docs-site/content/services.md` - From service YAML files
- `docs-site/content/configuration.md` - From service schemas
- `docs-site/content/_index.md` - From root README.md
- `docs-site/content/contributing.md` - From root CONTRIBUTING.md

**Manual** (edit directly):
- `docs-site/content/setup.md` - Installation instructions
- `docs-site/content/troubleshooting.md` - Common issues
- Other custom documentation pages

**To regenerate documentation:**
```bash
# Full regeneration
task docs

# Skip CLI build (faster for quick iterations)
cd docs-site && npm run generate -- --skip-build

# Generate specific doc only
cd docs-site && npm run generate -- --generator=services-guide

# Preview changes
task docs-serve
```

## Pull Request Process

1. Fork and create a feature branch
2. Make changes following project conventions
3. Write/update tests as needed
4. Run `task test` and `task lint`
5. Submit PR with clear description

### Commit Message Format

Follow conventional commits (see `docs/RELEASING.md` for details):

```
type(scope): description
```

**Types**: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`

**Examples:**
- `feat(services): add MySQL service`
- `fix(cli): resolve network creation issue`
- `docs: update setup instructions`

## Getting Help

- **Issues**: [Report bugs/request features](https://github.com/otto-nation/otto-stack/issues)
- **Documentation**: Check `docs-site/` or run `otto-stack help`
- **Troubleshooting**: Run `otto-stack doctor` for diagnostics

## Code Review

We review for:
- Code quality and tests
- Documentation updates
- Breaking changes
- Security considerations

---

We appreciate your contributions! 🙏

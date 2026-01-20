---
title: Contributing
description: Guide for contributing to otto-stack development
lead: Learn how to contribute to otto-stack development
date: "2025-10-01"
lastmod: "2026-01-20"
draft: false
weight: 60
toc: true
---

# Contributing to otto-stack

Thank you for contributing! This guide covers the essentials to get you started.

## Quick Start

```bash
# 1. Fork and clone
git clone https://github.com/your-username/otto-stack.git
cd otto-stack

# 2. Setup and build
task setup
task build

# 3. Run tests
task test

# 4. Verify
./build/otto-stack version
```

## Prerequisites

- **Go**: Check `.go-version` for the required version.
- **Task**: [Task runner](https://taskfile.dev/)
- **Docker**: For service management
- **Node.js**: For documentation generation

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

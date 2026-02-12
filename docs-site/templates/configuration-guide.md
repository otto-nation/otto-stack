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
{{{fileStructure}}}
```

## Main Configuration

**`.otto-stack/config.yaml`:**

```yaml
{{{configStructure}}}
```

{{{configSections}}}

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
# {{{serviceConfigExample}}}
```

### Customizing Services

Create a `.env` file in your project root to override defaults:

```bash
{{{customEnvExample}}}
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
{{{completeExample}}}
```

**`.env` (your customizations):**

```bash
{{{completeEnvExample}}}
```

## Next Steps

- **[Services Guide](/otto-stack/services/)** - Available services and environment variables
- **[CLI Reference](/otto-stack/cli-reference/)** - Command usage
- **[Troubleshooting](/otto-stack/troubleshooting/)** - Common issues

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

- **[Services Guide](services.md)** - Available services and environment variables
- **[CLI Reference](cli-reference.md)** - Command usage
- **[Troubleshooting](troubleshooting.md)** - Common issues

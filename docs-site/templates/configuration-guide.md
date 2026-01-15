# Configuration Guide

Otto-stack uses a single configuration file `otto-stack-config.yaml` to define your entire development stack.

## Configuration File Structure

```yaml
{{{structureExample}}}
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

{{#each serviceList}}

- {{this}}
  {{/each}}

For detailed information about each service, including configuration options and examples, see the [Services Guide](services.md).

## Service Configuration

Services can be configured in the `service_configuration` section. Each service has its own configuration schema with specific options for customization.

For complete service configuration details, examples, and available options, refer to the [Services Guide](services.md).

## Complete Example

Here's a complete configuration example with multiple services:

```yaml
{{{completeExample}}}
```

For more configuration examples and service-specific options, see the [Services Guide](services.md).

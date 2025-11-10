# Configuration Guide

Otto-stack uses a single configuration file `otto-stack-config.yaml` to define your entire development stack.

## Configuration File Structure

```yaml
{ { structureExample } }
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

For detailed information about each service, see the [Services Guide](services.md).

## Service Configuration Details

{{#each configurableServices}}

### {{name}}

{{#if description}}
{{description}}
{{/if}}

{{#each configSections}}

#### {{name}}

{{#if description}}
{{description}}
{{/if}}

- Type: `{{type}}`

{{#if properties}}
**Properties:**

{{#each properties}}

- **{{name}}** (`{{type}}`){{#if required}} _required_{{/if}}{{#if default}} = `{{default}}`{{/if}}: {{description}}
  {{/each}}
  {{/if}}

{{#if example}}

##### Example Configuration

```yaml
{ { example } }
```

{{/if}}

{{/each}}
{{/each}}

## Complete Example

```yaml
{ { completeExample } }
```

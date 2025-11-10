### {{name}}

{{description}}

{{#each details}}
**{{label}}:** {{value}}

{{/each}}

{{#if configurationSchema}}

#### Configuration Options

{{#each configurationSchema.fields}}

#### {{name}}

{{#if description}}
{{description}}

{{/if}}

- Type: `{{type}}`
  {{#if default}}
- Default: `{{default}}`
  {{/if}}
  {{#if required}}
- Required: Yes
  {{/if}}

{{#if items}}
**Items:**
{{#each items.properties}}

- **{{name}}** (`{{type}}`){{#if required}} _required_{{/if}}{{#if default}} = `{{default}}`{{/if}}{{#if description}}: {{description}}{{/if}}
  {{/each}}
  {{/if}}

{{#if properties}}
**Properties:**
{{#each properties}}

- **{{name}}** (`{{type}}`){{#if default}} = `{{default}}`{{/if}}{{#if description}}: {{description}}{{/if}}
  {{/each}}
  {{/if}}

{{/each}}

{{#if configurationSchema.examples}}

##### Example Configuration

```yaml
{{{toYaml configurationSchema.examples}}}
```

{{/if}}
{{/if}}

{{#if examples}}

#### Examples

{{#each examples}}

```bash
{{{this}}}
```

{{/each}}
{{/if}}

{{#if useCases}}

#### Use Cases

{{#each useCases}}

- {{this}}
  {{/each}}

{{/if}}

---

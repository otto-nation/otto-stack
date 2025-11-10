const yaml = require("js-yaml");

class SchemaParser {
  transformSchema(rawSchema) {
    if (!rawSchema) return null;

    return {
      fields: Object.entries(rawSchema).map(([key, config]) => ({
        name: key,
        type: config.type,
        description: config.description,
        required: config.required || false,
        default: config.default,
        items: config.items ? this.transformItems(config.items) : null,
        properties: config.properties
          ? this.transformProperties(config.properties)
          : null,
      })),
      examples: this.generateSchemaExamples(rawSchema),
    };
  }

  transformItems(items) {
    return {
      type: items.type,
      properties: items.properties
        ? this.transformProperties(items.properties)
        : null,
    };
  }

  transformProperties(properties) {
    return Object.entries(properties).map(([key, config]) => ({
      name: key,
      type: config.type,
      description: config.description,
      required: config.required || false,
      default: config.default,
    }));
  }

  generateSchemaExamples(schema) {
    const examples = {};

    Object.entries(schema).forEach(([key, config]) => {
      if (config.type === "string" && config.default) {
        examples[key] = config.default;
      } else if (config.type === "integer" && config.default !== undefined) {
        examples[key] = config.default;
      } else if (config.type === "boolean" && config.default !== undefined) {
        examples[key] = config.default;
      } else if (config.type === "array" && config.items) {
        examples[key] = [this.generateItemExample(config.items)];
      } else if (config.type === "object" && config.properties) {
        examples[key] = this.generateObjectExample(config.properties);
      }
    });

    return Object.keys(examples).length > 0 ? examples : null;
  }

  generateItemExample(items) {
    if (!items.properties) return {};

    const example = {};
    Object.entries(items.properties).forEach(([key, config]) => {
      if (config.default !== undefined) {
        example[key] = config.default;
      } else if (config.type === "string") {
        example[key] = `example-${key}`;
      } else if (config.type === "integer") {
        example[key] = 1;
      } else if (config.type === "boolean") {
        example[key] = true;
      }
    });

    return example;
  }

  generateObjectExample(properties) {
    const example = {};
    Object.entries(properties).forEach(([key, config]) => {
      if (config.default !== undefined) {
        example[key] = config.default;
      } else if (config.type === "string") {
        example[key] = `example-${key}`;
      } else if (config.type === "integer") {
        example[key] = 1;
      } else if (config.type === "boolean") {
        example[key] = true;
      }
    });
    return example;
  }

  toYaml(obj) {
    return yaml.dump(obj, {
      indent: 2,
      lineWidth: -1,
      noRefs: true,
      sortKeys: false,
    });
  }
}

module.exports = SchemaParser;

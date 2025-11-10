const ServiceAnalyzer = require("../utils/service-analyzer");
const SchemaParser = require("../utils/schema-parser");
const TemplateRenderer = require("../utils/template-renderer");

class ConfigurationGuideGenerator {
  constructor(config) {
    this.config = config;
    this.analyzer = new ServiceAnalyzer(config);
    this.schemaParser = new SchemaParser();
    this.templateRenderer = new TemplateRenderer();
  }

  generate() {
    console.log("Generating configuration guide...");

    const services = this.analyzer.loadAllServices();
    const configurableServices = this.getConfigurableServices(services);

    const templateData = {
      structureExample: this.generateStructureExample(),
      serviceList: this.generateServiceList(services),
      configurableServices:
        this.processConfigurableServices(configurableServices),
      completeExample: this.generateCompleteExample(configurableServices),
    };

    const today = new Date().toISOString().split("T")[0];
    const frontmatter = {
      title: "Configuration Guide",
      description: "Complete guide to configuring otto-stack services",
      lead: "Learn how to configure services for your specific needs",
      date: "2025-10-01",
      lastmod: today,
      draft: false,
      weight: 25,
      toc: true,
    };

    return this.templateRenderer.render(
      "configuration-guide.md",
      templateData,
      frontmatter,
    );
  }

  getConfigurableServices(services) {
    const configurable = {};
    Object.entries(services).forEach(([name, config]) => {
      if (config.configuration_schema) {
        configurable[name] = config;
      }
    });
    return configurable;
  }

  processConfigurableServices(configurableServices) {
    return Object.entries(configurableServices)
      .sort(([a], [b]) => a.localeCompare(b))
      .map(([name, config]) => ({
        name,
        description: config.description,
        configSections: this.processConfigSections(config.configuration_schema),
      }));
  }

  processConfigSections(schema) {
    if (!schema || !schema.properties) return [];

    return Object.entries(schema.properties).map(([key, prop]) => ({
      name: key,
      description: prop.description,
      type: prop.type,
      properties: this.processProperties(prop),
      example: this.generateSectionExample(key, prop),
    }));
  }

  processProperties(prop) {
    if (prop.type === "object" && prop.properties) {
      return Object.entries(prop.properties).map(([key, subProp]) => ({
        name: key,
        type: subProp.type,
        required: prop.required?.includes(key),
        default: subProp.default,
        description: subProp.description,
      }));
    }
    if (prop.type === "array" && prop.items?.properties) {
      return Object.entries(prop.items.properties).map(([key, subProp]) => ({
        name: key,
        type: subProp.type,
        required: prop.items.required?.includes(key),
        default: subProp.default,
        description: subProp.description,
      }));
    }
    return null;
  }

  generateSectionExample(key, prop) {
    return this.schemaParser.generateExample({ [key]: prop });
  }

  generateStructureExample() {
    return `project:
  name: "my-app"

stack:
  enabled:
    - postgres
    - redis
    - kafka

service_configuration:
  postgres:
    database: "my_app_db"
    password: "secure_password"
  redis:
    password: "redis_password"
    max_memory: "512m"`.trim();
  }

  generateServiceList(services) {
    return Object.keys(services)
      .sort()
      .map((name) => {
        const config = services[name];
        const description = config.description || config.name || "Service";
        return `**${name}** - ${description}`;
      });
  }

  generateBasicExample() {
    return `
project:
  name: "my-app"

stack:
  enabled:
    - postgres
    - redis

service_configuration:
  postgres:
    database: "my_app_db"
    password: "secure_password"
  redis:
    password: "redis_password"
    max_memory: "512m"`;
  }

  generateCompleteExample(configurableServices) {
    const examples = {};

    Object.entries(configurableServices).forEach(([name, config]) => {
      const schema = this.schemaParser.transformSchema(
        config.configuration_schema,
      );
      if (schema?.examples) {
        examples[name] = schema.examples;
      }
    });

    const exampleConfig = {
      project: {
        name: "my-app",
      },
      stack: {
        enabled: ["postgres", "redis", "kafka", "localstack-sqs"],
      },
      service_configuration: examples,
    };

    return this.schemaParser.toYaml(exampleConfig);
  }
}

module.exports = ConfigurationGuideGenerator;

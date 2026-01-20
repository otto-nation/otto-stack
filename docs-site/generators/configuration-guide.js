const fs = require("fs");
const path = require("path");
const yaml = require("js-yaml");
const ServiceAnalyzer = require("../utils/service-analyzer");
const TemplateRenderer = require("../utils/template-renderer");

class ConfigurationGuideGenerator {
  constructor(config) {
    this.config = config;
    this.analyzer = new ServiceAnalyzer(config);
    this.templateRenderer = new TemplateRenderer();
    this.schemaPath = path.join(__dirname, "../../internal/config/schema.yaml");
    this.exampleServices = ["postgres", "redis"];
    this.completeExampleServices = ["postgres", "redis", "kafka"];
  }

  generate() {
    console.log("Generating configuration guide...");

    const schema = this.loadSchema();
    const services = this.analyzer.loadAllServices();

    const templateData = {
      fileStructure: this.generateFileStructure(),
      configStructure: this.generateConfigStructure(schema),
      configSections: this.generateConfigSections(schema),
      serviceConfigExample: this.generateServiceConfigExample(services),
      customEnvExample: this.generateCustomEnvExample(services),
      completeExample: this.generateCompleteExample(schema),
      completeEnvExample: this.generateCompleteEnvExample(services),
    };

    const today = new Date().toISOString().split("T")[0];
    const frontmatter = {
      title: "Configuration Guide",
      description: "Configure your otto-stack development environment",
      lead: "Learn how to configure your development stack",
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

  loadSchema() {
    const content = fs.readFileSync(this.schemaPath, "utf8");
    return yaml.load(content);
  }

  generateFileStructure() {
    return `.otto-stack/
├── config.yaml              # Main configuration
├── generated/
│   ├── .env.generated       # Available environment variables
│   └── docker-compose.yml   # Generated Docker Compose
├── services/                # Service metadata
│   ├── postgres.yml
│   └── redis.yml
├── .gitignore
└── README.md`;
  }

  generateConfigStructure(schema) {
    const example = {};
    const schemaObj = schema.schema;

    // Build example from schema properties
    for (const [section, config] of Object.entries(schemaObj)) {
      if (!config.properties) continue;

      example[section] = {};
      for (const [key, prop] of Object.entries(config.properties)) {
        if (
          prop.default &&
          typeof prop.default === "string" &&
          !prop.default.startsWith("{{")
        ) {
          example[section][key] = prop.default;
        } else if (prop.type === "array") {
          example[section][key] =
            section === "stack" ? this.exampleServices : [];
        } else if (prop.type === "boolean") {
          example[section][key] = false;
        } else if (prop.type === "string") {
          example[section][key] =
            section === "project" && key === "name" ? "my-app" : "";
        }
      }
    }

    return yaml.dump(example, { indent: 2, lineWidth: 80 });
  }

  generateProjectSection(schema) {
    const projectSchema = schema.schema.project;
    if (!projectSchema) return "";

    const props = projectSchema.properties || {};
    return Object.entries(props)
      .map(([key, value]) => `- **${key}**: ${value.description || ""}`)
      .join("\n");
  }

  generateStackSection(schema) {
    const stackSchema = schema.schema.stack;
    if (!stackSchema) return "";

    const props = stackSchema.properties || {};
    return Object.entries(props)
      .map(([key, value]) => `- **${key}**: ${value.description || ""}`)
      .join("\n");
  }

  generateValidationSection(schema) {
    const validationSchema = schema.schema.validation;
    if (!validationSchema) return "";

    const props = validationSchema.properties || {};
    return Object.entries(props)
      .map(([key, value]) => `- **${key}**: ${value.description || ""}`)
      .join("\n");
  }

  generateConfigSections(schema) {
    const sections = [];
    const schemaObj = schema.schema;

    for (const [sectionName, config] of Object.entries(schemaObj)) {
      if (!config.properties) continue;

      const title = sectionName
        .split("_")
        .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
        .join(" ");

      const props = Object.entries(config.properties)
        .map(([key, value]) => `- **${key}**: ${value.description || ""}`)
        .join("\n");

      sections.push(`### ${title}\n\n${config.description || ""}\n\n${props}`);
    }

    return sections.join("\n\n");
  }

  generateServiceConfigExample(services) {
    // Get actual environment variables from service YAMLs
    const examples = [];

    this.completeExampleServices.forEach((serviceName) => {
      if (services[serviceName]) {
        const envVars = services[serviceName].environment || {};
        const envKeys = Object.keys(envVars);

        if (envKeys.length > 0) {
          examples.push(`# ${serviceName.toUpperCase()}`);
          // Show first 3-4 most relevant env vars
          envKeys.slice(0, 4).forEach((key) => {
            examples.push(`${key}=${envVars[key]}`);
          });
        }
      }
    });

    return examples.join("\n");
  }

  generateCompleteExample(schema) {
    const example = {
      project: {
        name: "my-fullstack-app",
        type: "docker",
      },
      stack: {
        enabled: this.completeExampleServices,
      },
    };

    if (schema.schema.validation) {
      example.validation = {
        options: {
          "config-syntax": true,
          docker: true,
        },
      };
    }

    return yaml.dump(example, { indent: 2, lineWidth: 80 });
  }

  generateCustomEnvExample(services) {
    const examples = [];

    this.exampleServices.forEach((serviceName) => {
      if (services[serviceName]) {
        const envVars = services[serviceName].environment || {};
        const envKeys = Object.keys(envVars);

        if (envKeys.length > 0) {
          examples.push(
            `# ${serviceName.charAt(0).toUpperCase() + serviceName.slice(1)}`,
          );
          envKeys.slice(0, 2).forEach((key) => {
            examples.push(`${key}=my_custom_value`);
          });
        }
      }
    });

    return examples.join("\n");
  }

  generateCompleteEnvExample(services) {
    const examples = [];

    this.exampleServices.forEach((serviceName) => {
      if (services[serviceName]) {
        const envVars = services[serviceName].environment || {};
        const envKeys = Object.keys(envVars);

        if (envKeys.length > 0) {
          examples.push(
            `# ${serviceName.charAt(0).toUpperCase() + serviceName.slice(1)}`,
          );
          envKeys.slice(0, 2).forEach((key) => {
            examples.push(`${key}=production_value`);
          });
        }
      }
    });

    return examples.join("\n");
  }
}

module.exports = ConfigurationGuideGenerator;

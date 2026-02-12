const fs = require("fs");
const path = require("path");
const yaml = require("js-yaml");
const ServiceAnalyzer = require("../utils/service-analyzer");
const TemplateRenderer = require("../utils/template-renderer");
const BaseGenerator = require("./base-generator");

const SCHEMA_PATH = path.join(__dirname, "../../internal/config/schema.yaml");
const EXAMPLE_SERVICES = ["postgres", "redis"];
const COMPLETE_EXAMPLE_SERVICES = ["postgres", "redis", "kafka"];

class ConfigurationGuideGenerator extends BaseGenerator {
  constructor(config) {
    super(config);
    this.analyzer = new ServiceAnalyzer(config);
    this.templateRenderer = new TemplateRenderer();
  }

  generate() {
    console.log("Generating configuration guide...");

    try {
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

      const frontmatter = this.createFrontmatter(
        "Configuration Guide",
        "Configure your otto-stack development environment",
        "Learn how to configure your development stack",
        25,
      );

      return this.templateRenderer.render(
        "configuration-guide.md",
        templateData,
        frontmatter,
      );
    } catch (error) {
      this.handleError("generate configuration guide", error);
    }
  }

  loadSchema() {
    try {
      const content = fs.readFileSync(SCHEMA_PATH, "utf8");
      return yaml.load(content);
    } catch (error) {
      this.handleError("load schema", error);
    }
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
          example[section][key] = section === "stack" ? EXAMPLE_SERVICES : [];
        } else if (prop.type === "boolean") {
          example[section][key] = false;
        } else if (prop.type === "string") {
          example[section][key] =
            section === "project" && key === "name" ? "my-app" : "";
        }
      }
    }

    return this.dumpYaml(example);
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
    const examples = [];

    COMPLETE_EXAMPLE_SERVICES.forEach((serviceName) => {
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
        enabled: COMPLETE_EXAMPLE_SERVICES,
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

    return this.dumpYaml(example);
  }

  generateCustomEnvExample(services) {
    const examples = [];

    EXAMPLE_SERVICES.forEach((serviceName) => {
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

    EXAMPLE_SERVICES.forEach((serviceName) => {
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

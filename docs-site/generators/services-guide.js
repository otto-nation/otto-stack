const ServiceAnalyzer = require("../utils/service-analyzer");
const SchemaParser = require("../utils/schema-parser");
const TemplateRenderer = require("../utils/template-renderer");
const BaseGenerator = require("./base-generator");

class ServicesGuideGenerator extends BaseGenerator {
  constructor(config) {
    super(config);
    this.analyzer = new ServiceAnalyzer(config);
    this.schemaParser = new SchemaParser();
    this.templateRenderer = new TemplateRenderer();
  }

  generate() {
    console.log("Generating services guide...");

    try {
      const services = this.analyzer.loadAllServices();
      const categories = this.analyzer.categorizeServices(services);

      const frontmatter = this.createFrontmatter(
        "Services",
        "Available services and configuration options",
        "Explore all the services you can use with otto-stack",
        30,
      );

      const content = this.generateContent(services, categories);
      return this.formatDocument(frontmatter, content);
    } catch (error) {
      this.handleError("generate services guide", error);
    }
  }

  generateContent(services, categories) {
    let content = `# Available Services

${Object.keys(services).length} services available for your development stack.

Each service can be configured through the \`service_configuration\` section in your \`otto-stack-config.yaml\` file. For detailed configuration instructions, see the [Configuration Guide](/otto-stack/configuration/).

`;

    this.getSortedCategories(categories).forEach(
      ([categoryName, categoryServices]) => {
        const categoryConfig = this.analyzer.getCategoryConfig(categoryName);
        const categoryTitle = `${categoryConfig.icon} ${categoryName.charAt(0).toUpperCase() + categoryName.slice(1)}`;

        content += `## ${categoryTitle}\n\n`;

        categoryServices
          .sort(([a], [b]) => a.localeCompare(b))
          .forEach(([name, config]) => {
            const serviceData = this.processServiceForTemplate(name, config);
            content += this.renderServiceSection(serviceData);
          });
      },
    );

    return content;
  }

  renderServiceSection(serviceData) {
    // Render the service template with proper Handlebars processing
    return this.templateRenderer.render("service.md", serviceData);
  }

  processServiceForTemplate(name, config) {
    const details = [];

    if (config.ports?.length > 0) {
      const portStrings = config.ports.map((p) => `${p.host}:${p.container}`);
      details.push({ label: "Ports", value: portStrings.join(", ") });
    }

    if (config.web_interface) {
      const webInterface = `[${config.web_interface.name}](${config.web_interface.url})`;
      details.push({ label: "Web Interface", value: webInterface });
    }

    if (config.provides?.length > 0) {
      details.push({ label: "Provides", value: config.provides.join(", ") });
    }

    if (config.requires?.length > 0) {
      details.push({ label: "Requires", value: config.requires.join(", ") });
    }

    let configurationSchema = null;
    if (config.configuration_schema) {
      const transformedSchema = this.schemaParser.transformSchema(
        config.configuration_schema,
      );
      if (transformedSchema) {
        configurationSchema = {
          fields: transformedSchema.fields,
          examples: transformedSchema.examples,
        };
      }
    }

    // Add examples and use cases from documentation
    const examples = config.documentation?.examples || [];
    const useCases = config.documentation?.use_cases || [];

    return {
      name,
      description: config.description,
      details,
      configurationSchema,
      examples,
      useCases,
    };
  }

  getSortedCategories(categories) {
    return Object.entries(categories).sort(([a], [b]) => {
      const orderA = this.analyzer.getCategoryConfig(a).order;
      const orderB = this.analyzer.getCategoryConfig(b).order;
      return orderA - orderB;
    });
  }
}

module.exports = ServicesGuideGenerator;

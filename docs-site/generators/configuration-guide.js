const ServiceAnalyzer = require("../utils/service-analyzer");
const TemplateRenderer = require("../utils/template-renderer");

class ConfigurationGuideGenerator {
  constructor(config) {
    this.config = config;
    this.analyzer = new ServiceAnalyzer(config);
    this.templateRenderer = new TemplateRenderer();
  }

  generate() {
    console.log("Generating configuration guide...");

    const services = this.analyzer.loadAllServices();

    const templateData = {
      structureExample: this.generateStructureExample(),
      serviceList: this.generateServiceList(services),
      completeExample: this.generateCompleteExample(services),
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

  generateCompleteExample(services) {
    // Create a realistic example with common services
    const exampleServices = ["postgres", "redis", "kafka"];
    const serviceConfig = {};

    exampleServices.forEach((serviceName) => {
      if (services[serviceName]) {
        // Add basic configuration for each service
        switch (serviceName) {
          case "postgres":
            serviceConfig[serviceName] = {
              database: "my_app_db",
              password: "secure_password",
              port: 5432,
            };
            break;
          case "redis":
            serviceConfig[serviceName] = {
              password: "redis_password",
              max_memory: "512m",
            };
            break;
          case "kafka":
            serviceConfig[serviceName] = {
              topics: ["events", "notifications"],
            };
            break;
        }
      }
    });

    const exampleConfig = {
      project: {
        name: "my-fullstack-app",
        type: "web",
      },
      stack: {
        enabled: exampleServices,
      },
      service_configuration: serviceConfig,
    };

    return this.toYaml(exampleConfig);
  }

  toYaml(obj) {
    const yaml = require("js-yaml");
    return yaml.dump(obj, {
      indent: 2,
      lineWidth: 80,
      noRefs: true,
    });
  }
}

module.exports = ConfigurationGuideGenerator;

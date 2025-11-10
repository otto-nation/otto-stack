const fs = require("fs");
const path = require("path");
const yaml = require("js-yaml");

class ServiceAnalyzer {
  constructor(config) {
    this.config = config;
  }

  loadAllServices() {
    const services = {};
    this.readServicesFromDir(this.config.servicesDir, services);
    return services;
  }

  readServicesFromDir(dir, services) {
    const items = fs.readdirSync(dir);
    items.forEach((item) => {
      const fullPath = path.join(dir, item);
      const stat = fs.statSync(fullPath);

      if (stat.isDirectory()) {
        this.readServicesFromDir(fullPath, services);
      } else if (item.endsWith(".yaml") || item.endsWith(".yml")) {
        const serviceName = path.basename(item, path.extname(item));
        const content = fs.readFileSync(fullPath, "utf8");
        const serviceConfig = yaml.load(content);

        if (!serviceConfig.hidden) {
          // Add file path for categorization
          serviceConfig._filePath = fullPath;
          services[serviceName] = serviceConfig;
        }
      }
    });
  }

  categorizeServices(services) {
    const categories = {};

    Object.entries(services).forEach(([name, config]) => {
      const category = this.detectCategory(name, config);
      if (!categories[category]) {
        categories[category] = [];
      }
      categories[category].push([name, config]);
    });

    return categories;
  }

  detectCategory(name, config) {
    // Use the folder structure as the category
    // Service files are organized like: services/database/postgres.yaml, services/cache/redis.yaml
    const servicePath = config._filePath || "";
    const pathParts = servicePath.split("/");

    // Find the category folder (should be the parent of the service file)
    const categoryIndex = pathParts.findIndex((part) => part === "services");
    if (categoryIndex !== -1 && pathParts[categoryIndex + 1]) {
      return pathParts[categoryIndex + 1];
    }

    return "other";
  }

  getCategoryConfig(categoryName) {
    const categoryConfigs = {
      database: { icon: "🗄️", order: 1 },
      cache: { icon: "⚡", order: 2 },
      messaging: { icon: "📨", order: 3 },
      cloud: { icon: "☁️", order: 4 },
      observability: { icon: "🔍", order: 5 },
      other: { icon: "🔧", order: 99 },
    };

    return categoryConfigs[categoryName] || categoryConfigs.other;
  }

  extractServiceDetails(config) {
    const details = [];

    if (config.container?.ports?.length > 0) {
      details.push({
        label: "Default Port",
        value: config.container.ports[0].external,
      });
    }

    if (config.service?.connection?.type) {
      details.push({
        label: "Connection Type",
        value: config.service.connection.type,
      });
    }

    if (config.service?.dependencies?.provides?.length > 0) {
      details.push({
        label: "Provides",
        value: config.service.dependencies.provides.join(", "),
      });
    }

    if (config.service?.dependencies?.required?.length > 0) {
      details.push({
        label: "Requires",
        value: config.service.dependencies.required.join(", "),
      });
    }

    if (config.documentation?.web_interfaces?.length > 0) {
      const interfaces = config.documentation.web_interfaces
        .map((iface) => `[${iface.name}](${iface.url})`)
        .join(", ");
      details.push({ label: "Web Interface", value: interfaces });
    }

    return details;
  }
}

module.exports = ServiceAnalyzer;

const yaml = require("js-yaml");

class BaseGenerator {
  constructor(config) {
    this.config = config;
  }

  createFrontmatter(title, description, lead, weight = 50) {
    return {
      title,
      description,
      lead,
      date: "2025-10-01",
      lastmod: this.getToday(),
      draft: false,
      weight,
      toc: true,
    };
  }

  getToday() {
    return new Date().toISOString().split("T")[0];
  }

  formatDocument(frontmatter, content) {
    const frontmatterYaml = yaml.dump(frontmatter);
    return `---\n${frontmatterYaml}---\n\n${content}`;
  }

  dumpYaml(data, options = {}) {
    return yaml.dump(data, { indent: 2, lineWidth: 80, ...options });
  }

  handleError(operation, error) {
    throw new Error(`Failed to ${operation}: ${error.message}`);
  }
}

module.exports = BaseGenerator;

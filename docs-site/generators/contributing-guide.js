const fs = require("fs");
const path = require("path");
const yaml = require("js-yaml");

class ContributingGuideGenerator {
  constructor(config) {
    this.config = config;
  }

  generate() {
    console.log("Generating contributing guide...");

    const rootContributingPath = path.join(__dirname, "../../CONTRIBUTING.md");

    if (!fs.existsSync(rootContributingPath)) {
      throw new Error("Root CONTRIBUTING.md file not found");
    }

    const content = fs.readFileSync(rootContributingPath, "utf8");

    const frontmatter = {
      title: "Contributing",
      description: "Guide for contributing to otto-stack development",
      lead: "Learn how to contribute to otto-stack development",
      date: "2025-10-01",
      lastmod: new Date().toISOString().split("T")[0],
      draft: false,
      weight: 60,
      toc: true,
    };

    return `---\n${yaml.dump(frontmatter)}---\n\n${content}`;
  }
}

module.exports = ContributingGuideGenerator;

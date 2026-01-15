const fs = require("fs");
const path = require("path");

class ContributingGuideGenerator {
  constructor(config) {
    this.config = config;
  }

  generate() {
    console.log("Generating contributing guide...");

    // Read the root CONTRIBUTING.md file
    const rootContributingPath = path.join(__dirname, "../../CONTRIBUTING.md");

    if (!fs.existsSync(rootContributingPath)) {
      throw new Error("Root CONTRIBUTING.md file not found");
    }

    const contributingContent = fs.readFileSync(rootContributingPath, "utf8");

    // Add Hugo frontmatter
    const today = new Date().toISOString().split("T")[0];
    const frontmatter = {
      title: "Contributing",
      description: "Guide for contributing to otto-stack development",
      lead: "Learn how to contribute to otto-stack development",
      date: "2025-10-01",
      lastmod: today,
      draft: false,
      weight: 60,
      toc: true,
    };

    const frontmatterYaml = require("js-yaml").dump(frontmatter);
    return `---\n${frontmatterYaml}---\n\n${contributingContent}`;
  }
}

module.exports = ContributingGuideGenerator;

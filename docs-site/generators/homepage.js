const fs = require("fs");
const path = require("path");
const yaml = require("js-yaml");

class HomepageGenerator {
  constructor(config) {
    this.config = config;
  }

  generate() {
    console.log("Generating homepage from README...");

    try {
      const readmePath = path.join("..", "README.md");
      const readmeContent = fs.readFileSync(readmePath, "utf8");

      // Extract content after the first heading
      const lines = readmeContent.split("\n");
      let contentStart = 0;
      for (let i = 0; i < lines.length; i++) {
        if (lines[i].startsWith("# ")) {
          contentStart = i;
          break;
        }
      }

      let content = lines.slice(contentStart).join("\n");

      // Fix links for Hugo with baseURL subdirectory
      // Convert docs-site/content/file.md -> file/ (relative format)
      content = content.replace(/docs-site\/content\/([^)]+)\.md/g, "$1/");

      // Convert any remaining .md links to Hugo format (relative)
      content = content.replace(/\]\(([^)]+)\.md\)/g, "]($1/)");

      // Fix docs-site root links
      content = content.replace(
        /\[([^\]]+)\]\(docs-site\/\)/g,
        '[$1]({{< ref "/" >}})',
      );

      // Fix LICENSE link
      content = content.replace(
        /\[([^\]]+)\]\(LICENSE\)/g,
        "[$1](https://github.com/otto-nation/otto-stack/blob/main/LICENSE)",
      );

      const today = new Date().toISOString().split("T")[0];
      const frontmatter = {
        title: "otto-stack",
        description:
          "A powerful development stack management tool built in Go for streamlined local development automation",
        lead: "Streamline your local development with powerful CLI tools and automated service management",
        date: "2025-10-01",
        lastmod: today,
        draft: false,
        weight: 50,
        toc: true,
      };

      const frontmatterYaml = yaml.dump(frontmatter);
      return `---\n${frontmatterYaml}---\n\n${content}`;
    } catch (error) {
      console.warn("Could not generate homepage:", error.message);
      return "";
    }
  }
}

module.exports = HomepageGenerator;

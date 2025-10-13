#!/usr/bin/env node

const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");
const yaml = require("js-yaml");

// Build the CLI first
console.log("Building otto-stack CLI...");
execSync("cd .. && go build -o docs-site/otto-stack ./cmd/otto-stack", {
  stdio: "inherit",
});

// Generate CLI reference
function generateCLIReference() {
  console.log("Generating CLI reference...");

  const helpOutput = execSync("./otto-stack --help", { encoding: "utf8" });
  const commands = extractCommands(helpOutput);

  const today = new Date().toISOString().split("T")[0];
  let markdown = `---
title: "CLI Reference"
description: "Complete command reference for otto-stack CLI"
lead: "Comprehensive reference for all otto-stack CLI commands and their usage"
date: "2025-10-01"
lastmod: "${today}"
draft: false
weight: 50
toc: true
---

# otto-stack CLI Reference

${helpOutput.split("\n").slice(0, 3).join("\n")}

## Commands

`;

  // Generate detailed help for each command
  commands.forEach((cmd) => {
    try {
      const cmdHelp = execSync(`./otto-stack ${cmd} --help`, {
        encoding: "utf8",
      });
      markdown += `### ${cmd}\n\n\`\`\`\n${cmdHelp}\`\`\`\n\n`;
    } catch (error) {
      console.warn(`Could not get help for command: ${cmd}`);
    }
  });

  fs.writeFileSync("content/reference.md", markdown);
  console.log("✅ Generated content/reference.md");
}

// Generate services guide
function generateServicesGuide() {
  console.log("Generating services guide...");

  try {
    const servicesDir = "../internal/config/services";
    const services = {};

    // Read all service YAML files recursively
    function readServicesFromDir(dir) {
      const items = fs.readdirSync(dir);
      items.forEach((item) => {
        const fullPath = path.join(dir, item);
        const stat = fs.statSync(fullPath);

        if (stat.isDirectory()) {
          readServicesFromDir(fullPath);
        } else if (item.endsWith(".yaml") || item.endsWith(".yml")) {
          const serviceName = path.basename(item, path.extname(item));
          const content = fs.readFileSync(fullPath, "utf8");
          services[serviceName] = yaml.load(content);
        }
      });
    }

    readServicesFromDir(servicesDir);

    const today = new Date().toISOString().split("T")[0];
    let markdown = `---
title: "Services"
description: "Available services and configuration options"
lead: "Explore all the services you can use with otto-stack"
date: "2025-10-01"
lastmod: "${today}"
draft: false
weight: 30
toc: true
---

# Available Services

${Object.keys(services).length} services available for your development stack.

`;

    // Sort services by name
    const sortedServices = Object.entries(services).sort(([a], [b]) =>
      a.localeCompare(b),
    );

    sortedServices.forEach(([name, config]) => {
      markdown += `## ${name}\n\n`;
      if (config.description) {
        markdown += `${config.description}\n\n`;
      }
      if (config.defaults && config.defaults.port) {
        markdown += `**Default Port:** ${config.defaults.port}\n\n`;
      }
      if (config.docker && config.docker.services) {
        const serviceNames = Object.keys(config.docker.services);
        markdown += `**Services:** ${serviceNames.join(", ")}\n\n`;
      }
      markdown += "---\n\n";
    });

    fs.writeFileSync("content/services.md", markdown);
    console.log("✅ Generated content/services.md");
  } catch (error) {
    console.warn("Could not generate services guide:", error.message);
  }
}

// Generate homepage from README
function generateHomepage() {
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

    // Fix links: convert docs-site/content/file.md to file.md
    content = content.replace(/docs-site\/content\/([^)]+\.md)/g, "$1");

    // Fix other problematic links
    content = content.replace(/\[([^\]]+)\]\(docs-site\/\)/g, "[$1](/)");
    content = content.replace(
      /\[([^\]]+)\]\(LICENSE\)/g,
      "[$1](https://github.com/otto-nation/otto-stack/blob/main/LICENSE)",
    );

    // Create frontmatter for Hugo with yyyy-MM-dd format
    const today = new Date().toISOString().split("T")[0];
    const frontmatter = `---
title: "otto-stack"
description: "A powerful development stack management tool built in Go for streamlined local development automation"
lead: "Streamline your local development with powerful CLI tools and automated service management"
date: "2025-10-01"
lastmod: "${today}"
draft: false
weight: 50
toc: true
---

`;

    const fullContent = frontmatter + content;
    fs.writeFileSync("content/_index.md", fullContent);
    console.log("✅ Generated content/_index.md from README.md");
  } catch (error) {
    console.warn("Could not generate homepage:", error.message);
  }
}

// Extract command names from help output
function extractCommands(helpOutput) {
  const lines = helpOutput.split("\n");
  const commandsStart = lines.findIndex((line) =>
    line.includes("Available Commands:"),
  );
  if (commandsStart === -1) return [];

  const commands = [];
  for (let i = commandsStart + 1; i < lines.length; i++) {
    const line = lines[i].trim();
    if (!line || line.startsWith("Flags:") || line.startsWith("Use ")) break;

    const match = line.match(/^\s*(\w+)\s+/);
    if (match && match[1] !== "help" && match[1] !== "completion") {
      commands.push(match[1]);
    }
  }
  return commands;
}

// Create content directory if it doesn't exist
if (!fs.existsSync("content")) {
  fs.mkdirSync("content", { recursive: true });
}

// Generate all docs
generateHomepage();
generateCLIReference();
generateServicesGuide();

// Format generated files with prettier
console.log("Formatting generated documentation...");
try {
  execSync("npm run format", { stdio: "inherit" });
  console.log("✅ Documentation formatted");
} catch (error) {
  console.warn("Could not format documentation:", error.message);
}

// Cleanup
fs.unlinkSync("./otto-stack");
console.log("🎉 Documentation generation complete!");

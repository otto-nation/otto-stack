const yaml = require("js-yaml");
const fs = require("fs");
const path = require("path");

class CLIReferenceGenerator {
  constructor(config) {
    this.config = config;
  }

  generate() {
    console.log("Generating CLI reference...");

    const commandsYaml = this.loadCommandsYaml();

    const today = new Date().toISOString().split("T")[0];
    const frontmatter = {
      title: "CLI Reference",
      description: "Complete command reference for otto-stack CLI",
      lead: "Comprehensive reference for all otto-stack CLI commands and their usage",
      date: "2025-10-01",
      lastmod: today,
      draft: false,
      weight: 50,
      toc: true,
    };

    let content = `# otto-stack CLI Reference

${commandsYaml.metadata?.description || "A powerful development stack management tool for streamlined local development automation"}

## Command Categories

`;

    // Add categories section
    if (commandsYaml.categories) {
      Object.entries(commandsYaml.categories).forEach(([key, category]) => {
        content += `### ${category.icon} ${category.name}\n\n`;
        content += `${category.description}\n\n`;
        content += `**Commands:** ${category.commands.map((cmd) => `\`${cmd}\``).join(", ")}\n\n`;
      });
    }

    content += `## Commands\n\n`;

    // Add detailed command information
    if (commandsYaml.commands) {
      Object.entries(commandsYaml.commands).forEach(
        ([commandName, command]) => {
          content += this.generateCommandSection(commandName, command);
        },
      );
    }

    // Add global flags section
    if (commandsYaml.global_flags) {
      content += `## Global Flags\n\n`;
      content += `These flags are available for all commands:\n\n`;

      Object.entries(commandsYaml.global_flags).forEach(([flagName, flag]) => {
        content += `- \`--${flagName}\``;
        if (flag.short) content += `, \`-${flag.short}\``;
        content += `: ${flag.description}`;
        if (flag.default !== undefined)
          content += ` (default: \`${flag.default}\`)`;
        content += `\n`;
      });
      content += `\n`;
    }

    const frontmatterYaml = yaml.dump(frontmatter);
    return `---\n${frontmatterYaml}---\n\n${content}`;
  }

  generateCommandSection(commandName, command) {
    let section = `### \`${commandName}\`\n\n`;

    section += `${command.description}\n\n`;

    if (command.long_description) {
      section += `${command.long_description.trim()}\n\n`;
    }

    if (command.usage) {
      section += `**Usage:** \`otto-stack ${command.usage}\`\n\n`;
    }

    if (command.aliases && command.aliases.length > 0) {
      section += `**Aliases:** ${command.aliases.map((alias) => `\`${alias}\``).join(", ")}\n\n`;
    }

    if (command.examples && command.examples.length > 0) {
      section += `**Examples:**\n\n`;
      command.examples.forEach((example) => {
        section += `\`\`\`bash\n${example.command}\n\`\`\`\n`;
        if (example.description) {
          section += `${example.description}\n\n`;
        }
      });
    }

    if (command.flags && Object.keys(command.flags).length > 0) {
      section += `**Flags:**\n\n`;
      Object.entries(command.flags).forEach(([flagName, flag]) => {
        section += `- \`--${flagName}\``;
        if (flag.short) section += `, \`-${flag.short}\``;
        section += ` (\`${flag.type}\`)`;
        section += `: ${flag.description}`;
        if (flag.default !== undefined)
          section += ` (default: \`${flag.default}\`)`;
        if (flag.options)
          section += ` (options: ${flag.options.map((opt) => `\`${opt}\``).join(", ")})`;
        section += `\n`;
      });
      section += `\n`;
    }

    if (command.related_commands && command.related_commands.length > 0) {
      section += `**Related Commands:** ${command.related_commands.map((cmd) => `[\`${cmd}\`](#${cmd})`).join(", ")}\n\n`;
    }

    if (command.tips && command.tips.length > 0) {
      section += `**Tips:**\n\n`;
      command.tips.forEach((tip) => {
        section += `- ${tip}\n`;
      });
      section += `\n`;
    }

    return section;
  }

  loadCommandsYaml() {
    const commandsPath = path.join(
      process.cwd(),
      "..",
      "internal",
      "config",
      "commands.yaml",
    );

    if (!fs.existsSync(commandsPath)) {
      throw new Error(`Commands configuration file not found: ${commandsPath}`);
    }

    try {
      const commandsContent = fs.readFileSync(commandsPath, "utf8");
      return yaml.load(commandsContent);
    } catch (error) {
      throw new Error(`Failed to load commands.yaml: ${error.message}`);
    }
  }
}

module.exports = CLIReferenceGenerator;

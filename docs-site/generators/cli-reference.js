const yaml = require("js-yaml");
const fs = require("fs");
const path = require("path");
const BaseGenerator = require("./base-generator");

const COMMANDS_PATH = path.join(
  process.cwd(),
  "..",
  "internal",
  "config",
  "commands.yaml",
);

class CLIReferenceGenerator extends BaseGenerator {
  generate() {
    console.log("Generating CLI reference...");

    try {
      const commandsYaml = this.loadCommandsYaml();
      const frontmatter = this.createFrontmatter(
        "CLI Reference",
        "Complete command reference for otto-stack CLI",
        "Comprehensive reference for all otto-stack CLI commands and their usage",
        50,
      );

      const content = this.generateContent(commandsYaml);
      return this.formatDocument(frontmatter, content);
    } catch (error) {
      this.handleError("generate CLI reference", error);
    }
  }

  generateContent(commandsYaml) {
    let content = `<!-- 
  ⚠️  AUTO-GENERATED FILE - DO NOT EDIT DIRECTLY
  This file is generated from internal/config/commands.yaml
  To make changes, edit the source file and run: task generate:docs
-->

# otto-stack CLI Reference

${commandsYaml.metadata?.description || "A powerful development stack management tool for streamlined local development automation"}

## Command Categories

`;

    if (commandsYaml.categories) {
      Object.entries(commandsYaml.categories).forEach(([key, category]) => {
        content += `### ${category.icon} ${category.name}\n\n`;
        content += `${category.description}\n\n`;
        content += `**Commands:** ${category.commands.map((cmd) => `\`${cmd}\``).join(", ")}\n\n`;
      });
    }

    content += `## Commands\n\n`;

    if (commandsYaml.commands) {
      Object.entries(commandsYaml.commands).forEach(
        ([commandName, command]) => {
          content += this.generateCommandSection(commandName, command);
        },
      );
    }

    if (commandsYaml.global_flags) {
      content += this.generateGlobalFlagsSection(commandsYaml.global_flags);
    }

    return content;
  }

  generateGlobalFlagsSection(globalFlags) {
    let section = `## Global Flags\n\n`;
    section += `These flags are available for all commands:\n\n`;

    Object.entries(globalFlags).forEach(([flagName, flag]) => {
      section += `- \`--${flagName}\``;
      if (flag.short) section += `, \`-${flag.short}\``;
      section += `: ${flag.description}`;
      if (flag.default !== undefined)
        section += ` (default: \`${flag.default}\`)`;
      section += `\n`;
    });
    section += `\n`;

    return section;
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
    if (!fs.existsSync(COMMANDS_PATH)) {
      throw new Error(
        `Commands configuration file not found: ${COMMANDS_PATH}`,
      );
    }

    try {
      const commandsContent = fs.readFileSync(COMMANDS_PATH, "utf8");
      return yaml.load(commandsContent);
    } catch (error) {
      this.handleError("load commands.yaml", error);
    }
  }
}

module.exports = CLIReferenceGenerator;

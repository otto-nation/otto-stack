const { execSync } = require("child_process");
const yaml = require("js-yaml");

class CLIReferenceGenerator {
  constructor(config) {
    this.config = config;
  }

  generate() {
    console.log("Generating CLI reference...");

    const helpOutput = execSync("./otto-stack --help", { encoding: "utf8" });
    const commands = this.extractCommands(helpOutput);

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

${helpOutput.split("\n").slice(0, 3).join("\n")}

## Available Commands

`;

    commands.forEach((command) => {
      content += `### ${command.name}\n\n${command.description}\n\n`;
      if (command.usage) {
        content += `**Usage:** \`${command.usage}\`\n\n`;
      }
      if (command.flags?.length > 0) {
        content += "**Flags:**\n\n";
        command.flags.forEach((flag) => {
          content += `- \`${flag.name}\`: ${flag.description}\n`;
        });
        content += "\n";
      }
    });

    const frontmatterYaml = yaml.dump(frontmatter);
    return `---\n${frontmatterYaml}---\n\n${content}`;
  }

  extractCommands(helpOutput) {
    const lines = helpOutput.split("\n");
    const commandsStart = lines.findIndex((line) =>
      line.includes("Available Commands:"),
    );
    if (commandsStart === -1) return [];

    const commands = [];
    for (let i = commandsStart + 1; i < lines.length; i++) {
      const line = lines[i].trim();
      if (!line || line.startsWith("Flags:") || line.startsWith("Use ")) break;

      const match = line.match(/^\s*(\w+)\s+(.+)/);
      if (match && match[1] !== "help" && match[1] !== "completion") {
        commands.push({
          name: match[1],
          description: match[2] || "",
          usage: `otto-stack ${match[1]}`,
          flags: [],
        });
      }
    }
    return commands;
  }
}

module.exports = CLIReferenceGenerator;

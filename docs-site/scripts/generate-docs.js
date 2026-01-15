#!/usr/bin/env node

const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");

// Load configuration
const config = require("../docs-config");

// Load generators
const ServicesGuideGenerator = require("../generators/services-guide");
const CLIReferenceGenerator = require("../generators/cli-reference");
const HomepageGenerator = require("../generators/homepage");
const ConfigurationGuideGenerator = require("../generators/configuration-guide");
const ContributingGuideGenerator = require("../generators/contributing-guide");

// Load utilities
const ServiceAnalyzer = require("../utils/service-analyzer");
const SchemaValidator = require("../utils/schema-validator");

const generators = {
  "services-guide": ServicesGuideGenerator,
  "cli-reference": CLIReferenceGenerator,
  homepage: HomepageGenerator,
  "configuration-guide": ConfigurationGuideGenerator,
  "contributing-guide": ContributingGuideGenerator,
};

async function validateServices() {
  if (!config.validation?.enabled) return;

  console.log("Validating service configurations...");

  const analyzer = new ServiceAnalyzer(config);
  const validator = new SchemaValidator();
  const services = analyzer.loadAllServices();

  const { errors, warnings } = validator.validateAllServices(services);

  if (warnings.length > 0) {
    console.warn("⚠️  Validation warnings:");
    warnings.forEach((warning) => console.warn(`   ${warning}`));
  }

  if (errors.length > 0) {
    console.error("❌ Validation errors:");
    errors.forEach((error) => console.error(`   ${error}`));

    if (config.validation.strict) {
      process.exit(1);
    }
  } else {
    console.log("✅ Service validation passed");
  }
}

async function main() {
  // Build the CLI first
  console.log("Building otto-stack CLI...");
  execSync("cd .. && go build -o docs-site/otto-stack ./cmd/otto-stack", {
    stdio: "inherit",
  });

  // Validate services
  await validateServices();

  // Create content directory if it doesn't exist
  if (!fs.existsSync(config.outputDir)) {
    fs.mkdirSync(config.outputDir, { recursive: true });
  }

  // Run enabled generators
  for (const generatorConfig of config.generators) {
    if (!generatorConfig.enabled) continue;

    const GeneratorClass = generators[generatorConfig.name];
    if (!GeneratorClass) {
      console.warn(`Unknown generator: ${generatorConfig.name}`);
      continue;
    }

    try {
      const generator = new GeneratorClass(config);
      const content = generator.generate();

      if (content) {
        const outputPath = path.join(config.outputDir, generatorConfig.output);
        fs.writeFileSync(outputPath, content);
        console.log(`✅ Generated ${generatorConfig.output}`);
      }
    } catch (error) {
      console.error(
        `❌ Failed to generate ${generatorConfig.name}:`,
        error.message,
      );
    }
  }

  // Format generated files with prettier
  console.log("Formatting generated documentation...");
  try {
    execSync("npm run format", { stdio: "inherit" });
    console.log("✅ Documentation formatted");
  } catch (error) {
    console.warn("Could not format documentation:", error.message);
  }

  // Cleanup
  if (fs.existsSync("./otto-stack")) {
    fs.unlinkSync("./otto-stack");
  }

  console.log("🎉 Documentation generation complete!");
}

main().catch(console.error);

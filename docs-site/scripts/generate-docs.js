#!/usr/bin/env node

const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");

// Parse CLI arguments
const args = process.argv.slice(2);
const flags = {
  skipBuild: args.includes("--skip-build"),
  skipFormat: args.includes("--skip-format"),
  generator: args.find((arg) => arg.startsWith("--generator="))?.split("=")[1],
};

// Load configuration
const config = require("../docs-config");

// Load generators
const generators = {
  "services-guide": require("../generators/services-guide"),
  "cli-reference": require("../generators/cli-reference"),
  homepage: require("../generators/homepage"),
  "configuration-guide": require("../generators/configuration-guide"),
  "contributing-guide": require("../generators/contributing-guide"),
};

// Load utilities
const ServiceAnalyzer = require("../utils/service-analyzer");
const SchemaValidator = require("../utils/schema-validator");

function validateServices() {
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

function buildCLI() {
  if (flags.skipBuild) {
    console.log("⏭️  Skipping CLI build");
    return;
  }

  const binaryPath = path.join(__dirname, "../otto-stack");
  const projectRoot = path.join(__dirname, "../..");

  // Check if binary exists and is recent
  if (fs.existsSync(binaryPath)) {
    const binaryAge = Date.now() - fs.statSync(binaryPath).mtimeMs;
    if (binaryAge < 60000) {
      // Less than 1 minute old
      console.log("⏭️  Using existing CLI binary");
      return;
    }
  }

  console.log("Building otto-stack CLI...");
  try {
    execSync(`go build -o ${binaryPath} ./cmd/otto-stack`, {
      cwd: projectRoot,
      stdio: "inherit",
    });
  } catch (error) {
    console.error("❌ Failed to build CLI");
    process.exit(1);
  }
}

function runGenerators() {
  const results = { success: [], failed: [] };

  // Filter generators if specific one requested
  const generatorsToRun = flags.generator
    ? config.generators.filter((g) => g.name === flags.generator)
    : config.generators.filter((g) => g.enabled);

  if (generatorsToRun.length === 0) {
    console.error(
      `❌ No generators found matching: ${flags.generator || "enabled"}`,
    );
    process.exit(1);
  }

  // Create output directory
  if (!fs.existsSync(config.outputDir)) {
    fs.mkdirSync(config.outputDir, { recursive: true });
  }

  // Run generators
  for (const generatorConfig of generatorsToRun) {
    const GeneratorClass = generators[generatorConfig.name];
    if (!GeneratorClass) {
      console.warn(`⚠️  Unknown generator: ${generatorConfig.name}`);
      results.failed.push(generatorConfig.name);
      continue;
    }

    try {
      const generator = new GeneratorClass(config);
      const content = generator.generate();

      if (content) {
        const outputPath = path.join(config.outputDir, generatorConfig.output);
        fs.writeFileSync(outputPath, content);
        console.log(`✅ Generated ${generatorConfig.output}`);
        results.success.push(generatorConfig.name);
      }
    } catch (error) {
      console.error(`❌ Failed to generate ${generatorConfig.name}:`);
      console.error(`   ${error.message}`);
      if (process.env.DEBUG) {
        console.error(error.stack);
      }
      results.failed.push(generatorConfig.name);
    }
  }

  return results;
}

function formatDocs() {
  if (flags.skipFormat) {
    console.log("⏭️  Skipping formatting");
    return;
  }

  console.log("Formatting generated documentation...");
  try {
    execSync("npm run format", { stdio: "inherit" });
    console.log("✅ Documentation formatted");
  } catch (error) {
    console.warn("⚠️  Could not format documentation");
  }
}

function cleanup() {
  const binaryPath = path.join(__dirname, "../otto-stack");
  if (fs.existsSync(binaryPath)) {
    fs.unlinkSync(binaryPath);
  }
}

function main() {
  console.log("📚 Otto-Stack Documentation Generator\n");

  buildCLI();
  validateServices();
  const results = runGenerators();
  formatDocs();
  cleanup();

  // Summary
  console.log("\n📊 Generation Summary:");
  console.log(`   ✅ Success: ${results.success.length}`);
  console.log(`   ❌ Failed: ${results.failed.length}`);

  if (results.failed.length > 0) {
    console.log(`\n❌ Failed generators: ${results.failed.join(", ")}`);
    process.exit(1);
  }

  console.log("\n🎉 Documentation generation complete!");
}

main();

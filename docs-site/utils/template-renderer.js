const fs = require("fs");
const path = require("path");
const yaml = require("js-yaml");
const Handlebars = require("handlebars");

class TemplateRenderer {
  constructor() {
    this.templatesDir = path.join(__dirname, "../templates");
    this.registerHelpers();
  }

  registerHelpers() {
    // Register toYaml helper for converting objects to YAML
    Handlebars.registerHelper("toYaml", function (obj) {
      if (!obj) return "";
      return yaml.dump(obj, { indent: 2, lineWidth: 80, noRefs: true });
    });
  }

  render(templateName, data, frontmatter = {}) {
    const templatePath = path.join(this.templatesDir, templateName);
    const templateSource = fs.readFileSync(templatePath, "utf8");
    const template = Handlebars.compile(templateSource);

    let rendered = template(data);

    // Add frontmatter if provided
    if (Object.keys(frontmatter).length > 0) {
      const frontmatterYaml = yaml.dump(frontmatter);
      rendered = `---\n${frontmatterYaml}---\n\n${rendered}`;
    }

    return rendered;
  }
}

module.exports = TemplateRenderer;

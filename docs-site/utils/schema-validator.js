const Ajv = require("ajv");

class SchemaValidator {
  constructor() {
    this.ajv = new Ajv({ allErrors: true });
  }

  validateAllServices(services) {
    const errors = [];
    const warnings = [];

    Object.entries(services).forEach(([name, config]) => {
      if (config.configuration_schema) {
        try {
          this.ajv.compile(config.configuration_schema);
        } catch (error) {
          errors.push(`${name}: Invalid schema - ${error.message}`);
        }
      }

      if (!config.description) {
        warnings.push(`${name}: Missing description`);
      }
    });

    return { errors, warnings };
  }

  validateServiceConfig(serviceName, config, schema) {
    const validate = this.ajv.compile(schema);
    const valid = validate(config);

    if (!valid) {
      return {
        valid: false,
        errors: validate.errors.map(
          (err) => `${serviceName}: ${err.instancePath} ${err.message}`,
        ),
      };
    }

    return { valid: true, errors: [] };
  }
}

module.exports = SchemaValidator;

# Contributing to otto-stack

Thank you for your interest in contributing to otto-stack! This guide will help you get started with development and explain how to contribute effectively.

## 📝 Quick Contributor Checklist

- [ ] Fork and clone the repository
- [ ] Install dependencies: `task setup`
- [ ] Build the project: `task build`
- [ ] Run tests: `task test`
- [ ] Edit YAML manifests (`internal/config/commands.yaml`, `internal/config/services/services.yaml`) for changes
- [ ] Run `task docs` to update documentation
- [ ] Commit both the manifest and generated docs
- [ ] Follow the contributing guide and PR template

## 📋 Overview

We welcome contributions from the development community! Whether you're fixing bugs, adding new services, improving documentation, or enhancing existing features, your contributions help make the framework better for everyone.

## 🤝 Community & Support

### GitHub Repository

- **Main Repository**: [otto-nation/otto-stack](https://github.com/otto-nation/otto-stack)
- **Issues**: [Report bugs and request features](https://github.com/otto-nation/otto-stack/issues)
- **Discussions**: [Ask questions and share ideas](https://github.com/otto-nation/otto-stack/issues)
- **Releases**: [Latest versions and changelog](https://github.com/otto-nation/otto-stack/releases)

### Getting Help

**Before opening an issue:**
1. Check the [Troubleshooting Guide](troubleshooting.md)
2. Review [existing issues](https://github.com/otto-nation/otto-stack/issues)
3. Search [issues](https://github.com/otto-nation/otto-stack/issues)

**For support requests:**
- Use [GitHub Issues](https://github.com/otto-nation/otto-stack/issues) for questions
- Check the [CLI Reference](cli-reference.md) for command help
- Run `otto-stack doctor` for system diagnostics

## 📚 Automated Documentation & YAML Manifests

This project uses automated documentation generation from YAML service configurations. Here's how it works:

### Documentation Generation (Node.js-based)

Documentation is automatically generated from YAML service configurations using a Node.js-based system:

```bash
# Generate all documentation
task docs

# Serve documentation locally for development
task docs-serve
```

This will update `docs-site/content/cli-reference.md` and `docs-site/content/services.md` based on the latest service configurations.

**Contributor Workflow Checklist:**
1. Ensure you have Node.js installed and dependencies set up (`cd docs-site && npm install`).
2. Edit service YAML files in `internal/config/services/` to add or update services.
3. Run `task docs` to regenerate documentation from service configurations.
4. Commit both the service configurations and the updated docs.
5. Never manually edit auto-generated docs (`docs-site/content/cli-reference.md`, `docs-site/content/services.md`).
6. Use `task docs-serve` to preview documentation changes locally.

Documentation for CLI commands (`docs-site/content/cli-reference.md`) and services (`docs-site/content/services.md`) is auto-generated from service YAML configurations using the Node.js documentation generator.

## 🚀 Getting Started

### Prerequisites

- **Go**: 1.24+ (managed via `.go-version` file)
- **Task**: [Task runner](https://taskfile.dev/) for build automation
- **Docker & Docker Compose**: For service management
- **Git**: For version control

### Development Setup

1. **Fork and clone the repository**:
   ```bash
   git clone https://github.com/your-username/otto-stack.git
   cd otto-stack
   ```

2. **Install Go version manager and dependencies**:
   ```bash
   # Install Go version specified in .go-version
   task setup
   ```

3. **Build the project**:
   ```bash
   task build
   ```

4. **Run tests**:
   ```bash
   task test
   ```

5. **Verify installation**:
   ```bash
   ./build/otto-stack version
   ./build/otto-stack help
   ```

## 🏗️ Architecture Overview

### Directory Structure

```
otto-stack/
├── cmd/                          # CLI command implementations
│   └── otto-stack/              # Main CLI application
├── internal/                     # Internal Go packages
│   ├── cli/                     # CLI command handlers
│   ├── config/                  # Configuration management
│   │   ├── commands.yaml        # CLI commands configuration
│   │   └── services/            # Service configurations by category
│   │       ├── database/        # Database services (postgres, mysql)
│   │       ├── cache/           # Cache services (redis)
│   │       ├── messaging/       # Messaging services (kafka, rabbitmq)
│   │       ├── observability/   # Monitoring services (prometheus, jaeger)
│   │       └── cloud/           # Cloud services (localstack-*)
│   └── pkg/                     # Reusable packages
├── docs-site/                   # Hugo documentation site
│   ├── content/                 # Markdown content files
│   ├── config/                  # Hugo configuration
│   ├── generators/              # Node.js documentation generators
│   ├── templates/               # Handlebars templates
│   ├── utils/                   # Documentation utilities
│   └── themes/                  # Hugo themes
├── scripts/                     # Build and utility scripts
└── otto-stack-config.sample.yaml # Sample configuration
```

## 🛠️ Adding New Services

### Step 1: Choose Service Category

Determine the appropriate category for your service:
- `database/` - Database services (postgres, mysql)
- `cache/` - Caching services (redis)
- `messaging/` - Message queues (kafka, rabbitmq)
- `observability/` - Monitoring tools (prometheus, jaeger)
- `cloud/` - Cloud service emulations (localstack-*)

### Step 2: Create Service Configuration

Create a new YAML file in the appropriate category directory:

```bash
# Example: Adding a new database service
touch internal/config/services/database/my-service.yaml
```

### Step 3: Define Service Configuration

Create the service YAML configuration following this structure:

```yaml
name: my-service
description: A brief description of what this service does
service_type: container

environment:
  MY_SERVICE_HOST: ${MY_SERVICE_HOST:-localhost}
  MY_SERVICE_PORT: ${MY_SERVICE_PORT:-8080}
  MY_SERVICE_PASSWORD: ${MY_SERVICE_PASSWORD:-password}

container:
  image: my-service:latest
  ports:
    - external: "8080"
      internal: "8080"
      protocol: tcp
  environment:
    MY_SERVICE_PASSWORD: ${MY_SERVICE_PASSWORD:-password}
  restart: unless-stopped
  memory_limit: 256m
  command:
    - my-service
    - --config=/etc/my-service/config.yml
  health_check:
    test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
    interval: 30s
    timeout: 10s
    retries: 3

service:
  connection:
    type: cli
    default_port: 8080
    client: my-service-cli
    host_flag: --host
    port_flag: --port
  dependencies:
    provides:
      - my-service-type
  management:
    connect:
      type: command
      command: ["my-service-cli"]
      args:
        default: ["--host", "localhost", "--port", "8080"]

configuration_schema:
  type: object
  properties:
    password:
      type: string
      default: "password"
      description: Service password
    max_connections:
      type: integer
      default: 100
      description: Maximum number of connections

documentation:
  examples:
    - my-service-cli --host localhost --port 8080 status
    - curl http://localhost:8080/health
  usage_notes: Brief usage notes and tips for the service
  links:
    - https://my-service.example.com/docs
  use_cases:
    - Primary use case
    - Secondary use case
    - Development testing
```

### Step 4: Test Your Service

```bash
# Generate updated documentation
task docs

# Verify the service appears in generated docs
cat docs-site/content/services.md | grep -A 10 "my-service"

# Test documentation generation
task docs-serve
```

### Step 5: Validate Configuration

Ensure your service configuration follows JSON Schema standards and includes:

- **Required fields**: `name`, `description`, `service_type`
- **Proper schema**: Use `type: object` with `properties` for `configuration_schema`
- **Documentation**: Include `examples`, `use_cases`, and `links`
- **Environment variables**: Use `${VAR:-default}` pattern for defaults

## 📝 Documentation Contributions

### Auto-Generated Documentation

Most documentation is automatically generated from templates and configurations:

1. **Service documentation** (`docs-site/content/services.md`):
   - Auto-generated from service YAML files in `internal/config/services/`
   - Uses Handlebars templates in `docs-site/templates/service.md`
   - Includes configuration schemas, examples, and use cases from YAML

2. **CLI reference** (`docs-site/content/cli-reference.md`):
   - Auto-generated from commands configuration
   - Uses Node.js generators in `docs-site/generators/`

3. **Configuration guide** (`docs-site/content/configuration.md`):
   - Auto-generated from service schemas and templates

4. **Homepage** (`docs-site/content/_index.md`):
   - Auto-generated from root `README.md`
   - Processes links and adds Hugo frontmatter

5. **Contributing guide** (`docs-site/content/contributing.md`):
   - Auto-generated from root `CONTRIBUTING.md`
   - Adds Hugo frontmatter for proper site integration

### Manual Documentation

These files can be edited directly:

1. **Setup guide** (`docs-site/content/setup.md`) - Installation instructions
2. **Usage guide** (`docs-site/content/usage.md`) - Basic usage examples
3. **Troubleshooting** (`docs-site/content/troubleshooting.md`) - Common issues

### Template System

Documentation uses Handlebars templates with these helpers:

- `{{{toYaml obj}}}` - Renders YAML configuration examples
- `{{#each items}}` - Iterates over collections
- `{{#if condition}}` - Conditional rendering

**Template locations:**
- `docs-site/templates/service.md` - Service documentation template
- `docs-site/utils/template-renderer.js` - Template rendering engine

### Contributing to Documentation

**For service documentation:**
- Edit service YAML files in `internal/config/services/`
- Add `documentation` section with `examples`, `use_cases`, `links`
- Run `task docs` to regenerate

**For templates:**
- Edit templates in `docs-site/templates/`
- Test with `task docs-serve`
- Ensure templates are in `.prettierignore` to prevent formatting issues

**For manual pages:**
- Edit Markdown files directly in `docs-site/content/`
- Follow existing structure and formatting
- Test locally with `task docs-serve`

### Documentation Standards

- Use clear, concise language
- Include working code examples
- Test all examples before submitting
- Follow existing Markdown structure
- Use proper frontmatter for Hugo pages

## 📦 Submitting Contributions

### Pull Request Process

1. Fork the repository and create a feature branch
2. Make your changes following the coding standards
3. Write or update tests as needed
4. Update documentation if necessary
5. Run `task test` to ensure all tests pass
6. Run `task lint` to check code style
7. Submit a pull request with a clear description

### Commit Message Format

Follow the conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

Examples:
- `feat(services): add MySQL service configuration`
- `fix(cli): resolve Docker network creation issue`
- `docs(contributing): update setup instructions`

## 📝 How to Update Commands and Services Documentation

### Services Documentation

**To update service documentation:**

1. **Edit service YAML files** in `internal/config/services/{category}/service-name.yaml`
2. **Add documentation section** to your service YAML:
   ```yaml
   documentation:
     examples:
       - service-cli --host localhost --port 5432 status
       - curl http://localhost:8080/health
     usage_notes: Brief usage notes and configuration tips
     links:
       - https://service-docs.example.com
     use_cases:
       - Primary application database
       - Development testing
       - Data persistence
   ```
3. **Update configuration schema** with proper JSON Schema format:
   ```yaml
   configuration_schema:
     type: object
     properties:
       password:
         type: string
         default: "password"
         description: Service password
   ```
4. **Regenerate documentation**: `task docs`

### CLI Reference Documentation

**CLI reference is auto-generated from:**
- Commands configuration in `internal/config/commands.yaml`
- Built-in command help text and descriptions
- Command usage patterns and examples

**To update CLI documentation:**
1. Edit command definitions in Go code
2. Update `internal/config/commands.yaml` if needed
3. Run `task docs` to regenerate

### Template Modifications

**To modify documentation templates:**
1. Edit templates in `docs-site/templates/`
2. Test changes with `task docs-serve`
3. Ensure templates remain in `.prettierignore`

**Important Notes:**
- **Never manually edit** `docs-site/content/cli-reference.md` or `docs-site/content/services.md`
- These files are **always regenerated** from source configurations
- Use `task docs-serve` to preview changes locally
- Automate with pre-commit hooks: `task docs` in your workflow

## 🤖 GitHub Workflows & CI/CD

The project uses GitHub Actions for continuous integration:

- **Build and Test**: Runs on every push and PR
- **Security Scanning**: CodeQL analysis and dependency checks
- **Documentation**: Builds and deploys Hugo site
- **Release**: Automated releases with release-please

## 📞 Getting Help

### Development Questions

- **GitHub Issues**: For general questions and ideas
- **Issues**: For bugs and feature requests
- **Documentation**: Check the comprehensive docs in `docs-site/`

### Code Review

All contributions go through code review. We look for:
- Code quality and style
- Test coverage
- Documentation updates
- Breaking change considerations

## 🔄 Maintenance

### Regular Maintenance Tasks

- Update dependencies regularly
- Review and update documentation
- Monitor security vulnerabilities
- Performance optimization

## 📚 See Also

- **Existing services**: Use as implementation examples
- **Framework scripts**: Study build and management scripts for patterns
- **Documentation**: Hugo-based documentation site in `docs-site/` directory

We appreciate your contributions! 🙏

---
*This file syncs to docs-site/content/contributing.md automatically via `otto-stack docs`*

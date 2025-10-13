# Contributing to otto-stack

Thank you for your interest in contributing to otto-stack! This guide will help you get started with development and explain how to contribute effectively.

## 📝 Quick Contributor Checklist

- [ ] Fork and clone the repository
- [ ] Install dependencies: `task setup`
- [ ] Build the project: `task build`
- [ ] Run tests: `task test`
- [ ] Edit YAML manifests (`internal/config/commands.yaml`, `internal/config/services/services.yaml`) for changes
- [ ] Run `otto-stack docs` to update documentation
- [ ] Commit both the manifest and generated docs
- [ ] Follow the contributing guide and PR template

## 📋 Overview

We welcome contributions from the development community! Whether you're fixing bugs, adding new services, improving documentation, or enhancing existing features, your contributions help make the framework better for everyone.

## 🤝 Community & Support

### GitHub Repository

- **Main Repository**: [otto-nation/otto-stack](https://github.com/otto-nation/otto-stack)
- **Issues**: [Report bugs and request features](https://github.com/otto-nation/otto-stack/issues)
- **Discussions**: [Ask questions and share ideas](https://github.com/otto-nation/otto-stack/discussions)
- **Releases**: [Latest versions and changelog](https://github.com/otto-nation/otto-stack/releases)

### Getting Help

**Before opening an issue:**
1. Check the [Troubleshooting Guide](docs-site/content/troubleshooting.md)
2. Review [existing issues](https://github.com/otto-nation/otto-stack/issues)
3. Search [discussions](https://github.com/otto-nation/otto-stack/discussions)

**For support requests:**
- Use [GitHub Discussions](https://github.com/otto-nation/otto-stack/discussions) for questions
- Check the [CLI Reference](docs-site/content/reference.md) for command help
- Run `otto-stack doctor` for system diagnostics

## 📚 Automated Documentation & YAML Manifests

This project uses automated documentation generation from YAML manifests. Here's how it works:

### Documentation Generation (Go-based)

Documentation for CLI commands and services is automatically generated from YAML manifests:

```bash
# Generate all documentation
otto-stack docs

# Generate only command reference
otto-stack docs --commands-only

# Generate only services guide
otto-stack docs --services-only

# Preview changes without writing files
otto-stack docs --dry-run
```

This will update `docs-site/content/reference.md` and `docs-site/content/services.md` based on the latest YAML manifests.

**Contributor Workflow Checklist:**
1. Ensure you have Go 1.21+ installed and the project built (`task build`).
2. Edit `internal/config/commands.yaml` and/or `internal/config/services/services.yaml` to add or update commands/services.
3. Run `otto-stack docs` to regenerate documentation from YAML manifests.
4. Commit both the manifest and the updated docs.
5. Never manually edit auto-generated docs (`docs-site/content/reference.md`, `docs-site/content/services.md`).
6. Optionally, set up CI or pre-commit hooks to automate doc generation.

Documentation for commands (`docs-site/content/reference.md`) and services (`docs-site/content/services.md`) is auto-generated from these manifests using the native Go `otto-stack docs` command.

## 🚀 Getting Started

### Prerequisites

- **Go**: 1.21+ (managed via `.go-version` file)
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
├── internal/                     # Internal Go packages
│   ├── cli/                      # CLI command handlers
│   ├── pkg/                      # Reusable packages
│   └── ...
├── services/                     # Service definitions and configs
│   ├── postgres/                 # PostgreSQL service
│   ├── redis/                    # Redis service
│   ├── prometheus/               # Prometheus service
│   ├── kafka/                    # Kafka service
│   └── services.yaml             # YAML manifest for all services
├── scripts/                      # Build and utility scripts
│   └── commands.yaml             # YAML manifest for all commands
├── docs-site/                    # Hugo documentation site
│   ├── content/                  # Markdown content files
│   ├── config/                   # Hugo configuration
│   └── themes/                   # Hugo themes
└── otto-stack-config.sample.yaml # Sample configuration
```

## 🛠️ Adding New Services

### Step 1: Create Service Directory

Create a new directory under `services/` for your service:

```bash
mkdir -p services/my-service
```

### Step 2: Create Service Metadata

Create `services/my-service/service.yaml`:

```yaml
name: my-service
description: "A brief description of what this service does"
category: database  # database|cache|observability|messaging|cloud-services
version: "1.0.0"
maintainer: "Your Name <your.email@example.com>"

image:
  name: my-service
  tag: latest
  registry: docker.io

defaults:
  port: 8080
  memory_limit: 256m
  cpu_limit: 0.5
  restart_policy: unless-stopped
  environment:
    MY_SERVICE_HOST: localhost
    MY_SERVICE_PORT: 8080
    MY_SERVICE_PASSWORD: changeme

overrides:
  port:
    type: integer
    default: 8080
    description: "Port for the service to listen on"
  memory_limit:
    type: string
    default: 256m
    description: "Memory limit for the container"
  password:
    type: string
    default: changeme
    description: "Password for service authentication"

health_check:
  enabled: true
  command: "curl -f http://localhost:8080/health || exit 1"
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s

dependencies: []

spring_boot:
  enabled: true
  config_template: my-service-config.yml
  dependencies:
    - spring-boot-starter-web
    - spring-boot-starter-actuator

tags:
  - database
  - sql
  - development
```

### Step 3: Create Docker Compose Definition

Create `services/my-service/docker-compose.yml`:

```yaml
services:
  my-service:
    image: ${MY_SERVICE_IMAGE:-my-service:latest}
    container_name: ${PROJECT_NAME}_my-service
    ports:
      - "${MY_SERVICE_PORT:-8080}:8080"
    environment:
      - MY_SERVICE_PASSWORD=${MY_SERVICE_PASSWORD:-changeme}
    volumes:
      - my-service-data:/data
      - ${MY_SERVICE_CONFIG_FILE:-./config/my-service.conf}:/etc/my-service/my-service.conf:ro
    networks:
      - otto-stack-framework
    restart: ${MY_SERVICE_RESTART_POLICY:-unless-stopped}
    mem_limit: ${MY_SERVICE_MEMORY_LIMIT:-256m}
    cpus: ${MY_SERVICE_CPU_LIMIT:-0.5}
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

volumes:
  my-service-data:
    driver: local

networks:
  otto-stack-framework:
    external: true
```

### Step 4: Update Framework Configuration

Add your service to `internal/config/services/services.yaml`:

```yaml
services:
  my-service:
    enabled: true
    category: database
    description: "My custom service for development"
    ports: [8080]
    dependencies: []
    resource_tier: small
```

### Step 5: Test Your Service

```bash
# Build the project
task build

# Generate updated documentation
./build/otto-stack docs

# Test the service
./build/otto-stack start my-service
./build/otto-stack status my-service
./build/otto-stack stop my-service
```

## 📝 Documentation Contributions

### Adding Documentation

1. **Service documentation**: Document new services in `docs-site/content/services.md`
2. **Configuration examples**: Add examples to `docs-site/content/configuration.md`
3. **Troubleshooting**: Add common issues to `docs-site/content/troubleshooting.md`
4. **Quick reference**: Update command references in `docs-site/content/reference.md`

### Documentation Standards

- Use clear, concise language
- Include code examples where applicable
- Follow the existing documentation structure
- Use proper Markdown formatting
- Test all code examples before submitting

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

## 🎯 Coding Standards

### Go Standards

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Handle errors appropriately
- Write unit tests for new functionality

### YAML Standards

- Use 2-space indentation
- Quote string values when necessary
- Use descriptive keys and comments
- Validate YAML syntax before committing

### Documentation Standards

- Use Markdown format
- Include code examples
- Keep language clear and concise
- Update table of contents when adding sections

## 📝 How to Update Commands and Services Documentation

**Note:**
Do not manually edit `docs-site/content/reference.md` or `docs-site/content/services.md`—these files are always generated from the manifests.

**Tip:**
You can automate this process with a pre-commit hook or CI workflow using `./otto-stack docs`.

## 🤖 GitHub Workflows & CI/CD

The project uses GitHub Actions for continuous integration:

- **Build and Test**: Runs on every push and PR
- **Security Scanning**: CodeQL analysis and dependency checks
- **Documentation**: Builds and deploys Hugo site
- **Release**: Automated releases with release-please

## 📞 Getting Help

### Development Questions

- **GitHub Discussions**: For general questions and ideas
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

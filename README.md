# otto-stack

A powerful development stack management tool built in Go for streamlined local development automation

## What is otto-stack?

**otto-stack** is a modern CLI tool that provides quick setup, Docker integration, configuration management, and built-in monitoring for development environments.

## Quick Start

### Installation

#### Quick Install (Recommended)

```bash
# Install latest version
curl -fsSL https://raw.githubusercontent.com/otto-nation/otto-stack/main/install.sh | bash
```

#### Manual Installation

**macOS/Linux:**
```bash
# Download the latest release
curl -L -o otto-stack "https://github.com/otto-nation/otto-stack/releases/latest/download/otto-stack-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)"
chmod +x otto-stack
sudo mv otto-stack /usr/local/bin/
```

**Windows:**
```powershell
# Download from releases page
# https://github.com/otto-nation/otto-stack/releases/latest
```

#### Package Managers

**Homebrew (macOS):**
```bash
# Coming soon
brew install isaacgarza/tap/otto-stack
```

### Basic Usage

```bash
# Initialize a new project
otto-stack init

# Start your development stack
otto-stack up
```

## Key Features

- **Project Templates**: Go, Node.js, Python, and full-stack setups
- **Service Management**: Databases, message queues, monitoring tools
- **Health Monitoring**: Built-in health checks and status monitoring
- **Docker Integration**: Seamless container management

## Documentation

- **[Setup & Installation](docs-site/content/setup.md)**
- **[Usage Guide](docs-site/content/usage.md)**
- **[Services Guide](docs-site/content/services.md)**
- **[Configuration](docs-site/content/configuration.md)**
- **[CLI Reference](docs-site/content/reference.md)**
- **[Contributing](docs-site/content/contributing.md)**

## Get Started

1. **[Complete installation guide](docs-site/content/setup.md)**
2. **[Learn basic usage](docs-site/content/usage.md)**
3. **[Explore available services](docs-site/content/services.md)**

## Contributing

We welcome contributions! Please see our [Contributing Guide](docs-site/content/contributing.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 📖 [Documentation](docs-site/)
- 🐛 [Issues](https://github.com/otto-nation/otto-stack/issues)
- 💬 [Discussions](https://github.com/otto-nation/otto-stack/discussions)

---

> **Built with ❤️ by the otto-stack team**
> Making local development environments simple, consistent, and powerful.

# otto-stack

A powerful development stack management tool built in Go for streamlined local development automation

📚 **[Documentation](https://otto-nation.github.io/otto-stack/)**

## What is otto-stack?

**otto-stack** is a modern CLI tool that provides quick setup, Docker integration, configuration management, and built-in monitoring for development environments.

## Quick Start

### Installation

#### Homebrew (Recommended)

```bash
brew install otto-nation/tap/otto-stack
```

#### Script Install

```bash
# Install latest version
curl -fsSL https://raw.githubusercontent.com/otto-nation/otto-stack/main/scripts/install.sh | bash

# Install to custom directory
curl -fsSL https://raw.githubusercontent.com/otto-nation/otto-stack/main/scripts/install.sh | bash -s -- --dir ~/.local/bin
```

#### Manual Install

```bash
# Download the latest release
curl -L -o otto-stack "https://github.com/otto-nation/otto-stack/releases/latest/download/otto-stack-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)"
chmod +x otto-stack
sudo mv otto-stack /usr/local/bin/
```

#### Uninstall

```bash
# Homebrew
brew uninstall otto-stack

# Script install
curl -fsSL https://raw.githubusercontent.com/otto-nation/otto-stack/main/scripts/uninstall.sh | bash
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
- **[Services Guide](docs-site/content/services.md)**
- **[Configuration](docs-site/content/configuration.md)**
- **[CLI Reference](docs-site/content/cli-reference.md)**
- **[Scripts Reference](docs-site/content/scripts.md)**
- **[Contributing](docs-site/content/contributing.md)**

## Get Started

1. **[Complete installation guide](docs-site/content/setup.md)**
2. **[Explore available services](docs-site/content/services.md)**

## Git Workflow

For AI-powered git automation, install [otto-workbench](https://github.com/otto-nation/otto-workbench):

```bash
git clone https://github.com/otto-nation/otto-workbench ~/workbench
cd ~/workbench && ./install.sh
```

Then use from any project:
```bash
task --global commit      # AI-generated commit messages
task --global create-pr   # AI-generated pull requests
task --global update-pr   # Update PR descriptions
```

**otto-workbench** also provides shell aliases, development utilities, and AWS/Kubernetes helpers.

## Contributing

We welcome contributions! Please see our [Contributing Guide](docs-site/content/contributing.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 📖 [Documentation](docs-site/)
- 🐛 [Issues](https://github.com/otto-nation/otto-stack/issues)

---

> **Built with ❤️ by the otto-stack team**
> Making local development environments simple, consistent, and powerful.

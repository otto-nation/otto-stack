---
title: "otto-stack"
description: "A powerful development stack management tool built in Go for streamlined local development automation"
lead: "Streamline your local development with powerful CLI tools and automated service management"
date: "2025-10-01"
lastmod: "2025-10-14"
draft: false
weight: 50
toc: true
---

# otto-stack

A powerful development stack management tool built in Go for streamlined local development automation

## What is otto-stack?

**otto-stack** is a modern CLI tool that provides quick setup, Docker integration, configuration management, and built-in monitoring for development environments.

## Quick Start

### Installation

#### Quick Install (Recommended)

```bash
# Install latest version
curl -fsSL https://raw.githubusercontent.com/otto-nation/otto-stack/main/scripts/install.sh | bash
```

#### Custom Installation

```bash
# Install to custom directory
curl -fsSL https://raw.githubusercontent.com/otto-nation/otto-stack/main/scripts/install.sh | bash -s -- --dir ~/.local/bin

# Or download and run locally
curl -fsSL -o install.sh https://raw.githubusercontent.com/otto-nation/otto-stack/main/scripts/install.sh
chmod +x install.sh
./install.sh --help
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

#### Uninstallation

```bash
# Remove otto-stack from your system
curl -fsSL https://raw.githubusercontent.com/otto-nation/otto-stack/main/scripts/uninstall.sh | bash
```

#### Package Managers

**Homebrew (macOS):**

```bash
# Coming soon
brew install otto-nation/tap/otto-stack
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

- **[Setup & Installation](setup.md)**
- **[Usage Guide](usage.md)**
- **[Services Guide](services.md)**
- **[Configuration](configuration.md)**
- **[CLI Reference](reference.md)**
- **[Scripts Reference](scripts.md)**
- **[Contributing](contributing.md)**

## Get Started

1. **[Complete installation guide](setup.md)**
2. **[Learn basic usage](usage.md)**
3. **[Explore available services](services.md)**

## Contributing

We welcome contributions! Please see our [Contributing Guide](contributing.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/otto-nation/otto-stack/blob/main/LICENSE) file for details.

## Support

- 📖 [Documentation](/)
- 🐛 [Issues](https://github.com/otto-nation/otto-stack/issues)

---

> **Built with ❤️ by the otto-stack team**
> Making local development environments simple, consistent, and powerful.

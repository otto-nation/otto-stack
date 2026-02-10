---
title: "Setup"
description: "Installation guide for otto-stack"
lead: "Get otto-stack installed and running on your system"
date: "2025-10-01"
lastmod: "2025-11-10"
draft: false
weight: 10
toc: true
---

# Setup & Installation

## 📦 Installation

### Quick Install (Recommended)

```bash
# Install latest version
curl -fsSL https://raw.githubusercontent.com/otto-nation/otto-stack/main/scripts/install.sh | bash
```

### Custom Installation

```bash
# Install to custom directory
curl -fsSL https://raw.githubusercontent.com/otto-nation/otto-stack/main/scripts/install.sh | bash -s -- --dir ~/.local/bin
```

### Manual Installation

**macOS/Linux:**

```bash
# Download the latest release
curl -L -o otto-stack "https://github.com/otto-nation/otto-stack/releases/latest/download/otto-stack-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)"
chmod +x otto-stack
sudo mv otto-stack /usr/local/bin/
```

**Windows:**
Download from [releases page](https://github.com/otto-nation/otto-stack/releases/latest)

### Package Managers

**Homebrew (macOS):**

```bash
# Coming soon
brew install otto-nation/tap/otto-stack
```

## 🐳 Prerequisites

### Docker

**macOS:**

```bash
# Option 1: Docker Desktop
# Download from https://docker.com

# Option 2: Colima (lightweight)
brew install colima docker
colima start
```

**Linux:**

```bash
# Install Docker Engine
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
```

**Windows:**

- Install [Docker Desktop](https://docker.com)
- Enable WSL2 integration

### System Requirements

- **RAM**: 4GB minimum, 8GB recommended
- **Disk**: 2GB free space
- **CPU**: Any modern processor

## 🚀 Quick Start

```bash
# Verify installation
otto-stack version

# Initialize new project
otto-stack init

# Start development environment
otto-stack up

# Check status
otto-stack status
```

## ⚙️ Configuration

Otto-stack creates configuration files during initialization:

**Main config** (`.otto-stack/config.yaml`):

```yaml
project:
  name: my-app
  type: docker

stack:
  enabled:
    - postgres
    - redis
```

**Environment variables** (`.env.generated`):

Shows available configuration options with defaults:

```bash
POSTGRES_DB=${POSTGRES_DB:-local_dev}
POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-password}
REDIS_PASSWORD=${REDIS_PASSWORD:-password}
```

Create a `.env` file to customize:

```bash
POSTGRES_DB=my_custom_db
POSTGRES_PASSWORD=secure_password
```

See [Configuration Guide](/configuration) for details.

## 🔧 Verification

```bash
# Check system health
otto-stack doctor

# Test Docker connectivity
docker info

# Verify services start
otto-stack up
otto-stack status
```

## 🧹 Uninstallation

```bash
# Remove otto-stack from your system
curl -fsSL https://raw.githubusercontent.com/otto-nation/otto-stack/main/scripts/uninstall.sh | bash
```

## 📚 Next Steps

- **[Configuration](/configuration)** - Configure your development stack
- **[Services Guide](/services)** - Available services and options
- **[CLI Reference](/cli-reference)** - Complete command reference
- **[Troubleshooting](/troubleshooting)** - Common issues and solutions

## 🆘 Need Help?

- **Issues**: [GitHub Issues](https://github.com/otto-nation/otto-stack/issues)
- **Questions**: [GitHub Issues](https://github.com/otto-nation/otto-stack/issues)
- **Documentation**: [Full Documentation](/)

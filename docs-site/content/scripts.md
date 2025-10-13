---
title: "Scripts Reference"
description: "Development and installation scripts for otto-stack"
date: 2024-10-13
lastmod: 2024-10-13
draft: false
weight: 80
---

# Scripts Reference

This reference covers all utility scripts available in the otto-stack project for installation, development, and maintenance.

## Installation Scripts

### install.sh

Universal installation script that downloads and installs the latest otto-stack release from GitHub.

**Usage:**

```bash
# Install to default location (/usr/local/bin)
curl -fsSL https://raw.githubusercontent.com/otto-nation/otto-stack/main/scripts/install.sh | bash

# Install to custom directory
./scripts/install.sh --dir ~/.local/bin

# Install using environment variable
INSTALL_DIR=~/bin ./scripts/install.sh
```

**Options:**

- `-d, --dir DIR` - Installation directory (default: /usr/local/bin)
- `-h, --help` - Show help message
- `-v, --version` - Show version

**Environment Variables:**

- `INSTALL_DIR` - Installation directory (overrides --dir)

### uninstall.sh

Removes otto-stack installations from the system.

**Usage:**

```bash
# Interactive uninstall
./scripts/uninstall.sh

# Force uninstall without confirmation
./scripts/uninstall.sh --force
```

**Options:**

- `-f, --force` - Force removal without confirmation
- `-h, --help` - Show help message
- `-v, --version` - Show version

## Development Scripts

### common.sh

Shared utilities and functions used by other scripts. Not meant to be run directly.

**Provides:**

- Colored output functions (`print_status`, `print_success`, `print_error`, etc.)
- Platform detection (`detect_platform`)
- Dependency checking (`check_dependencies`)
- File operations (`install_file`, `verify_executable`)
- GitHub API interactions (`get_latest_version`)

### setup-hooks.sh

Sets up Git hooks for the development workflow.

**Usage:**

```bash
./scripts/setup-hooks.sh
```

**What it does:**

- Installs pre-commit and pre-push Git hooks
- Configures automatic code quality checks
- Sets up commit message validation

### sync-go-version.sh

Synchronizes Go version across different configuration files in the project.

**Usage:**

```bash
./scripts/sync-go-version.sh
```

**Files updated:**

- `.go-version`
- `go.mod`
- GitHub Actions workflows
- Docker configurations

### get-go-version.sh

Retrieves the current Go version from project configuration.

**Usage:**

```bash
./scripts/get-go-version.sh
```

## Build & Release Scripts

### generate-release-configs.sh

Generates release configuration files from central configuration.

**Usage:**

```bash
./scripts/generate-release-configs.sh
```

### test-cross-platform.sh

Tests builds across different platforms to ensure compatibility.

**Usage:**

```bash
./scripts/test-cross-platform.sh
```

**Platforms tested:**

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

### validate-project.sh

Validates project structure and configuration files.

**Usage:**

```bash
./scripts/validate-project.sh
```

**Checks:**

- Project structure integrity
- Configuration file validity
- Documentation consistency
- Build requirements

## Security & Quality Scripts

### run-security-scan.sh

Runs comprehensive security scans on the codebase.

**Usage:**

```bash
./scripts/run-security-scan.sh
```

**Scans performed:**

- Dependency vulnerability scanning
- Static code analysis
- Secret detection
- License compliance

### update-homebrew-formula.sh

Updates the Homebrew formula with the latest release checksums.

**Usage:**

```bash
# Update with latest release
./scripts/update-homebrew-formula.sh

# Update with specific version
./scripts/update-homebrew-formula.sh --version v1.2.3

# Validate existing formula
./scripts/update-homebrew-formula.sh --validate-only
```

**Options:**

- `-v, --version VER` - Use specific version instead of latest
- `--validate-only` - Only validate existing formula syntax
- `-h, --help` - Show help message

### update-docs-lastmod.sh

Updates documentation modification timestamps.

**Usage:**

```bash
./scripts/update-docs-lastmod.sh
```

**Usage:**

```bash
./scripts/update-docs-lastmod.sh
```

## Running Scripts

All scripts are designed to be run from the project root directory:

```bash
# Make scripts executable (if needed)
chmod +x scripts/*.sh

# Run any script
./scripts/script-name.sh [options]
```

## Dependencies

Most scripts require:

- `bash` (version 4.0+)
- `curl` (for GitHub API interactions)
- `git` (for repository operations)

Additional dependencies are checked by individual scripts as needed.

## Script Development

When creating new scripts:

1. Use the common utilities from `scripts/common.sh`
2. Follow the established patterns for error handling
3. Include proper help text and option parsing
4. Add appropriate documentation to this page
5. Make scripts executable: `chmod +x scripts/new-script.sh`

## Troubleshooting

### Permission Denied

If you get permission errors:

```bash
chmod +x scripts/*.sh
```

### Command Not Found

Ensure you're running scripts from the project root:

```bash
cd /path/to/otto-stack
./scripts/script-name.sh
```

### Missing Dependencies

Install required tools:

```bash
# macOS
brew install curl git

# Ubuntu/Debian
sudo apt-get install curl git

# CentOS/RHEL
sudo yum install curl git
```

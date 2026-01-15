# Scripts Reference

Development and maintenance scripts for otto-stack contributors.

## Installation Scripts

### install.sh

Downloads and installs the latest otto-stack release from GitHub.

```bash
# Install to default location (/usr/local/bin)
curl -fsSL https://raw.githubusercontent.com/otto-nation/otto-stack/main/scripts/install.sh | bash

# Install to custom directory
./scripts/install.sh --dir ~/.local/bin
```

**Options:**
- `-d, --dir DIR` - Installation directory (default: /usr/local/bin)
- `-h, --help` - Show help message

### uninstall.sh

Removes otto-stack installations from the system.

```bash
# Interactive uninstall
./scripts/uninstall.sh

# Force uninstall without confirmation
./scripts/uninstall.sh --force
```

## Development Scripts

### common.sh

Shared utilities and functions used by other scripts. Not meant to be run directly.

**Provides:**
- Colored output functions (`print_status`, `print_success`, `print_error`)
- Platform detection (`detect_platform`)
- Dependency checking (`check_dependencies`)
- File operations (`install_file`, `verify_executable`)
- GitHub API interactions (`get_latest_version`)
- AI helper functions (`check_ai_available`, `get_ai_response`)

### setup-hooks.sh

Sets up Git hooks for the development workflow.

```bash
./scripts/setup-hooks.sh
```

Installs pre-commit and pre-push hooks for automatic code quality checks.

### go-version.sh

Unified Go version management script.

**Get version:**

```bash
./scripts/go-version.sh get [--major-minor|--full|--env|--github-matrix]
```

**Sync version across files:**

```bash
./scripts/go-version.sh sync [--check|--fix]
```

Synchronizes Go version across `.go-version`, `go.mod`, and `.golangci.yml`.

See [go-version-management.md](go-version-management.md) for details.

### homebrew.sh

Unified Homebrew formula management script.

**Update formula with checksums:**

```bash
./scripts/homebrew.sh update -v v1.2.3
```

**Deploy formula to tap repository:**

```bash
./scripts/homebrew.sh deploy
```

See [HOMEBREW_TAP.md](HOMEBREW_TAP.md) for details.

## Build & Release Scripts

### generate-release-configs.sh

Generates release configuration files from central configuration.

```bash
./scripts/generate-release-configs.sh
```

### validate-project.sh

Validates project structure and configuration files.

```bash
./scripts/validate-project.sh
```

Checks project structure integrity, configuration validity, and documentation consistency.

### update-docs-lastmod.sh

Updates documentation modification timestamps.

```bash
./scripts/update-docs-lastmod.sh
```

## Running Scripts

All scripts should be run from the project root:

```bash
./scripts/script-name.sh [options]
```

## Dependencies

Most scripts require:
- `bash` (version 4.0+)
- `curl` (for GitHub API interactions)
- `git` (for repository operations)

Additional dependencies are checked by individual scripts as needed.

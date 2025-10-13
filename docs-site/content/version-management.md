---
title: "Version Management System"
description: "Multi-version support and automatic version switching for otto-stack"
lead: "Install, manage, and switch between different otto-stack versions"
date: "2025-10-01"
lastmod: "2025-10-11"
draft: false
weight: 75
toc: true
---

# Version Management System

The otto-stack version management system allows you to install, manage, and automatically switch between different versions of otto-stack based on project requirements.

## Quick Start

```bash
# Set version requirement for current project
otto-stack versions set ">=1.0.0"

# Install a specific version
otto-stack versions install 1.2.3

# List installed versions
otto-stack versions list

# Automatically use the right version for your project
otto-stack up  # Uses version based on .otto-stack-version file
```

## Core Concepts

### Version Files

Projects specify their otto-stack version requirements using `.otto-stack-version` files:

**Simple text format:**

```
1.2.3
```

**Version constraint:**

```
>=1.0.0
```

**YAML format with metadata:**

```yaml
version: "^1.2.0"
metadata:
  created_by: "otto-stack"
  project: "my-awesome-app"
```

### Automatic Version Detection

otto-stack automatically detects version requirements by searching for version files in:

1. Current directory (`.`)
2. `.otto-stack/` subdirectory
3. `.config/` subdirectory
4. `config/` subdirectory

Supported file names:

- `.otto-stack-version`
- `.otto-stack-version.yaml`
- `.otto-stack-version.yml`
- `otto-stack-version`
- `otto-stack-version.yaml`
- `otto-stack-version.yml`

### Project Root Detection

When determining which version to use, otto-stack finds the project root by looking for:

- Version files (`.otto-stack-version`)
- Git repositories (`.git` directory)
- Common project files (`go.mod`, `package.json`, `Cargo.toml`, etc.)

## Version Constraints

otto-stack supports semantic versioning constraints:

| Constraint | Description           | Example                              |
| ---------- | --------------------- | ------------------------------------ |
| `1.2.3`    | Exact version         | Must be exactly 1.2.3                |
| `>=1.2.3`  | Greater than or equal | 1.2.3, 1.2.4, 1.3.0, 2.0.0           |
| `>1.2.3`   | Greater than          | 1.2.4, 1.3.0, 2.0.0                  |
| `<=1.2.3`  | Less than or equal    | 1.0.0, 1.2.2, 1.2.3                  |
| `<1.2.3`   | Less than             | 1.0.0, 1.2.2                         |
| `~1.2.3`   | Tilde (patch changes) | 1.2.3, 1.2.4, 1.2.10 (but not 1.3.0) |
| `^1.2.3`   | Caret (minor changes) | 1.2.3, 1.3.0, 1.9.9 (but not 2.0.0)  |
| `*`        | Any version           | Any available version                |

## Commands

### `versions list`

List all installed versions of otto-stack.

```bash
otto-stack versions list

# Example output:
VERSION  ACTIVE  INSTALLED   SOURCE
1.1.0             2024-01-15  github
1.2.0    *        2024-02-01  github
1.2.3             2024-02-15  github
```

**Flags:**

- `--json` - Output in JSON format

### `versions install`

Install a specific version of otto-stack.

```bash
# Install specific version
otto-stack versions install 1.2.3

# Install latest version
otto-stack versions install latest
```

The command will:

1. Download the binary from GitHub releases
2. Verify checksums (if available)
3. Extract and install to version-specific directory
4. Register the version in the local registry

### `versions uninstall`

Remove a specific version of otto-stack.

```bash
otto-stack versions uninstall 1.2.3
```

**Note:** Cannot uninstall the currently active version.

### `versions use`

Set the global default version of otto-stack.

```bash
otto-stack versions use 1.2.3
```

This sets the version to use when no project-specific version is found.

### `versions available`

List all available versions from GitHub releases.

```bash
otto-stack versions available

# Limit results
otto-stack versions available --limit 10

# JSON output
otto-stack versions available --json
```

### `versions detect`

Detect version requirements for a project.

```bash
# Detect in current directory
otto-stack versions detect

# Detect in specific directory
otto-stack versions detect /path/to/project
```

**Example output:**

```
Project: /home/user/my-project
Required version: >=1.0.0
Resolved to installed version: 1.2.3
```

**Flags:**

- `--json` - Output in JSON format

### `versions set`

Set version requirement for a project.

```bash
# Set for current directory
otto-stack versions set ">=1.0.0"

# Set for specific directory
otto-stack versions set "^1.2.0" /path/to/project

# Use YAML format
otto-stack versions set "1.2.3" --format yaml
```

**Flags:**

- `--format` - File format (`text` or `yaml`)

### `versions cleanup`

Clean up old versions to save disk space.

```bash
# Keep 3 most recent versions
otto-stack versions cleanup --keep 3

# Dry run to see what would be removed
otto-stack versions cleanup --dry-run
```

## Multi-Project Workflow

### Example Setup

```bash
# Project A needs otto-stack 1.1.x
cd project-a
otto-stack versions set "~1.1.0"

# Project B needs otto-stack 1.2.x or higher
cd project-b
otto-stack versions set ">=1.2.0"

# Install required versions
otto-stack versions install 1.1.5
otto-stack versions install 1.2.3
```

### Automatic Switching

Once versions are installed and requirements are set:

```bash
# Automatically uses 1.1.5
cd project-a
otto-stack up

# Automatically uses 1.2.3
cd project-b
otto-stack up
```

### Version Resolution

When you run a otto-stack command, the system:

1. **Finds project root** by searching upward for version files or project markers
2. **Reads version constraint** from `.otto-stack-version` file
3. **Resolves to best match** among installed versions
4. **Delegates execution** to the correct version binary
5. **Falls back** to active version if no constraint found

## Storage Layout

Versions are stored in your home directory:

```
~/.otto-stack/
├── versions/                 # Installed version binaries
│   ├── 1.1.0/
│   │   └── otto-stack
│   ├── 1.2.0/
│   │   └── otto-stack
│   └── 1.2.3/
│       └── otto-stack
└── ...

~/.config/otto-stack/
├── installed_versions.json   # Registry of installed versions
└── project_configs.json     # Per-project configurations
```

## Advanced Usage

### CI/CD Integration

Pin specific versions in CI environments:

```bash
# In your CI script
echo "1.2.3" > .otto-stack-version
otto-stack up  # Will use exactly 1.2.3
```

### Team Consistency

Commit `.otto-stack-version` files to ensure team consistency:

```bash
# Set project requirement
otto-stack versions set "^1.2.0"

# Commit to repo
git add .otto-stack-version
git commit -m "Pin otto-stack version to ^1.2.0"
```

### Version Verification

Verify installed versions and project requirements:

```bash
# Check what version would be used
otto-stack versions detect

# List all installed versions
otto-stack versions list

# Check available updates
otto-stack versions available | head -5
```

### Cleanup Strategy

Regular maintenance to manage disk usage:

```bash
# Keep only 3 most recent versions
otto-stack versions cleanup --keep 3

# Preview cleanup without making changes
otto-stack versions cleanup --dry-run
```

## Troubleshooting

### Common Issues

**No compatible version found:**

```bash
Error: No installed version satisfies requirement: >=1.3.0
Run 'otto-stack versions install 1.3.0' to install a compatible version.
```

**Solution:** Install a compatible version:

```bash
otto-stack versions install 1.3.0
# or install latest
otto-stack versions install latest
```

**Version file not detected:**

```bash
otto-stack versions detect
# Output: No specific version requirement found
```

**Solution:** Create a version requirement:

```bash
otto-stack versions set ">=1.0.0"
```

**Binary delegation fails:**

- Check file permissions in `~/.otto-stack/versions/`
- Verify version registry: `otto-stack versions list`
- Reinstall problematic version: `otto-stack versions install X.Y.Z`

### Debug Information

Get detailed information about version resolution:

```bash
# Enable verbose output
otto-stack --verbose versions detect

# Check version manager state
ls -la ~/.otto-stack/versions/
cat ~/.config/otto-stack/installed_versions.json
```

## Best Practices

1. **Pin versions in production**: Use exact versions for deployments
2. **Use constraints for development**: Allow flexibility with `^` or `~`
3. **Commit version files**: Ensure team consistency
4. **Regular cleanup**: Manage disk space with periodic cleanup
5. **Document requirements**: Include version info in project README
6. **Test version changes**: Verify compatibility when updating constraints

## Security

- **Checksum verification**: Downloads are verified against GitHub release checksums when available
- **HTTPS downloads**: All downloads use secure HTTPS connections
- **No auto-execution**: Downloaded binaries are not executed during installation
- **User-scoped storage**: Versions are stored in user directories, not system-wide

## Migration from Manual Management

If you're currently managing otto-stack versions manually:

1. **Install version manager**: Upgrade to a version with version management
2. **Set project requirements**: Add `.otto-stack-version` files to your projects
3. **Install needed versions**: Use `otto-stack versions install` for required versions
4. **Remove manual installations**: Clean up old manual installations
5. **Update documentation**: Update team documentation with new workflow

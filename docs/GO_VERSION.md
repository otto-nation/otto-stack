# Go Version Management

otto-stack uses a centralized Go version management system to ensure consistency across all project files and CI/CD workflows.

## Single Source of Truth

The `.go-version` file at the project root is the single source of truth for the Go version:

```bash
# .go-version contains just the version number
1.24.11
```

## Automated Synchronization

The `scripts/go-version.sh` script synchronizes the Go version across:

- `go.mod` - Go module definition
- `.golangci.yml` - Linting configuration

### Get Version

```bash
# Get the raw version
./scripts/go-version.sh get                 # Output: 1.24.11

# Get major.minor only
./scripts/go-version.sh get --major-minor   # Output: 1.24

# Get as environment variable
./scripts/go-version.sh get --env           # Output: GO_VERSION=1.24.11

# Get as GitHub Actions matrix
./scripts/go-version.sh get --github-matrix # Output: ['1.24']
```

### Sync Version

```bash
# Check if all files are in sync
./scripts/go-version.sh sync --check

# Fix any version mismatches
./scripts/go-version.sh sync --fix
```

## GitHub Actions Integration

The `.github/actions/setup-go-version` composite action automatically:

1. Reads the Go version from `.go-version`
2. Sets up Go with the correct version
3. Configures module caching
4. Outputs the version for use in other steps

Usage in workflows:

```yaml
- name: Setup Go
  uses: ./.github/actions/setup-go-version
  id: setup-go

- name: Use Go version
  run: echo "Using Go ${{ steps.setup-go.outputs.go-version }}"
```

## Upgrading Go Version

To upgrade the Go version across the entire project:

1. Update `.go-version`:
   ```bash
   echo "1.24.11" > .go-version
   ```

2. Sync all configuration files:
   ```bash
   ./scripts/go-version.sh sync --fix
   ```

3. Verify the changes:
   ```bash
   ./scripts/go-version.sh sync --check
   ```

4. Commit the changes:
   ```bash
   git add .go-version go.mod .golangci.yml
   git commit -m "chore: upgrade Go to 1.24.11"
   ```

#!/bin/bash

# Sync Go version across all configuration files
# This script ensures all Go version references are consistent with .go-version

set -euo pipefail

# Get the directory of this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
GO_VERSION_FILE="$PROJECT_ROOT/.go-version"
GOLANGCI_CONFIG="$PROJECT_ROOT/.golangci.yml"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Function to show usage
show_usage() {
    cat << EOF
Usage: $0 [options]

Sync Go version across all configuration files using .go-version as source of truth.

Options:
    --check     Check if versions are in sync (exit 1 if not)
    --fix       Fix version mismatches (default action)
    --help      Show this help message

Files synchronized:
    - go.mod
    - .golangci.yml
    - Dockerfile (if exists)
    - CONTRIBUTING.md (if exists)
    - Any other Go version references

The script reads the version from .go-version and updates all configuration files
to use the same version.
EOF
}

# Function to check if .go-version file exists and is valid
validate_go_version_file() {
    if [[ ! -f "$GO_VERSION_FILE" ]]; then
        print_status "$RED" "Error: .go-version file not found at $GO_VERSION_FILE"
        exit 1
    fi

    local go_version=$(cat "$GO_VERSION_FILE" | tr -d '[:space:]')
    if [[ -z "$go_version" ]]; then
        print_status "$RED" "Error: .go-version file is empty"
        exit 1
    fi

    if [[ ! "$go_version" =~ ^[0-9]+\.[0-9]+(\.[0-9]+)?$ ]]; then
        print_status "$RED" "Error: Invalid Go version format in .go-version: $go_version"
        print_status "$RED" "Expected format: x.y or x.y.z (e.g., 1.21 or 1.21.0)"
        exit 1
    fi

    echo "$go_version"
}

# Function to update .golangci.yml
update_golangci_config() {
    local go_version=$1
    local check_only=${2:-false}

    if [[ ! -f "$GOLANGCI_CONFIG" ]]; then
        print_status "$YELLOW" "Warning: .golangci.yml not found, skipping"
        return 0
    fi

    # Check if the file has a go: field at all
    if ! grep -q -E '^\s*go:\s*"' "$GOLANGCI_CONFIG"; then
        # No go: field exists, skip this check as golangci-lint will use go.mod version
        return 0
    fi

    # Extract current version from golangci config
    local current_version=$(grep -E '^\s*go:\s*"' "$GOLANGCI_CONFIG" | sed 's/.*"\([^"]*\)".*/\1/' || echo "")

    if [[ "$current_version" == "$go_version" ]]; then
        print_status "$GREEN" "✓ .golangci.yml already has correct Go version: $go_version"
        return 0
    fi

    if [[ "$check_only" == "true" ]]; then
        print_status "$RED" "✗ .golangci.yml has wrong Go version: '$current_version' (expected: $go_version)"
        return 1
    fi

    # Update the version
    if command -v sed >/dev/null 2>&1; then
        # Use different sed syntax for macOS vs Linux
        if [[ "$OSTYPE" == "darwin"* ]]; then
            sed -i '' -E 's/go: "[^"]+"/go: "'"$go_version"'"/' "$GOLANGCI_CONFIG"
        else
            sed -i -E 's/go: "[^"]+"/go: "'"$go_version"'"/' "$GOLANGCI_CONFIG"
        fi
        print_status "$GREEN" "✓ Updated .golangci.yml Go version to: $go_version"
    else
        print_status "$RED" "Error: sed command not found"
        return 1
    fi
}

    # Function to update go.mod
    update_go_mod() {
        local go_version=$1
        local check_only=${2:-false}
        local go_mod="$PROJECT_ROOT/go.mod"

        if [[ ! -f "$go_mod" ]]; then
            print_status "$YELLOW" "Warning: go.mod not found, skipping"
            return 0
        fi

        # Extract current version from go.mod (full version including patch)
        local current_version=$(grep -E '^go [0-9]+\.[0-9]+' "$go_mod" | awk '{print $2}' || echo "")

        if [[ -z "$current_version" ]]; then
            print_status "$YELLOW" "Warning: No go directive found in go.mod"
            return 0
        fi

        if [[ "$current_version" == "$go_version" ]]; then
            print_status "$GREEN" "✓ go.mod already has correct Go version: $current_version"
            return 0
        fi

        if [[ "$check_only" == "true" ]]; then
            print_status "$RED" "✗ go.mod has wrong Go version: $current_version (expected: $go_version)"
            return 1
        fi

        # Update using go mod edit
        go mod edit -go="$go_version" "$go_mod"
        print_status "$GREEN" "✓ Updated go.mod Go version to: $go_version"
    }

# Function to update Dockerfile if it exists
update_dockerfile() {
    local go_version=$1
    local check_only=${2:-false}
    local dockerfile="$PROJECT_ROOT/Dockerfile"

    if [[ ! -f "$dockerfile" ]]; then
        return 0
    fi

    # Look for Go version in FROM statements
    local current_version=$(grep -E '^FROM.*golang:[0-9]+\.[0-9]+' "$dockerfile" | head -1 | sed -E 's/.*golang:([0-9]+\.[0-9]+).*/\1/' || echo "")

    if [[ -z "$current_version" ]]; then
        print_status "$BLUE" "ℹ Dockerfile found but no Go version detected in FROM statements"
        return 0
    fi

    # Extract major.minor version for comparison
    local version_major_minor=$(echo "$go_version" | sed -E 's/^([0-9]+\.[0-9]+).*/\1/')

    if [[ "$current_version" == "$version_major_minor" ]]; then
        print_status "$GREEN" "✓ Dockerfile already has correct Go version: $current_version"
        return 0
    fi

    if [[ "$check_only" == "true" ]]; then
        print_status "$RED" "✗ Dockerfile has wrong Go version: $current_version (expected: $version_major_minor)"
        return 1
    fi

    # Update the version
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' -E "s/(FROM.*golang:)[0-9]+\.[0-9]+/\1$version_major_minor/" "$dockerfile"
    else
        sed -i -E "s/(FROM.*golang:)[0-9]+\.[0-9]+/\1$version_major_minor/" "$dockerfile"
    fi
    print_status "$GREEN" "✓ Updated Dockerfile Go version to: $version_major_minor"
}

# Function to update CONTRIBUTING.md
update_contributing() {
    local go_version=$1
    local check_only=${2:-false}
    local contributing="$PROJECT_ROOT/CONTRIBUTING.md"

    if [[ ! -f "$contributing" ]]; then
        return 0
    fi

    # Look for Go version in prerequisites section (handles **Go**: format)
    local current_version=$(grep -E '\*\*Go\*\*:.*[0-9]+\.[0-9]+' "$contributing" | sed -E 's/.*([0-9]+\.[0-9]+)\+.*/\1/' || echo "")

    if [[ -z "$current_version" ]]; then
        print_status "$BLUE" "ℹ CONTRIBUTING.md found but no Go version detected"
        return 0
    fi

    # Extract major.minor version for comparison
    local version_major_minor=$(echo "$go_version" | sed -E 's/^([0-9]+\.[0-9]+).*/\1/')

    if [[ "$current_version" == "$version_major_minor" ]]; then
        print_status "$GREEN" "✓ CONTRIBUTING.md already has correct Go version: $current_version+"
        return 0
    fi

    if [[ "$check_only" == "true" ]]; then
        print_status "$RED" "✗ CONTRIBUTING.md has wrong Go version: $current_version+ (expected: $version_major_minor+)"
        return 1
    fi

    # Update the version (handles **Go**: format)
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' -E "s/(\*\*Go\*\*:.*)[0-9]+\.[0-9]+\+/\1$version_major_minor+/" "$contributing"
    else
        sed -i -E "s/(\*\*Go\*\*:.*)[0-9]+\.[0-9]+\+/\1$version_major_minor+/" "$contributing"
    fi
    print_status "$GREEN" "✓ Updated CONTRIBUTING.md Go version to: $version_major_minor+"
}

# Function to check all files
check_all_versions() {
    local go_version=$1
    local errors=0

    print_status "$BLUE" "Checking Go version consistency..."
    print_status "$BLUE" "Source version: $go_version"
    echo

    update_go_mod "$go_version" true || ((errors++))
    update_golangci_config "$go_version" true || ((errors++))
    update_dockerfile "$go_version" true || ((errors++))
    update_contributing "$go_version" true || ((errors++))

    echo
    if [[ $errors -eq 0 ]]; then
        print_status "$GREEN" "✅ All files have consistent Go versions"
        return 0
    else
        print_status "$RED" "❌ Found $errors version mismatch(es)"
        print_status "$YELLOW" "Run '$0 --fix' to fix the mismatches"
        return 1
    fi
}

# Function to fix all versions
fix_all_versions() {
    local go_version=$1
    local changes_made=false

    # Check if any changes are needed first
    local go_mod_needs_update=false
    local golangci_needs_update=false
    local dockerfile_needs_update=false
    local contributing_needs_update=false

    # Check go.mod
    if [[ -f "$PROJECT_ROOT/go.mod" ]]; then
        local current_go_version
        current_go_version=$(grep "^go " "$PROJECT_ROOT/go.mod" | awk '{print $2}' || echo "")
        if [[ "$current_go_version" != "$go_version" ]]; then
            go_mod_needs_update=true
        fi
    fi

    # Check .golangci.yml
    if [[ -f "$GOLANGCI_CONFIG" ]]; then
        # Only check if the file has a go: field
        if grep -q -E '^\s*go:\s*"' "$GOLANGCI_CONFIG"; then
            local current_golangci_version
            current_golangci_version=$(grep -E '^\s*go:\s*"' "$GOLANGCI_CONFIG" | sed 's/.*"\([^"]*\)".*/\1/' || echo "")
            if [[ "$current_golangci_version" != "$go_version" ]]; then
                golangci_needs_update=true
                changes_made=true
            fi
        fi
    fi

    # Check Dockerfile
    if [[ -f "$PROJECT_ROOT/Dockerfile" ]]; then
        local current_dockerfile_version
        current_dockerfile_version=$(grep "FROM golang:" "$PROJECT_ROOT/Dockerfile" | head -1 | sed 's/.*golang:\([0-9.]*\).*/\1/' || echo "")
        local expected_dockerfile_version
        expected_dockerfile_version=$(echo "$go_version" | sed 's/\([0-9]*\.[0-9]*\).*/\1/')
        if [[ "$current_dockerfile_version" != "$expected_dockerfile_version" ]]; then
            dockerfile_needs_update=true
            changes_made=true
        fi
    fi

    # Check CONTRIBUTING.md
    if [[ -f "$PROJECT_ROOT/CONTRIBUTING.md" ]]; then
        local current_contributing_version
        current_contributing_version=$(grep -E '\*\*Go\*\*:.*[0-9]+\.[0-9]+' "$PROJECT_ROOT/CONTRIBUTING.md" | sed -E 's/.*([0-9]+\.[0-9]+)\+.*/\1/' || echo "")
        local expected_contributing_version
        expected_contributing_version=$(echo "$go_version" | sed 's/\([0-9]*\.[0-9]*\).*/\1/')
        if [[ -n "$current_contributing_version" && "$current_contributing_version" != "$expected_contributing_version" ]]; then
            contributing_needs_update=true
            changes_made=true
        fi
    fi

    # Only show messages and make changes if needed
    if [[ "$go_mod_needs_update" == true || "$golangci_needs_update" == true || "$dockerfile_needs_update" == true || "$contributing_needs_update" == true ]]; then
        print_status "$BLUE" "Fixing Go version inconsistencies..."
        print_status "$BLUE" "Target version: $go_version"
        echo

        if [[ "$go_mod_needs_update" == true ]]; then
            update_go_mod "$go_version" false
        fi
        if [[ "$golangci_needs_update" == true ]]; then
            update_golangci_config "$go_version" false
        fi
        if [[ "$dockerfile_needs_update" == true ]]; then
            update_dockerfile "$go_version" false
        fi
        if [[ "$contributing_needs_update" == true ]]; then
            update_contributing "$go_version" false
        fi

        echo
        print_status "$GREEN" "✅ All files synchronized to Go version: $go_version"
    else
        print_status "$GREEN" "✅ All files already synchronized to Go version: $go_version"
    fi
}

# Main function
main() {
    local action="fix"

    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --check)
                action="check"
                shift
                ;;
            --fix)
                action="fix"
                shift
                ;;
            --help|-h)
                show_usage
                exit 0
                ;;
            *)
                print_status "$RED" "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done

    # Validate and get Go version
    local go_version
    go_version=$(validate_go_version_file)

    # Perform the requested action
    case $action in
        check)
            check_all_versions "$go_version"
            ;;
        fix)
            fix_all_versions "$go_version"
            ;;
        *)
            print_status "$RED" "Unknown action: $action"
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"

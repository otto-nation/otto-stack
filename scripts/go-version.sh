#!/bin/bash
# Unified Go version management script
# Handles both getting and syncing Go version across project files

set -euo pipefail

# Source common utilities
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=scripts/common.sh
source "$SCRIPT_DIR/common.sh"

PROJECT_ROOT=$(get_project_root)
GO_VERSION_FILE="$PROJECT_ROOT/.go-version"
GOLANGCI_CONFIG="$PROJECT_ROOT/.golangci.yml"

# Read and validate Go version from .go-version file
get_go_version() {
    if [[ ! -f "$GO_VERSION_FILE" ]]; then
        print_error ".go-version file not found at $GO_VERSION_FILE"
        exit 1
    fi

    local version
    version=$(cat "$GO_VERSION_FILE" | tr -d '[:space:]')

    if [[ -z "$version" ]]; then
        print_error ".go-version file is empty"
        exit 1
    fi

    if [[ ! "$version" =~ ^[0-9]+\.[0-9]+(\.[0-9]+)?$ ]]; then
        print_error "Invalid Go version format: $version"
        print_error "Expected format: x.y or x.y.z (e.g., 1.21 or 1.21.0)"
        exit 1
    fi

    echo "$version"
}

# Get subcommand - output version in various formats
cmd_get() {
    local version
    version=$(get_go_version)

    case "${1:-}" in
        --major-minor)
            echo "$version" | sed -E 's/^([0-9]+\.[0-9]+).*/\1/'
            ;;
        --full)
            echo "$version"
            ;;
        --env)
            echo "GO_VERSION=$version"
            ;;
        --github-matrix)
            local major_minor
            major_minor=$(echo "$version" | sed -E 's/^([0-9]+\.[0-9]+).*/\1/')
            echo "['$major_minor']"
            ;;
        --help|-h)
            cat << EOF
Usage: $0 get [option]

Get Go version from .go-version file.

Options:
    (no option)      Output version as-is
    --major-minor    Output major.minor only (e.g., 1.21)
    --full           Output full version
    --env            Output as GO_VERSION=x.y.z
    --github-matrix  Output as JSON array for GitHub Actions
    --help, -h       Show this help
EOF
            ;;
        "")
            echo "$version"
            ;;
        *)
            print_error "Unknown option: $1"
            exit 1
            ;;
    esac
}

# Update .golangci.yml
update_golangci() {
    local version=$1
    local check_only=${2:-false}

    [[ ! -f "$GOLANGCI_CONFIG" ]] && return 0

    if ! grep -q -E '^\s*go:\s*"' "$GOLANGCI_CONFIG"; then
        return 0
    fi

    local current
    current=$(grep -E '^\s*go:\s*"' "$GOLANGCI_CONFIG" | sed 's/.*"\([^"]*\)".*/\1/' || echo "")

    if [[ "$current" == "$version" ]]; then
        print_success ".golangci.yml already has correct version: $version"
        return 0
    fi

    if [[ "$check_only" == "true" ]]; then
        print_error ".golangci.yml has wrong version: '$current' (expected: $version)"
        return 1
    fi

    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' -E 's/go: "[^"]+"/go: "'"$version"'"/' "$GOLANGCI_CONFIG"
    else
        sed -i -E 's/go: "[^"]+"/go: "'"$version"'"/' "$GOLANGCI_CONFIG"
    fi
    print_success "Updated .golangci.yml to: $version"
}

# Update go.mod
update_go_mod() {
    local version=$1
    local check_only=${2:-false}
    local go_mod="$PROJECT_ROOT/go.mod"

    [[ ! -f "$go_mod" ]] && return 0

    local current
    current=$(grep -E '^go [0-9]+\.[0-9]+' "$go_mod" | awk '{print $2}' || echo "")

    [[ -z "$current" ]] && return 0

    if [[ "$current" == "$version" ]]; then
        print_success "go.mod already has correct version: $version"
        return 0
    fi

    if [[ "$check_only" == "true" ]]; then
        print_error "go.mod has wrong version: $current (expected: $version)"
        return 1
    fi

    go mod edit -go="$version" "$go_mod"
    print_success "Updated go.mod to: $version"
}

# Sync subcommand - synchronize version across files
cmd_sync() {
    local check_only=false
    local version

    case "${1:---fix}" in
        --check)
            check_only=true
            ;;
        --fix)
            check_only=false
            ;;
        --help|-h)
            cat << EOF
Usage: $0 sync [option]

Sync Go version across all configuration files.

Options:
    --check    Check if versions are in sync (exit 1 if not)
    --fix      Fix version mismatches (default)
    --help, -h Show this help

Files synchronized:
    - go.mod
    - .golangci.yml
EOF
            return 0
            ;;
        *)
            print_error "Unknown option: $1"
            exit 1
            ;;
    esac

    version=$(get_go_version)
    print_status "Syncing Go version: $version"

    local failed=0
    update_go_mod "$version" "$check_only" || failed=1
    update_golangci "$version" "$check_only" || failed=1

    if [[ $failed -eq 1 ]]; then
        if [[ "$check_only" == "true" ]]; then
            print_error "Version mismatch detected"
        else
            print_error "Failed to update some files"
        fi
        exit 1
    fi

    if [[ "$check_only" == "true" ]]; then
        print_success "All versions are in sync"
    else
        print_success "All versions synchronized"
    fi
}

# Main command dispatcher
show_usage() {
    cat << EOF
Usage: $0 <command> [options]

Unified Go version management for the project.

Commands:
    get     Get Go version from .go-version file
    sync    Sync Go version across configuration files
    help    Show this help message

Run '$0 <command> --help' for more information on a command.
EOF
}

case "${1:-}" in
    get)
        shift
        cmd_get "$@"
        ;;
    sync)
        shift
        cmd_sync "$@"
        ;;
    help|--help|-h)
        show_usage
        ;;
    "")
        print_error "No command specified"
        show_usage
        exit 1
        ;;
    *)
        print_error "Unknown command: $1"
        show_usage
        exit 1
        ;;
esac

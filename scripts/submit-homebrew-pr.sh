#!/bin/bash

# Prepare Homebrew Core PR for otto-stack
# This script prepares the branch and provides instructions for manual submission

set -euo pipefail

# Source common utilities
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=./common.sh
source "${SCRIPT_DIR}/common.sh"

# Configuration
HOMEBREW_CORE_DIR="../homebrew-core"
FORMULA_PATH="Formula/o/otto-stack.rb"

# Extract constants (cache username to avoid double prompt)
export GITHUB_USER="${GITHUB_USER:-}"
eval "$(${SCRIPT_DIR}/extract-constants.sh)"

usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Prepare a Homebrew Core PR for otto-stack

OPTIONS:
    -v, --version VERSION    Version to submit (default: latest from GitHub)
    -s, --sha256 SHA256      SHA256 hash of release tarball (default: auto-calculate)
    -h, --help              Show this help message

EXAMPLES:
    $0                                    # Use latest version and auto-calculate SHA
    $0 --version v1.0.0                  # Use specific version, auto-calculate SHA
    $0 --version v1.0.0 --sha256 abc123  # Use specific version and SHA
EOF
}

main() {
    local version=""
    local sha256=""

    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -v|--version)
                version="$2"
                shift 2
                ;;
            -s|--sha256)
                sha256="$2"
                shift 2
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                usage
                exit 1
                ;;
        esac
    done

    # Get latest version if not provided
    if [[ -z "$version" ]]; then
        print_info "Fetching latest version from GitHub..."
        version=$(get_latest_version "$GITHUB_ORG/$GITHUB_REPO")
        if [[ -z "$version" ]]; then
            print_error "Failed to fetch latest version"
            exit 1
        fi
        print_info "Using latest version: $version"
    fi

    # Calculate SHA256 if not provided
    if [[ -z "$sha256" ]]; then
        print_info "Calculating SHA256 for source tarball..."
        local tarball_url="https://github.com/$GITHUB_ORG/$GITHUB_REPO/archive/refs/tags/$version.tar.gz"
        local temp_file=$(mktemp)
        
        if ! curl -L -o "$temp_file" "$tarball_url"; then
            print_error "Failed to download source tarball"
            rm -f "$temp_file"
            exit 1
        fi
        
        sha256=$(sha256sum "$temp_file" | cut -d' ' -f1)
        rm -f "$temp_file"
        print_info "Calculated SHA256: $sha256"
    fi

    # Check if homebrew-core directory exists
    if [[ ! -d "$HOMEBREW_CORE_DIR" ]]; then
        print_error "Homebrew core directory not found: $HOMEBREW_CORE_DIR"
        print_info "Please clone your homebrew-core fork to: $HOMEBREW_CORE_DIR"
        exit 1
    fi

    print_info "Preparing Homebrew Core PR for $APP_NAME_TITLE $version"

    # Navigate to homebrew-core directory
    cd "$HOMEBREW_CORE_DIR"

    # Ensure we're on the right branch and up to date
    print_info "Updating homebrew-core fork..."
    git fetch upstream
    git checkout master
    git merge upstream/master
    git push origin master

    # Create new branch for this version
    local branch_name="otto-stack-${version#v}"
    git checkout -b "$branch_name"

    # Update the formula
    print_info "Updating formula with version $version and SHA256 $sha256"
    sed -i.bak \
        -e "s|url \".*\"|url \"https://github.com/$GITHUB_ORG/$GITHUB_REPO/archive/refs/tags/$version.tar.gz\"|" \
        -e "s|sha256 \".*\"|sha256 \"$sha256\"|" \
        "$FORMULA_PATH"

    # Remove backup file
    rm -f "${FORMULA_PATH}.bak"

    # Commit changes
    git add "$FORMULA_PATH"
    git commit -m "otto-stack: update to $version"

    # Push branch
    git push origin "$branch_name"

    print_success "Branch prepared successfully!"
    print_info ""
    print_info "📋 Next steps (REQUIRED before submitting PR):"
    print_info ""
    print_info "1. Test the formula locally:"
    print_info "   brew uninstall --force otto-stack"
    print_info "   HOMEBREW_NO_INSTALL_FROM_API=1 brew install --build-from-source $FORMULA_PATH"
    print_info "   brew test otto-stack"
    print_info "   brew audit --strict otto-stack"
    print_info "   brew style otto-stack"
    print_info ""
    print_info "2. If tests pass, create PR manually:"
    print_info "   https://github.com/Homebrew/homebrew-core/compare/master...$GITHUB_USER:homebrew-core:$branch_name"
    print_info ""
    print_info "3. Use this commit message format:"
    print_info "   otto-stack: update to $version"
    print_info ""
    print_warning "⚠️  Do NOT create automated PRs to Homebrew Core"
    print_warning "   Manual testing and human review are required"
}

main "$@"

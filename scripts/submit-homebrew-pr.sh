#!/bin/bash

# Prepare Homebrew Core PR for otto-stack
# This script prepares the branch and provides instructions for manual submission

set -euo pipefail

# Source common utilities
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/common.sh"

# Configuration
HOMEBREW_CORE_DIR="../homebrew-core"
FORMULA_PATH="Formula/o/otto-stack.rb"

# Extract constants
eval "$(extract_constants)"

usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Prepare a Homebrew Core PR for otto-stack

OPTIONS:
    -v, --version VERSION    Version to submit (required)
    -s, --sha256 SHA256      SHA256 hash of the release tarball (required)
    -h, --help              Show this help message

EXAMPLES:
    $0 --version v1.0.0 --sha256 abc123...
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
                error "Unknown option: $1"
                usage
                exit 1
                ;;
        esac
    done

    # Validate required arguments
    if [[ -z "$version" ]]; then
        error "Version is required"
        usage
        exit 1
    fi

    if [[ -z "$sha256" ]]; then
        error "SHA256 hash is required"
        usage
        exit 1
    fi

    # Check if homebrew-core directory exists
    if [[ ! -d "$HOMEBREW_CORE_DIR" ]]; then
        error "Homebrew core directory not found: $HOMEBREW_CORE_DIR"
        exit 1
    fi

    info "Preparing Homebrew Core PR for $APP_NAME_TITLE $version"

    # Navigate to homebrew-core directory
    cd "$HOMEBREW_CORE_DIR"

    # Ensure we're on the right branch and up to date
    info "Updating homebrew-core fork..."
    git fetch upstream
    git checkout master
    git merge upstream/master
    git push origin master

    # Create new branch for this version
    local branch_name="otto-stack-${version#v}"
    git checkout -b "$branch_name"

    # Update the formula
    info "Updating formula with version $version and SHA256 $sha256"
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

    success "Branch prepared successfully!"
    info ""
    info "📋 Next steps (REQUIRED before submitting PR):"
    info ""
    info "1. Test the formula locally:"
    info "   brew uninstall --force otto-stack"
    info "   HOMEBREW_NO_INSTALL_FROM_API=1 brew install --build-from-source $FORMULA_PATH"
    info "   brew test otto-stack"
    info "   brew audit --strict otto-stack"
    info "   brew style otto-stack"
    info ""
    info "2. If tests pass, create PR manually:"
    info "   https://github.com/Homebrew/homebrew-core/compare/master...$GITHUB_USER:homebrew-core:$branch_name"
    info ""
    info "3. Use this commit message format:"
    info "   otto-stack: update to $version"
    info ""
    warning "⚠️  Do NOT create automated PRs to Homebrew Core"
    warning "   Manual testing and human review are required"
}

main "$@"

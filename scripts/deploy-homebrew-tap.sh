#!/bin/bash
set -euo pipefail

# Deploy Homebrew formula to tap repository
# This script is designed to run in CI/CD after a release

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=scripts/common.sh
source "$SCRIPT_DIR/common.sh"

# Configuration
readonly TAP_REPO="otto-nation/homebrew-tap"
readonly TAP_URL="https://github.com/${TAP_REPO}.git"
readonly FORMULA_NAME="otto-stack"

show_usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Deploy Homebrew formula to tap repository

OPTIONS:
    -v, --version VER    Version to deploy (required)
    -t, --token TOKEN    GitHub token for pushing (required in CI)
    --dry-run            Show what would be done without pushing
    -h, --help           Show this help message

ENVIRONMENT VARIABLES:
    GITHUB_TOKEN         GitHub token (alternative to --token)
    HOMEBREW_TAP_TOKEN   Dedicated tap token (preferred)

EXAMPLES:
    $0 -v v1.2.3                    # Deploy version v1.2.3
    $0 -v v1.2.3 --dry-run          # Test without pushing
    HOMEBREW_TAP_TOKEN=xxx $0 -v v1.2.3  # Use dedicated token

EOF
}

parse_args() {
    VERSION=""
    TOKEN="${HOMEBREW_TAP_TOKEN:-${GITHUB_TOKEN:-}}"
    DRY_RUN=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -v|--version)
                VERSION="$2"
                shift 2
                ;;
            -t|--token)
                TOKEN="$2"
                shift 2
                ;;
            --dry-run)
                DRY_RUN=true
                shift
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
    
    if [[ -z "$VERSION" ]]; then
        print_error "Version is required"
        show_usage
        exit 1
    fi
    
    if [[ -z "$TOKEN" ]] && [[ "$DRY_RUN" == "false" ]]; then
        print_error "GitHub token is required (set HOMEBREW_TAP_TOKEN or GITHUB_TOKEN)"
        exit 1
    fi
}

clone_tap() {
    local tap_dir="$1"
    local auth_url
    
    if [[ -n "$TOKEN" ]]; then
        auth_url="https://x-access-token:${TOKEN}@github.com/${TAP_REPO}.git"
    else
        auth_url="$TAP_URL"
    fi
    
    print_status "Cloning tap repository..."
    git clone "$auth_url" "$tap_dir"
}

update_formula() {
    local tap_dir="$1"
    local version="$2"
    local formula_src="Formula/${FORMULA_NAME}.rb"
    local formula_dest="${tap_dir}/Formula/${FORMULA_NAME}.rb"
    
    if [[ ! -f "$formula_src" ]]; then
        print_error "Formula not found: $formula_src"
        print_info "Run: ./scripts/update-homebrew-formula.sh -v $version"
        return 1
    fi
    
    print_status "Copying formula to tap..."
    mkdir -p "${tap_dir}/Formula"
    cp "$formula_src" "$formula_dest"
    
    print_success "✓ Formula copied"
}

commit_and_push() {
    local tap_dir="$1"
    local version="$2"
    
    cd "$tap_dir"
    
    # Configure git
    git config user.name "github-actions[bot]"
    git config user.email "github-actions[bot]@users.noreply.github.com"
    
    # Check if there are changes
    if ! git diff --quiet Formula/${FORMULA_NAME}.rb; then
        print_status "Committing changes..."
        git add "Formula/${FORMULA_NAME}.rb"
        git commit -m "chore: update ${FORMULA_NAME} to ${version}"
        
        if [[ "$DRY_RUN" == "true" ]]; then
            print_info "🏃 DRY RUN: Would push to ${TAP_REPO}"
            git show HEAD
        else
            print_status "Pushing to ${TAP_REPO}..."
            git push origin main
            print_success "✅ Pushed to tap repository"
        fi
    else
        print_info "ℹ️  No changes to formula"
    fi
}

main() {
    parse_args "$@"
    
    print_header "🍺 Homebrew Tap Deployment"
    echo
    
    print_info "Version: $VERSION"
    print_info "Tap: $TAP_REPO"
    if [[ "$DRY_RUN" == "true" ]]; then
        print_warning "DRY RUN MODE - No changes will be pushed"
    fi
    echo
    
    # Create temp directory for tap
    local tap_dir
    tap_dir=$(mktemp -d)
    trap "rm -rf '$tap_dir'" EXIT
    
    # Clone tap repository
    clone_tap "$tap_dir" || exit 1
    
    # Update formula
    update_formula "$tap_dir" "$VERSION" || exit 1
    
    # Commit and push
    commit_and_push "$tap_dir" "$VERSION" || exit 1
    
    print_success "🎉 Homebrew tap deployment complete!"
    print_info "Users can now install with: brew install otto-nation/tap/${FORMULA_NAME}"
}

main "$@"

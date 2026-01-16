#!/bin/bash
# Unified Homebrew formula management script
# Handles both updating and deploying Homebrew formulas

set -euo pipefail

# Source common utilities
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=scripts/common.sh
source "$SCRIPT_DIR/common.sh"

# Configuration
readonly REPO="${GITHUB_ORG}/${GITHUB_REPO}"
readonly TAP_REPO="otto-nation/homebrew-tap"
readonly TAP_URL="https://github.com/${TAP_REPO}.git"
readonly FORMULA_FILE="Formula/${APP_NAME}.rb"

# Calculate SHA256 for a remote file
get_remote_sha256() {
    local url="$1"
    local temp_file
    
    temp_file=$(mktemp)
    trap "rm -f '$temp_file'" RETURN
    
    if ! download_file "$url" "$temp_file" false; then
        print_error "Failed to download $url"
        return 1
    fi
    
    shasum -a 256 "$temp_file" | cut -d' ' -f1
}

# Update formula with checksums
cmd_update() {
    local version=""
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -v|--version)
                version="$2"
                shift 2
                ;;
            -h|--help)
                cat << EOF
Usage: $0 update [OPTIONS]

Update Homebrew formula with checksums for a release.

OPTIONS:
    -v, --version VER    Version to update (required)
    -h, --help           Show this help

EXAMPLE:
    $0 update -v v1.2.3
EOF
                return 0
                ;;
            *)
                print_error "Unknown option: $1"
                return 1
                ;;
        esac
    done
    
    if [[ -z "$version" ]]; then
        print_error "Version is required. Use -v or --version"
        return 1
    fi
    
    local project_root
    project_root=$(get_project_root)
    
    print_status "Updating Homebrew formula for version $version..."
    
    # Define platforms and their download URLs
    declare -A platforms=(
        ["darwin-amd64"]="https://github.com/$REPO/releases/download/$version/$APP_NAME-darwin-amd64"
        ["darwin-arm64"]="https://github.com/$REPO/releases/download/$version/$APP_NAME-darwin-arm64"
    )
    
    # Calculate checksums
    declare -A checksums
    for platform in "${!platforms[@]}"; do
        local url="${platforms[$platform]}"
        print_status "Calculating SHA256 for $platform..."
        checksums[$platform]=$(get_remote_sha256 "$url")
    done
    
    # Update formula file
    local formula_path="$project_root/$FORMULA_FILE"
    if [[ ! -f "$formula_path" ]]; then
        print_error "Formula file not found: $formula_path"
        return 1
    fi
    
    # Update version and checksums in formula
    local version_clean="${version#v}"
    
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' "s/version \".*\"/version \"$version_clean\"/" "$formula_path"
        sed -i '' "s|url \"https://github.com/$REPO/releases/download/v[^/]*/|url \"https://github.com/$REPO/releases/download/$version/|" "$formula_path"
        sed -i '' "/darwin.*amd64/,/sha256/ s/sha256 \".*\"/sha256 \"${checksums[darwin-amd64]}\"/" "$formula_path"
        sed -i '' "/darwin.*arm64/,/sha256/ s/sha256 \".*\"/sha256 \"${checksums[darwin-arm64]}\"/" "$formula_path"
    else
        sed -i "s/version \".*\"/version \"$version_clean\"/" "$formula_path"
        sed -i "s|url \"https://github.com/$REPO/releases/download/v[^/]*/|url \"https://github.com/$REPO/releases/download/$version/|" "$formula_path"
        sed -i "/darwin.*amd64/,/sha256/ s/sha256 \".*\"/sha256 \"${checksums[darwin-amd64]}\"/" "$formula_path"
        sed -i "/darwin.*arm64/,/sha256/ s/sha256 \".*\"/sha256 \"${checksums[darwin-arm64]}\"/" "$formula_path"
    fi
    
    print_success "Formula updated successfully"
    print_info "Updated checksums:"
    for platform in "${!checksums[@]}"; do
        print_info "  $platform: ${checksums[$platform]}"
    done
}

# Deploy formula to tap repository
cmd_deploy() {
    local version=""
    local token="${HOMEBREW_TAP_TOKEN:-${GITHUB_TOKEN:-}}"
    local dry_run=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -v|--version)
                version="$2"
                shift 2
                ;;
            -t|--token)
                token="$2"
                shift 2
                ;;
            --dry-run)
                dry_run=true
                shift
                ;;
            -h|--help)
                cat << EOF
Usage: $0 deploy [OPTIONS]

Deploy Homebrew formula to tap repository.

OPTIONS:
    -v, --version VER    Version to deploy (required)
    -t, --token TOKEN    GitHub token for pushing
    --dry-run            Show what would be done without pushing
    -h, --help           Show this help

ENVIRONMENT VARIABLES:
    GITHUB_TOKEN         GitHub token (alternative to --token)
    HOMEBREW_TAP_TOKEN   Dedicated tap token (preferred)

EXAMPLE:
    $0 deploy -v v1.2.3
    $0 deploy -v v1.2.3 --dry-run
EOF
                return 0
                ;;
            *)
                print_error "Unknown option: $1"
                return 1
                ;;
        esac
    done
    
    if [[ -z "$version" ]]; then
        print_error "Version is required. Use -v or --version"
        return 1
    fi
    
    if [[ -z "$token" ]] && [[ "$dry_run" == false ]]; then
        print_error "GitHub token required. Set GITHUB_TOKEN or HOMEBREW_TAP_TOKEN"
        return 1
    fi
    
    local project_root
    project_root=$(get_project_root)
    
    local formula_source="$project_root/$FORMULA_FILE"
    if [[ ! -f "$formula_source" ]]; then
        print_error "Formula file not found: $formula_source"
        return 1
    fi
    
    print_status "Deploying formula for version $version to $TAP_REPO..."
    
    # Create temp directory for tap repo
    local tap_dir
    tap_dir=$(mktemp -d)
    trap "rm -rf '$tap_dir'" RETURN
    
    # Clone tap repository
    print_status "Cloning tap repository..."
    if [[ -n "$token" ]]; then
        git clone "https://x-access-token:${token}@github.com/${TAP_REPO}.git" "$tap_dir" 2>&1 | grep -v "x-access-token" || true
    else
        git clone "$TAP_URL" "$tap_dir"
    fi
    
    # Copy formula
    local formula_dest="$tap_dir/$FORMULA_FILE"
    mkdir -p "$(dirname "$formula_dest")"
    cp "$formula_source" "$formula_dest"
    
    # Commit and push
    cd "$tap_dir"
    git config user.name "GitHub Actions"
    git config user.email "actions@github.com"
    
    if git diff --quiet; then
        print_info "No changes to formula"
        return 0
    fi
    
    git add "$FORMULA_FILE"
    git commit -m "Update ${APP_NAME} to $version"
    
    if [[ "$dry_run" == true ]]; then
        print_warning "Dry run - would push:"
        git show --stat
        return 0
    fi
    
    print_status "Pushing to tap repository..."
    git push origin main
    
    print_success "Formula deployed successfully to $TAP_REPO"
}

# Main command dispatcher
show_usage() {
    cat << EOF
Usage: $0 <command> [options]

Unified Homebrew formula management.

Commands:
    update    Update formula with checksums for a release
    deploy    Deploy formula to tap repository
    help      Show this help message

Run '$0 <command> --help' for more information on a command.
EOF
}

case "${1:-}" in
    update)
        shift
        cmd_update "$@"
        ;;
    deploy)
        shift
        cmd_deploy "$@"
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

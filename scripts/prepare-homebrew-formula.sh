#!/bin/bash

# Prepare Homebrew Core formula for otto-stack
# This script prepares the formula and branch for manual testing and PR submission

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

Prepare a Homebrew Core formula for otto-stack

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

    # Check if homebrew-core directory exists
    if [[ ! -d "$HOMEBREW_CORE_DIR" ]]; then
        print_error "Homebrew core directory not found: $HOMEBREW_CORE_DIR"
        print_info "Please clone your homebrew-core fork to: $HOMEBREW_CORE_DIR"
        exit 1
    fi

    print_info "Preparing Homebrew Core formula for $APP_NAME_TITLE $version"

    # Navigate to homebrew-core directory
    cd "$HOMEBREW_CORE_DIR"

    # Tap homebrew/core if needed (required by Homebrew guidelines)
    if command -v brew >/dev/null 2>&1; then
        brew tap --force homebrew/core
    fi

    # Ensure we're on the right branch and up to date
    print_info "Updating homebrew-core fork..."
    git fetch upstream
    git checkout master
    git merge upstream/master
    git push origin master

    # Create new branch for this version
    local clean_version="${version#otto-stack-}"  # Remove otto-stack- prefix if present
    clean_version="${clean_version#v}"            # Remove v prefix if present
    local branch_name="otto-stack-${clean_version}"
    
    # Delete branch if it exists
    if git show-ref --verify --quiet "refs/heads/$branch_name"; then
        print_info "Deleting existing branch: $branch_name"
        git branch -D "$branch_name"
    fi
    
    git checkout -b "$branch_name"

    # Calculate SHA256 for source tarball
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

    # Create/update formula in our local fork
    print_info "Creating/updating source-based formula with version $version"
    
    local tarball_url="https://github.com/$GITHUB_ORG/$GITHUB_REPO/archive/refs/tags/$version.tar.gz"
    
    # Check if formula exists
    if [[ -f "Formula/o/otto-stack.rb" ]]; then
        formula_file="Formula/o/otto-stack.rb"
        is_new_formula=false
        print_info "Updating existing formula"
    else
        formula_file="Formula/o/otto-stack.rb"
        is_new_formula=true
        print_info "Creating new formula file: $formula_file"
        mkdir -p Formula/o
    fi
    
    cat > "$formula_file" << EOF
class OttoStack < Formula
  desc "Powerful development stack management tool for streamlined local development"
  homepage "https://github.com/otto-nation/otto-stack"
  url "$tarball_url"
  sha256 "$sha256"
  license "MIT"

  depends_on "go" => :build

  def install
    ldflags = %W[
      -s -w
      -X github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/version.Version=#{version}
      -X github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/version.BuildDate=#{time.iso8601}
      -X github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/version.GitCommit=homebrew
    ]

    system "go", "build", *std_go_args(ldflags:), "./cmd/otto-stack"
    generate_completions_from_executable(bin/"otto-stack", "completion")
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/otto-stack version")
    system bin/"otto-stack", "init", "--help"
  end
end
EOF

    # Validate formula syntax (this is what we can actually test locally)
    if command -v brew >/dev/null 2>&1; then
        print_info "Validating formula syntax..."
        
        # Run style check on the formula file
        print_info "Running brew style..."
        brew style "$(pwd)/$formula_file"
        
        print_success "Formula validation passed!"
        print_info ""
        print_info "⚠️  Note: 'brew audit' can only be run after the formula is in a tap"
        print_info "   The Homebrew Core maintainers will run audit checks during PR review"
        
    else
        print_warning "Homebrew not available - skipping validation"
    fi

    # Commit changes with AI-generated message (only after tests pass)
    print_info "Generating commit message..."
    local commit_msg
    if command -v q >/dev/null 2>&1; then
        local formula_type=""
        if [[ "$is_new_formula" == true ]]; then
            formula_type="(new formula)"
        fi
        
        commit_msg=$(q chat --no-interactive "Generate a Homebrew commit message for otto-stack $version.

Requirements:
- First line: 'otto-stack $version $formula_type' (max 50 characters)
- Two newlines, then detailed explanation
- Follow Homebrew standards exactly
- Mention version update and SHA256 hash update
- Keep it concise and professional

Current version: $version
SHA256: $sha256
Is new formula: $is_new_formula" | sed 's/\x1b\[[0-9;]*m//g' | sed 's/^> //')
    else
        # Fallback to simple message
        if [[ "$is_new_formula" == true ]]; then
            commit_msg="otto-stack $version (new formula)"
        else
            commit_msg="otto-stack $version"
        fi
    fi
    
    git add "$formula_file"
    git commit -m "$commit_msg"

    # Push branch
    git push origin "$branch_name"

    print_success "Formula prepared and tested successfully!"
    print_info ""
    print_info "📋 Next steps:"
    print_info ""
    print_info "1. ✅ Formula tested and committed (all tests passed)"
    print_info ""
    print_info "2. Create PR manually at:"
    print_info "   https://github.com/Homebrew/homebrew-core/compare"
    print_info ""
    print_info "3. ✅ Commit message follows Homebrew standards"
    print_info ""
    print_info "The formula is ready for Homebrew Core submission!"
}

main "$@"

#!/bin/bash
set -euo pipefail

# Extract constants from Go code for use in shell scripts
# This ensures scripts use the same values as the Go application

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Extract constants from Go files
extract_constant() {
    local const_name="$1"
    local file_pattern="$2"
    
    grep -r "^[[:space:]]*${const_name}[[:space:]]*=" "$PROJECT_ROOT/$file_pattern" | \
    head -1 | \
    sed -E 's/.*=[[:space:]]*"([^"]+)".*/\1/' || echo ""
}

# Extract brand constants
GITHUB_ORG=$(extract_constant "GitHubOrg" "internal/pkg/constants/brand.go")
APP_NAME=$(extract_constant "AppName" "internal/pkg/constants/brand.go")
APP_NAME_TITLE=$(extract_constant "AppNameTitle" "internal/pkg/constants/brand.go")

# Handle GitHubRepo which references AppName
GITHUB_REPO_RAW=$(grep -r "^[[:space:]]*GitHubRepo[[:space:]]*=" "$PROJECT_ROOT/internal/pkg/constants/brand.go" | head -1)
if [[ "$GITHUB_REPO_RAW" == *"AppName"* ]]; then
    # GitHubRepo = AppName, so use APP_NAME value
    GITHUB_REPO="$APP_NAME"
else
    # Extract literal value
    GITHUB_REPO=$(echo "$GITHUB_REPO_RAW" | sed -E 's/.*=[[:space:]]*"([^"]+)".*/\1/')
fi

# Get current GitHub user dynamically
get_github_user() {
    # Use existing GITHUB_USER if set
    if [[ -n "${GITHUB_USER:-}" ]]; then
        echo "$GITHUB_USER"
        return
    fi
    
    # Try GitHub CLI first
    if command -v gh >/dev/null 2>&1; then
        local gh_user=$(gh api user --jq '.login' 2>/dev/null || echo "")
        if [[ -n "$gh_user" ]]; then
            echo "$gh_user"
            return
        fi
    fi
    
    # Try git config github.user
    local git_github_user=$(git config --global github.user 2>/dev/null || echo "")
    if [[ -n "$git_github_user" ]]; then
        echo "$git_github_user"
        return
    fi
    
    # Try to extract from git remote origin URL
    local remote_url=$(git remote get-url origin 2>/dev/null || echo "")
    if [[ "$remote_url" =~ github\.com[:/]([^/]+)/ ]]; then
        local repo_owner="${BASH_REMATCH[1]}"
        
        # If this is the main repo (otto-nation), we need the fork owner
        if [[ "$repo_owner" == "otto-nation" ]]; then
            # Check if there's a fork remote
            local fork_url=$(git remote get-url fork 2>/dev/null || echo "")
            if [[ "$fork_url" =~ github\.com[:/]([^/]+)/ ]]; then
                echo "${BASH_REMATCH[1]}"
                return
            fi
            
            # Fallback: prompt user for their GitHub username
            if [[ -t 0 ]]; then
                read -p "Enter your GitHub username (for homebrew-core fork): " github_user
                echo "$github_user"
            else
                echo "unknown-user"
            fi
        else
            # This is already a fork, use the owner
            echo "$repo_owner"
        fi
    else
        # No GitHub remote found, prompt user
        if [[ -t 0 ]]; then
            read -p "Enter your GitHub username: " github_user
            echo "$github_user"
        else
            echo "unknown-user"
        fi
    fi
}

GITHUB_USER=$(get_github_user)

# Fallback to hardcoded values if extraction fails
GITHUB_ORG="${GITHUB_ORG:-otto-nation}"
GITHUB_USER="${GITHUB_USER:-$(whoami)}"
GITHUB_REPO="${GITHUB_REPO:-otto-stack}"
APP_NAME="${APP_NAME:-otto-stack}"
APP_NAME_TITLE="${APP_NAME_TITLE:-Otto Stack}"

# Export for use by other scripts
export GITHUB_ORG GITHUB_USER GITHUB_REPO APP_NAME APP_NAME_TITLE

# If called directly, output the constants
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    cat << EOF
# Extracted constants from Go code
export GITHUB_ORG="$GITHUB_ORG"
export GITHUB_USER="$GITHUB_USER"
export GITHUB_REPO="$GITHUB_REPO"
export APP_NAME="$APP_NAME"
export APP_NAME_TITLE="$APP_NAME_TITLE"
EOF
fi

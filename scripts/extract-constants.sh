#!/bin/bash
set -euo pipefail

# Constants for shell scripts (hardcoded for simplicity)
# These values match the Go constants in internal/pkg/constants/app.go

# Brand constants
export GITHUB_ORG="otto-nation"
export APP_NAME="otto-stack"
export APP_NAME_TITLE="Otto Stack"
export GITHUB_REPO="otto-stack"

# Get current GitHub user dynamically
get_github_user() {
    # Use existing GITHUB_USER if set
    if [[ -n "${GITHUB_USER:-}" ]]; then
        echo "$GITHUB_USER"
        return
    fi
    
    # Try to get from git config
    if command -v git >/dev/null 2>&1; then
        local git_user
        git_user=$(git config --get user.name 2>/dev/null || echo "")
        if [[ -n "$git_user" ]]; then
            echo "$git_user"
            return
        fi
    fi
    
    # Fallback to environment or prompt
    echo "${USER:-unknown}"
}

# Set dynamic values
export GITHUB_USER="${GITHUB_USER:-$(get_github_user)}"

# Export all constants for use in other scripts
export GITHUB_ORG APP_NAME APP_NAME_TITLE GITHUB_REPO GITHUB_USER

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

# Fallback to hardcoded values if extraction fails
GITHUB_ORG="${GITHUB_ORG:-otto-nation}"
GITHUB_REPO="${GITHUB_REPO:-otto-stack}"
APP_NAME="${APP_NAME:-otto-stack}"
APP_NAME_TITLE="${APP_NAME_TITLE:-Otto Stack}"

# Export for use by other scripts
export GITHUB_ORG GITHUB_REPO APP_NAME APP_NAME_TITLE

# If called directly, output the constants
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    cat << EOF
# Extracted constants from Go code
export GITHUB_ORG="$GITHUB_ORG"
export GITHUB_REPO="$GITHUB_REPO"
export APP_NAME="$APP_NAME"
export APP_NAME_TITLE="$APP_NAME_TITLE"
EOF
fi

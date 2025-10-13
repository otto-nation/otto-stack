#!/bin/bash
set -euo pipefail

# otto-stack uninstallation script
readonly BINARY_NAME="otto-stack"

# Source common utilities
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=scripts/common.sh
source "$SCRIPT_DIR/common.sh"

# Find otto-stack installations
find_installations() {
    local installations=()
    local common_paths=(
        "/usr/local/bin"
        "/usr/bin"
        "$HOME/.local/bin"
        "$HOME/bin"
    )
    
    # Check common installation paths
    for path in "${common_paths[@]}"; do
        if [[ -f "${path}/${BINARY_NAME}" ]]; then
            installations+=("${path}/${BINARY_NAME}")
        fi
    done
    
    # Check PATH
    if command_exists "$BINARY_NAME"; then
        local path_location
        path_location=$(command -v "$BINARY_NAME")
        # Add to list if not already found
        if [[ ! " ${installations[*]} " =~ " ${path_location} " ]]; then
            installations+=("$path_location")
        fi
    fi
    
    printf '%s\n' "${installations[@]}"
}

# Remove installation
remove_installation() {
    local file_path="$1"
    local dir_path
    dir_path=$(dirname "$file_path")
    
    if [[ -w "$dir_path" ]]; then
        rm -f "$file_path"
        print_success "Removed ${file_path}"
    else
        print_warning "Removing ${file_path} (requires sudo)..."
        sudo rm -f "$file_path"
        print_success "Removed ${file_path}"
    fi
}

# Show usage information
show_usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Uninstall otto-stack from the system

OPTIONS:
    -f, --force      Force removal without confirmation
    -h, --help       Show this help message
    -v, --version    Show version and exit

EXAMPLES:
    $0               # Interactive uninstall
    $0 --force       # Force uninstall without confirmation

EOF
}

# Parse command line arguments
parse_args() {
    FORCE=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -f|--force)
                FORCE=true
                shift
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            -v|--version)
                echo "otto-stack uninstaller v1.0.0"
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
}

# Confirm removal
confirm_removal() {
    local installations=("$@")
    
    if [[ "$FORCE" == "true" ]]; then
        return 0
    fi
    
    echo
    print_warning "The following otto-stack installations will be removed:"
    for installation in "${installations[@]}"; do
        echo "  - $installation"
    done
    echo
    
    confirm "Are you sure you want to continue?"
}

# Main uninstallation function
main() {
    parse_args "$@"
    
    print_header "🗑️  otto-stack uninstaller"
    echo
    
    # Find all installations
    local installations
    readarray -t installations < <(find_installations)
    
    if [[ ${#installations[@]} -eq 0 ]]; then
        print_warning "No otto-stack installations found"
        exit 0
    fi
    
    # Confirm removal
    if ! confirm_removal "${installations[@]}"; then
        print_info "Uninstallation cancelled"
        exit 0
    fi
    
    # Remove installations
    print_status "Removing otto-stack installations..."
    for installation in "${installations[@]}"; do
        if [[ -f "$installation" ]]; then
            remove_installation "$installation"
        fi
    done
    
    # Verify removal
    if command_exists "$BINARY_NAME"; then
        print_warning "⚠️  otto-stack may still be available in PATH"
        print_warning "You may need to restart your shell or check additional locations"
    else
        print_success "🎉 otto-stack has been completely removed from your system"
    fi
}

# Run main function with all arguments
main "$@"

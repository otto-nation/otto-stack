#!/bin/bash
set -euo pipefail

# otto-stack installation script
readonly REPO="otto-nation/otto-stack"
readonly BINARY_NAME="otto-stack"
readonly DEFAULT_INSTALL_DIR="/usr/local/bin"

# Source common utilities
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=scripts/common.sh
source "$SCRIPT_DIR/common.sh"

# Download and install binary
install_binary() {
    local platform version install_dir download_url temp_dir temp_file target_path
    
    platform=$(detect_platform) || exit 1
    version=$(get_latest_version "$REPO") || exit 1
    install_dir="${INSTALL_DIR:-$DEFAULT_INSTALL_DIR}"
    
    print_status "Installing ${BINARY_NAME} ${version} for ${platform}..."
    
    # Construct download URL
    download_url="https://github.com/${REPO}/releases/download/${version}/${BINARY_NAME}-${platform}"
    if [[ "$platform" == windows-* ]]; then
        download_url="${download_url}.exe"
    fi
    
    # Create temporary directory with cleanup
    temp_dir=$(create_temp_dir)
    temp_file="${temp_dir}/${BINARY_NAME}"
    
    # Download binary
    print_status "Downloading from ${download_url}..."
    if ! download_file "$download_url" "$temp_file"; then
        print_error "Failed to download binary from ${download_url}"
        print_error "Please check if the release exists for your platform"
        exit 1
    fi
    
    # Make executable and verify
    chmod +x "$temp_file"
    verify_executable "$temp_file" || exit 1
    
    # Install binary
    target_path=$(install_file "$temp_file" "$install_dir" "$BINARY_NAME")
    print_success "${BINARY_NAME} installed successfully to ${target_path}"
    
    # Verify installation
    if command_exists "$BINARY_NAME"; then
        print_success "🎉 Installation verified! Run '${BINARY_NAME} --help' to get started"
    else
        print_warning "⚠️  ${BINARY_NAME} installed but not in PATH"
        print_warning "Add ${install_dir} to your PATH or run: export PATH=\"${install_dir}:\$PATH\""
    fi
}

# Show usage information
show_usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Install otto-stack from GitHub releases

OPTIONS:
    -d, --dir DIR    Installation directory (default: ${DEFAULT_INSTALL_DIR})
    -h, --help       Show this help message
    -v, --version    Show version and exit

ENVIRONMENT VARIABLES:
    INSTALL_DIR      Installation directory (overrides -d/--dir)

EXAMPLES:
    $0                           # Install to ${DEFAULT_INSTALL_DIR}
    $0 --dir ~/.local/bin        # Install to custom directory
    INSTALL_DIR=~/bin $0         # Install using environment variable

EOF
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -d|--dir)
                INSTALL_DIR="$2"
                shift 2
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            -v|--version)
                echo "otto-stack installer v1.0.0"
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

# Main installation function
main() {
    parse_args "$@"
    
    print_header "🚀 otto-stack installer"
    echo
    
    check_dependencies curl || exit 1
    install_binary
}

# Run main function with all arguments
main "$@"

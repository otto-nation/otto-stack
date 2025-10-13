#!/bin/bash
set -euo pipefail

# Source common utilities
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=scripts/common.sh
source "$SCRIPT_DIR/common.sh"

# Constants from Go code (loaded by common.sh)
readonly REPO="${GITHUB_ORG}/${GITHUB_REPO}"
readonly FORMULA_FILE="Formula/${APP_NAME}.rb"

# Get the latest release version and download URLs
get_release_info() {
    local version
    version=$(get_latest_version "$REPO") || {
        print_error "Failed to get latest version"
        return 1
    }
    
    echo "$version"
}

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
update_formula() {
    local version="$1"
    local project_root
    project_root=$(get_project_root)
    
    print_status "Updating Homebrew formula for version $version..."
    
    # Define platforms and their download URLs
    declare -A platforms=(
        ["darwin-amd64"]="https://github.com/$REPO/releases/download/$version/$APP_NAME-darwin-amd64"
        ["darwin-arm64"]="https://github.com/$REPO/releases/download/$version/$APP_NAME-darwin-arm64"
        ["linux-amd64"]="https://github.com/$REPO/releases/download/$version/$APP_NAME-linux-amd64"
        ["linux-arm64"]="https://github.com/$REPO/releases/download/$version/$APP_NAME-linux-arm64"
    )
    
    # Calculate checksums for each platform
    declare -A checksums
    for platform in "${!platforms[@]}"; do
        local url="${platforms[$platform]}"
        print_status "Calculating checksum for $platform..."
        
        local checksum
        checksum=$(get_remote_sha256 "$url") || {
            print_error "Failed to get checksum for $platform"
            return 1
        }
        
        checksums[$platform]="$checksum"
        print_success "✓ $platform: $checksum"
    done
    
    # Create updated formula
    cat > "$project_root/$FORMULA_FILE" << EOF
class $(echo "$APP_NAME" | sed 's/-//g' | sed 's/\b\w/\U&/g') < Formula
  desc "Development stack management tool for streamlined local development automation"
  homepage "https://github.com/$REPO"
  version "$version"
  license "MIT"

  on_macos do
    on_intel do
      url "${platforms[darwin-amd64]}"
      sha256 "${checksums[darwin-amd64]}"
    end

    on_arm do
      url "${platforms[darwin-arm64]}"
      sha256 "${checksums[darwin-arm64]}"
    end
  end

  on_linux do
    on_intel do
      url "${platforms[linux-amd64]}"
      sha256 "${checksums[linux-amd64]}"
    end

    on_arm do
      url "${platforms[linux-arm64]}"
      sha256 "${checksums[linux-arm64]}"
    end
  end

  def install
    case Hardware::CPU.arch
    when :x86_64
      arch_suffix = "amd64"
    when :arm64
      arch_suffix = "arm64"
    else
      raise "Unsupported architecture: #{Hardware::CPU.arch}"
    end

    os_name = OS.mac? ? "darwin" : "linux"
    binary_name = "$APP_NAME-#{os_name}-#{arch_suffix}"
    
    bin.install binary_name => "$APP_NAME"
  end

  test do
    assert_match "$APP_NAME", shell_output("#{bin}/$APP_NAME --version")
    
    # Test basic functionality
    system bin/"$APP_NAME", "--help"
    
    # Test that the binary is properly linked
    assert_predicate bin/"$APP_NAME", :exist?
    assert_predicate bin/"$APP_NAME", :executable?
  end

  def caveats
    <<~EOS
      To get started with $APP_NAME_TITLE:
        $APP_NAME init

      For more information:
        $APP_NAME --help
        
      Documentation: https://github.com/$REPO/tree/main/docs-site
    EOS
  end
end
EOF

    print_success "✅ Updated $FORMULA_FILE with version $version"
}

# Validate formula syntax
validate_formula() {
    local project_root
    project_root=$(get_project_root)
    
    if command_exists brew; then
        print_status "Validating formula syntax..."
        if brew formula "$project_root/$FORMULA_FILE" >/dev/null 2>&1; then
            print_success "✅ Formula syntax is valid"
        else
            print_warning "⚠️  Formula syntax validation failed (brew not available or formula has issues)"
        fi
    else
        print_info "ℹ️  Skipping formula validation (brew not installed)"
    fi
}

# Show usage information
show_usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Update Homebrew formula with latest release checksums

OPTIONS:
    -v, --version VER    Use specific version instead of latest
    -h, --help           Show this help message
    --validate-only      Only validate existing formula

EXAMPLES:
    $0                   # Update with latest release
    $0 -v v1.2.3         # Update with specific version
    $0 --validate-only   # Just validate current formula

EOF
}

# Parse command line arguments
parse_args() {
    VERSION=""
    VALIDATE_ONLY=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -v|--version)
                VERSION="$2"
                shift 2
                ;;
            --validate-only)
                VALIDATE_ONLY=true
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
}

# Main function
main() {
    parse_args "$@"
    
    print_header "🍺 Homebrew Formula Updater"
    echo
    
    check_dependencies curl shasum || exit 1
    
    if [[ "$VALIDATE_ONLY" == "true" ]]; then
        validate_formula
        exit 0
    fi
    
    # Get version to use
    local version
    if [[ -n "$VERSION" ]]; then
        version="$VERSION"
        print_info "Using specified version: $version"
    else
        version=$(get_release_info) || exit 1
        print_info "Using latest version: $version"
    fi
    
    # Update formula
    update_formula "$version" || exit 1
    
    # Validate result
    validate_formula
    
    print_success "🎉 Homebrew formula updated successfully!"
    print_info "Formula location: $(get_project_root)/$FORMULA_FILE"
}

# Run main function with all arguments
main "$@"

#!/bin/bash
# Common utilities for otto-stack scripts

# Get script directory
get_script_dir() {
    cd "$(dirname "${BASH_SOURCE[1]}")" && pwd
}

# Get project root directory
get_project_root() {
    local script_dir
    script_dir=$(get_script_dir)
    cd "$script_dir/.." && pwd
}

# Load constants from Go code
load_constants() {
    local script_dir
    script_dir=$(get_script_dir)
    
    # Source the extracted constants
    # shellcheck source=scripts/extract-constants.sh
    source "$script_dir/extract-constants.sh"
}

# Initialize constants
load_constants

# Colors
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly CYAN='\033[0;36m'
readonly BOLD='\033[1m'
readonly NC='\033[0m'

# Print functions with consistent formatting
print_status() { echo -e "${BLUE}🔧${NC} $1"; }
print_success() { echo -e "${GREEN}✅${NC} $1"; }
print_warning() { echo -e "${YELLOW}⚠️${NC} $1"; }
print_error() { echo -e "${RED}❌${NC} $1" >&2; }
print_info() { echo -e "${CYAN}ℹ️${NC} $1"; }
print_header() { echo -e "${BOLD}${BLUE}$1${NC}"; }

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Validate required commands
check_dependencies() {
    local missing_deps=()
    
    for cmd in "$@"; do
        if ! command_exists "$cmd"; then
            missing_deps+=("$cmd")
        fi
    done
    
    if [[ ${#missing_deps[@]} -gt 0 ]]; then
        print_error "Missing required dependencies: ${missing_deps[*]}"
        print_error "Please install them and try again"
        return 1
    fi
    
    return 0
}

# Detect platform (OS and architecture)
detect_platform() {
    local os arch
    
    case "$(uname -s)" in
        Linux*)   os="linux" ;;
        Darwin*)  os="darwin" ;;
        CYGWIN*|MINGW*|MSYS*) os="windows" ;;
        *) print_error "Unsupported OS: $(uname -s)"; return 1 ;;
    esac
    
    case "$(uname -m)" in
        x86_64|amd64) arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        *) print_error "Unsupported architecture: $(uname -m)"; return 1 ;;
    esac
    
    echo "${os}-${arch}"
}

# Get latest release version from GitHub API
get_latest_version() {
    local repo="$1"
    local version
    
    if ! command_exists curl; then
        print_error "curl is required to fetch version information"
        return 1
    fi
    
    version=$(curl -s "https://api.github.com/repos/${repo}/releases/latest" | \
              grep '"tag_name":' | \
              sed -E 's/.*"([^"]+)".*/\1/')
    
    if [[ -z "$version" ]]; then
        print_error "Failed to fetch latest version for ${repo}"
        return 1
    fi
    
    echo "$version"
}

# Create temporary directory with cleanup trap
create_temp_dir() {
    local temp_dir
    temp_dir=$(mktemp -d)
    
    # Set up cleanup trap
    trap "rm -rf '$temp_dir'" EXIT
    
    echo "$temp_dir"
}

# Download file with progress
download_file() {
    local url="$1"
    local output="$2"
    local show_progress="${3:-true}"
    
    if [[ "$show_progress" == "true" ]]; then
        curl -fsSL --progress-bar -o "$output" "$url"
    else
        curl -fsSL -o "$output" "$url"
    fi
}

# Verify file exists and is executable
verify_executable() {
    local file="$1"
    
    if [[ ! -f "$file" ]]; then
        print_error "File not found: $file"
        return 1
    fi
    
    if [[ ! -x "$file" ]]; then
        print_error "File is not executable: $file"
        return 1
    fi
    
    return 0
}

# Install file to directory with proper permissions
install_file() {
    local source="$1"
    local target_dir="$2"
    local filename="$3"
    local target_path="${target_dir}/${filename}"
    
    if [[ -w "$target_dir" ]]; then
        cp "$source" "$target_path"
        chmod +x "$target_path"
    else
        print_warning "Installing to ${target_dir} (requires sudo)..."
        sudo cp "$source" "$target_path"
        sudo chmod +x "$target_path"
    fi
    
    echo "$target_path"
}

# Show confirmation prompt
confirm() {
    local message="$1"
    local default="${2:-n}"
    local prompt
    
    case "$default" in
        y|Y) prompt="[Y/n]" ;;
        n|N) prompt="[y/N]" ;;
        *) prompt="[y/n]" ;;
    esac
    
    read -p "$message $prompt " -n 1 -r
    echo
    
    case "$default" in
        y|Y) [[ $REPLY =~ ^[Nn]$ ]] && return 1 || return 0 ;;
        n|N) [[ $REPLY =~ ^[Yy]$ ]] && return 0 || return 1 ;;
        *) [[ $REPLY =~ ^[Yy]$ ]] && return 0 || return 1 ;;
    esac
}

# Export functions for use in other scripts
export -f print_status print_success print_warning print_error print_info print_header
export -f command_exists get_script_dir get_project_root check_dependencies
export -f detect_platform get_latest_version create_temp_dir download_file
export -f verify_executable install_file confirm load_constants

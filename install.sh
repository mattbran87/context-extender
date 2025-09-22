#!/bin/bash

# Context-Extender Installation Script
# Auto-detects platform and architecture for macOS and Linux
# Usage: curl -fsSL https://raw.githubusercontent.com/mattbran87/context-extender/master/install.sh | sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Configuration
REPO="mattbran87/context-extender"
VERSION="1.2.0"
BINARY_NAME="context-extender"

# Print colored output
print_status() {
    echo -e "${BLUE}ℹ️  ${NC}$1"
}

print_success() {
    echo -e "${GREEN}✅ ${NC}$1"
}

print_warning() {
    echo -e "${YELLOW}⚠️  ${NC}$1"
}

print_error() {
    echo -e "${RED}❌ ${NC}$1"
}

print_header() {
    echo -e "${BOLD}🎉 Context-Extender Installation Script${NC}"
    echo -e "${BOLD}======================================${NC}"
    echo
}

# Detect operating system
detect_os() {
    case "$(uname -s)" in
        Darwin)
            OS="darwin"
            OS_NAME="macOS"
            ;;
        Linux)
            OS="linux"
            OS_NAME="Linux"
            ;;
        *)
            print_error "Unsupported operating system: $(uname -s)"
            print_error "This script supports macOS and Linux only."
            print_error "For Windows, please download the .exe file from:"
            print_error "https://github.com/${REPO}/releases/tag/v${VERSION}"
            exit 1
            ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)
            ARCH="amd64"
            ARCH_NAME="x86_64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ARCH_NAME="ARM64"
            ;;
        *)
            print_error "Unsupported architecture: $(uname -m)"
            print_error "Supported architectures: x86_64, ARM64"
            exit 1
            ;;
    esac
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."

    if ! command_exists curl; then
        print_error "curl is required but not installed."
        print_error "Please install curl and try again."
        exit 1
    fi

    if ! command_exists chmod; then
        print_error "chmod is required but not installed."
        exit 1
    fi

    print_success "Prerequisites check passed"
}

# Get installation directory
get_install_dir() {
    # Try to find a good installation directory
    if [ -w "/usr/local/bin" ]; then
        INSTALL_DIR="/usr/local/bin"
    elif [ -w "$HOME/.local/bin" ]; then
        INSTALL_DIR="$HOME/.local/bin"
        # Create directory if it doesn't exist
        mkdir -p "$INSTALL_DIR"
    elif [ -w "$HOME/bin" ]; then
        INSTALL_DIR="$HOME/bin"
        # Create directory if it doesn't exist
        mkdir -p "$INSTALL_DIR"
    else
        INSTALL_DIR="$PWD"
        print_warning "No writable PATH directory found, installing to current directory"
    fi
}

# Download binary
download_binary() {
    BINARY_FILE="${BINARY_NAME}-${VERSION}-${OS}-${ARCH}"
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/v${VERSION}/${BINARY_FILE}"
    TARGET_PATH="${INSTALL_DIR}/${BINARY_NAME}"

    print_status "Downloading Context-Extender v${VERSION} for ${OS_NAME} ${ARCH_NAME}..."
    print_status "URL: ${DOWNLOAD_URL}"

    if curl -fsSL -o "${TARGET_PATH}" "${DOWNLOAD_URL}"; then
        print_success "Downloaded successfully to ${TARGET_PATH}"
    else
        print_error "Failed to download binary from ${DOWNLOAD_URL}"
        print_error "Please check your internet connection and try again."
        exit 1
    fi

    # Make binary executable
    chmod +x "${TARGET_PATH}"
    print_success "Made binary executable"
}

# Check if binary is in PATH
check_path() {
    if [ "$INSTALL_DIR" != "$PWD" ]; then
        if echo "$PATH" | grep -q "$INSTALL_DIR"; then
            print_success "Installation directory is in PATH"
            BINARY_CMD="context-extender"
        else
            print_warning "Installation directory ${INSTALL_DIR} is not in PATH"
            print_status "You may need to add it to your PATH or use the full path"
            BINARY_CMD="${TARGET_PATH}"
        fi
    else
        print_warning "Installed to current directory"
        BINARY_CMD="./context-extender"
    fi
}

# Verify installation
verify_installation() {
    print_status "Verifying installation..."

    if "${BINARY_CMD}" version >/dev/null 2>&1; then
        VERSION_OUTPUT=$("${BINARY_CMD}" version)
        print_success "Installation verified!"
        print_status "Version: $(echo "$VERSION_OUTPUT" | head -n1)"
    else
        print_error "Installation verification failed"
        print_error "Binary may not be working correctly"
        exit 1
    fi
}

# Run installation wizard
run_wizard() {
    echo
    print_status "Ready to run Context-Extender installation wizard!"
    print_status "The wizard will:"
    print_status "  • Check Claude Code installation"
    print_status "  • Install conversation capture hooks"
    print_status "  • Initialize database"
    print_status "  • Test the setup"
    echo

    read -p "Do you want to run the installation wizard now? (y/N): " -r
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo
        print_status "Starting Context-Extender installation wizard..."
        echo

        if "${BINARY_CMD}" install; then
            echo
            print_success "Context-Extender installation completed successfully!"
        else
            echo
            print_warning "Installation wizard encountered issues"
            print_status "You can run it again later with: ${BINARY_CMD} install"
        fi
    else
        echo
        print_success "Context-Extender binary installed successfully!"
        print_status "To run the installation wizard later:"
        print_status "  ${BINARY_CMD} install"
        print_status ""
        print_status "For quick setup without wizard:"
        print_status "  ${BINARY_CMD} configure"
    fi
}

# Show completion message
show_completion() {
    echo
    print_success "Installation Summary:"
    print_status "  • Binary: ${TARGET_PATH}"
    print_status "  • Command: ${BINARY_CMD}"
    print_status "  • Version: ${VERSION}"
    print_status "  • Platform: ${OS_NAME} ${ARCH_NAME}"
    echo
    print_status "Next steps:"
    print_status "  • Run '${BINARY_CMD} --help' to see all commands"
    print_status "  • Run '${BINARY_CMD} configure --status' to check setup"
    print_status "  • Visit: https://github.com/${REPO} for documentation"
    echo
    print_success "Thank you for using Context-Extender! 🎉"
}

# Main installation function
main() {
    print_header

    print_status "Auto-detecting your platform..."
    detect_os
    detect_arch
    print_success "Detected: ${OS_NAME} ${ARCH_NAME}"
    echo

    check_prerequisites
    echo

    get_install_dir
    print_status "Installing to: ${INSTALL_DIR}"
    echo

    download_binary
    echo

    check_path
    echo

    verify_installation
    echo

    run_wizard

    show_completion
}

# Handle script interruption
trap 'echo; print_error "Installation interrupted by user"; exit 1' INT

# Run main function
main "$@"
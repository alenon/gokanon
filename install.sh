#!/bin/bash

# GoKanon Installation Script
# This script installs the latest version of gokanon

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REPO="alenon/gokanon"
BINARY_NAME="gokanon"
INSTALL_DIR="/usr/local/bin"
USE_SUDO=true

# Print colored output
print_info() {
    echo -e "${GREEN}==>${NC} $1"
}

print_error() {
    echo -e "${RED}Error:${NC} $1" >&2
}

print_warning() {
    echo -e "${YELLOW}Warning:${NC} $1"
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Darwin*)
            OS="darwin"
            ;;
        Linux*)
            OS="linux"
            ;;
        MINGW*|MSYS*|CYGWIN*)
            print_error "Windows is not supported by this install script. Please download the binary manually from:"
            print_error "https://github.com/${REPO}/releases/latest"
            exit 1
            ;;
        *)
            print_error "Unsupported operating system: $(uname -s)"
            exit 1
            ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac
}

# Check if running as root
check_sudo() {
    if [ "$EUID" -eq 0 ]; then
        USE_SUDO=false
    elif [ ! -w "$INSTALL_DIR" ]; then
        if ! command -v sudo >/dev/null 2>&1; then
            print_warning "sudo not found and $INSTALL_DIR is not writable"
            print_info "Installing to ~/.local/bin instead"
            INSTALL_DIR="$HOME/.local/bin"
            USE_SUDO=false
            mkdir -p "$INSTALL_DIR"
        fi
    else
        USE_SUDO=false
    fi
}

# Get latest release version
get_latest_version() {
    print_info "Fetching latest release version..."

    # Try using GitHub API
    if command -v curl >/dev/null 2>&1; then
        VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command -v wget >/dev/null 2>&1; then
        VERSION=$(wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        print_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi

    if [ -z "$VERSION" ]; then
        print_error "Failed to fetch latest version"
        exit 1
    fi

    print_info "Latest version: $VERSION"
}

# Download binary
download_binary() {
    BINARY_FILE="${BINARY_NAME}-${OS}-${ARCH}"
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_FILE}.tar.gz"
    TMP_DIR=$(mktemp -d)

    print_info "Downloading from $DOWNLOAD_URL"

    cd "$TMP_DIR"

    if command -v curl >/dev/null 2>&1; then
        curl -L -o "${BINARY_FILE}.tar.gz" "$DOWNLOAD_URL"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "${BINARY_FILE}.tar.gz" "$DOWNLOAD_URL"
    fi

    if [ ! -f "${BINARY_FILE}.tar.gz" ]; then
        print_error "Failed to download binary"
        cd - > /dev/null
        rm -rf "$TMP_DIR"
        exit 1
    fi

    print_info "Extracting binary..."
    tar -xzf "${BINARY_FILE}.tar.gz"

    if [ ! -f "$BINARY_FILE" ]; then
        print_error "Failed to extract binary"
        cd - > /dev/null
        rm -rf "$TMP_DIR"
        exit 1
    fi

    chmod +x "$BINARY_FILE"
}

# Install binary
install_binary() {
    print_info "Installing to $INSTALL_DIR/$BINARY_NAME"

    if [ "$USE_SUDO" = true ]; then
        sudo mv "$TMP_DIR/$BINARY_FILE" "$INSTALL_DIR/$BINARY_NAME"
    else
        mv "$TMP_DIR/$BINARY_FILE" "$INSTALL_DIR/$BINARY_NAME"
    fi

    # Clean up
    cd - > /dev/null
    rm -rf "$TMP_DIR"
}

# Verify installation
verify_installation() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        print_info "Installation successful!"
        print_info "Version: $($BINARY_NAME --version 2>/dev/null || echo 'installed')"
        echo ""
        print_info "Run '$BINARY_NAME --help' to get started"
    else
        print_warning "Installation complete, but $BINARY_NAME is not in PATH"
        print_info "Add $INSTALL_DIR to your PATH:"
        echo ""
        echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
        echo ""
        print_info "Add this to your shell profile (~/.bashrc, ~/.zshrc, etc.)"
    fi
}

# Main installation flow
main() {
    echo ""
    print_info "GoKanon Installer"
    echo ""

    # Detect system
    detect_os
    detect_arch
    print_info "Detected system: $OS-$ARCH"

    # Check installation permissions
    check_sudo

    # Get latest version
    get_latest_version

    # Download and install
    download_binary
    install_binary

    # Verify
    verify_installation

    echo ""
    print_info "For more information, visit: https://github.com/${REPO}"
    echo ""
}

# Run main function
main

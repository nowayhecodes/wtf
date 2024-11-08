#!/bin/bash

set -e # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Installing WTF...${NC}"

# Detect architecture and OS
ARCH=$(uname -m)
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

# Convert architecture names
case ${ARCH} in
x86_64)
    ARCH="x86_64"
    ;;
aarch64)
    ARCH="aarch64"
    ;;
*)
    echo -e "${RED}Unsupported architecture: ${ARCH}${NC}"
    exit 1
    ;;
esac

# Allow custom installation directory
DEFAULT_INSTALL_DIR="/usr/local/bin"
INSTALL_DIR=${RPM_INSTALL_DIR:-$DEFAULT_INSTALL_DIR}
RPM_HOME="${HOME}/.rpm"

# Create temporary directory
TMP_DIR=$(mktemp -d)
cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

# Download binary
echo "Downloading WTF..."
BINARY_URL="https://github.com/nowayhecodes/wtf/releases/latest/download/wtf-${OS}-${ARCH}.tar.gz"
if ! curl -L --progress-bar "$BINARY_URL" -o "$TMP_DIR/wtf.tar.gz"; then
    echo -e "${RED}Failed to download WTF${NC}"
    exit 1
fi

# Extract binary
cd "$TMP_DIR"
tar xzf wtf.tar.gz

# Create WTF home directory
mkdir -p "$WTF_HOME"
mkdir -p "$WTF_HOME/bin"

# Install binary
if [ ! -w "$INSTALL_DIR" ]; then
    echo "Escalating privileges to install to $INSTALL_DIR"
    sudo mv "$TMP_DIR/wtf" "$INSTALL_DIR/wtf"
    sudo chmod +x "$INSTALL_DIR/wtf"
else
    mv "$TMP_DIR/wtf" "$INSTALL_DIR/wtf"
    chmod +x "$INSTALL_DIR/wtf"
fi

# Update shell configuration
update_shell_config() {
    local shell_config="$1"
    local updated=false

    # Create config file if it doesn't exist
    touch "$shell_config"

    # Check if RPM_HOME is already in the config
    if ! grep -q "export WTF_HOME=" "$shell_config"; then
        echo -e "\n# WTF Configuration" >>"$shell_config"
        echo "export WTF_HOME=\"$WTF_HOME\"" >>"$shell_config"
        updated=true
    fi

    # Add to PATH if custom installation directory is used
    if [ "$INSTALL_DIR" != "$DEFAULT_INSTALL_DIR" ]; then
        if ! grep -q "export PATH=.*$INSTALL_DIR" "$shell_config"; then
            echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >>"$shell_config"
            updated=true
        fi
    fi

    # Add RPM_HOME/bin to PATH
    if ! grep -q "export PATH=.*\$WTF_HOME/bin" "$shell_config"; then
        echo "export PATH=\"\$PATH:\$WTF_HOME/bin\"" >>"$shell_config"
        updated=true
    fi

    if [ "$updated" = true ]; then
        echo -e "${GREEN}Updated ${shell_config}${NC}"
    fi
}

# Update shell configurations
if [ -n "$BASH_VERSION" ]; then
    update_shell_config "$HOME/.bashrc"
elif [ -n "$ZSH_VERSION" ]; then
    update_shell_config "$HOME/.zshrc"
fi

# Verify installation
if command -v rpm >/dev/null; then
    echo -e "${GREEN}WTF has been successfully installed!${NC}"
    echo -e "You can now use WTF by running: ${BLUE}wtf install <package>${NC}"
    echo -e "\nEnvironment variables set:"
    echo -e "${BLUE}WTF_HOME=${WTF_HOME}${NC}"
    echo -e "${BLUE}PATH includes: ${INSTALL_DIR}${NC}"
    echo -e "\nPlease restart your shell or run:"
    echo -e "${BLUE}source ~/.bashrc${NC} (for bash)"
    echo -e "${BLUE}source ~/.zshrc${NC} (for zsh)"
else
    echo -e "${RED}Installation failed. Please try again or install manually.${NC}"
    exit 1
fi

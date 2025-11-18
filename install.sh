#!/bin/bash
#
# LineSense Installation Script
#
# This script automates the installation and setup of LineSense.
# It handles building, installation, shell integration, and initial configuration.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/Traves-Theberge/LineSense/main/install.sh | bash
#   # or
#   ./install.sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Configuration
REPO_URL="https://github.com/Traves-Theberge/LineSense.git"
INSTALL_DIR="${HOME}/.local/bin"
CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/linesense"

# Print banner
print_banner() {
    echo -e "${CYAN}"
    echo "   __    _            ____                       "
    echo "  / /   (_)_ __   ___/ ___|  ___ _ __  ___  ___  "
    echo " / /    | | '_ \ / _ \___ \ / _ \ '_ \/ __|/ _ \ "
    echo "/ /___  | | | | |  __/___) |  __/ | | \__ \  __/ "
    echo "\____/  |_|_| |_|\___|____/ \___|_| |_|___/\___| "
    echo -e "${NC}"
    echo -e "${BOLD}AI-Powered Shell Assistant${NC}"
    echo "----------------------------------------"
}

# Print colored messages
info() {
    echo -e "${BLUE}ℹ️  [INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}✅ [SUCCESS]${NC} $1"
}

warn() {
    echo -e "${YELLOW}⚠️  [WARN]${NC} $1"
}

error() {
    echo -e "${RED}❌ [ERROR]${NC} $1"
    exit 1
}

step() {
    echo -e "\n${BOLD}👉 $1${NC}"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Detect user's shell
detect_shell() {
    if [ -n "$BASH_VERSION" ]; then
        echo "bash"
    elif [ -n "$ZSH_VERSION" ]; then
        echo "zsh"
    else
        # Fallback to $SHELL
        basename "$SHELL"
    fi
}

# Get shell RC file
get_shell_rc() {
    local shell_name="$1"
    case "$shell_name" in
        bash)
            if [ -f "$HOME/.bashrc" ]; then
                echo "$HOME/.bashrc"
            else
                echo "$HOME/.bash_profile"
            fi
            ;;
        zsh)
            echo "$HOME/.zshrc"
            ;;
        *)
            warn "Unknown shell: $shell_name"
            echo ""
            ;;
    esac
}

# Check prerequisites
check_prerequisites() {
    step "Checking prerequisites..."

    # Check for Go
    if ! command_exists go; then
        error "Go is not installed. Please install Go 1.21 or later from https://golang.org/dl/"
    fi

    # Check Go version
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    GO_MAJOR=$(echo "$GO_VERSION" | cut -d. -f1)
    GO_MINOR=$(echo "$GO_VERSION" | cut -d. -f2)

    if [ "$GO_MAJOR" -lt 1 ] || ([ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -lt 21 ]); then
        error "Go version 1.21 or later is required. Found: $GO_VERSION"
    fi

    success "Go $GO_VERSION detected"

    # Check for git
    if ! command_exists git; then
        error "Git is not installed. Please install git first."
    fi

    success "Git detected"
}

# Clone or update repository
setup_repository() {
    info "Setting up LineSense repository..."

    local tmp_dir=$(mktemp -d)
    cd "$tmp_dir"

    info "Cloning LineSense from $REPO_URL..."
    if git clone --depth 1 "$REPO_URL" linesense; then
        cd linesense
        success "Repository cloned successfully"
    else
        error "Failed to clone repository"
    fi
}

# Build LineSense
build_linesense() {
    info "Building LineSense..."

    if go build -o linesense ./cmd/linesense; then
        success "Build completed successfully"
    else
        error "Build failed"
    fi
}

# Install binary
install_binary() {
    info "Installing LineSense to $INSTALL_DIR..."

    # Create install directory if it doesn't exist
    mkdir -p "$INSTALL_DIR"

    # Copy binary
    if cp linesense "$INSTALL_DIR/linesense"; then
        chmod +x "$INSTALL_DIR/linesense"
        success "Binary installed to $INSTALL_DIR/linesense"
    else
        error "Failed to install binary"
    fi

    # Check if install dir is in PATH
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        warn "$INSTALL_DIR is not in your PATH"
        info "Add the following to your shell RC file:"
        echo "    export PATH=\"\$PATH:$INSTALL_DIR\""
    fi
}

# Install shell integration
install_shell_integration() {
    local shell_name=$(detect_shell)
    local rc_file=$(get_shell_rc "$shell_name")

    if [ -z "$rc_file" ]; then
        warn "Could not determine shell RC file. Skipping shell integration."
        return
    fi

    info "Installing shell integration for $shell_name..."

    # Copy shell integration script
    local integration_dir="$CONFIG_DIR/shell"
    mkdir -p "$integration_dir"

    if [ -f "scripts/linesense.$shell_name" ]; then
        cp "scripts/linesense.$shell_name" "$integration_dir/"
        chmod +x "$integration_dir/linesense.$shell_name"
        success "Shell integration script copied to $integration_dir"
    else
        warn "Shell integration script not found for $shell_name"
        return
    fi

    # Add sourcing to RC file
    local source_line="[ -f \"$integration_dir/linesense.$shell_name\" ] && source \"$integration_dir/linesense.$shell_name\""

    if grep -q "linesense.$shell_name" "$rc_file" 2>/dev/null; then
        info "Shell integration already present in $rc_file"
    else
        info "Adding shell integration to $rc_file..."
        echo "" >> "$rc_file"
        echo "# LineSense AI Shell Assistant" >> "$rc_file"
        echo "$source_line" >> "$rc_file"
        success "Shell integration added to $rc_file"
    fi
}

# Initialize configuration
init_config() {
    info "Initializing configuration..."

    if "$INSTALL_DIR/linesense" config init; then
        success "Configuration initialized"
    else
        warn "Configuration initialization failed. You can run 'linesense config init' manually."
    fi
}

# Print next steps
print_next_steps() {
    echo ""
    echo -e "${GREEN}╔════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║          LineSense Installation Complete! 🎉                   ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    info "Next steps:"
    echo ""
    echo "  1. Restart your shell or run:"
    echo -e "     ${BLUE}source $(get_shell_rc $(detect_shell))${NC}"
    echo ""
    echo "  2. Set your OpenRouter API key:"
    echo -e "     ${BLUE}linesense config set-key${NC}"
    echo ""
    echo "  3. Try it out:"
    echo -e "     ${BLUE}linesense suggest \"list files sorted by size\"${NC}"
    echo -e "     ${BLUE}linesense explain \"git rebase -i HEAD~3\"${NC}"
    echo ""
    echo "  4. Use shell integration (Ctrl+Space for suggestions)"
    echo ""
    echo "  📚 Documentation: https://github.com/Traves-Theberge/LineSense"
    echo "  🐛 Report issues: https://github.com/Traves-Theberge/LineSense/issues"
    echo ""
}

# Cleanup
cleanup() {
    if [ -n "$tmp_dir" ] && [ -d "$tmp_dir" ]; then
        cd /
        rm -rf "$tmp_dir"
    fi
}

# Main installation flow
main() {
    print_banner

    # Set trap for cleanup
    trap cleanup EXIT

    # Run installation steps
    check_prerequisites
    setup_repository
    build_linesense
    install_binary
    install_shell_integration
    init_config
    print_next_steps

    success "Installation completed successfully!"
}

# Run main function
main "$@"

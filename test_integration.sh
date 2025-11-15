#!/bin/bash
#
# LineSense Shell Integration Test
# Run this script in your terminal to test the shell integration
#

echo "======================================"
echo "LineSense Shell Integration Test"
echo "======================================"
echo ""

# Check if linesense binary exists
if ! command -v linesense &> /dev/null; then
    echo "❌ ERROR: linesense binary not found in PATH"
    echo "   Run: go install ./cmd/linesense"
    exit 1
fi

echo "✅ linesense binary found: $(which linesense)"
echo "   Version: $(linesense --version)"
echo ""

# Check if shell integration file exists
INTEGRATION_FILE="$HOME/.config/linesense/shell/linesense.bash"
if [ ! -f "$INTEGRATION_FILE" ]; then
    echo "❌ ERROR: Shell integration file not found"
    echo "   Expected: $INTEGRATION_FILE"
    exit 1
fi

echo "✅ Shell integration file found: $INTEGRATION_FILE"
echo ""

# Check if it's sourced in .bashrc
if grep -q "linesense" ~/.bashrc; then
    echo "✅ Shell integration is in .bashrc"
    grep "linesense" ~/.bashrc
else
    echo "⚠️  WARNING: Shell integration not found in .bashrc"
    echo "   Add this line to your .bashrc:"
    echo '   [ -f "$HOME/.config/linesense/shell/linesense.bash" ] && source "$HOME/.config/linesense/shell/linesense.bash"'
fi
echo ""

# Check if API key is set
if [ -f "$HOME/.config/linesense/config.toml" ]; then
    echo "✅ Configuration file exists"
    if grep -q "api_key" "$HOME/.config/linesense/config.toml"; then
        echo "✅ API key is configured"
    else
        echo "⚠️  WARNING: No API key found"
        echo "   Run: linesense config set-key"
    fi
else
    echo "⚠️  WARNING: No configuration file"
    echo "   Run: linesense config init"
fi
echo ""

echo "======================================"
echo "Testing CLI Commands"
echo "======================================"
echo ""

# Test suggest command
echo "Testing: linesense suggest --line \"list files\""
linesense suggest --line "list files" 2>&1 | head -10
echo ""

echo "======================================"
echo "Shell Integration Instructions"
echo "======================================"
echo ""
echo "The shell integration ONLY works in interactive bash sessions."
echo "It will NOT work in:"
echo "  - Script execution (non-interactive)"
echo "  - This test script output"
echo "  - Claude Code terminal interface"
echo ""
echo "To test the keybindings:"
echo "  1. Open a NEW bash terminal (or run: bash)"
echo "  2. The integration should auto-load (you'll see a message)"
echo "  3. Type: npm"
echo "  4. Press: Ctrl+Space"
echo "  5. You should see AI suggestions!"
echo ""
echo "Alternative test:"
echo "  1. Type: docker ps -a"
echo "  2. Press: Ctrl+X then Ctrl+E"
echo "  3. You should see command explanation!"
echo ""
echo "If it still doesn't work, manually source it:"
echo "  source ~/.config/linesense/shell/linesense.bash"
echo ""

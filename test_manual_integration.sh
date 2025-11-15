#!/bin/bash
#
# Manual Shell Integration Test
#
# This script helps you verify that LineSense shell integration works correctly.

set -e

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║       LineSense Manual Shell Integration Test                 ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}Step 1: Verifying shell integration scripts...${NC}"
if [ -f "$HOME/.config/linesense/shell/linesense.bash" ]; then
    echo -e "${GREEN}✓${NC} Bash integration script found"
else
    echo "✗ Bash integration script NOT found"
    exit 1
fi

if [ -f "$HOME/.config/linesense/shell/linesense.zsh" ]; then
    echo -e "${GREEN}✓${NC} Zsh integration script found"
else
    echo "✗ Zsh integration script NOT found (optional)"
fi

echo ""
echo -e "${BLUE}Step 2: Testing bash script syntax...${NC}"
if bash -n "$HOME/.config/linesense/shell/linesense.bash"; then
    echo -e "${GREEN}✓${NC} Bash script syntax is valid"
else
    echo "✗ Bash script has syntax errors"
    exit 1
fi

echo ""
echo -e "${BLUE}Step 3: Sourcing bash integration...${NC}"
source "$HOME/.config/linesense/shell/linesense.bash"
echo -e "${GREEN}✓${NC} Bash integration sourced successfully"

echo ""
echo -e "${BLUE}Step 4: Checking if functions are defined...${NC}"

# Check if the main functions exist
if type _linesense_suggest_widget &>/dev/null; then
    echo -e "${GREEN}✓${NC} _linesense_suggest_widget function defined"
else
    echo "✗ _linesense_suggest_widget function NOT defined"
    exit 1
fi

if type _linesense_explain_widget &>/dev/null; then
    echo -e "${GREEN}✓${NC} _linesense_explain_widget function defined"
else
    echo "✗ _linesense_explain_widget function NOT defined"
    exit 1
fi

echo ""
echo -e "${BLUE}Step 5: Checking keybindings...${NC}"

# Check if keybindings are set (bash uses bind command)
if bind -p | grep -q "_linesense_suggest_widget"; then
    echo -e "${GREEN}✓${NC} Suggest keybinding is set"
    SUGGEST_KEY=$(bind -p | grep "_linesense_suggest_widget" | awk '{print $1}' | head -1)
    echo "  Bound to: $SUGGEST_KEY"
else
    echo "✗ Suggest keybinding NOT set"
fi

if bind -p | grep -q "_linesense_explain_widget"; then
    echo -e "${GREEN}✓${NC} Explain keybinding is set"
    EXPLAIN_KEY=$(bind -p | grep "_linesense_explain_widget" | awk '{print $1}' | head -1)
    echo "  Bound to: $EXPLAIN_KEY"
else
    echo "✗ Explain keybinding NOT set"
fi

echo ""
echo -e "${BLUE}Step 6: Testing linesense command availability...${NC}"
if command -v linesense &>/dev/null; then
    echo -e "${GREEN}✓${NC} linesense command is available"
    echo "  Version: $(linesense --version 2>&1)"
else
    echo "✗ linesense command NOT found in PATH"
    exit 1
fi

echo ""
echo "════════════════════════════════════════════════════════════════"
echo -e "${GREEN}All automated checks passed!${NC}"
echo "════════════════════════════════════════════════════════════════"
echo ""

echo -e "${YELLOW}Manual Testing Instructions:${NC}"
echo ""
echo "To test the interactive features:"
echo ""
echo "1. Open a new terminal or run:"
echo "   source ~/.bashrc"
echo ""
echo "2. Add this line to your ~/.bashrc:"
echo "   [ -f \"\$HOME/.config/linesense/shell/linesense.bash\" ] && source \"\$HOME/.config/linesense/shell/linesense.bash\""
echo ""
echo "3. Test the suggest feature:"
echo "   - Type: list files"
echo "   - Press: Ctrl+Space"
echo "   - You should see AI suggestions"
echo ""
echo "4. Test the explain feature:"
echo "   - Type: docker ps -a"
echo "   - Press: Ctrl+X then Ctrl+E"
echo "   - You should see an explanation"
echo ""
echo "5. Try typing various commands and test the keybindings!"
echo ""
echo "════════════════════════════════════════════════════════════════"

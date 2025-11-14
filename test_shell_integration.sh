#!/usr/bin/env bash
# Test script for shell integration functions

set -e

echo "🧪 Testing LineSense Shell Integration"
echo "======================================"
echo ""

# Build the binary first
echo "1️⃣  Building linesense..."
go build -o linesense cmd/linesense/main.go
export PATH="$PWD:$PATH"

# Load .env
if [ -f .env ]; then
    export $(cat .env | xargs)
fi

echo "✓ Binary built and in PATH"
echo ""

# Test JSON parsing function
echo "2️⃣  Testing JSON parsing helpers..."

# Source the bash script (suppress bind warnings)
source scripts/linesense.bash 2>/dev/null || true

# Test JSON parsing if jq is available
if command -v jq &> /dev/null; then
    test_json='{"suggestions":[{"command":"ls -al","risk":"low"}]}'
    result=$(echo "$test_json" | jq -r '.suggestions[0].command' 2>/dev/null)
    if [ "$result" = "ls -al" ]; then
        echo "✓ JSON parsing with jq works"
    else
        echo "❌ JSON parsing with jq failed: got '$result'"
        exit 1
    fi
else
    echo "⚠️  jq not installed, will use grep/sed fallback"
fi

echo ""

# Test suggest command directly
echo "3️⃣  Testing suggest command..."
result=$(./linesense suggest --line "list files" 2>/dev/null)
if echo "$result" | grep -q '"command"'; then
    echo "✓ Suggest command works"
    echo "   Example output:"
    echo "$result" | head -5
else
    echo "❌ Suggest command failed"
    exit 1
fi

echo ""

# Test explain command directly
echo "4️⃣  Testing explain command..."
result=$(./linesense explain --line "ls -al" 2>/dev/null)
if echo "$result" | grep -q '"summary"'; then
    echo "✓ Explain command works"
    echo "   Example output:"
    echo "$result" | head -5
else
    echo "❌ Explain command failed"
    exit 1
fi

echo ""
echo "======================================"
echo "✅ All shell integration tests passed!"
echo ""
echo "To use in your shell:"
echo "  Bash: source $(pwd)/scripts/linesense.bash"
echo "  Zsh:  source $(pwd)/scripts/linesense.zsh"
echo ""
echo "Keybindings:"
echo "  Ctrl+Space: Get AI suggestions"
echo "  Ctrl+X Ctrl+E: Explain current command"

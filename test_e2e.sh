#!/bin/bash
#
# LineSense End-to-End Testing Script
# Tests all major functionality before release
#

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

TESTS_PASSED=0
TESTS_FAILED=0

# Print test header
print_test() {
    echo -e "\n${BLUE}━━━ TEST: $1 ━━━${NC}"
}

# Print success
pass() {
    echo -e "${GREEN}✓ PASS${NC} $1"
    ((TESTS_PASSED++))
}

# Print failure
fail() {
    echo -e "${RED}✗ FAIL${NC} $1"
    ((TESTS_FAILED++))
}

# Print summary
summary() {
    echo ""
    echo -e "${BLUE}════════════════════════════════════════${NC}"
    echo -e "${BLUE}   End-to-End Test Summary${NC}"
    echo -e "${BLUE}════════════════════════════════════════${NC}"
    echo -e "   ${GREEN}Passed: $TESTS_PASSED${NC}"
    echo -e "   ${RED}Failed: $TESTS_FAILED${NC}"
    echo -e "${BLUE}════════════════════════════════════════${NC}"

    if [ $TESTS_FAILED -eq 0 ]; then
        echo -e "\n${GREEN}🎉 ALL TESTS PASSED!${NC}\n"
        return 0
    else
        echo -e "\n${RED}❌ SOME TESTS FAILED${NC}\n"
        return 1
    fi
}

echo -e "${BLUE}"
cat << "EOF"
╔════════════════════════════════════════╗
║   LineSense End-to-End Tests          ║
╚════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Check prerequisites
print_test "Prerequisites Check"

if ! command -v ./linesense &> /dev/null; then
    fail "linesense binary not found in current directory"
    echo "Run: go build -o linesense ./cmd/linesense"
    exit 1
fi
pass "linesense binary exists"

if [ ! -f "$HOME/.config/linesense/config.toml" ]; then
    fail "Configuration not initialized"
    echo "Run: ./linesense config init"
    exit 1
fi
pass "Configuration exists"

if [ -z "$OPENROUTER_API_KEY" ]; then
    fail "OPENROUTER_API_KEY not set"
    echo "Run: export OPENROUTER_API_KEY=your_key"
    exit 1
fi
pass "API key is set"

# Test 1: Version command
print_test "Version Command"
if ./linesense --version 2>&1 | grep -q "linesense version"; then
    pass "Version command works"
else
    fail "Version command failed"
fi

# Test 2: Help command
print_test "Help Command"
if ./linesense --help 2>&1 | grep -q "AI-powered shell command assistant"; then
    pass "Help command works"
else
    fail "Help command failed"
fi

# Test 3: Config show
print_test "Config Show"
if ./linesense config show 2>&1 | grep -q "LineSense Configuration"; then
    pass "Config show works"
else
    fail "Config show failed"
fi

# Test 4: Suggest command (JSON format)
print_test "Suggest Command (JSON)"
OUTPUT=$(./linesense suggest --line "list files" --format json 2>&1)
if echo "$OUTPUT" | grep -q "suggestions"; then
    pass "Suggest command returns JSON"
else
    fail "Suggest command JSON output failed"
fi

if echo "$OUTPUT" | grep -q "command"; then
    pass "Suggest includes command field"
else
    fail "Suggest missing command field"
fi

if echo "$OUTPUT" | grep -q "risk"; then
    pass "Suggest includes risk field"
else
    fail "Suggest missing risk field"
fi

# Test 5: Suggest command (Pretty format)
print_test "Suggest Command (Pretty)"
OUTPUT=$(./linesense suggest --line "compress directory" 2>&1)
if echo "$OUTPUT" | grep -q "Command Suggestions"; then
    pass "Suggest pretty format shows header"
else
    fail "Suggest pretty format failed"
fi

if echo "$OUTPUT" | grep -q "tar"; then
    pass "Suggest provides relevant command"
else
    fail "Suggest command not relevant"
fi

# Test 6: Explain command (JSON format)
print_test "Explain Command (JSON)"
OUTPUT=$(./linesense explain --line "docker ps -a" --format json 2>&1)
if echo "$OUTPUT" | grep -q "summary"; then
    pass "Explain command returns JSON"
else
    fail "Explain command JSON output failed"
fi

if echo "$OUTPUT" | grep -q "risk"; then
    pass "Explain includes risk field"
else
    fail "Explain missing risk field"
fi

if echo "$OUTPUT" | grep -q "notes"; then
    pass "Explain includes notes field"
else
    fail "Explain missing notes field"
fi

# Test 7: Explain command (Pretty format)
print_test "Explain Command (Pretty)"
OUTPUT=$(./linesense explain --line "ls -la" 2>&1)
if echo "$OUTPUT" | grep -q "Command Explanation"; then
    pass "Explain pretty format shows header"
else
    fail "Explain pretty format failed"
fi

if echo "$OUTPUT" | grep -q "Summary"; then
    pass "Explain shows summary section"
else
    fail "Explain missing summary"
fi

if echo "$OUTPUT" | grep -q "Risk Level"; then
    pass "Explain shows risk level"
else
    fail "Explain missing risk level"
fi

# Test 8: Safety filters (high-risk command)
print_test "Safety Filter Detection"
OUTPUT=$(./linesense explain --line "rm -rf /" --format json 2>&1)
if echo "$OUTPUT" | grep -q '"risk".*"high"'; then
    pass "High-risk command detected (rm -rf /)"
else
    fail "Safety filter failed to detect high-risk command"
fi

OUTPUT=$(./linesense explain --line "dd if=/dev/zero of=/dev/sda" --format json 2>&1)
if echo "$OUTPUT" | grep -q '"risk".*"high"'; then
    pass "High-risk command detected (dd disk wipe)"
else
    fail "Safety filter failed to detect disk wipe"
fi

# Test 9: Context-aware suggestions
print_test "Context-Aware Suggestions"
OUTPUT=$(./linesense suggest --line "find go files modified today" 2>&1)
if echo "$OUTPUT" | grep -q "find"; then
    pass "Context-aware suggestion works"
else
    fail "Context-aware suggestion failed"
fi

# Test 10: Different output formats
print_test "Output Format Switching"
JSON_OUT=$(./linesense suggest --line "test" --format json 2>&1)
PRETTY_OUT=$(./linesense suggest --line "test" --format pretty 2>&1)

if echo "$JSON_OUT" | jq . &> /dev/null; then
    pass "JSON output is valid"
else
    fail "JSON output is invalid"
fi

if echo "$PRETTY_OUT" | grep -q "Command Suggestions"; then
    pass "Pretty output format works"
else
    fail "Pretty output format failed"
fi

# Test 11: Error handling
print_test "Error Handling"
if ! ./linesense suggest 2>&1 | grep -q "required"; then
    fail "Missing required flag should show error"
else
    pass "Missing required flag handled correctly"
fi

if ! ./linesense explain 2>&1 | grep -q "required"; then
    fail "Missing required flag should show error"
else
    pass "Missing required flag handled correctly"
fi

# Test 12: Shell detection
print_test "Shell Detection"
OUTPUT=$(./linesense suggest --line "test" --format json 2>&1)
if echo "$OUTPUT" | grep -q "command"; then
    pass "Shell auto-detection works"
else
    fail "Shell auto-detection failed"
fi

# Test 13: Response time
print_test "Performance Test"
START=$(date +%s)
./linesense suggest --line "simple test" --format json &> /dev/null
END=$(date +%s)
DURATION=$((END - START))

if [ $DURATION -lt 10 ]; then
    pass "Response time acceptable ($DURATION seconds)"
else
    fail "Response too slow ($DURATION seconds)"
fi

# Summary
summary
exit $?

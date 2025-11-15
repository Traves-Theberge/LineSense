#!/usr/bin/env bash
# LineSense bash integration
# Source this file in your ~/.bashrc:
#   source /path/to/linesense.bash

# Check if linesense is available
if ! command -v linesense &> /dev/null; then
    echo "Warning: linesense not found in PATH. Shell integration disabled." >&2
    return 1
fi

# Function to parse JSON (uses jq if available, falls back to grep/sed)
_linesense_parse_json() {
    local json="$1"
    local key="$2"

    if command -v jq &> /dev/null; then
        echo "$json" | jq -r "$key" 2>/dev/null
    else
        # Fallback: basic grep/sed parsing (works for simple cases)
        echo "$json" | grep -o "\"$key\":\"[^\"]*\"" | head -1 | sed "s/\"$key\":\"\(.*\)\"/\1/"
    fi
}

# Function to request a suggestion from linesense
_linesense_request() {
    local current_line="$READLINE_LINE"
    local cwd="$PWD"

    # Don't suggest for empty lines
    if [ -z "$current_line" ]; then
        return
    fi

    # Call linesense suggest and capture JSON output
    local result
    result=$(linesense suggest --shell bash --line "$current_line" --cwd "$cwd" 2>/dev/null)

    if [ $? -eq 0 ] && [ -n "$result" ]; then
        # Parse JSON and extract first suggestion
        local suggestion risk

        if command -v jq &> /dev/null; then
            suggestion=$(echo "$result" | jq -r '.suggestions[0].command' 2>/dev/null)
            risk=$(echo "$result" | jq -r '.suggestions[0].risk' 2>/dev/null)
        else
            suggestion=$(echo "$result" | grep -o '"command":"[^"]*"' | head -1 | sed 's/"command":"\(.*\)"/\1/')
            risk=$(echo "$result" | grep -o '"risk":"[^"]*"' | head -1 | sed 's/"risk":"\(.*\)"/\1/')
        fi

        if [ -n "$suggestion" ] && [ "$suggestion" != "null" ]; then
            # Show risk indicator for high-risk commands
            if [ "$risk" = "high" ]; then
                echo -e "\n⚠️  WARNING: High-risk command detected!" >&2
            fi

            # Replace current line with suggestion
            READLINE_LINE="$suggestion"
            READLINE_POINT=${#READLINE_LINE}
        fi
    fi
}

# Function to explain the current command
_linesense_explain() {
    local current_line="$READLINE_LINE"
    local cwd="$PWD"

    # Don't explain empty lines
    if [ -z "$current_line" ]; then
        return
    fi

    # Save cursor position and clear line
    echo -ne "\r\033[K" >&2

    # Call linesense explain and capture JSON output
    local result
    result=$(linesense explain --shell bash --line "$current_line" --cwd "$cwd" 2>/dev/null)

    if [ $? -eq 0 ] && [ -n "$result" ]; then
        local summary risk

        if command -v jq &> /dev/null; then
            summary=$(echo "$result" | jq -r '.summary' 2>/dev/null)
            risk=$(echo "$result" | jq -r '.risk' 2>/dev/null)
        else
            summary=$(echo "$result" | grep -o '"summary":"[^"]*"' | head -1 | sed 's/"summary":"\(.*\)"/\1/')
            risk=$(echo "$result" | grep -o '"risk":"[^"]*"' | head -1 | sed 's/"risk":"\(.*\)"/\1/')
        fi

        # Display formatted explanation
        echo "" >&2
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" >&2
        echo "📝 LineSense Explanation" >&2
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" >&2
        echo "" >&2
        echo "Command: $current_line" >&2
        echo "" >&2

        # Color-code risk level
        case "$risk" in
            high)
                echo -e "⚠️  Risk: \033[1;31mHIGH\033[0m" >&2
                ;;
            medium)
                echo -e "⚠️  Risk: \033[1;33mMEDIUM\033[0m" >&2
                ;;
            low)
                echo -e "✓ Risk: \033[1;32mLOW\033[0m" >&2
                ;;
        esac

        echo "" >&2
        echo "$summary" >&2
        echo "" >&2
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" >&2
        echo "" >&2
    else
        echo "❌ Failed to get explanation" >&2
    fi
}

# Default keybindings
# Override these by setting environment variables before sourcing this file:
#   export LINESENSE_SUGGEST_KEY="\C-s"  # Ctrl+S for suggest
#   export LINESENSE_EXPLAIN_KEY="\C-h"  # Ctrl+H for explain

# Suggest keybinding (default: Ctrl+Space)
LINESENSE_SUGGEST_KEY="${LINESENSE_SUGGEST_KEY:-"\C- "}"
bind -x "\"${LINESENSE_SUGGEST_KEY}\": _linesense_request"

# Explain keybinding (default: Ctrl+X Ctrl+E to avoid conflicts)
# Note: Ctrl+E conflicts with readline's end-of-line, so we use Ctrl+X Ctrl+E
LINESENSE_EXPLAIN_KEY="${LINESENSE_EXPLAIN_KEY:-"\C-x\C-e"}"
bind -x "\"${LINESENSE_EXPLAIN_KEY}\": _linesense_explain"


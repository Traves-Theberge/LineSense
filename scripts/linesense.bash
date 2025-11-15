#!/usr/bin/env bash
# LineSense bash integration
# Source this file in your ~/.bashrc:
#   source /path/to/linesense.bash

# Check if linesense is available
if ! command -v linesense &> /dev/null; then
    echo "Warning: linesense not found in PATH. Shell integration disabled." >&2
    return 1
fi

# Function to request a suggestion from linesense
_linesense_request() {
    local current_line="$READLINE_LINE"
    local cwd="$PWD"

    # Don't suggest for empty lines
    if [ -z "$current_line" ]; then
        return
    fi

    # Clear line and show the pretty output
    echo "" >&2

    # Call linesense suggest with pretty format
    linesense suggest --shell bash --line "$current_line" --cwd "$cwd" --format pretty

    # Note: We display the suggestions but don't auto-replace the line
    # User can manually copy the command they want
}

# Function to explain the current command
_linesense_explain() {
    local current_line="$READLINE_LINE"
    local cwd="$PWD"

    # Don't explain empty lines
    if [ -z "$current_line" ]; then
        return
    fi

    # Clear line and show the pretty output
    echo "" >&2

    # Call linesense explain with pretty format
    linesense explain --shell bash --line "$current_line" --cwd "$cwd" --format pretty
}

# Default keybindings
# Override these by setting environment variables before sourcing this file:
#   export LINESENSE_SUGGEST_KEY="\C-s"  # Ctrl+S for suggest
#   export LINESENSE_EXPLAIN_KEY="\C-h"  # Ctrl+H for explain

# Suggest keybinding (default: Ctrl+Space)
# Note: \C-@ is the readline notation for Ctrl+Space (ASCII NUL)
LINESENSE_SUGGEST_KEY="${LINESENSE_SUGGEST_KEY:-\C-@}"
bind -x "\"${LINESENSE_SUGGEST_KEY}\": _linesense_request"

# Explain keybinding (default: Ctrl+X)
LINESENSE_EXPLAIN_KEY="${LINESENSE_EXPLAIN_KEY:-\C-x}"
bind -x "\"${LINESENSE_EXPLAIN_KEY}\": _linesense_explain"


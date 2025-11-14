#!/usr/bin/env zsh
# LineSense zsh integration
# Source this file in your ~/.zshrc:
#   source /path/to/linesense.zsh

# Check if linesense is available
if ! command -v linesense &> /dev/null; then
    print "Warning: linesense not found in PATH. Shell integration disabled." >&2
    return 1
fi

# ZLE widget for linesense suggestions
linesense-widget() {
    local current_buffer="$BUFFER"
    local cwd="$PWD"

    # Don't suggest for empty buffers
    if [[ -z "$current_buffer" ]]; then
        return
    fi

    # Call linesense suggest and capture JSON output
    local result
    result=$(linesense suggest --shell zsh --line "$current_buffer" --cwd "$cwd" 2>/dev/null)

    if [[ $? -eq 0 && -n "$result" ]]; then
        # Parse JSON and extract first suggestion
        local suggestion risk

        if command -v jq &> /dev/null; then
            suggestion=$(echo "$result" | jq -r '.suggestions[0].command' 2>/dev/null)
            risk=$(echo "$result" | jq -r '.suggestions[0].risk' 2>/dev/null)
        else
            suggestion=$(echo "$result" | grep -o '"command":"[^"]*"' | head -1 | sed 's/"command":"\(.*\)"/\1/')
            risk=$(echo "$result" | grep -o '"risk":"[^"]*"' | head -1 | sed 's/"risk":"\(.*\)"/\1/')
        fi

        if [[ -n "$suggestion" && "$suggestion" != "null" ]]; then
            # Show risk indicator for high-risk commands
            if [[ "$risk" == "high" ]]; then
                print "\n⚠️  WARNING: High-risk command detected!" >&2
            fi

            # Replace buffer with suggestion
            BUFFER="$suggestion"
            CURSOR=${#BUFFER}
        fi
    fi

    # Redraw the line
    zle reset-prompt
}

# ZLE widget for command explanation
linesense-explain-widget() {
    local current_buffer="$BUFFER"
    local cwd="$PWD"

    # Don't explain empty buffers
    if [[ -z "$current_buffer" ]]; then
        return
    fi

    # Call linesense explain and capture JSON output
    local result
    result=$(linesense explain --shell zsh --line "$current_buffer" --cwd "$cwd" 2>/dev/null)

    if [[ $? -eq 0 && -n "$result" ]]; then
        local summary risk

        if command -v jq &> /dev/null; then
            summary=$(echo "$result" | jq -r '.summary' 2>/dev/null)
            risk=$(echo "$result" | jq -r '.risk' 2>/dev/null)
        else
            summary=$(echo "$result" | grep -o '"summary":"[^"]*"' | head -1 | sed 's/"summary":"\(.*\)"/\1/')
            risk=$(echo "$result" | grep -o '"risk":"[^"]*"' | head -1 | sed 's/"risk":"\(.*\)"/\1/')
        fi

        # Display formatted explanation
        print "" >&2
        print "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" >&2
        print "📝 LineSense Explanation" >&2
        print "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" >&2
        print "" >&2
        print "Command: $current_buffer" >&2
        print "" >&2

        # Color-code risk level
        case "$risk" in
            high)
                print -P "⚠️  Risk: %F{red}%BHIGH%b%f" >&2
                ;;
            medium)
                print -P "⚠️  Risk: %F{yellow}%BMEDIUM%b%f" >&2
                ;;
            low)
                print -P "✓ Risk: %F{green}%BLOW%b%f" >&2
                ;;
        esac

        print "" >&2
        print "$summary" >&2
        print "" >&2
        print "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" >&2
        print "" >&2
    else
        print "❌ Failed to get explanation" >&2
    fi

    # Redraw the line
    zle reset-prompt
}

# Register ZLE widgets
zle -N linesense-widget
zle -N linesense-explain-widget

# Default keybindings
# Override these by setting environment variables before sourcing this file:
#   export LINESENSE_SUGGEST_KEY="^S"  # Ctrl+S for suggest
#   export LINESENSE_EXPLAIN_KEY="^H"  # Ctrl+H for explain

# Suggest keybinding (default: Ctrl+Space)
LINESENSE_SUGGEST_KEY="${LINESENSE_SUGGEST_KEY:-"^ "}"
bindkey "${LINESENSE_SUGGEST_KEY}" linesense-widget

# Explain keybinding (default: Ctrl+X Ctrl+E to avoid conflicts)
# Note: ^E conflicts with end-of-line in many zsh configs, so we use ^X^E
LINESENSE_EXPLAIN_KEY="${LINESENSE_EXPLAIN_KEY:-"^X^E"}"
bindkey "${LINESENSE_EXPLAIN_KEY}" linesense-explain-widget

# Print keybinding information
print "LineSense zsh integration loaded:" >&2
print "  Suggest: ${LINESENSE_SUGGEST_KEY} (default: Ctrl+Space)" >&2
print "  Explain: ${LINESENSE_EXPLAIN_KEY} (default: Ctrl+X Ctrl+E)" >&2

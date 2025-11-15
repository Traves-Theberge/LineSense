package ai

import (
	"fmt"
	"strings"

	"github.com/traves/linesense/internal/core"
)

// buildSuggestSystemPrompt creates the system prompt for command suggestions
func buildSuggestSystemPrompt() string {
	return `You are an expert shell command assistant. Your job is to suggest 3-5 complete, correct shell commands based on the user's partial input and context.

IMPORTANT RULES:
1. Provide 3-5 alternative command suggestions
2. Order suggestions from most likely to least likely
3. Make commands safe and appropriate
4. Use the context (git info, history, cwd) to make intelligent suggestions
5. If the input is already complete, suggest improvements or alternatives
6. For ambiguous or typo inputs, interpret user intent and suggest corrections
7. Prefer standard Unix/Linux commands
8. Keep commands concise but complete

RESPONSE FORMAT:
One suggestion per line in this exact format:
COMMAND | brief explanation (5-10 words max)

Example:
ls -la | List all files with details
find . -type f -name "*.txt" | Find all text files recursively
tree -L 2 | Show directory tree 2 levels deep`
}

// buildSuggestUserPrompt creates the user prompt with context
func buildSuggestUserPrompt(ctx *core.ContextEnvelope) string {
	var parts []string

	// Add the current line (what the user is typing)
	parts = append(parts, fmt.Sprintf("Current input: %s", ctx.Line))

	// Add shell and working directory
	parts = append(parts, fmt.Sprintf("\nShell: %s", ctx.Shell))
	parts = append(parts, fmt.Sprintf("Working directory: %s", ctx.CWD))

	// Add git context if available
	if ctx.Git != nil && ctx.Git.IsRepo {
		parts = append(parts, "\nGit context:")
		parts = append(parts, fmt.Sprintf("- Branch: %s", ctx.Git.Branch))
		parts = append(parts, fmt.Sprintf("- Status: %s", ctx.Git.StatusSummary))
		if len(ctx.Git.Remotes) > 0 {
			parts = append(parts, fmt.Sprintf("- Remotes: %s", strings.Join(ctx.Git.Remotes, ", ")))
		}
	}

	// Add recent history if available
	if len(ctx.History) > 0 {
		parts = append(parts, "\nRecent commands (last 5):")
		start := len(ctx.History) - 5
		if start < 0 {
			start = 0
		}
		for _, entry := range ctx.History[start:] {
			parts = append(parts, fmt.Sprintf("- %s", entry.Command))
		}
	}

	parts = append(parts, "\nSuggest the complete command:")

	return strings.Join(parts, "\n")
}

// buildExplainSystemPrompt creates the system prompt for command explanations
func buildExplainSystemPrompt() string {
	return `You are an expert shell command explainer. Your job is to explain what a command does, its risks, and potential side effects.

IMPORTANT RULES:
1. Be concise but thorough
2. Explain what the command does in plain English
3. Identify potential risks (low/medium/high)
4. Warn about destructive operations
5. Mention important flags and options
6. Note common pitfalls or mistakes

RESPONSE FORMAT:
Summary: [one-sentence explanation]
Risk: [low|medium|high]
Details: [detailed explanation]`
}

// buildExplainUserPrompt creates the user prompt for explanations
func buildExplainUserPrompt(ctx *core.ContextEnvelope) string {
	var parts []string

	// Add the command to explain
	parts = append(parts, fmt.Sprintf("Explain this command: %s", ctx.Line))

	// Add context
	parts = append(parts, fmt.Sprintf("\nShell: %s", ctx.Shell))
	parts = append(parts, fmt.Sprintf("Working directory: %s", ctx.CWD))

	// Add git context if relevant
	if ctx.Git != nil && ctx.Git.IsRepo {
		parts = append(parts, fmt.Sprintf("\nGit repository: branch=%s, status=%s", ctx.Git.Branch, ctx.Git.StatusSummary))
	}

	return strings.Join(parts, "\n")
}

// parseSuggestions extracts command suggestions from AI response
func parseSuggestions(response string, originalLine string) []core.Suggestion {
	// Clean up the response
	cleaned := strings.TrimSpace(response)

	// Remove markdown code blocks if present
	cleaned = strings.TrimPrefix(cleaned, "```bash")
	cleaned = strings.TrimPrefix(cleaned, "```sh")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)

	// Split by lines to get multiple suggestions
	lines := strings.Split(cleaned, "\n")

	var suggestions []core.Suggestion
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines
		if line == "" {
			continue
		}

		// Skip numbered lines (in case AI added numbering)
		line = strings.TrimPrefix(line, "1. ")
		line = strings.TrimPrefix(line, "2. ")
		line = strings.TrimPrefix(line, "3. ")
		line = strings.TrimPrefix(line, "4. ")
		line = strings.TrimPrefix(line, "5. ")
		line = strings.TrimSpace(line)

		// Skip if still empty after cleanup
		if line == "" {
			continue
		}

		// Split on pipe to separate command from explanation
		parts := strings.SplitN(line, "|", 2)
		command := strings.TrimSpace(parts[0])
		explanation := ""
		if len(parts) > 1 {
			explanation = strings.TrimSpace(parts[1])
		}

		// Skip if command is empty
		if command == "" {
			continue
		}

		// Create suggestion with risk assessment
		risk := assessRisk(command)

		// Use the extracted explanation if available, otherwise use default
		if explanation == "" {
			explanation = fmt.Sprintf("Suggested based on: %s", originalLine)
		}

		suggestions = append(suggestions, core.Suggestion{
			Command:     command,
			Risk:        risk,
			Explanation: explanation,
			Source:      "llm",
		})

		// Limit to 5 suggestions max
		if len(suggestions) >= 5 {
			break
		}
	}

	return suggestions
}

// parseExplanation extracts explanation from AI response
func parseExplanation(response string) core.Explanation {
	lines := strings.Split(response, "\n")

	explanation := core.Explanation{
		Risk:  core.RiskMedium, // Default
		Notes: []string{},
	}

	// Parse structured response
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "Summary:") {
			explanation.Summary = strings.TrimSpace(strings.TrimPrefix(line, "Summary:"))
		} else if strings.HasPrefix(line, "Risk:") {
			riskStr := strings.TrimSpace(strings.TrimPrefix(line, "Risk:"))
			switch strings.ToLower(riskStr) {
			case "low":
				explanation.Risk = core.RiskLow
			case "medium":
				explanation.Risk = core.RiskMedium
			case "high":
				explanation.Risk = core.RiskHigh
			}
		} else if strings.HasPrefix(line, "Details:") {
			explanation.Notes = append(explanation.Notes, strings.TrimSpace(strings.TrimPrefix(line, "Details:")))
		} else if !strings.HasPrefix(line, "Summary:") && !strings.HasPrefix(line, "Risk:") {
			// Add other non-empty lines as notes
			explanation.Notes = append(explanation.Notes, line)
		}
	}

	// If no summary was parsed, use the full response
	if explanation.Summary == "" {
		explanation.Summary = strings.TrimSpace(response)
	}

	return explanation
}

// assessRisk performs basic risk assessment on a command
func assessRisk(command string) core.RiskLevel {
	commandLower := strings.ToLower(command)

	// High risk patterns
	highRiskPatterns := []string{
		"rm -rf", "dd if=", "mkfs", "> /dev/", "rm /",
		"chmod 777", "chown -R", "sudo rm",
	}
	for _, pattern := range highRiskPatterns {
		if strings.Contains(commandLower, pattern) {
			return core.RiskHigh
		}
	}

	// Medium risk patterns
	mediumRiskPatterns := []string{
		"sudo", "rm ", "mv ", "chmod", "chown",
		"kill", "pkill", "systemctl stop",
	}
	for _, pattern := range mediumRiskPatterns {
		if strings.Contains(commandLower, pattern) {
			return core.RiskMedium
		}
	}

	// Default to low risk
	return core.RiskLow
}

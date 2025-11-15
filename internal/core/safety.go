package core

import (
	"regexp"
	"strings"

	"github.com/traves/linesense/internal/config"
)

// RiskLevel represents the risk level of a command
type RiskLevel string

const (
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"
)

// Built-in high-risk patterns (always checked)
var builtinHighRiskPatterns = []string{
	`rm\s+-rf\s+/`,     // rm -rf /
	`dd\s+if=`,         // dd commands
	`mkfs`,             // filesystem formatting
	`>\s*/dev/`,        // writing to device files
	`chmod\s+777`,      // overly permissive permissions
	`chmod\s+-R\s+777`, // recursive 777
	`curl.*\|\s*bash`,  // curl to bash
	`wget.*\|\s*sh`,    // wget to shell
	`:\(\)\{.*\};:`,    // fork bomb (escaped parens and braces)
	`killall\s+-9`,     // force kill all
}

// Built-in medium-risk patterns
var builtinMediumRiskPatterns = []string{
	`sudo`,             // elevated privileges
	`rm\s+`,            // file removal
	`mv\s+`,            // file move
	`chmod`,            // permission changes
	`chown`,            // ownership changes
	`kill`,             // process termination
	`pkill`,            // process killing
	`systemctl`,        // system service management
	`reboot`,           // system reboot
	`shutdown`,         // system shutdown
	`iptables`,         // firewall changes
	`apt-get\s+remove`, // package removal
	`yum\s+remove`,     // package removal
}

// ApplySafetyFilters filters and classifies suggestions based on safety rules
func ApplySafetyFilters(suggestions []Suggestion, cfg *config.SafetyConfig) []Suggestion {
	var filtered []Suggestion

	for _, suggestion := range suggestions {
		// Check if command is blocked
		if IsBlocked(suggestion.Command, cfg) {
			continue // Skip blocked commands
		}

		// Classify risk for remaining commands
		risk := ClassifyRisk(suggestion.Command, cfg)
		suggestion.Risk = risk

		filtered = append(filtered, suggestion)
	}

	return filtered
}

// ClassifyRisk determines the risk level of a command
func ClassifyRisk(command string, cfg *config.SafetyConfig) RiskLevel {
	commandLower := strings.ToLower(command)

	// Check high-risk patterns first
	for _, pattern := range builtinHighRiskPatterns {
		matched, err := regexp.MatchString(pattern, commandLower)
		if err == nil && matched {
			return RiskHigh
		}
	}

	// Check config-defined high-risk patterns
	if cfg != nil {
		for _, pattern := range cfg.RequireConfirmPatterns {
			matched, err := regexp.MatchString(pattern, commandLower)
			if err == nil && matched {
				return RiskHigh
			}
		}
	}

	// Check medium-risk patterns
	for _, pattern := range builtinMediumRiskPatterns {
		matched, err := regexp.MatchString(pattern, commandLower)
		if err == nil && matched {
			return RiskMedium
		}
	}

	// Default to low risk
	return RiskLow
}

// IsBlocked returns true if the command should be blocked entirely
func IsBlocked(command string, cfg *config.SafetyConfig) bool {
	if cfg == nil {
		return false
	}

	commandLower := strings.ToLower(command)

	// Check against denylist patterns
	for _, pattern := range cfg.Denylist {
		matched, err := regexp.MatchString(pattern, commandLower)
		if err == nil && matched {
			return true
		}
	}

	return false
}

// ValidateCommand performs additional validation on commands
func ValidateCommand(command string) error {
	// Check for empty commands
	if strings.TrimSpace(command) == "" {
		return nil // Empty is okay, just means no suggestion
	}

	// Check for suspiciously long commands (potential injection)
	if len(command) > 10000 {
		return nil // Just ignore, don't error
	}

	return nil
}

package core

import (
	"os"
	"strings"

	"github.com/traves/linesense/internal/config"
)

// ContextEnvelope is collected before each suggestion / explanation
type ContextEnvelope struct {
	Shell        string            `json:"shell"`                  // "bash" | "zsh"
	Line         string            `json:"line"`                   // current input line
	CWD          string            `json:"cwd"`
	Git          *GitInfo          `json:"git,omitempty"`
	Env          map[string]string `json:"env,omitempty"`          // filtered env (if enabled)
	History      []HistoryEntry    `json:"history,omitempty"`      // last N commands
	UsageSummary *UsageSummary     `json:"usage_summary,omitempty"`
}

// GitInfo contains git repository information
type GitInfo struct {
	IsRepo        bool     `json:"is_repo"`
	Branch        string   `json:"branch,omitempty"`
	StatusSummary string   `json:"status_summary,omitempty"`
	Remotes       []string `json:"remotes,omitempty"`
}

// HistoryEntry represents a shell history entry
type HistoryEntry struct {
	Command   string  `json:"command"`
	Timestamp *string `json:"timestamp,omitempty"`
	ExitCode  *int    `json:"exit_code,omitempty"`
}

// UsageSummary contains usage pattern information
type UsageSummary struct {
	FrequentlyUsedCommands []string `json:"frequently_used_commands"` // top N commands in this cwd
}

// BuildContext gathers all contextual information
func BuildContext(shell, line, cwd string, cfg *config.Config) (*ContextEnvelope, error) {
	ctx := &ContextEnvelope{
		Shell: shell,
		Line:  line,
		CWD:   cwd,
	}

	// Collect git context if enabled
	if cfg.Context.IncludeGit {
		gitInfo, err := CollectGitInfo(cwd)
		if err == nil && gitInfo != nil {
			ctx.Git = gitInfo
		}
		// Silently ignore errors - git info is optional
	}

	// Collect shell history if enabled
	if cfg.Context.HistoryLength > 0 {
		history, err := CollectHistory(shell, cfg.Context.HistoryLength)
		if err == nil && len(history) > 0 {
			ctx.History = history
		}
		// Silently ignore errors - history is optional
	}

	// Collect environment variables if enabled
	if cfg.Context.IncludeEnv {
		ctx.Env = collectFilteredEnv()
	}

	// TODO: Build usage summary from usage log
	// This will be implemented when usage logging is complete

	return ctx, nil
}

// collectFilteredEnv returns a filtered map of environment variables
// Filters out sensitive variables like API keys and passwords
func collectFilteredEnv() map[string]string {
	// List of environment variable patterns to exclude (case-insensitive)
	sensitivePatterns := []string{
		"KEY", "SECRET", "PASSWORD", "PASS", "TOKEN", "AUTH",
		"CREDENTIAL", "PRIVATE", "API_KEY",
	}

	env := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) != 2 {
			continue
		}

		key := pair[0]
		value := pair[1]

		// Check if key contains sensitive patterns
		keyUpper := strings.ToUpper(key)
		sensitive := false
		for _, pattern := range sensitivePatterns {
			if strings.Contains(keyUpper, pattern) {
				sensitive = true
				break
			}
		}

		if !sensitive {
			env[key] = value
		}
	}

	return env
}

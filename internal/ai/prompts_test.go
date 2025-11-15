package ai

import (
	"strings"
	"testing"

	"github.com/traves/linesense/internal/core"
)

func TestBuildSuggestSystemPrompt(t *testing.T) {
	prompt := buildSuggestSystemPrompt()

	if prompt == "" {
		t.Error("System prompt should not be empty")
	}

	// Verify it contains key instructions
	requiredPhrases := []string{
		"shell command",
		"3-5",
		"alternative",
		"safe",
		"context",
	}

	for _, phrase := range requiredPhrases {
		if !strings.Contains(prompt, phrase) {
			t.Errorf("System prompt should contain %q", phrase)
		}
	}
}

func TestBuildSuggestUserPrompt(t *testing.T) {
	ctx := &core.ContextEnvelope{
		Shell:          "bash",
		Line:           "git com",
		CWD:            "/home/user/project",
		OS:             "linux",
		Distribution:   "ubuntu",
		PackageManager: "apt",
		Git: &core.GitInfo{
			IsRepo:        true,
			Branch:        "main",
			StatusSummary: "clean",
			Remotes:       []string{"origin"},
		},
		History: []core.HistoryEntry{
			{Command: "ls -la"},
			{Command: "cd project"},
			{Command: "git status"},
		},
	}

	prompt := buildSuggestUserPrompt(ctx)

	// Verify all context is included
	if !strings.Contains(prompt, "git com") {
		t.Error("Should contain current input")
	}
	if !strings.Contains(prompt, "linux") {
		t.Error("Should contain OS")
	}
	if !strings.Contains(prompt, "ubuntu") {
		t.Error("Should contain distribution")
	}
	if !strings.Contains(prompt, "apt") {
		t.Error("Should contain package manager")
	}
	if !strings.Contains(prompt, "bash") {
		t.Error("Should contain shell")
	}
	if !strings.Contains(prompt, "/home/user/project") {
		t.Error("Should contain working directory")
	}
	if !strings.Contains(prompt, "main") {
		t.Error("Should contain git branch")
	}
	if !strings.Contains(prompt, "clean") {
		t.Error("Should contain git status")
	}
	if !strings.Contains(prompt, "git status") {
		t.Error("Should contain recent history")
	}
}

func TestBuildSuggestUserPrompt_NoGit(t *testing.T) {
	ctx := &core.ContextEnvelope{
		Shell:          "zsh",
		Line:           "ls",
		CWD:            "/tmp",
		OS:             "darwin",
		PackageManager: "brew",
	}

	prompt := buildSuggestUserPrompt(ctx)

	// Should still work without git context
	if !strings.Contains(prompt, "ls") {
		t.Error("Should contain current input")
	}
	if !strings.Contains(prompt, "zsh") {
		t.Error("Should contain shell")
	}
	if strings.Contains(prompt, "Git context") {
		t.Error("Should not contain git context when not in repo")
	}
}

func TestBuildExplainSystemPrompt(t *testing.T) {
	prompt := buildExplainSystemPrompt()

	if prompt == "" {
		t.Error("Explain system prompt should not be empty")
	}

	// Verify key instructions
	requiredPhrases := []string{
		"explain",
		"risks",
		"Summary",
		"Risk",
		"Details",
	}

	for _, phrase := range requiredPhrases {
		if !strings.Contains(prompt, phrase) {
			t.Errorf("Explain prompt should contain %q", phrase)
		}
	}
}

func TestBuildExplainUserPrompt(t *testing.T) {
	ctx := &core.ContextEnvelope{
		Shell:        "bash",
		Line:         "rm -rf /tmp/test",
		CWD:          "/home/user",
		OS:           "linux",
		Distribution: "arch",
		Git: &core.GitInfo{
			IsRepo:        true,
			Branch:        "develop",
			StatusSummary: "modified",
		},
	}

	prompt := buildExplainUserPrompt(ctx)

	// Verify command and context
	if !strings.Contains(prompt, "rm -rf /tmp/test") {
		t.Error("Should contain command to explain")
	}
	if !strings.Contains(prompt, "linux") {
		t.Error("Should contain OS")
	}
	if !strings.Contains(prompt, "arch") {
		t.Error("Should contain distribution")
	}
	if !strings.Contains(prompt, "bash") {
		t.Error("Should contain shell")
	}
	if !strings.Contains(prompt, "develop") {
		t.Error("Should contain git branch")
	}
}

func TestParseSuggestions(t *testing.T) {
	tests := []struct {
		name         string
		response     string
		originalLine string
		wantCommand  string
		wantRisk     core.RiskLevel
	}{
		{
			name:         "simple command",
			response:     "git commit -m 'Update README'",
			originalLine: "git com",
			wantCommand:  "git commit -m 'Update README'",
			wantRisk:     core.RiskLow,
		},
		{
			name:         "command with markdown",
			response:     "```bash\nls -la\n```",
			originalLine: "ls",
			wantCommand:  "ls -la",
			wantRisk:     core.RiskLow,
		},
		{
			name:         "high risk command",
			response:     "sudo rm -rf /var/log/old",
			originalLine: "del logs",
			wantCommand:  "sudo rm -rf /var/log/old",
			wantRisk:     core.RiskHigh,
		},
		{
			name:         "medium risk command",
			response:     "sudo systemctl restart nginx",
			originalLine: "restart nginx",
			wantCommand:  "sudo systemctl restart nginx",
			wantRisk:     core.RiskMedium,
		},
		{
			name:         "multiline response (now returns multiple suggestions)",
			response:     "git status\ngit st",
			originalLine: "git st",
			wantCommand:  "git status", // First suggestion
			wantRisk:     core.RiskLow,
		},
		{
			name:         "command with explanation",
			response:     "ls -la | List all files with details",
			originalLine: "ls",
			wantCommand:  "ls -la",
			wantRisk:     core.RiskLow,
		},
		{
			name:         "multiple commands with explanations",
			response:     "git status | Show repository status\ngit diff | Show uncommitted changes\ngit log | Show commit history",
			originalLine: "git",
			wantCommand:  "git status", // First suggestion
			wantRisk:     core.RiskLow,
		},
		{
			name:         "empty response",
			response:     "",
			originalLine: "???",
			wantCommand:  "",
			wantRisk:     core.RiskLow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions := parseSuggestions(tt.response, tt.originalLine)

			if tt.wantCommand == "" {
				if len(suggestions) != 0 {
					t.Errorf("Expected no suggestions for empty response, got %d", len(suggestions))
				}
				return
			}

			if len(suggestions) < 1 {
				t.Fatalf("Expected at least 1 suggestion, got %d", len(suggestions))
			}

			// Check the first suggestion matches expected
			suggestion := suggestions[0]
			if suggestion.Command != tt.wantCommand {
				t.Errorf("First command = %q, want %q", suggestion.Command, tt.wantCommand)
			}
			if suggestion.Risk != tt.wantRisk {
				t.Errorf("First risk = %v, want %v", suggestion.Risk, tt.wantRisk)
			}
			if suggestion.Source != "llm" {
				t.Errorf("Source = %q, want llm", suggestion.Source)
			}
		})
	}
}

func TestParseSuggestions_Explanations(t *testing.T) {
	tests := []struct {
		name            string
		response        string
		wantExplanation string
	}{
		{
			name:            "command with explanation",
			response:        "ls -la | List all files with details",
			wantExplanation: "List all files with details",
		},
		{
			name:            "command without explanation",
			response:        "git status",
			wantExplanation: "Suggested based on: git", // Default explanation
		},
		{
			name:            "multiple commands with explanations",
			response:        "git status | Show repository status\ngit diff | Show uncommitted changes",
			wantExplanation: "Show repository status", // First suggestion
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions := parseSuggestions(tt.response, "git")

			if len(suggestions) < 1 {
				t.Fatalf("Expected at least 1 suggestion, got %d", len(suggestions))
			}

			if suggestions[0].Explanation != tt.wantExplanation {
				t.Errorf("Explanation = %q, want %q", suggestions[0].Explanation, tt.wantExplanation)
			}
		})
	}
}

func TestParseExplanation(t *testing.T) {
	tests := []struct {
		name        string
		response    string
		wantSummary string
		wantRisk    core.RiskLevel
		checkNotes  bool
		wantNotes   int
	}{
		{
			name: "structured response",
			response: `Summary: Lists all files in the current directory
Risk: low
Details: The ls command is safe and only reads directory information`,
			wantSummary: "Lists all files in the current directory",
			wantRisk:    core.RiskLow,
			checkNotes:  true,
			wantNotes:   1,
		},
		{
			name: "high risk response",
			response: `Summary: Deletes all files recursively
Risk: high
Details: This is extremely dangerous and can cause data loss`,
			wantSummary: "Deletes all files recursively",
			wantRisk:    core.RiskHigh,
			checkNotes:  true,
			wantNotes:   1,
		},
		{
			name: "medium risk response",
			response: `Summary: Restarts the web server
Risk: medium
Details: Will cause temporary service interruption`,
			wantSummary: "Restarts the web server",
			wantRisk:    core.RiskMedium,
		},
		{
			name:        "unstructured response",
			response:    "This command lists files in the directory",
			wantSummary: "This command lists files in the directory",
			wantRisk:    core.RiskMedium, // Default
		},
		{
			name: "case insensitive risk",
			response: `Summary: Test command
Risk: HIGH
Details: Test details`,
			wantSummary: "Test command",
			wantRisk:    core.RiskHigh,
		},
		{
			name: "response with additional lines",
			response: `Summary: Git status check
Risk: low
This command shows the status of your repository
It will not make any changes
Only reads repository information`,
			wantSummary: "Git status check",
			wantRisk:    core.RiskLow,
			checkNotes:  true,
			wantNotes:   3,
		},
		{
			name: "empty response",
			response: `


`,
			wantSummary: "",
			wantRisk:    core.RiskMedium,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			explanation := parseExplanation(tt.response)

			if explanation.Summary != tt.wantSummary {
				t.Errorf("Summary = %q, want %q", explanation.Summary, tt.wantSummary)
			}
			if explanation.Risk != tt.wantRisk {
				t.Errorf("Risk = %v, want %v", explanation.Risk, tt.wantRisk)
			}
			if tt.checkNotes && len(explanation.Notes) != tt.wantNotes {
				t.Errorf("Notes count = %d, want %d", len(explanation.Notes), tt.wantNotes)
			}
		})
	}
}

func TestAssessRisk(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		wantRisk core.RiskLevel
	}{
		// High risk commands
		{"rm -rf root", "rm -rf /", core.RiskHigh},
		{"dd dangerous", "dd if=/dev/zero of=/dev/sda", core.RiskHigh},
		{"mkfs format", "mkfs.ext4 /dev/sda", core.RiskHigh},
		{"chmod 777", "chmod 777 /etc/passwd", core.RiskHigh},
		{"sudo rm recursive", "sudo rm -rf /var", core.RiskHigh},

		// Medium risk commands
		{"sudo command", "sudo apt-get update", core.RiskMedium},
		{"rm single file", "rm file.txt", core.RiskMedium},
		{"mv file", "mv old.txt new.txt", core.RiskMedium},
		{"chmod", "chmod 644 file.txt", core.RiskMedium},
		{"kill process", "kill -9 1234", core.RiskMedium},
		{"systemctl stop", "systemctl stop nginx", core.RiskMedium},

		// Low risk commands
		{"ls", "ls -la", core.RiskLow},
		{"cat", "cat file.txt", core.RiskLow},
		{"echo", "echo hello", core.RiskLow},
		{"pwd", "pwd", core.RiskLow},
		{"git status", "git status", core.RiskLow},
		{"grep", "grep pattern file.txt", core.RiskLow},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			risk := assessRisk(tt.command)
			if risk != tt.wantRisk {
				t.Errorf("assessRisk(%q) = %v, want %v", tt.command, risk, tt.wantRisk)
			}
		})
	}
}

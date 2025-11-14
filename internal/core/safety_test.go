package core

import (
	"testing"

	"github.com/traves/linesense/internal/config"
)

func TestClassifyRisk(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected RiskLevel
	}{
		// High-risk commands
		{"rm -rf root", "rm -rf /", RiskHigh},
		{"dd command", "dd if=/dev/zero of=/dev/sda", RiskHigh},
		{"mkfs", "mkfs.ext4 /dev/sda1", RiskHigh},
		{"chmod 777", "chmod 777 file.sh", RiskHigh},
		{"curl to bash", "curl http://example.com | bash", RiskHigh},
		{"fork bomb", ":(){ :|:& };:", RiskHigh},

		// Medium-risk commands
		{"sudo command", "sudo apt-get update", RiskMedium},
		{"rm file", "rm myfile.txt", RiskMedium},
		{"mv file", "mv old.txt new.txt", RiskMedium},
		{"chmod", "chmod 644 file.txt", RiskMedium},
		{"kill", "kill -9 1234", RiskMedium},
		{"systemctl", "systemctl restart nginx", RiskMedium},

		// Low-risk commands
		{"ls", "ls -la", RiskLow},
		{"cat", "cat file.txt", RiskLow},
		{"echo", "echo hello", RiskLow},
		{"pwd", "pwd", RiskLow},
		{"grep", "grep pattern file.txt", RiskLow},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ClassifyRisk(tt.command, nil)
			if result != tt.expected {
				t.Errorf("ClassifyRisk(%q) = %v, want %v", tt.command, result, tt.expected)
			}
		})
	}
}

func TestIsBlocked(t *testing.T) {
	cfg := &config.SafetyConfig{
		Denylist: []string{
			`rm\s+-rf\s+/`,
			`dd\s+if=`,
		},
	}

	tests := []struct {
		name     string
		command  string
		expected bool
	}{
		{"blocked rm -rf /", "rm -rf /", true},
		{"blocked dd", "dd if=/dev/zero of=/dev/sda", true},
		{"allowed ls", "ls -la", false},
		{"allowed rm file", "rm myfile.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsBlocked(tt.command, cfg)
			if result != tt.expected {
				t.Errorf("IsBlocked(%q) = %v, want %v", tt.command, result, tt.expected)
			}
		})
	}
}

func TestApplySafetyFilters(t *testing.T) {
	cfg := &config.SafetyConfig{
		Denylist: []string{
			`rm\s+-rf\s+/`,
		},
		RequireConfirmPatterns: []string{
			`format`,
		},
	}

	suggestions := []Suggestion{
		{Command: "ls -la", Risk: RiskLow},
		{Command: "rm -rf /", Risk: RiskLow},     // Should be blocked
		{Command: "sudo rm file.txt", Risk: RiskLow}, // Should be medium risk
		{Command: "format disk", Risk: RiskLow},  // Should be high risk
	}

	filtered := ApplySafetyFilters(suggestions, cfg)

	// Should have 3 suggestions (rm -rf / is blocked)
	if len(filtered) != 3 {
		t.Errorf("Expected 3 suggestions, got %d", len(filtered))
	}

	// Check ls is low risk
	if filtered[0].Risk != RiskLow {
		t.Errorf("Expected ls to be low risk, got %v", filtered[0].Risk)
	}

	// Check sudo rm is medium risk
	if filtered[1].Risk != RiskMedium {
		t.Errorf("Expected sudo rm to be medium risk, got %v", filtered[1].Risk)
	}

	// Check format is high risk
	if filtered[2].Risk != RiskHigh {
		t.Errorf("Expected format to be high risk, got %v", filtered[2].Risk)
	}
}

func TestIsBlocked_NilConfig(t *testing.T) {
	// Should not block anything with nil config
	result := IsBlocked("rm -rf /", nil)
	if result {
		t.Error("IsBlocked should return false with nil config")
	}
}

func TestValidateCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		wantErr  bool
	}{
		{"valid command", "ls -la", false},
		{"empty command", "", false},
		{"very long command", string(make([]byte, 20000)), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCommand(tt.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfigDir(t *testing.T) {
	// Save original env
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	tests := []struct {
		name        string
		xdgHome     string
		expectMatch string
	}{
		{
			name:        "with XDG_CONFIG_HOME set",
			xdgHome:     "/tmp/test-config",
			expectMatch: "/tmp/test-config/linesense",
		},
		{
			name:        "without XDG_CONFIG_HOME",
			xdgHome:     "",
			expectMatch: "/.config/linesense", // Will have home prefix
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.xdgHome != "" {
				os.Setenv("XDG_CONFIG_HOME", tt.xdgHome)
			} else {
				os.Unsetenv("XDG_CONFIG_HOME")
			}

			result := GetConfigDir()

			if tt.xdgHome != "" {
				if result != tt.expectMatch {
					t.Errorf("GetConfigDir() = %v, want %v", result, tt.expectMatch)
				}
			} else {
				// Should contain .config/linesense
				if !filepath.IsAbs(result) {
					t.Errorf("GetConfigDir() should return absolute path, got %v", result)
				}
				if !filepath.HasPrefix(result, string(filepath.Separator)) {
					t.Errorf("GetConfigDir() should start with /, got %v", result)
				}
			}
		})
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	// Save original env
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	// Set to non-existent directory
	os.Setenv("XDG_CONFIG_HOME", "/tmp/linesense-test-nonexistent")

	_, err := LoadConfig()
	if err == nil {
		t.Error("LoadConfig() should return error for missing file")
	}
}

func TestLoadConfig_ValidFile(t *testing.T) {
	// Create temporary config
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "linesense")
	os.MkdirAll(configDir, 0755)

	configContent := `
[ai]
provider_profile = "test-profile"

[context]
history_length = 25
include_git = false
include_env = false

[safety]
denylist = ["test-pattern"]
`

	configPath := filepath.Join(configDir, "config.toml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Set XDG to temp dir
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)
	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Verify values
	if cfg.AI.ProviderProfile != "test-profile" {
		t.Errorf("ProviderProfile = %v, want test-profile", cfg.AI.ProviderProfile)
	}
	if cfg.Context.HistoryLength != 25 {
		t.Errorf("HistoryLength = %v, want 25", cfg.Context.HistoryLength)
	}
	if cfg.Context.IncludeGit != false {
		t.Errorf("IncludeGit = %v, want false", cfg.Context.IncludeGit)
	}
	if len(cfg.Safety.Denylist) == 0 {
		t.Error("Denylist should be populated")
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	// Create minimal config
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "linesense")
	os.MkdirAll(configDir, 0755)

	configContent := `
# Minimal config - should use defaults
`

	configPath := filepath.Join(configDir, "config.toml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)
	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Check defaults
	if cfg.AI.ProviderProfile == "" {
		t.Error("ProviderProfile should have default value")
	}
	if cfg.Context.HistoryLength == 0 {
		t.Error("HistoryLength should have default value")
	}
}

func TestGetConfigDir_WithHome(t *testing.T) {
	// Test when HOME is set but XDG_CONFIG_HOME is not
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	originalHOME := os.Getenv("HOME")
	defer func() {
		os.Setenv("XDG_CONFIG_HOME", originalXDG)
		os.Setenv("HOME", originalHOME)
	}()

	os.Unsetenv("XDG_CONFIG_HOME")
	os.Setenv("HOME", "/tmp/test-home")

	dir := GetConfigDir()
	expected := "/tmp/test-home/.config/linesense"
	if dir != expected {
		t.Errorf("GetConfigDir() = %v, want %v", dir, expected)
	}
}

func TestLoadConfig_InvalidTOML(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "linesense")
	os.MkdirAll(configDir, 0755)

	// Invalid TOML syntax
	configContent := `
[ai
provider_profile = "broken
`

	configPath := filepath.Join(configDir, "config.toml")
	os.WriteFile(configPath, []byte(configContent), 0644)

	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)
	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	_, err := LoadConfig()
	if err == nil {
		t.Error("LoadConfig() should error on invalid TOML")
	}
}

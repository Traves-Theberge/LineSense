package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadProvidersConfig_ValidFile(t *testing.T) {
	// Create temporary config
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "linesense")
	os.MkdirAll(configDir, 0755)

	providersContent := `
[default]
provider = "openrouter"
model = "test/model"
temperature = 0.5
max_tokens = 100

[profile.fast]
provider = "openrouter"
model = "test/fast-model"
temperature = 0.2
max_tokens = 50
`

	providersPath := filepath.Join(configDir, "providers.toml")
	if err := os.WriteFile(providersPath, []byte(providersContent), 0644); err != nil {
		t.Fatalf("Failed to write test providers config: %v", err)
	}

	// Set XDG to temp dir
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)
	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	cfg, err := LoadProvidersConfig()
	if err != nil {
		t.Fatalf("LoadProvidersConfig() error = %v", err)
	}

	// Verify default profile
	defaultProfile, err := cfg.GetProfile("default")
	if err != nil {
		t.Fatalf("GetProfile(default) error = %v", err)
	}

	if defaultProfile.Provider != "openrouter" {
		t.Errorf("Provider = %v, want openrouter", defaultProfile.Provider)
	}
	if defaultProfile.Model != "test/model" {
		t.Errorf("Model = %v, want test/model", defaultProfile.Model)
	}
	if defaultProfile.Temperature != 0.5 {
		t.Errorf("Temperature = %v, want 0.5", defaultProfile.Temperature)
	}
	if defaultProfile.MaxTokens != 100 {
		t.Errorf("MaxTokens = %v, want 100", defaultProfile.MaxTokens)
	}

	// Verify fast profile
	fastProfile, err := cfg.GetProfile("fast")
	if err != nil {
		t.Fatalf("GetProfile(fast) error = %v", err)
	}

	if fastProfile.Model != "test/fast-model" {
		t.Errorf("Fast Model = %v, want test/fast-model", fastProfile.Model)
	}
	if fastProfile.MaxTokens != 50 {
		t.Errorf("Fast MaxTokens = %v, want 50", fastProfile.MaxTokens)
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	// Create minimal config
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "linesense")
	os.MkdirAll(configDir, 0755)

	providersContent := `
[default]
provider = "openrouter"
model = "test/model"
temperature = 0.3
max_tokens = 500
`

	providersPath := filepath.Join(configDir, "providers.toml")
	if err := os.WriteFile(providersPath, []byte(providersContent), 0644); err != nil {
		t.Fatalf("Failed to write test providers config: %v", err)
	}

	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)
	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	cfg, err := LoadProvidersConfig()
	if err != nil {
		t.Fatalf("LoadProvidersConfig() error = %v", err)
	}

	// Try to get non-existent profile
	_, err = cfg.GetProfile("nonexistent")
	if err == nil {
		t.Error("GetProfile(nonexistent) should return error")
	}
}

func TestLoadProvidersConfig_MissingFile(t *testing.T) {
	// Set to non-existent directory
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)
	os.Setenv("XDG_CONFIG_HOME", "/tmp/linesense-test-providers-nonexistent")

	_, err := LoadProvidersConfig()
	if err == nil {
		t.Error("LoadProvidersConfig() should return error for missing file")
	}
}

func TestProviderProfile_Defaults(t *testing.T) {
	// Create config with minimal profile
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "linesense")
	os.MkdirAll(configDir, 0755)

	providersContent := `
[profile.minimal]
provider = "openrouter"
model = "test/model"
`

	providersPath := filepath.Join(configDir, "providers.toml")
	if err := os.WriteFile(providersPath, []byte(providersContent), 0644); err != nil {
		t.Fatalf("Failed to write test providers config: %v", err)
	}

	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)
	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	cfg, err := LoadProvidersConfig()
	if err != nil {
		t.Fatalf("LoadProvidersConfig() error = %v", err)
	}

	profile, err := cfg.GetProfile("minimal")
	if err != nil {
		t.Fatalf("GetProfile(minimal) error = %v", err)
	}

	// Verify required fields are present
	if profile.Provider != "openrouter" {
		t.Errorf("Provider = %v, want openrouter", profile.Provider)
	}
	if profile.Model != "test/model" {
		t.Errorf("Model = %v, want test/model", profile.Model)
	}
	// Note: Temperature and MaxTokens will be 0 if not specified in TOML
	// The AI provider is responsible for setting defaults
}

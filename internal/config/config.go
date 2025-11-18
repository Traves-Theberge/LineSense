package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config represents the global configuration from ~/.config/linesense/config.toml
type Config struct {
	Shell       ShellConfig       `toml:"shell"`
	Keybindings KeybindingsConfig `toml:"keybindings"`
	Context     ContextConfig     `toml:"context"`
	Safety      SafetyConfig      `toml:"safety"`
	AI          AIConfig          `toml:"ai"`
}

// ShellConfig controls which shells are enabled
type ShellConfig struct {
	EnableBash bool `toml:"enable_bash"`
	EnableZsh  bool `toml:"enable_zsh"`
}

// KeybindingsConfig defines keybindings for shell actions
type KeybindingsConfig struct {
	Suggest      string `toml:"suggest"`      // e.g. "ctrl+space"
	Explain      string `toml:"explain"`      // e.g. "ctrl+e"
	Alternatives string `toml:"alternatives"` // e.g. "alt+a"
}

// ContextConfig controls what context is gathered
type ContextConfig struct {
	HistoryLength      int    `toml:"history_length"` // how many recent commands to use
	IncludeGit         bool   `toml:"include_git"`
	IncludeFiles       bool   `toml:"include_files"`
	IncludeEnv         bool   `toml:"include_env"`
	GlobalInstructions string `toml:"global_instructions"` // User-defined global context/rules
}

// SafetyConfig defines safety rules
type SafetyConfig struct {
	RequireConfirmPatterns []string `toml:"require_confirm_patterns"`
	Denylist               []string `toml:"denylist"`
	DefaultExecution       string   `toml:"default_execution"` // "paste_only" (v0.1)
}

// AIConfig controls AI provider settings
type AIConfig struct {
	ProviderProfile string `toml:"provider_profile"` // "default" | "fast" | "smart" | etc.
}

// LoadConfig loads the global config from standard locations
func LoadConfig() (*Config, error) {
	configPath, err := getConfigPath("config.toml")
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config path: %w", err)
	}

	var cfg Config
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	// Set defaults if not specified
	if cfg.Context.HistoryLength == 0 {
		cfg.Context.HistoryLength = 100
	}
	if cfg.AI.ProviderProfile == "" {
		cfg.AI.ProviderProfile = "default"
	}

	return &cfg, nil
}

// GetConfigDir returns the configuration directory path
func GetConfigDir() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return filepath.Join(".config", "linesense")
		}
		configDir = filepath.Join(home, ".config")
	}
	return filepath.Join(configDir, "linesense")
}

// getConfigPath resolves the full path to a config file
func getConfigPath(filename string) (string, error) {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		configDir = filepath.Join(home, ".config")
	}

	return filepath.Join(configDir, "linesense", filename), nil
}

package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// ProvidersConfig represents ~/.config/linesense/providers.toml
type ProvidersConfig struct {
	Default    ProfileConfig            `toml:"default"`
	Profiles   map[string]ProfileConfig `toml:"profile"`
	OpenRouter OpenRouterConfig         `toml:"openrouter"`
}

// ProfileConfig defines a provider profile
type ProfileConfig struct {
	Provider    string  `toml:"provider"`    // "openrouter"
	Model       string  `toml:"model"`       // e.g. "openrouter/openai/gpt-4.1-mini"
	Temperature float64 `toml:"temperature"`
	MaxTokens   int     `toml:"max_tokens"`
}

// OpenRouterConfig contains OpenRouter-specific settings
type OpenRouterConfig struct {
	Type       string `toml:"type"`        // "openrouter"
	APIKeyEnv  string `toml:"api_key_env"` // e.g. "OPENROUTER_API_KEY"
	BaseURL    string `toml:"base_url"`    // e.g. "https://openrouter.ai/api/v1"
	TimeoutMs  int    `toml:"timeout_ms"`
}

// LoadProvidersConfig loads the providers config
func LoadProvidersConfig() (*ProvidersConfig, error) {
	configPath, err := getConfigPath("providers.toml")
	if err != nil {
		return nil, fmt.Errorf("failed to resolve providers config path: %w", err)
	}

	var cfg ProvidersConfig
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse providers config file %s: %w", configPath, err)
	}

	// Set defaults for OpenRouter if not specified
	if cfg.OpenRouter.Type == "" {
		cfg.OpenRouter.Type = "openrouter"
	}
	if cfg.OpenRouter.APIKeyEnv == "" {
		cfg.OpenRouter.APIKeyEnv = "OPENROUTER_API_KEY"
	}
	if cfg.OpenRouter.BaseURL == "" {
		cfg.OpenRouter.BaseURL = "https://openrouter.ai/api/v1"
	}
	if cfg.OpenRouter.TimeoutMs == 0 {
		cfg.OpenRouter.TimeoutMs = 30000
	}

	return &cfg, nil
}

// GetProfile returns the profile configuration for the given profile name
func (p *ProvidersConfig) GetProfile(profileName string) (*ProfileConfig, error) {
	if profileName == "" || profileName == "default" {
		return &p.Default, nil
	}

	profile, ok := p.Profiles[profileName]
	if !ok {
		return nil, fmt.Errorf("profile %q not found", profileName)
	}

	return &profile, nil
}

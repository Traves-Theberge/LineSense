package ai

import (
	"fmt"

	"github.com/traves/linesense/internal/config"
	"github.com/traves/linesense/internal/core"
)

// NewProvider creates a provider instance based on configuration
func NewProvider(cfg *config.ProvidersConfig, profileName string) (core.Provider, error) {
	// Get the profile configuration
	profile, err := cfg.GetProfile(profileName)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile %q: %w", profileName, err)
	}

	// Currently only OpenRouter is supported
	// Future: Add support for other providers (Anthropic, OpenAI direct, etc.)
	switch profile.Provider {
	case "openrouter", "":
		return NewOpenRouterProvider(cfg.OpenRouter, *profile)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", profile.Provider)
	}
}

package ai

import (
	"testing"

	"github.com/traves/linesense/internal/config"
)

func TestNewProvider_OpenRouter(t *testing.T) {
	cfg := &config.ProvidersConfig{
		Default: config.ProfileConfig{
			Provider:    "openrouter",
			Model:       "test/model",
			Temperature: 0.7,
			MaxTokens:   500,
		},
		OpenRouter: config.OpenRouterConfig{
			Type:       "openrouter",
			APIKeyEnv:  "TEST_API_KEY",
			BaseURL:    "https://api.test.com",
			TimeoutMs:  5000,
		},
	}

	// Set test API key
	t.Setenv("TEST_API_KEY", "test-key-123")

	provider, err := NewProvider(cfg, "default")
	if err != nil {
		t.Fatalf("NewProvider() error = %v", err)
	}

	if provider == nil {
		t.Fatal("Provider should not be nil")
	}

	if provider.Name() != "openrouter" {
		t.Errorf("Provider name = %v, want openrouter", provider.Name())
	}
}

func TestNewProvider_MissingAPIKey(t *testing.T) {
	cfg := &config.ProvidersConfig{
		Default: config.ProfileConfig{
			Provider: "openrouter",
			Model:    "test/model",
		},
		OpenRouter: config.OpenRouterConfig{
			APIKeyEnv: "NONEXISTENT_API_KEY",
			BaseURL:   "https://api.test.com",
		},
	}

	// Don't set the API key
	_, err := NewProvider(cfg, "default")
	if err == nil {
		t.Error("NewProvider() should error when API key is missing")
	}
}

func TestNewProvider_InvalidProfile(t *testing.T) {
	cfg := &config.ProvidersConfig{
		Default: config.ProfileConfig{
			Provider: "openrouter",
			Model:    "test/model",
		},
		OpenRouter: config.OpenRouterConfig{
			APIKeyEnv: "TEST_API_KEY",
			BaseURL:   "https://api.test.com",
		},
	}

	t.Setenv("TEST_API_KEY", "test-key")

	_, err := NewProvider(cfg, "nonexistent-profile")
	if err == nil {
		t.Error("NewProvider() should error for nonexistent profile")
	}
}

func TestNewProvider_UnsupportedProvider(t *testing.T) {
	cfg := &config.ProvidersConfig{
		Default: config.ProfileConfig{
			Provider: "unsupported-provider",
			Model:    "test/model",
		},
		OpenRouter: config.OpenRouterConfig{
			APIKeyEnv: "TEST_API_KEY",
			BaseURL:   "https://api.test.com",
		},
	}

	t.Setenv("TEST_API_KEY", "test-key")

	_, err := NewProvider(cfg, "default")
	if err == nil {
		t.Error("NewProvider() should error for unsupported provider")
	}
}

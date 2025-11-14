package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/traves/linesense/internal/config"
	"github.com/traves/linesense/internal/core"
)

// OpenRouterProvider implements the Provider interface for OpenRouter
type OpenRouterProvider struct {
	config  config.OpenRouterConfig
	profile config.ProfileConfig
	apiKey  string
}

// NewOpenRouterProvider creates a new OpenRouter provider
func NewOpenRouterProvider(cfg config.OpenRouterConfig, profile config.ProfileConfig) (*OpenRouterProvider, error) {
	// Get API key from environment
	apiKey := os.Getenv(cfg.APIKeyEnv)
	if apiKey == "" {
		return nil, fmt.Errorf("API key not found in environment variable %s", cfg.APIKeyEnv)
	}

	return &OpenRouterProvider{
		config:  cfg,
		profile: profile,
		apiKey:  apiKey,
	}, nil
}

// Name returns the provider name
func (p *OpenRouterProvider) Name() string {
	return "openrouter"
}

// Suggest generates command suggestions using OpenRouter
func (p *OpenRouterProvider) Suggest(ctx context.Context, input core.SuggestInput) ([]core.Suggestion, error) {
	// Build the prompt
	systemPrompt := buildSuggestSystemPrompt()
	userPrompt := buildSuggestUserPrompt(input.Context)

	// Make API request
	response, err := p.callOpenRouter(ctx, input.ModelID, systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("OpenRouter API call failed: %w", err)
	}

	// Parse suggestions from response
	suggestions := parseSuggestions(response, input.Context.Line)

	return suggestions, nil
}

// Explain generates an explanation using OpenRouter
func (p *OpenRouterProvider) Explain(ctx context.Context, input core.ExplainInput) (core.Explanation, error) {
	// Build the prompt
	systemPrompt := buildExplainSystemPrompt()
	userPrompt := buildExplainUserPrompt(input.Context)

	// Make API request
	response, err := p.callOpenRouter(ctx, input.ModelID, systemPrompt, userPrompt)
	if err != nil {
		return core.Explanation{}, fmt.Errorf("OpenRouter API call failed: %w", err)
	}

	// Parse explanation from response
	explanation := parseExplanation(response)

	return explanation, nil
}

// OpenRouter API types
type openRouterRequest struct {
	Model       string                   `json:"model"`
	Messages    []openRouterMessage      `json:"messages"`
	Temperature float64                  `json:"temperature,omitempty"`
	MaxTokens   int                      `json:"max_tokens,omitempty"`
}

type openRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// callOpenRouter makes an API call to OpenRouter
func (p *OpenRouterProvider) callOpenRouter(ctx context.Context, modelID, systemPrompt, userPrompt string) (string, error) {
	// Use configured model if not overridden
	if modelID == "" {
		modelID = p.profile.Model
	}

	// Build request
	reqBody := openRouterRequest{
		Model: modelID,
		Messages: []openRouterMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: p.profile.Temperature,
		MaxTokens:   p.profile.MaxTokens,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/chat/completions", strings.TrimSuffix(p.config.BaseURL, "/"))
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))

	// Set timeout
	client := &http.Client{
		Timeout: time.Duration(p.config.TimeoutMs) * time.Millisecond,
	}

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var apiResp openRouterResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API error
	if apiResp.Error != nil {
		return "", fmt.Errorf("API error: %s", apiResp.Error.Message)
	}

	// Extract content
	if len(apiResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	return apiResp.Choices[0].Message.Content, nil
}

package core

import (
	"context"

	"github.com/traves/linesense/internal/config"
)

// Suggestion represents a suggested command
type Suggestion struct {
	Command     string    `json:"command"`
	Risk        RiskLevel `json:"risk"` // "low" | "medium" | "high"
	Explanation string    `json:"explanation"`
	Source      string    `json:"source"` // "llm" | "preset"
}

// Explanation represents an explanation of a command
type Explanation struct {
	Summary string    `json:"summary"`
	Risk    RiskLevel `json:"risk"`
	Notes   []string  `json:"notes,omitempty"`
}

// Provider is the interface for AI providers
type Provider interface {
	Name() string
	Suggest(ctx context.Context, input SuggestInput) ([]Suggestion, error)
	Explain(ctx context.Context, input ExplainInput) (Explanation, error)
}

// SuggestInput contains input for suggestion generation
type SuggestInput struct {
	ModelID string           `json:"model_id"`
	Prompt  string           `json:"prompt"`
	Context *ContextEnvelope `json:"context"`
}

// ExplainInput contains input for command explanation
type ExplainInput struct {
	ModelID string           `json:"model_id"`
	Prompt  string           `json:"prompt"`
	Context *ContextEnvelope `json:"context"`
}

// Engine is the main engine for suggestions and explanations
type Engine struct {
	config   *config.Config
	provider Provider
}

// NewEngine creates a new engine instance
func NewEngine(cfg *config.Config, provider Provider) *Engine {
	return &Engine{
		config:   cfg,
		provider: provider,
	}
}

// Suggest generates command suggestions
func (e *Engine) Suggest(_ context.Context, shell, line, cwd string) ([]Suggestion, error) {
	// TODO: Build context envelope
	// TODO: Build prompt for suggestions
	// TODO: Call provider.Suggest
	// TODO: Apply safety filters
	// TODO: Return filtered suggestions
	panic("not implemented")
}

// Explain generates an explanation for a command
func (e *Engine) Explain(_ context.Context, shell, line, cwd string) (Explanation, error) {
	// TODO: Build context envelope
	// TODO: Build prompt for explanation
	// TODO: Call provider.Explain
	// TODO: Return explanation
	panic("not implemented")
}

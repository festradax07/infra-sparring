// Package llm abstracts the language model that generates challenges.
//
// The model only returns data: it never receives tools that execute commands
// (ADR 0004). Adapters for each provider implement Provider.
package llm

import "context"

// Provider sends one request to a language model and returns its answer.
type Provider interface {
	Generate(ctx context.Context, req Request) (Response, error)
}

// Request is a single, stateless call to the model.
type Request struct {
	// System sets the model's role and rules.
	System string
	// Prompt is the user message, for example the competency to generate.
	Prompt string
	// MaxTokens caps the length of the answer.
	MaxTokens int
}

// Response is the model's answer plus the tokens it cost (RNF-09).
type Response struct {
	Text  string
	Usage Usage
}

// Usage counts the tokens billed for one request.
type Usage struct {
	InputTokens  int
	OutputTokens int
}

// Total returns input plus output tokens.
func (u Usage) Total() int {
	return u.InputTokens + u.OutputTokens
}

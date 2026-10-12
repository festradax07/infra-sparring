// Package llmtest provides a scripted llm.Provider for unit tests, so they
// never call a real model.
package llmtest

import (
	"context"
	"errors"
	"sync"

	"github.com/festradax07/infra-sparring/internal/llm"
)

// Compile-time check that FakeProvider satisfies llm.Provider.
var _ llm.Provider = (*FakeProvider)(nil)

// ErrNoMoreResponses is returned when the fake runs out of queued responses.
var ErrNoMoreResponses = errors.New("llmtest: no more queued responses")

// FakeProvider returns queued responses in order and records every request.
type FakeProvider struct {
	mu        sync.Mutex
	responses []llm.Response
	requests  []llm.Request
}

// NewFakeProvider returns a fake that answers with responses, in order.
func NewFakeProvider(responses ...llm.Response) *FakeProvider {
	return &FakeProvider{responses: responses}
}

// Generate returns the next queued response, or ErrNoMoreResponses.
func (f *FakeProvider) Generate(_ context.Context, req llm.Request) (llm.Response, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.requests = append(f.requests, req)
	if len(f.responses) == 0 {
		return llm.Response{}, ErrNoMoreResponses
	}

	next := f.responses[0]
	f.responses = f.responses[1:]
	return next, nil
}

// Requests returns a copy of the requests received so far, in order.
func (f *FakeProvider) Requests() []llm.Request {
	f.mu.Lock()
	defer f.mu.Unlock()

	return append([]llm.Request(nil), f.requests...)
}

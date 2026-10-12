package llmtest

import (
	"context"
	"errors"
	"testing"

	"github.com/festradax07/infra-sparring/internal/llm"
)

func TestFakeProvider_ReturnsResponsesInOrderThenFails(t *testing.T) {
	ctx := context.Background()
	fake := NewFakeProvider(
		llm.Response{Text: "first"},
		llm.Response{Text: "second"},
	)

	for _, want := range []string{"first", "second"} {
		got, err := fake.Generate(ctx, llm.Request{Prompt: want})
		if err != nil {
			t.Fatalf("generate: %v", err)
		}
		if got.Text != want {
			t.Errorf("got %q, want %q", got.Text, want)
		}
	}

	if _, err := fake.Generate(ctx, llm.Request{}); !errors.Is(err, ErrNoMoreResponses) {
		t.Errorf("got error %v, want %v", err, ErrNoMoreResponses)
	}
	if n := len(fake.Requests()); n != 3 {
		t.Errorf("recorded %d requests, want 3", n)
	}
}

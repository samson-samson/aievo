package ollama

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTestClient(t *testing.T, opts ...Option) *LLM {
	t.Helper()
	var ollamaModel string
	if ollamaModel = os.Getenv("OLLAMA_TEST_MODEL"); ollamaModel == "" {
		t.Skip("OLLAMA_TEST_MODEL not set")
		return nil
	}

	opts = append([]Option{WithModel(ollamaModel)}, opts...)

	c, err := New(opts...)
	require.NoError(t, err)
	return c
}

func TestGenerateContent(t *testing.T) {
	t.Parallel()
	llm := newTestClient(t)

	ctx := context.Background()
	resp, err := llm.Generate(ctx, "你是谁?")
	require.NoError(t, err)
	require.NotEmpty(t, resp)
	t.Logf("response: %s", resp.Content)
}

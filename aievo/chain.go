package aievo

import (
	"context"

	"github.com/antgroup/aievo/llm"
)

type Handler func(ctx context.Context, prompt string, _ ...llm.GenerateOption) (string, error)

// Chain connect stage into one stage.
func Chain(chains ...Handler) Handler {
	return func(ctx context.Context, prompt string, opts ...llm.GenerateOption) (string, error) {
		var result string
		var err error
		for i := 0; i < len(chains); i++ {
			result, err = chains[i](ctx, prompt, opts...)
			if err != nil {
				return "", err
			}
		}
		return result, err
	}
}

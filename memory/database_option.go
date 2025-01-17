package memory

import (
	"context"

	"github.com/antgroup/aievo/schema"
)

// DatabaseOption is a function for creating new db storage
// with other default values.
type DatabaseOption func(b *Database)

// WithWindow is an option for providing max item for user message buf.
func WithWindow(window int) DatabaseOption {
	return func(b *Database) {
		b.window = window
	}
}

// WithSaveFunc is an option for providing the save func.
func WithSaveFunc(fun func(ctx context.Context,
	msg *schema.Message) error) DatabaseOption {
	return func(b *Database) {
		b.saveFunc = fun
	}
}

// WithLoadFunc is an option for providing the save func.
func WithLoadFunc(fun func(ctx context.Context) []schema.Message) DatabaseOption {
	return func(b *Database) {
		b.loadFunc = fun
	}
}

package memory

import (
	"context"

	"github.com/antgroup/aievo/schema"
)

type Database struct {
	buffer *Buffer

	window   int
	saveFunc func(ctx context.Context, msg *schema.Message) error
	loadFunc func(ctx context.Context) []schema.Message
}

func NewDatabaseMemory(opts ...DatabaseOption) *Database {
	dm := &Database{
		window: -1,
	}
	for _, opt := range opts {
		opt(dm)
	}
	dm.buffer = NewBufferWindowMemory(dm.window)
	return dm
}

func (d *Database) Load(ctx context.Context, filter func(index, consumption int, message schema.Message) bool) []schema.Message {
	d.load(ctx)
	return d.buffer.Load(ctx, filter)
}

func (d *Database) LoadNext(ctx context.Context, filter func(message schema.Message) bool) *schema.Message {
	d.load(ctx)
	return d.buffer.LoadNext(ctx, filter)
}

func (d *Database) Save(ctx context.Context, msg schema.Message) error {
	if err := d.save(ctx, &msg); err != nil {
		return err
	}
	return d.buffer.Save(ctx, msg)
}

func (d *Database) Clear(_ context.Context) error {
	return nil
}

func (d *Database) load(ctx context.Context) {
	if d.loadFunc != nil {
		_ = d.buffer.Clear(ctx)
		messages := d.loadFunc(ctx)
		d.buffer.Messages = append(d.buffer.Messages, messages...)
	}
}

func (d *Database) save(ctx context.Context, msg *schema.Message) error {
	if d.saveFunc != nil {
		return d.saveFunc(ctx, msg)
	}
	return nil
}

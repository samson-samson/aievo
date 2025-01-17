package memory

import (
	"context"

	"github.com/antgroup/aievo/schema"
)

type Buffer struct {
	Messages []schema.Message
	index    int
	window   int
}

func NewBufferMemory() *Buffer {
	return &Buffer{}
}

func NewBufferWindowMemory(window int) *Buffer {
	return &Buffer{window: window}
}

func (c *Buffer) Load(ctx context.Context, filter func(index, consumption int, message schema.Message) bool) []schema.Message {
	msgs := make([]schema.Message, 0, len(c.Messages))
	for i, message := range c.Messages {
		if filter == nil || filter(i, c.index, message) {
			msgs = append(msgs, message)
		}
	}
	if len(msgs) > c.window && c.window > 0 {
		msgs = msgs[len(msgs)-c.window:]
	}
	return msgs
}

func (c *Buffer) LoadNext(ctx context.Context, filter func(message schema.Message) bool) *schema.Message {
	if c.index >= len(c.Messages) {
		return nil
	}
	for ; c.index < len(c.Messages); c.index++ {
		if c.Messages[c.index].IsMsg() || c.Messages[c.index].IsEnd() ||
			c.Messages[c.index].IsCreative() {
			if c.Messages[c.index].Sender != c.Messages[c.index].Receiver {
				if filter != nil && !filter(c.Messages[c.index]) {
					return nil
				}
				c.index++
				return &c.Messages[c.index-1]
			}
		}
	}
	return nil
}

func (c *Buffer) Save(ctx context.Context, msg schema.Message) error {
	c.Messages = append(c.Messages, msg)
	return nil
}

func (c *Buffer) Clear(ctx context.Context) error {
	c.Messages = c.Messages[:0]
	return nil
}
